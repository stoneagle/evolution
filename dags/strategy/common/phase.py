from quota.util import action
from library import conf


def now_and_shake_before(trend_df):
    """
    获取trend不同时机的row
    """
    # 最新bar
    now = trend_df.iloc[-1]
    # 次新bar
    pre = trend_df.iloc[-2]
    # 当前macd趋势开始震荡前的bar
    if now[action.INDEX_STATUS] != action.STATUS_SHAKE and pre[action.INDEX_STATUS] == action.STATUS_SHAKE:
        trend_df = trend_df.head(len(trend_df) - 1)
    trend_no_shake_df = trend_df[trend_df[action.INDEX_STATUS] != action.STATUS_SHAKE]
    shake_before = trend_no_shake_df.iloc[-1]

    # 当前macd趋势的开始bar的数据
    if shake_before.name - shake_before[action.INDEX_TREND_COUNT] > 0:
        start = trend_df.loc[shake_before.name - shake_before[action.INDEX_TREND_COUNT]]
    else:
        start = None
    return start, shake_before, pre, now


def now(trend_df):
    """
    获取trend最近两个phase的起点与终点
    """
    trend_no_shake_df = trend_df[trend_df[action.INDEX_STATUS] != action.STATUS_SHAKE]
    now = trend_no_shake_df.iloc[-1]
    if now.name - now[action.INDEX_TREND_COUNT] > 0:
        start = trend_df.loc[now.name - now[action.INDEX_TREND_COUNT]]
    else:
        start = None
    return start, now


def latest_dict(trend_df, phase_df):
    """
    获取trend的phase列表
    """
    trend_no_shake_df = trend_df[trend_df[action.INDEX_PHASE_STATUS] != action.STATUS_SHAKE]

    phase_row = dict()
    update_flag = False
    if len(phase_df) != 0:
        latest_date = phase_df.iloc[-1][conf.HDF5_SHARE_DATE_INDEX]
        trend_no_shake_df = trend_no_shake_df[trend_no_shake_df[conf.HDF5_SHARE_DATE_INDEX] >= latest_date]
        update_flag = True
    trend_no_shake_df = trend_no_shake_df.reset_index(drop=True)

    for index, row in trend_no_shake_df.iterrows():
        if index < len(trend_no_shake_df) - 1:
            next_row = trend_no_shake_df.iloc[index + 1]

        if len(phase_row) == 0:
            phase_row[action.INDEX_STATUS] = row[action.INDEX_PHASE_STATUS]
            phase_row[conf.HDF5_SHARE_DATE_INDEX] = row[conf.HDF5_SHARE_DATE_INDEX]
            t_count = 1
            m_start = row["macd"]
            c_start = row["close"]
            continue

        if phase_row[action.INDEX_STATUS] != next_row[action.INDEX_PHASE_STATUS]:
            if update_flag is True:
                phase_len = len(phase_df)
                phase_df.loc[phase_len - 1:phase_len, "macd"] = row["macd"] - m_start
                phase_df.loc[phase_len - 1:phase_len, "close"] = row["close"] - c_start
                phase_df.loc[phase_len - 1:phase_len, "count"] = t_count
                update_flag = False
            else:
                phase_row["macd"] = row["macd"] - m_start
                phase_row["close"] = row["close"] - c_start
                phase_row["count"] = t_count
                phase_df = phase_df.append(phase_row, ignore_index=True)
            phase_row = dict()
        else:
            t_count += 1
    phase_df = phase_df.append(phase_row, ignore_index=True)
    return phase_df
