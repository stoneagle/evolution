import h5py
from library import conf, console, tool, tradetime
from quota import macd, kline
from quota.util import action
import warnings
warnings.filterwarnings("ignore")
INDEX_MACD_TREND_STATUS = "_m_status"
INDEX_MACD_TREND_COUNT = "_m_count"
INDEX_MACD_TREND_CROSS = "_m_cross"
INDEX_MACD_TREND_VALUE = "_m_macd"
INDEX_MACD_TREND_DIF = "_m_dif"
INDEX_MACD_DIVERSE_COUNT = "_d_count"
INDEX_MACD_DIVERSE_PRICE_MIN = "_d_p_min"
INDEX_MACD_DIVERSE_PRICE_MAX = "_d_p_max"
INDEX_MACD_DIVERSE_PRICE_START = "_m_p_start"


def share_filter(omit_list):
    """
    每日选股，筛选出存在背离的股票
    """
    # 筛选记录内容如下:
    # 1. 月线macd趋势
    # 2. 周线macd趋势
    # 3. 日线macd趋势，是否背离，数值差距
    # 4. 30min的macd趋势，是否背离，数值差距，震荡中枢数量
    # 5. 5min的macd趋势，是否背离，数值差距，震荡中枢数量
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    filter_df = tool.init_empty_df(_ini_filter_columns())
    for code_prefix in f:
        if code_prefix in omit_list:
            continue
        console.write_head(
            conf.HDF5_OPERATE_SCREEN,
            conf.HDF5_RESOURCE_TUSHARE,
            code_prefix
        )
        for code in f[code_prefix]:
            code_group_path = '/' + code_prefix + '/' + code
            if f.get(code_group_path) is None:
                console.write_blank()
                console.write_msg(code + "的tushare数据不存在")
                continue

            # 忽略停牌、退市、无法获取的情况
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
                console.write_blank()
                console.write_msg(code + "已退市或停牌")
                continue

            try:
                code_dict = _daily_code(f, code)
                if code_dict is None:
                    console.write_pass()
                    continue
                else:
                    console.write_exec()
                    filter_df = filter_df.append(code_dict, ignore_index=True)
            except Exception as er:
                console.write_msg("[" + code + "]" + str(er))
        console.write_blank()
        console.write_tail()
    f.close()
    f_screen = h5py.File(conf.HDF5_FILE_SCREEN, 'a')
    if f_screen.get(conf.SCREEN_SHARE_FILTER) is None:
        f_screen.create_group(conf.SCREEN_SHARE_FILTER)
    today_str = tradetime.get_today()
    tool.delete_dataset(f_screen[conf.SCREEN_SHARE_FILTER], today_str)
    tool.merge_df_dataset(f_screen[conf.SCREEN_SHARE_FILTER], today_str, filter_df)
    f_screen.close()
    return


def mark_grade(today_str=None):
    """
    对筛选结果进行打分
    """
    console.write_head(
        conf.HDF5_OPERATE_SCREEN,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.SCREEN_SHARE_GRADE
    )
    f = h5py.File(conf.HDF5_FILE_SCREEN, 'a')
    f_share = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    if today_str is None:
        today_str = tradetime.get_today()
    if f[conf.SCREEN_SHARE_FILTER].get(today_str) is None:
        console.write_msg(today_str + "个股筛选结果不存在")
        return
    grade_df = tool.init_empty_df(["code", "status", "d_price_space", "d_price_per", "30_price_space", "30_price_per", "d_macd", "30_macd"])
    screen_df = tool.df_from_dataset(f[conf.SCREEN_SHARE_FILTER], today_str, None)
    screen_df["d_m_status"] = screen_df["d_m_status"].str.decode("utf-8")
    screen_df["w_m_status"] = screen_df["w_m_status"].str.decode("utf-8")
    screen_df["m_m_status"] = screen_df["m_m_status"].str.decode("utf-8")
    screen_df["code"] = screen_df["code"].str.decode("utf-8")
    for index, row in screen_df.iterrows():
        code = row["code"]
        grade_dict = dict()
        grade_dict["code"] = code
        grade_dict["status"] = 0
        grade_dict["status"] += _macd_status_grade(row["d_m_status"])
        grade_dict["status"] += _macd_status_grade(row["w_m_status"])
        grade_dict["status"] += _macd_status_grade(row["m_m_status"])
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        for ktype in ["D", "30"]:
            detail_ds_name = ktype
            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if f_share[code_group_path].get(detail_ds_name) is None:
                console.write_msg(code + "的detail数据不存在")
                continue
            if f_share[code_group_path].get(index_ds_name) is None:
                console.write_msg(code + "的index数据不存在")
                continue
            detail_df = tool.df_from_dataset(f_share[code_group_path], detail_ds_name, None)
            index_df = tool.df_from_dataset(f_share[code_group_path], index_ds_name, None)
            latest_price = detail_df.tail(1)["close"].values[0]
            latest_macd = index_df.tail(1)["macd"].values[0]
            diverse_price_start = row[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_START]
            if diverse_price_start == 0:
                grade_dict[str.lower(ktype) + "_price_space"] = 0
                grade_dict[str.lower(ktype) + "_price_per"] = 0
            else:
                grade_dict[str.lower(ktype) + "_price_space"] = round(diverse_price_start - latest_price, 2)
                grade_dict[str.lower(ktype) + "_price_per"] = round(grade_dict[str.lower(ktype) + "_price_space"] * 100 / diverse_price_start, 2)
            grade_dict[str.lower(ktype) + "_macd"] = latest_macd
        grade_df = grade_df.append(grade_dict, ignore_index=True)
    if f.get(conf.SCREEN_SHARE_GRADE) is None:
        f.create_group(conf.SCREEN_SHARE_GRADE)
    tool.delete_dataset(f[conf.SCREEN_SHARE_GRADE], today_str)
    tool.merge_df_dataset(f[conf.SCREEN_SHARE_GRADE], today_str, grade_df)
    f_share.close()
    f.close()
    console.write_tail()
    return


def _daily_code(f, code):
    MACD_DIVERSE_LIMIT = 5

    code_prefix = code[0:3]
    code_dict = dict()
    code_dict["code"] = code
    omit_flag = False
    for ktype in ["D", "W", "M", "30", "5"]:
        if f[code_prefix][code].get(ktype) is None:
            console.write_blank()
            console.write_msg(code + "阶段" + ktype + "的detail数据不存在")
            return None
        share_df = tool.df_from_dataset(f[code_prefix][code], ktype, None)

        index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
        if f[code_prefix][code].get(index_ds_name) is None:
            console.write_blank()
            console.write_msg(code + "阶段" + ktype + "的macd与均线数据不存在")
            return None
        index_df = tool.df_from_dataset(f[code_prefix][code], index_ds_name, None)
        share_df = share_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')
        share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        # 检查macd的趋势
        if ktype in ["D", "M", "W"]:
            code_dict = _filter_macd_trend(share_df.tail(50), code_dict, ktype)
            if code_dict is None:
                omit_flag = True
                break

        # 检查macd的背离
        if ktype in ["D", "30"]:
            code_dict = _filter_macd_diverse(share_df.tail(100), code_dict, ktype)

        # 检查震荡中枢
        # if ktype in ["30", "5"]:
        #     code_dict = _filter_wrap_central(share_df, code_dict, ktype)

    # 如果日线不存在macd的上涨趋势，则忽略该股票
    if omit_flag is True:
        return None

    # 如果日线、30min一个阶段内都不存在背离，忽略
    if code_dict["30" + INDEX_MACD_DIVERSE_COUNT] <= MACD_DIVERSE_LIMIT and code_dict["d" + INDEX_MACD_DIVERSE_COUNT] <= MACD_DIVERSE_LIMIT:
        return None
    return code_dict


def _filter_macd_trend(share_df, code_dict, ktype):
    trend_df = macd.trend(share_df)
    tail_row = trend_df.tail(1)
    tail_macd = tail_row["macd"].values[0]
    tail_dif = tail_row["dif"].values[0]
    tail_phase_status = tail_row[action.INDEX_PHASE_STATUS].values[0]
    if ktype == "D":
        # 忽略日线macd下降
        if tail_phase_status == action.STATUS_DOWN:
            return None
        # 忽略dif已上穿零轴的情况
        elif tail_dif > 0:
            return None
        # 忽略macd已启动的股票
        elif tail_macd > 0.1:
            return None
        # 忽略下降中shake的股票
        elif tail_phase_status == action.STATUS_SHAKE:
            for i in range(2, 20):
                pre_phase_status = trend_df.tail(i)[action.INDEX_PHASE_STATUS].values[0]
                if pre_phase_status != action.STATUS_SHAKE:
                    break
            if tail_phase_status == action.STATUS_DOWN:
                return None
    code_dict[str.lower(ktype) + INDEX_MACD_TREND_STATUS] = tail_phase_status
    code_dict[str.lower(ktype) + INDEX_MACD_TREND_COUNT] = tail_row[action.INDEX_TREND_COUNT].values[0]
    code_dict[str.lower(ktype) + INDEX_MACD_TREND_CROSS] = tail_row[macd.INDEX_CROSS_COUNT].values[0]
    code_dict[str.lower(ktype) + INDEX_MACD_TREND_VALUE] = tail_macd
    code_dict[str.lower(ktype) + INDEX_MACD_TREND_DIF] = tail_dif
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
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_COUNT] = diverse_count
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_MIN] = lp_diverse_df["close"].min()
        # TODO 获取背离开始期间，价格的最大值
        # code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_MAX] = lp_diverse_df["close"].max()
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_START] = lp_diverse_df.head(1)["close"].values[0]
    else:
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_COUNT] = 0
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_MIN] = 0
        code_dict[str.lower(ktype) + INDEX_MACD_DIVERSE_PRICE_START] = 0
    return code_dict


def _filter_wrap_central(share_df, code_dict, ktype):
    central_df = kline.central(share_df)
    print(central_df)
    return code_dict


def watch(code_list):
    """
    监听筛选出的股票
    """
    return


def _ini_filter_columns():
    columns = []
    columns.append("code")
    for ktype in ["D", "W", "M"]:
        ktype = str.lower(ktype)
        columns.append(ktype + INDEX_MACD_TREND_STATUS)
        columns.append(ktype + INDEX_MACD_TREND_VALUE)
        columns.append(ktype + INDEX_MACD_TREND_DIF)
        columns.append(ktype + INDEX_MACD_TREND_COUNT)
        columns.append(ktype + INDEX_MACD_TREND_CROSS)
    for ktype in ["30", "D"]:
        ktype = str.lower(ktype)
        columns.append(ktype + INDEX_MACD_DIVERSE_COUNT)
        columns.append(ktype + INDEX_MACD_DIVERSE_PRICE_MIN)
        columns.append(ktype + INDEX_MACD_DIVERSE_PRICE_START)
    return columns


def _macd_status_grade(status):
    if status == action.STATUS_UP:
        ret = 2
    elif status == action.STATUS_SHAKE:
        ret = 1
    else:
        ret = 0
    return ret
