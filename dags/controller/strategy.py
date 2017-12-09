import h5py
from library import conf, console, tool
from strategy.macd import trend
from strategy.util import action


def wrap_box():
    """
    根据缠论聚合k线，判断中枢位置，并定位买卖点
    """
    return


def macd_divergence(code_list, ktype, start_date):
    """
    计算MACD与价格的背离情况，判断盘中波动走向
    macd代表柱子，正为红，负为绿
    dea代表黄线
    diff代表白线
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_STRATEGY,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_INDEX_PHASE
    )
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        # TODO, 立刻根据日期获取股票数据
        if f.get(code_group_path) is None:
            continue
        # 忽略停牌、退市、无法获取的情况
        if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
            continue
        ds_name = conf.HDF5_INDEX_MACD_TREND + "_" + ktype
        if f[code_prefix][code].get(ds_name) is None:
            continue
        if f[code_prefix][code].get(ktype) is None:
            continue

        trend_df = tool.df_from_dataset(f[code_prefix][code], ds_name, None)
        share_df = tool.df_from_dataset(f[code_prefix][code], ktype, None)

        trend_df[conf.HDF5_SHARE_DATE_INDEX] = trend_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        trend_df[action.INDEX_DIRECTION] = trend_df[action.INDEX_DIRECTION].str.decode("utf-8")
        trend_df[action.INDEX_ACTION] = trend_df[action.INDEX_ACTION].str.decode("utf-8")
        trend_df[action.INDEX_STATUS] = trend_df[action.INDEX_STATUS].str.decode("utf-8")
        trend_df[action.INDEX_PHASE_STATUS] = trend_df[action.INDEX_PHASE_STATUS].str.decode("utf-8")
        share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        share_df = share_df[["high", "low", "date"]]

        trend.get_macd_phases(trend_df, share_df)
    console.write_tail()
    f.close()
    return
