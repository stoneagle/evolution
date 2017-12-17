import h5py
from library import conf, console, tool
from strategy import macd, kline


def filter_share(code_list, start_date):
    """
    整理筛选的股票缠论k线
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_WRAP,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_INDEX_WRAP
    )
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            continue
        # 忽略停牌、退市、无法获取的情况
        if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
            continue

        for ktype in conf.HDF5_SHARE_WRAP_KTYPE:
            ds_name = ktype
            if f[code_prefix][code].get(ds_name) is None:
                continue
            share_df = tool.df_from_dataset(f[code_prefix][code], ds_name, None)
            wrap_df = one_df(share_df)
            if wrap_df is not None:
                ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
                if f[code_prefix][code].get(ds_name) is not None:
                    tool.delete_dataset(f[code_prefix][code], ds_name)
                tool.create_df_dataset(f[code_prefix][code], ds_name, wrap_df)
                console.write_exec()
            else:
                console.write_pass()
    console.write_blank()
    console.write_tail()
    f.close()
    return


def all_classify(classify_list, init_flag=True):
    """
    整理分类的缠论k线
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f_classify[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_WRAP,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )
            for ktype in conf.HDF5_SHARE_WRAP_KTYPE:
                ds_name = conf.HDF5_CLASSIFY_DS_DETAIL + "_" + ktype
                if f_classify[ctype][classify_name].get(ds_name) is None:
                    continue
                share_df = tool.df_from_dataset(f_classify[ctype][classify_name], ds_name, None)
                wrap_df = one_df(share_df)
                wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
                if init_flag is True:
                    tool.delete_dataset(f_classify[ctype][classify_name], wrap_ds_name)
                if wrap_df is not None:
                    tool.merge_df_dataset(f_classify[ctype][classify_name], wrap_ds_name, wrap_df)
            console.write_tail()
    f_classify.close()
    f.close()
    return


def all_index(init_flag=True):
    """
    整理指数的缠论k线
    """
    f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
    for code in f:
        console.write_head(
            conf.HDF5_OPERATE_WRAP,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        for ktype in conf.HDF5_SHARE_WRAP_KTYPE:
            if f[code].get(ktype) is None:
                continue
            share_df = tool.df_from_dataset(f[code], ktype, None)
            wrap_df = one_df(share_df)
            wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
            if init_flag is True:
                tool.delete_dataset(f[code], wrap_ds_name)
            if wrap_df is not None:
                tool.merge_df_dataset(f[code], wrap_ds_name, wrap_df)
        console.write_tail()
    f.close()
    return


def one_df(share_df, index_flag=False):
    """
    整理某只股票的缠论k线
    """
    share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    share_df = share_df[["date", "high", "low", "open", "close"]]

    # 如果数据集过少则直接返回
    if len(share_df) <= 3:
        return None
    wrap_df = kline.trans_wrap(share_df, False)
    wrap_df = wrap_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    index_df = macd.value(wrap_df)
    wrap_df = wrap_df.merge(index_df, left_index=True, right_index=True, how='left')
    return wrap_df.reset_index()
