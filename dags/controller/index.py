import h5py
import talib
import pandas as pd
import numpy as np
from library import conf, console, tool


def all_index(gem_flag, start_date):
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_INDEX,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_SHARE_DETAIL
    )
    for code_prefix in f:
        if gem_flag is True and code_prefix == "300":
            continue
        for code in f[code_prefix]:
            # 忽略停牌、退市、无法获取的情况
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
                continue
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
                continue
            for ktype in conf.HDF5_SHARE_KTYPE:
                index_df = one_index(f, code, ktype, start_date)
                ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if start_date is None:
                    tool.delete_dataset(f[code_prefix][code], ds_name)
                tool.merge_df_dataset(f[code_prefix][code], ds_name, index_df)
    console.write_tail()
    f.close()
    return


def one_index(f, code, ktype, start_date):
    code_prefix = code[0:3]
    code_group_path = '/' + code_prefix + '/' + code
    if f.get(code_group_path) is not None and f[code_prefix][code].get(ktype) is None:
        return
    df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
    df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    df = df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    # 计算macd
    macd = talib.MACD(df["close"].values, fastperiod=12, slowperiod=26, signalperiod=9)
    # 计算均值
    df_index = pd.DataFrame(np.column_stack(macd), index=df.index, columns=conf.HDF5_INDEX_COLUMN)
    df_index["macd"] = df_index["macd"] * 2
    for i in [5, 10, 30, 60, 240, 480, 365]:
        df_index["ma_" + str(i)] = df["close"].rolling(window=i).mean()
    if start_date is not None:
        df_index = df_index.ix[start_date:]
    return df_index
