import h5py
import talib
import pandas as pd
import numpy as np
from library import conf, console, tool


def all_share(gem_flag, start_date):
    """
    获取所有股票的macd与所处均线等指标
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code_prefix in f:
        if gem_flag is True and code_prefix == "300":
            continue
        for code in f[code_prefix]:
            console.write_head(
                conf.HDF5_OPERATE_INDEX,
                conf.HDF5_RESOURCE_TUSHARE,
                code
            )
            # 忽略停牌、退市、无法获取的情况
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
                continue
            if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
                continue

            code_group_path = '/' + code_prefix + '/' + code
            for ktype in conf.HDF5_SHARE_KTYPE:
                if f.get(code_group_path) is None or f[code_prefix][code].get(ktype) is None:
                    continue
                df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
                index_df = one_df(df, start_date)
                ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if start_date is None:
                    tool.delete_dataset(f[code_prefix][code], ds_name)
                tool.merge_df_dataset(f[code_prefix][code], ds_name, index_df.reset_index())
            console.write_tail()
    f.close()
    return


def all_classify(classify_list, start_date):
    """
    获取所有分类的macd与所处均线等指标(依赖分类数据聚合)
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_ARRANGE,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )

            for ktype in conf.HDF5_SHARE_KTYPE:
                ds_name = conf.HDF5_CLASSIFY_DS_DETAIL + "_" + ktype
                if f[ctype][classify_name].get(ds_name) is None:
                    continue

                df = tool.df_from_dataset(f[ctype][classify_name], ds_name, None)
                df["close"] = df["close"].apply(lambda x: round(x, 2))
                index_df = one_df(df, start_date)
                index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if start_date is None:
                    tool.delete_dataset(f[ctype][classify_name], index_ds_name)
                tool.merge_df_dataset(f[ctype][classify_name], index_ds_name, index_df.reset_index())
            console.write_tail()
            break
    f.close()
    return


def all_index(start_date):
    """
    处理所有指数的均线与macd
    """
    f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
    for code in f:
        console.write_head(
            conf.HDF5_OPERATE_ARRANGE,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        for ktype in conf.HDF5_SHARE_KTYPE:
            if f[code].get(ktype) is None:
                continue
            df = tool.df_from_dataset(f[code], ktype, None)
            index_df = one_df(df, start_date)
            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if start_date is None:
                tool.delete_dataset(f[code], index_ds_name)
            tool.merge_df_dataset(f[code], index_ds_name, index_df.reset_index())
        console.write_tail()
    f.close()
    return


def one_df(df, start_date):
    """
    处理单组df的均线与macd
    """
    df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    df = df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    # 计算macd
    macd = talib.MACD(df["close"].values, fastperiod=12, slowperiod=26, signalperiod=9)

    # 计算均值
    tmp_mean_dict = dict()
    for i in range(0, 481):
        tmp_mean_dict[i] = df["close"].rolling(window=i).mean()

    df_index = pd.DataFrame(np.column_stack(macd), index=df.index, columns=conf.HDF5_INDEX_COLUMN)
    df_index["close"] = df["close"]
    df_index["macd"] = df_index["macd"] * 2

    if start_date is not None:
        df_index = df_index.ix[start_date:]

    # 计算close价格所处均线位置
    for index, row in df_index.iterrows():
        # 先确认均线是上行还是下行
        ma5_price = tmp_mean_dict[5].loc[index]
        ma10_price = tmp_mean_dict[10].loc[index]
        if np.isnan(ma5_price) or np.isnan(ma10_price):
            df_index.loc[index, "ma_border"] = np.NaN
            continue

        if ma5_price >= ma10_price:
            up_flag = True
        else:
            up_flag = False

        border = 0
        # 1. 如果下行，从小于找起直到出现大于;如果上行，从大于找起直到出现小于
        for i in [30, 60, 240, 480]:
            compare_price = tmp_mean_dict[i].loc[index]
            if compare_price is np.NaN:
                break

            if (ma5_price > compare_price and up_flag is False) or (ma5_price < compare_price and up_flag is True):
                border = i
                break

        # 2. 根据边界均线，按比例调整，找出离close最贴近的均线
        if border != 0:
            border_price = tmp_mean_dict[border].loc[index]
            df_index.loc[index, "ma_border"] = border
            min_diff = 0
            min_diff_ma_num = 0
            if (up_flag is False and row["close"] > border_price) or (up_flag is True and row["close"] < border_price):
                for i in range(0, border - 5):
                    tmp_mean = tmp_mean_dict[border - i]
                    dynamic_price = tmp_mean.loc[index, ]

                    # 寻找close最贴近的均线
                    if min_diff == 0 or abs(row["close"] - dynamic_price) <= min_diff:
                        min_diff = abs(row["close"] - dynamic_price)
                        min_diff_ma_num = border - i
                df_index.loc[index, "ma_border"] = min_diff_ma_num
        else:
            df_index.loc[index, "ma_border"] = 5
    df_index = df_index.drop("close", axis=1)
    return df_index
