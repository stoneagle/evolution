from strategy.util import action, phase
from library import conf
import talib
import pandas as pd
import numpy as np
INDEX_CROSS_COUNT = "cross"
INDEX_PRICE_DIFF = "price_diff"
INDEX_MACD_DIFF = "macd_diff"
INDEX_DIVERSE = "diverse"


def value(detail_df):
    # 计算macd
    macd = talib.MACD(detail_df["close"].values, fastperiod=12, slowperiod=26, signalperiod=9)
    index_df = pd.DataFrame(np.column_stack(macd), index=detail_df.index, columns=conf.HDF5_INDEX_COLUMN)
    index_df["macd"] = index_df["macd"] * 2
    return index_df


def trend(index_df):
    """
    获取macd的趋势状况
    """
    index_df = index_df[index_df["dif"].notnull()]
    index_df = index_df.set_index(conf.HDF5_SHARE_DATE_INDEX)[["dif", "dea", "macd", "close"]]
    # 如果数据集过少则直接返回
    if len(index_df) <= 3:
        return None
    trend_df = action.Action(index_df.reset_index()).all(date_column=conf.HDF5_SHARE_DATE_INDEX, value_column="macd")
    trend_df = trend_df.drop(action.INDEX_VALUE, axis=1)
    index_df = index_df.reset_index()
    trend_df = trend_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='left')
    # 计算阶段中，dif与dea交叉次数
    for index, row in trend_df.iterrows():
        # 如果dif与dea的差值出现方向变化，记作交叉，当交叉点数值穿越0轴时重新计数
        diff = row["dif"] - row["dea"]
        if index == 0:
            pre_diff = diff
            pre_dif = row["dif"]
            trend_df.loc[index, INDEX_CROSS_COUNT] = 0
            continue

        pre_dif = trend_df.iloc[index - 1]["dif"]
        pre_dea = trend_df.iloc[index - 1]["dea"]
        pre_diff = pre_dif - pre_dea
        # 零轴上，下沉零轴下
        if (pre_diff >= 0 and diff < 0) or (pre_diff < 0 and diff >= 0):
            if (pre_dif > 0 and row["dif"] <= 0) or (pre_dif < 0 and row["dif"] >= 0):
                trend_df.loc[index, INDEX_CROSS_COUNT] = 1
            else:
                trend_df.loc[index, INDEX_CROSS_COUNT] = trend_df.iloc[index - 1][INDEX_CROSS_COUNT] + 1
        else:
            if (pre_dif > 0 and row["dif"] <= 0) or (pre_dif < 0 and row["dif"] >= 0):
                trend_df.loc[index, INDEX_CROSS_COUNT] = 0
            else:
                trend_df.loc[index, INDEX_CROSS_COUNT] = trend_df.iloc[index - 1][INDEX_CROSS_COUNT]
    return trend_df


def price_diverse(share_df):
    """
    获取价格背离情况
    当前点与起点相比，价格与macd的背离情况，包括波动幅度与比例
    """
    trend_df = trend(share_df)
    if trend is None:
        return None
    trend_df = trend_df.merge(share_df[["date", "high", "low", "open"]], left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='left')
    # 计算阶段中，macd、price与起点的差值
    phase_start_index = 0
    for index, row in trend_df.iterrows():
        if index == 0 or phase_start_index == 0:
            phase_start_index = index

        phase_start_price = trend_df.loc[phase_start_index, "close"]
        phase_start_macd = trend_df.loc[phase_start_index, "macd"]
        price_diff = row["close"] - phase_start_price
        macd_diff = row["macd"] - phase_start_macd
        trend_df.loc[index, INDEX_PRICE_DIFF] = price_diff
        trend_df.loc[index, INDEX_MACD_DIFF] = macd_diff
        if (macd_diff >= 0 and price_diff < 0) or (macd_diff < 0 and price_diff >= 0):
            trend_df.loc[index, INDEX_DIVERSE] = True
        else:
            trend_df.loc[index, INDEX_DIVERSE] = False

        if row[action.INDEX_PHASE_STATUS] != trend_df.loc[phase_start_index, action.INDEX_PHASE_STATUS]:
            phase_start_index = 0

    return trend_df


def phase_static(share_df):
    """
    获取各阶段的背离情况
    按照up、down与shake，记录各个阶段的红柱、绿柱面积
    """
    trend_df = trend(share_df)
    if trend is None:
        return None
    trend_df = trend_df.merge(share_df[["date", "high", "low", "open"]], left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='left')
    phase_df = phase.Phase(trend_df).merge(phase.PTYPE_MACD)
    return phase_df
