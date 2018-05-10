from quota.util import action
from library import conf
MACD_DIFF = "macd_diff"
PRICE_START = "price_start"
PRICE_END = "price_end"
DIF_END = "dif_end"
COUNT = "count"


def pre(trend_df):
    """
    获取trend上一phase的起点与终点
    """
    now_start, now_end = now(trend_df)
    pre_end = now_start
    if pre_end.name - pre_end[action.INDEX_TREND_COUNT] > 0:
        pre_start = trend_df.loc[pre_end.name - pre_end[action.INDEX_TREND_COUNT]]
    else:
        pre_start = None
    return pre_start, pre_end


def now(trend_df):
    """
    获取trend当前phase的起点与终点
    """
    trend_no_shake_df = trend_df[trend_df[action.INDEX_STATUS] != action.STATUS_SHAKE]
    end = trend_no_shake_df.iloc[-1]
    if end.name - end[action.INDEX_TREND_COUNT] > 0:
        start = trend_df.loc[end.name - end[action.INDEX_TREND_COUNT]]
    else:
        start = None
    return start, end


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
        else:
            break

        if len(phase_row) == 0:
            phase_row[action.INDEX_PHASE_STATUS] = row[action.INDEX_PHASE_STATUS]
            phase_row[conf.HDF5_SHARE_DATE_INDEX] = row[conf.HDF5_SHARE_DATE_INDEX]
            t_count = 1
            m_start = row["macd"]
            if len(phase_df) > 1:
                phase_row[PRICE_START] = phase_df.iloc[-1][PRICE_END]
            else:
                phase_row[PRICE_START] = int(row["close"])
        else:
            t_count += 1
            if phase_row[action.INDEX_PHASE_STATUS] != next_row[action.INDEX_PHASE_STATUS]:
                if update_flag is True:
                    phase_len = len(phase_df)
                    phase_df.loc[phase_len - 1:phase_len, MACD_DIFF] = row["macd"] - m_start
                    phase_df.loc[phase_len - 1:phase_len, PRICE_END] = row["close"]
                    phase_df.loc[phase_len - 1:phase_len, DIF_END] = row["dif"]
                    phase_df.loc[phase_len - 1:phase_len, COUNT] = int(t_count)
                    update_flag = False
                else:
                    phase_row[MACD_DIFF] = row["macd"] - m_start
                    phase_row[PRICE_END] = row["close"]
                    phase_row[DIF_END] = row["dif"]
                    phase_row[COUNT] = t_count
                    phase_df = phase_df.append(phase_row, ignore_index=True)
                phase_row = dict()

    if len(phase_df) == 0:
        phase_df = phase_df.append(phase_row, ignore_index=True)
    elif len(phase_df) > 0 and phase_df.iloc[-1][conf.HDF5_SHARE_DATE_INDEX] != phase_row[conf.HDF5_SHARE_DATE_INDEX]:
        phase_df = phase_df.append(phase_row, ignore_index=True)
    return phase_df
