from strategy.util import wrap
from library import tool
import numpy as np
POINT_UP = "up"
POINT_DOWN = "down"
INDEX_POINT_TYPE = "point"
INDEX_POINT_DIRECTION = "point_direction"
INDEX_POINT_INDEX = "point_index"
INDEX_POINT_CHECK_COUNT = "point_check_count"


def trans_wrap(share_df, merge_price):
    """
    聚合k线成为缠论模式
    """
    return wrap.Wrap(share_df).merge_line("high", "low", merge_price)


def central(share_df):
    """
    聚合k线后，标记顶与底，统计震荡中枢数量
    """
    wrap_df = trans_wrap(share_df, False)
    phase_df = point_phase(wrap_df)
    # 如果聚合区间不足三条，无法进行分析
    if len(phase_df) <= 3:
        return None

    # 根据聚合的up与down的值，判断是否存在中枢
    central_dict = dict()
    central_df = tool.init_empty_df(["high", "low", "count", "start_date", "end_date"])
    for index, row in phase_df.iterrows():
        if index == 0:
            central_dict["high"] = row["high"]
            central_dict["low"] = row["low"]
            central_dict["count"] = 1
            central_dict["out_count"] = 0
            central_dict["start_date"] = row["start_date"]
        else:
            compare_high = min(row["high"], central_dict["high"])
            compare_low = max(row["low"], central_dict["low"])
            if compare_high < compare_low:
                if central_dict["count"] < 3:
                    # 不存在重叠，重置数据
                    pre_row = phase_df.iloc[index - 1]
                    central_dict["count"] = 1
                    central_dict["out_count"] = 0
                    central_dict["high"] = min(row["high"], pre_row["high"])
                    central_dict["low"] = max(row["low"], pre_row["low"])
                    central_dict["start_date"] = row["start_date"]
                else:
                    # 已重叠，并出现脱离中枢范围的情况
                    central_dict["out_count"] += 1
                    if central_dict["out_count"] >= 2:
                        central_dict["end_date"] = phase_df.iloc[index - 2]["end_date"]
                        central_df = central_df.append(central_dict, ignore_index=True)
                        pre_row = phase_df.iloc[index - 1]
                        central_dict["high"] = min(row["high"], pre_row["high"])
                        central_dict["low"] = max(row["low"], pre_row["low"])
                        central_dict["start_date"] = pre_row["start_date"]
                        central_dict["count"] = 2
                        central_dict["out_count"] = 0
            else:
                # 存在重叠
                central_dict["high"] = compare_high
                central_dict["low"] = compare_low
                central_dict["count"] += 1
        # 前三根，判断是否存在重叠范围，如果存在则出现中枢
        # 连续两根超出中枢范围，进入新的趋势判断
    central_dict = central_dict.drop("out_count", axis=1)
    return central_dict


def point_phase(wrap_df):
    point_df = point(wrap_df[["high", "low", "date"]])
    # 取up的high，down的low值，组成一个阶段的up与down值，逐个阶段判断是否存在重合
    phase_df = tool.init_empty_df(["high", "low", "start_date", "end_date"])
    row_dict = dict()
    for index, row in point_df.iterrows():
        # 如果第一个尖点是底则忽略
        if index == 0 and row[INDEX_POINT_TYPE] == POINT_UP:
            continue

        if row[INDEX_POINT_TYPE] == POINT_DOWN:
            row_dict["low"] = row["low"]
            row_dict["start_date"] = row["date"]
        else:
            row_dict["high"] = row["high"]
            row_dict["end_date"] = row["date"]
            phase_df = phase_df.append(row_dict, ignore_index=True)
            row_dict = dict()
    return phase_df


def point(wrap_df):
    # 标记顶与底
    direction = "None"
    point_index = 0
    check_count = 0
    for index, row in wrap_df.iterrows():
        # 如果连续出现4根反方向的k线，之前的尖点视为顶或底(中间如果出现震荡，判断是否超出之前的尖点)
        if index == 0:
            point_index = index
            continue

        pre_high = wrap_df.iloc[index - 1]["high"]
        pre_low = wrap_df.iloc[index - 1]["low"]

        if index == 1 and direction == "None":
            if row["high"] > pre_high:
                direction = POINT_UP
            else:
                direction = POINT_DOWN
            point_index = index
            wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = 0
            wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
            wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
            continue

        # 如果处于check状态
        if check_count > 0:
            point_high = wrap_df.iloc[point_index]["high"]
            point_low = wrap_df.iloc[point_index]["low"]
            # 转折继续延续，则更新check数量
            if (direction == POINT_UP and row["high"] < pre_high) or (direction == POINT_DOWN and row["low"] > pre_low):
                check_count += 1
                wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = check_count
                wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
                wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
                if check_count >= wrap.TURN_LIMIT:
                    # 如果转折满足数量，则记录尖点，并更新新的尖点
                    if direction == POINT_UP:
                        direction = POINT_DOWN
                        wrap_df.loc[point_index, INDEX_POINT_TYPE] = POINT_UP
                    else:
                        direction = POINT_UP
                        wrap_df.loc[point_index, INDEX_POINT_TYPE] = POINT_DOWN
                    check_count = 0
                    point_index = index
                    wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = 0
                    wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
                    wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
            # 转折震荡后恢复原有走势
            elif (direction == POINT_UP and row["high"] > point_high) or (direction == POINT_DOWN and row["low"] < point_low):
                check_count = 0
                point_index = index
                wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = 0
                wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
                wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
            else:
                wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = check_count
                wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
                wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
            continue

        if (direction == POINT_UP and row["high"] >= pre_high) or (direction == POINT_DOWN and row["low"] <= pre_low):
            # 如果k线延续，更新尖点的index
            point_index = index
            wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = 0
            wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
            wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
        else:
            # 如果k线出现转折，进入check状态
            check_count = 1
            wrap_df.loc[index, INDEX_POINT_CHECK_COUNT] = 1
            wrap_df.loc[index, INDEX_POINT_DIRECTION] = direction
            wrap_df.loc[index, INDEX_POINT_INDEX] = point_index
    if np.isnan(wrap_df.iloc[len(wrap_df) - 1][INDEX_POINT_TYPE]):
        wrap_df.loc[len(wrap_df) - 1, INDEX_POINT_TYPE] = direction
    wrap_df = wrap_df[wrap_df[INDEX_POINT_TYPE].notnull()]
    return wrap_df[["date", "high", "low", INDEX_POINT_TYPE]].reset_index(drop=True)
