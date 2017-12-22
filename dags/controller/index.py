import h5py
from library import conf, console, tool
from quota import macd, index


def all_share(omit_list, init_flag=True):
    """
    获取所有股票的macd与所处均线等指标
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code_prefix in f:
        if code_prefix in omit_list:
            continue
        console.write_head(
            conf.HDF5_OPERATE_INDEX,
            conf.HDF5_RESOURCE_TUSHARE,
            code_prefix
        )
        for code in f[code_prefix]:
            # 忽略停牌、退市、无法获取的情况
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
                continue
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
                continue

            code_group_path = '/' + code_prefix + '/' + code
            for ktype in conf.HDF5_SHARE_KTYPE:
                try:
                    if f.get(code_group_path) is None or f[code_prefix][code].get(ktype) is None:
                        console.write_msg(code + "-" + ktype + "的detail不存在")
                        continue
                    df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
                    index_df = one_df(df, init_flag)
                    ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                    if init_flag is True:
                        tool.delete_dataset(f[code_prefix][code], ds_name)
                    tool.merge_df_dataset(f[code_prefix][code], ds_name, index_df.reset_index())
                except Exception as er:
                    print(str(er))
            console.write_exec()
        console.write_blank()
        console.write_tail()
    f.close()
    return


def all_classify(classify_list, init_flag=True):
    """
    获取所有分类的macd与所处均线等指标(依赖分类数据聚合)
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_INDEX,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )
            for ktype in conf.HDF5_SHARE_KTYPE:
                ds_name = ktype
                if f[ctype][classify_name].get(ds_name) is None:
                    console.write_msg(classify_name + "分类聚合detail不存在")
                    continue

                df = tool.df_from_dataset(f[ctype][classify_name], ds_name, None)
                df["close"] = df["close"].apply(lambda x: round(x, 2))
                try:
                    index_df = one_df(df, init_flag, True)
                except Exception as er:
                    console.write_msg("[" + classify_name + "]" + str(er))
                    continue
                index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if init_flag is True:
                    tool.delete_dataset(f[ctype][classify_name], index_ds_name)
                tool.merge_df_dataset(f[ctype][classify_name], index_ds_name, index_df.reset_index())
            console.write_tail()
    f.close()
    return


def all_index(init_flag=True):
    """
    处理所有指数的均线与macd
    """
    f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
    for code in f:
        console.write_head(
            conf.HDF5_OPERATE_INDEX,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        for ktype in conf.HDF5_SHARE_KTYPE:
            if f[code].get(ktype) is None:
                console.write_msg(code + "-" + ktype + "的detail不存在")
                continue
            df = tool.df_from_dataset(f[code], ktype, None)
            index_df = one_df(df, init_flag)
            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if init_flag is True:
                tool.delete_dataset(f[code], index_ds_name)
            tool.merge_df_dataset(f[code], index_ds_name, index_df.reset_index())
        console.write_tail()
    f.close()
    return


def one_df(detail_df, init_flag, ma_pos_flag=False):
    """
    处理单组df的均线与macd
    """
    detail_df[conf.HDF5_SHARE_DATE_INDEX] = detail_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    detail_df = detail_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    if init_flag is not True:
        detail_df = detail_df.tail(50)
    index_df = macd.value(detail_df)
    if ma_pos_flag is True:
        index_df["close"] = detail_df["close"]
        index_df = index.mean_position(index_df, "close")
    else:
        index_df[index.INDEX_MA_BORDER] = 0
    return index_df
