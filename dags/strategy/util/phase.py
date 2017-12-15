from strategy.util import action
from library import tool

PTYPE_MACD = "macd"
PTYPE_WRAP = "wrap"

POSITION_NEGATIVE = "negative"
POSITION_POSITIVE = "positive"

MACD_START_DATE = "start_date"
MACD_END_DATE = "end_date"
MACD_START = "start"
MACD_END = "end"
MACD_AREA = "area"
MACD_VALUE_OPEN = "open"
MACD_VALUE_CLOSE = "close"
MACD_VALUE_HIGH = "high"
MACD_VALUE_LOW = "low"
MACD_POSITION = "position"
MACD_STATUS = "status"
MACD_COLUMNS = [MACD_STATUS, MACD_POSITION, MACD_START_DATE, MACD_END_DATE, MACD_START, MACD_END, MACD_AREA, MACD_VALUE_OPEN, MACD_VALUE_CLOSE, MACD_VALUE_HIGH, MACD_VALUE_LOW]


class Phase(object):
    # share数据
    share_df = None
    # 聚合类型
    ptype = ""
    # 所有阶段聚合结果
    phase_df = None

    def __init__(self, share_df):
        self.share_df = share_df
        return

    def merge(self, ptype):
        if ptype == PTYPE_MACD:
            self.phase_df = tool.init_empty_df(MACD_COLUMNS)
            ret = self.merge_macd("close")
        else:
            ret = None
        return ret

    def merge_macd(self, value_column):
        one = dict()
        trend_len = len(self.share_df)
        for index, row in self.share_df.iterrows():
            share_row = self.share_df[self.share_df[action.INDEX_DATE] == row[action.INDEX_DATE]]
            share_high = share_row["high"].values[0]
            share_low = share_row["low"].values[0]
            if index == trend_len - 1:
                break
            next_row = self.share_df.iloc[index + 1]

            if len(one) == 0:
                one[MACD_STATUS] = row[action.INDEX_PHASE_STATUS]
                one[MACD_START] = row[value_column]
                one[MACD_AREA] = abs(row[value_column])
                one[MACD_START_DATE] = row[action.INDEX_DATE]
                if one[MACD_STATUS] == action.STATUS_UP:
                    one[MACD_POSITION] = self._macd_position(share_high)
                    one[MACD_VALUE_OPEN] = share_high
                else:
                    one[MACD_POSITION] = self._macd_position(share_low)
                    one[MACD_VALUE_OPEN] = share_low
                one[MACD_VALUE_CLOSE] = 0
                one[MACD_VALUE_HIGH] = share_high
                one[MACD_VALUE_LOW] = share_low
                continue
            one[MACD_VALUE_HIGH] = max(one[MACD_VALUE_HIGH], share_high)
            one[MACD_VALUE_LOW] = min(one[MACD_VALUE_LOW], share_low)
            one[MACD_AREA] += abs(row[value_column])

            # 如果当前phase_status与该阶段的phase_status不一致，或者macd穿越0轴，则统计阶段数据
            if next_row[action.INDEX_PHASE_STATUS] != one[MACD_STATUS] or (self._check_position(one[MACD_START], next_row[value_column])):
                one[MACD_END] = row[value_column]
                one[MACD_END_DATE] = row[action.INDEX_DATE]
                if one[MACD_STATUS] == action.STATUS_UP:
                    one[MACD_VALUE_CLOSE] = share_high
                else:
                    one[MACD_VALUE_CLOSE] = share_low
                self.phase_df = self.phase_df.append(one, ignore_index=True)
                one = dict()
                continue
        return self.phase_df

    def _macd_position(self, value):
        if value >= 0:
            position = POSITION_POSITIVE
        else:
            position = POSITION_NEGATIVE
        return position

    def _check_position(self, pre, now):
        if pre > 0:
            if now <= 0:
                ret = True
            else:
                ret = False
        else:
            if now > 0:
                ret = True
            else:
                ret = False
        return ret
