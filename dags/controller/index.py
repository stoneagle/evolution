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
                tool.merge_df_dataset(f[code_prefix][code], ds_name, index_df.reset_index())
    console.write_tail()
    f.close()
    return


def one_index(f, code, ktype, start_date):
    code_prefix = code[0:3]
    code_group_path = '/' + code_prefix + '/' + code
    if f.get(code_group_path) is None or f[code_prefix][code].get(ktype) is None:
        return
    df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
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
    for i in [5, 10, 30, 60, 240, 480, 365]:
        df_index["ma_" + str(i)] = tmp_mean_dict[i]
    if start_date is not None:
        df_index = df_index.ix[start_date:]

    # 计算close价格所处均线位置
    for index, row in df_index.iterrows():
        # 先确认均线是上行还是下行
        ma5_price = row["ma_5"]
        ma10_price = row["ma_10"]
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
            compare_price = row["ma_" + str(i)]
            if compare_price is np.NaN:
                break

            if (ma5_price > compare_price and up_flag is False) or (ma5_price < compare_price and up_flag is True):
                border = i
                break

        # 2. 根据边界均线，按比例调整，找出离close最贴近的均线
        if border != 0:
            border_price = row["ma_" + str(border)]
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
    return df_index
