from strategy.util import wrap
import numpy as np
POINT_UP = "up"
POINT_DOWN = "down"
INDEX_POINT_TYPE = "point"


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
    point_df = point(wrap_df.head(25))
    # 判断是否存在中枢
    print(point_df[["high", "low", "point"]])
    return


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
            continue

        # 如果处于check状态
        if check_count > 0:
            point_high = wrap_df.iloc[point_index]["high"]
            point_low = wrap_df.iloc[point_index]["low"]
            # 转折继续延续，则更新check数量
            if (direction == POINT_UP and row["high"] < pre_high) or (direction == POINT_DOWN and row["low"] > pre_low):
                check_count += 1
                if check_count >= wrap.TURN_LIMIT:
                    # 如果转折满足数量，则记录尖点，并更新新的尖点
                    check_count = 0
                    point_index = index
                    if direction == POINT_UP:
                        direction = POINT_DOWN
                        wrap_df.loc[point_index, INDEX_POINT_TYPE] = POINT_UP
                    else:
                        direction = POINT_UP
                        wrap_df.loc[point_index, INDEX_POINT_TYPE] = POINT_DOWN
            # 转折震荡后恢复原有走势
            if (direction == POINT_UP and row["high"] > point_high) or (direction == POINT_DOWN and row["low"] < point_low):
                check_count = 0
                point_index = index
            continue

        if (direction == POINT_UP and row["high"] > pre_high) or (direction == POINT_DOWN and row["low"] < pre_low):
            # 如果k线延续，更新尖点的index
            point_index = index
        else:
            # 如果k线出现转折，进入check状态
            check_count = 1
    print(wrap_df[["high", "low"]])
    if np.isnan(wrap_df.iloc[len(wrap_df) - 1][INDEX_POINT_TYPE]):
        wrap_df.loc[len(wrap_df) - 1, INDEX_POINT_TYPE] = direction
    return wrap_df[wrap_df[INDEX_POINT_TYPE].notnull()]
