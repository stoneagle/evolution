import h5py
from library import conf, console, tool
from strategy.macd import trend


def wrap_k_line():
    """
    将k线整理成缠论模式
    """
    return


def wrap_k_box():
    """
    根据缠论k先获取箱体范围
    """
    return


def macd_all_divergence(code_list, start_date):
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_TACTICS,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_SHARE_DETAIL
    )
    for code in code_list:

        for ktype in conf.HDF5_SHARE_KTYPE:
            macd_code_divergence(f, code, ktype, start_date)
    console.write_tail()
    f.close()
    return


def macd_code_divergence(f, code, ktype, start_date):
    """
    计算MACD的背离指标
    macd代表柱子，正为红，负为绿
    dea代表黄线
    diff代表白线
    """
    code_prefix = code[0:3]
    code_group_path = '/' + code_prefix + '/' + code
    # TODO, 立刻根据日期获取股票数据
    if f.get(code_group_path) is None:
        return

    # 忽略停牌、退市、无法获取的情况
    if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
        return

    ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
    if f[code_prefix][code].get(ds_name) is None:
        return

    index_df = tool.df_from_dataset(f[code_prefix][code], ds_name, None)
    index_df[conf.HDF5_SHARE_DATE_INDEX] = index_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    index_df = index_df[index_df["dif"].notnull()]
    index_df = index_df.set_index(conf.HDF5_SHARE_DATE_INDEX)[["dif", "dea", "macd", "close"]]

    # 如果数据集过少则直接返回
    # TODO 异常报错
    if len(index_df) <= 3:
        return

    return trend.get_macd_trends(index_df)
