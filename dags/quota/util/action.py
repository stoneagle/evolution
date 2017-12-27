from library import tool

# 行为类型
TREND_STILL = "trend_still"
TREND_TURN = "trend_turn"
STILL_OUT = "still_out"
STILL_IN = "still_in"
TURN_IN = "turn_in"
TURN_OUT = "turn_out"

# 走势方向
DIRECTION_UP = "up"
DIRECTION_DOWN = "down"

# 趋势状态
STATUS_UP = "up"
STATUS_DOWN = "down"
STATUS_SHAKE = "shake"

INDEX_DATE = "date"
INDEX_VALUE = "value"
INDEX_DIRECTION = "direction"
INDEX_TURN_COUNT = "turn_count"
INDEX_TREND_COUNT = "trend_count"
INDEX_ACTION = "action"
INDEX_STATUS = "status"
INDEX_PHASE_STATUS = "phase_status"
# 进入震荡状态所需最小turn次数
SHAKE_MIN_NUM = 3
# 单次有效突破最低次数
TURN_MIN_NUM = 3
# 默认macd的range因子，取10%数值
DEFAULT_FACTOR_MACD_RANGE = 0.1


class Action(object):
    # 存储index_df数据
    index_df = None
    # 存储所有时间节点的走向情况
    df = tool.init_empty_df([INDEX_DATE, INDEX_VALUE, INDEX_DIRECTION, INDEX_TURN_COUNT, INDEX_TREND_COUNT, INDEX_ACTION, INDEX_STATUS, INDEX_PHASE_STATUS])
    # 转折后的上边界
    up_border = 0
    # 转折后的下边界
    down_border = 0
    # 方案A，绝对值(废弃):不同价格量级的数据(例如比特币)，波动幅度不同，不能用固定值
    # 方案B，该段趋势波动差，差值越大越快说明趋势越急，对转折要求更高
    factor_macd_range = None

    def __init__(self, index_df, factor_macd_range=None):
        self.index_df = index_df
        if factor_macd_range is not None:
            self.factor_macd_range = factor_macd_range
        else:
            self.factor_macd_range = DEFAULT_FACTOR_MACD_RANGE
        return

    def all(self, date_column, value_column):
        """
        获取从头到尾的趋势
        """
        first = self.index_df.iloc[1]
        second = self.index_df.iloc[2]
        diff = second[value_column] - first[value_column]
        if diff > 0 or (diff == 0 and second >= 0):
            direction = DIRECTION_UP
            status = STATUS_UP
        elif diff < 0 or (diff == 0 and second < 0):
            direction = DIRECTION_DOWN
            status = STATUS_DOWN
        first_row = {
            INDEX_DATE: first[date_column],
            INDEX_VALUE: first[value_column],
            INDEX_DIRECTION: direction,
            INDEX_STATUS: status,
            INDEX_PHASE_STATUS: status,
            INDEX_TURN_COUNT: 0,
            INDEX_TREND_COUNT: 1,
            INDEX_ACTION: TREND_STILL,
        }
        second_row = {
            INDEX_DATE: second[date_column],
            INDEX_VALUE: second[value_column],
            INDEX_DIRECTION: direction,
            INDEX_STATUS: status,
            INDEX_PHASE_STATUS: status,
            INDEX_TURN_COUNT: 0,
            INDEX_TREND_COUNT: 2,
            INDEX_ACTION: TREND_STILL
        }
        self.df = self.df.append(first_row, ignore_index=True)
        self.df = self.df.append(second_row, ignore_index=True)

        for index, row in self.index_df[3:].iterrows():
            date = row[date_column]
            value = row[value_column]
            pre_turn_count = self.df.iloc[-1][INDEX_TURN_COUNT]
            # 1.单边趋势
            if pre_turn_count == 0:
                one = self.trend(value, date)
            # 2.第N次转折后，根据选择方向进行判断:
            elif pre_turn_count >= 1:
                one = self.turn(value, date)
            self.df = self.df.append(one, ignore_index=True)
        return self.df

    def trend(self, value, date):
        """
        单边趋势的方向选择
        """
        one = dict()
        one[INDEX_DATE] = date
        one[INDEX_VALUE] = value
        pre_row = self.df.iloc[-1]
        # 趋势延续
        if self.compare_border(TREND_STILL, value) or (value == pre_row[INDEX_VALUE]):
            self.reset_border()
            one[INDEX_ACTION] = TREND_STILL
            one[INDEX_TREND_COUNT] = pre_row[INDEX_TREND_COUNT] + 1
            one[INDEX_TURN_COUNT] = 0
            one[INDEX_STATUS] = pre_row[INDEX_STATUS]
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
            one[INDEX_DIRECTION] = pre_row[INDEX_DIRECTION]
        # 出现转折
        elif self.compare_border(TREND_TURN, value):
            pre_value = pre_row[INDEX_VALUE]
            if len(self.df) > pre_row[INDEX_TREND_COUNT]:
                start_row = self.df.iloc[-1 - pre_row[INDEX_TREND_COUNT]]
            else:
                start_row = self.df.iloc[-pre_row[INDEX_TREND_COUNT] + 1]
            macd_range = abs(pre_value - start_row[INDEX_VALUE])
            if pre_row[INDEX_DIRECTION] == DIRECTION_UP:
                self.set_border(pre_value - self.factor_macd_range * macd_range, pre_value)
            else:
                self.set_border(pre_value + self.factor_macd_range * macd_range, pre_value)
            one[INDEX_ACTION] = TREND_TURN
            one[INDEX_TREND_COUNT] = 1
            one[INDEX_TURN_COUNT] = 1
            one[INDEX_STATUS] = STATUS_SHAKE
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
            one[INDEX_DIRECTION] = self.get_op_d(pre_row[INDEX_DIRECTION])
        return one

    def turn(self, value, date):
        """
        翻转后的方向选择
        """
        one = dict()
        one[INDEX_DATE] = date
        one[INDEX_VALUE] = value
        pre_row = self.df.iloc[-1]
        border_row = None
        for i in range(1, len(self.df)):
            border_row = self.df.iloc[-i]
            if border_row[INDEX_STATUS] != STATUS_SHAKE:
                break
        # 如果未出现转折，并且超出转折边界
        if self.compare_border(STILL_OUT, value):
            one[INDEX_ACTION] = STILL_OUT
            one[INDEX_DIRECTION] = pre_row[INDEX_DIRECTION]
            one[INDEX_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
            one[INDEX_TREND_COUNT] = 1 + pre_row[INDEX_TREND_COUNT]
            one[INDEX_TURN_COUNT] = 0
            pre_turn_count = pre_row[INDEX_TURN_COUNT]
            # 如果只转折一次
            if pre_turn_count == 1:
                if one[INDEX_TREND_COUNT] >= TURN_MIN_NUM:
                    self.df.loc[len(self.df) - pre_row[INDEX_TREND_COUNT]:len(self.df), INDEX_PHASE_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
                    self.reset_border()
                else:
                    one[INDEX_STATUS] = STATUS_SHAKE
                    one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
                    one[INDEX_TURN_COUNT] = 1
            elif pre_turn_count % 2 == 0:
                if pre_turn_count < SHAKE_MIN_NUM:
                    # 视作上一阶段延续
                    for i in range(1, len(self.df)):
                        if self.df.iloc[len(self.df) - i][INDEX_PHASE_STATUS] != STATUS_SHAKE:
                            self.df.loc[len(self.df) - i:len(self.df), INDEX_PHASE_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
                            break
                    one[INDEX_TREND_COUNT] = border_row[INDEX_TREND_COUNT] + i
                    self.reset_border()
                else:
                    # 视作新的相反阶段开始
                    self.reset_border()
            elif pre_turn_count % 2 == 1:
                # 小于shake次数，趋势逆转，大于shake次数，震荡结束
                if pre_turn_count < SHAKE_MIN_NUM:
                    self.reset_border()
                else:
                    self.reset_border()
        # 如果未出现转折，延伸长度未超出转折边界
        elif self.compare_border(STILL_IN, value):
            one[INDEX_ACTION] = STILL_IN
            one[INDEX_TREND_COUNT] = pre_row[INDEX_TREND_COUNT] + 1
            one[INDEX_TURN_COUNT] = pre_row[INDEX_TURN_COUNT]
            one[INDEX_DIRECTION] = pre_row[INDEX_DIRECTION]
            one[INDEX_STATUS] = pre_row[INDEX_STATUS]
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
        # 如果出现转折，延伸长度未超出上一阶段极值
        elif self.compare_border(TURN_IN, value):
            one[INDEX_ACTION] = TURN_IN
            one[INDEX_TREND_COUNT] = 1
            one[INDEX_TURN_COUNT] = pre_row[INDEX_TURN_COUNT] + 1
            one[INDEX_DIRECTION] = self.get_op_d(pre_row[INDEX_DIRECTION])
            one[INDEX_STATUS] = STATUS_SHAKE
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
        # 如果出现转折，超出上一阶段极值
        elif self.compare_border(TURN_OUT, value):
            one[INDEX_ACTION] = TURN_OUT
            one[INDEX_TURN_COUNT] = 0
            one[INDEX_DIRECTION] = self.get_op_d(pre_row[INDEX_DIRECTION])
            one[INDEX_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
            one[INDEX_PHASE_STATUS] = one[INDEX_STATUS]
            one[INDEX_TREND_COUNT] = 1
            pre_turn_count = pre_row[INDEX_TURN_COUNT]
            if pre_turn_count < SHAKE_MIN_NUM:
                if pre_turn_count % 2 == 1:
                    # 上一阶段延续
                    self.df.loc[len(self.df) - pre_row[INDEX_TREND_COUNT]:len(self.df), INDEX_PHASE_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
                    one[INDEX_TREND_COUNT] = border_row[INDEX_TREND_COUNT] + i
                else:
                    # 趋势逆转
                    for i in range(1, 100):
                        if self.df.iloc[len(self.df) - i][INDEX_PHASE_STATUS] != STATUS_SHAKE:
                            self.df.loc[len(self.df) - i:len(self.df), INDEX_PHASE_STATUS] = self.get_d_s(one[INDEX_DIRECTION])
                            break
            self.reset_border()
        return one

    def compare_border(self, atype, value):
        """
        比较value与边界
        """
        pre_row = self.df.iloc[-1]
        d = pre_row[INDEX_DIRECTION]
        pre = pre_row[INDEX_VALUE]
        up_border = self.up_border
        down_border = self.down_border

        up_switch = {
            TREND_STILL: value >= pre,
            TREND_TURN: value < pre,
            STILL_OUT: value > up_border,
            STILL_IN: pre <= value <= up_border,
            TURN_IN: down_border <= value < pre,
            TURN_OUT: value < down_border,
        }
        down_switch = {
            TREND_STILL: value <= pre,
            TREND_TURN: value > pre,
            STILL_OUT: value < down_border,
            STILL_IN: pre >= value >= down_border,
            TURN_IN: up_border >= value > pre,
            TURN_OUT: value > up_border,
        }
        if d == DIRECTION_UP:
            return up_switch[atype]
        elif d == DIRECTION_DOWN:
            return down_switch[atype]
        else:
            raise Exception("边界获取异常")

    def set_border(self, a, b):
        """
        设置初始边界
        """
        self.up_border = max(a, b)
        self.down_border = min(a, b)
        return

    def reset_border(self):
        self.up_border = 0
        self.down_border = 0
        return

    def get_op_d(self, d):
        """
        获取相反方向
        """
        if d == DIRECTION_UP:
            return DIRECTION_DOWN
        elif d == DIRECTION_DOWN:
            return DIRECTION_UP

    def get_d_s(self, d):
        """
        获取方向的对应状态
        """
        if d == DIRECTION_UP:
            return STATUS_UP
        elif d == DIRECTION_DOWN:
            return STATUS_DOWN
