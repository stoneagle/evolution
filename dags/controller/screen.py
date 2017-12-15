import h5py
from library import conf, console, tool
from strategy import macd, kline
from strategy.util import action
import warnings
warnings.filterwarnings("ignore")


def daily_filter(code_list, gem_flag):
    """
    每日选股，对个股进行排名与打分
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    filter_dict = dict()
    for code in code_list:
        console.write_head(
            conf.HDF5_OPERATE_STRATEGY,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        code_prefix = code[0:3]
        # 判断是否跳过创业板
        if gem_flag is True and code_prefix == "300":
            continue

        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            console.write_msg(code + "的tushare数据不存在")
            continue

        # 忽略停牌、退市、无法获取的情况
        if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
            console.write_msg(code + "已退市或停牌")
            continue

        code_dict = dict()
        omit_flag = False
        for ktype in ["D", "W", "M", "30", "5"]:
            if f[code_prefix][code].get(ktype) is None:
                console.write_msg(code + "阶段" + ktype + "的股票数据不存在")
                continue
            share_df = tool.df_from_dataset(f[code_prefix][code], ktype, None)

            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if f[code_prefix][code].get(index_ds_name) is None:
                console.write_msg(code + "阶段" + ktype + "的macd与均线数据不存在")
                continue
            index_df = tool.df_from_dataset(f[code_prefix][code], index_ds_name, None)
            share_df = share_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')
            share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            # 检查macd的趋势
            # if ktype in ["D", "M", "W"]:
            #     try:
            #         code_dict = _filter_macd_trend(share_df.tail(50), code_dict, ktype)
            #     except Exception as er:
            #         console.write_msg(str(er))
            #         omit_flag = True
            #         break

            # 检查macd的背离
            # if ktype in ["D", "30"]:
            #     code_dict = _filter_macd_diverse(share_df.tail(100), code_dict, ktype)

            # 检查震荡中枢
            # if ktype in ["30", "5"]:
            if ktype in ["30"]:
                code_dict = _filter_wrap_central(share_df, code_dict, ktype)

        if omit_flag is True:
            continue

        # 如果日线、30min一个阶段内都不存在背离，忽略
        if code_dict["30_diverse_count"] <= 3 and code_dict["D_diverse_count"] <= 3:
            continue
        filter_dict[code] = code_dict
        console.write_tail()
    print(filter_dict)
    f.close()
    # 筛选记录内容如下:
    # 1. 月线macd趋势
    # 2. 周线macd趋势
    # 3. 日线macd趋势，是否背离，数值差距
    # 4. 30min的macd趋势，是否背离，数值差距，震荡中枢数量
    # 5. 5min的macd趋势，是否背离，数值差距，震荡中枢数量
    return


def _filter_macd_trend(share_df, code_dict, ktype):
    result = macd.trend(share_df.tail(50))
    tail_row = result.tail(1)
    tail_phase_status = tail_row[action.INDEX_PHASE_STATUS].values[0]
    # 忽略日线macd下降，或者在下降中shake的股票
    if ktype == "D":
        if tail_phase_status == action.STATUS_DOWN:
            raise Exception("macd处于下降通道")
        elif tail_phase_status == action.STATUS_SHAKE:
            for i in range(2, 20):
                pre_phase_status = result.tail(i)[action.INDEX_PHASE_STATUS].values[0]
                if pre_phase_status != action.STATUS_SHAKE:
                    break
            if tail_phase_status == action.STATUS_DOWN:
                raise Exception("macd处于下降通道")
    code_dict[ktype + "_phase_status"] = tail_phase_status
    code_dict[ktype + "_trend_count"] = tail_row[action.INDEX_TREND_COUNT].values[0]
    code_dict[ktype + "_cross_count"] = tail_row[macd.INDEX_CROSS_COUNT].values[0]
    code_dict[ktype + "_macd"] = tail_row["macd"].values[0]
    code_dict[ktype + "_dif"] = tail_row["dif"].values[0]
    return code_dict


def _filter_macd_diverse(share_df, code_dict, ktype):
    diverse_df = macd.price_diverse(share_df)
    latest_macd = diverse_df.iloc[-1]["macd"]
    # 根据当前macd所处零轴位置，判断该周期下，是否存在背离
    for i in range(2, len(diverse_df)):
        pre_macd = diverse_df.iloc[-i]["macd"]
        if (latest_macd > 0 and pre_macd <= 0) or (latest_macd <= 0 and pre_macd > 0):
            border = i
            break
    lp_df = diverse_df.tail(border)
    lp_df = lp_df[diverse_df[action.INDEX_PHASE_STATUS] == action.STATUS_UP]
    lp_diverse_df = lp_df[lp_df[macd.INDEX_DIVERSE] == 1]
    diverse_count = lp_diverse_df[macd.INDEX_DIVERSE].count()
    if diverse_count > 0:
        code_dict[ktype + "_diverse_count"] = diverse_count
        code_dict[ktype + "_diverse_price_min"] = lp_diverse_df["close"].min()
        code_dict[ktype + "_diverse_price_start"] = lp_diverse_df.head(1)["close"].values[0]
    else:
        code_dict[ktype + "_diverse_count"] = 0
        code_dict[ktype + "_diverse_price_min"] = 0
        code_dict[ktype + "_diverse_price_start"] = 0
    return code_dict


def _filter_wrap_central(share_df, code_dict, ktype):
    kline.central(share_df)
    return code_dict


def watch(code_list):
    """
    监听筛选出的股票
    """
    return
