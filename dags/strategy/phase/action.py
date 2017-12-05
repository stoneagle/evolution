from library import conf
from strategy.phase import phase
# 行为类型
TREND_STILL = "trend_still"
TREND_TURN = "trend_turn"
TURN_FIRST_STILL = "turn_first_still"
TURN_FIRST_TURN_OUT = "turn_first_turn_out"
TURN_FIRST_TURN_IN = "turn_first_turn_in"
TURN_OVEN_STILL_OUT = "turn_even_still_out"
TURN_OVEN_STILL_IN = "turn_even_still_in"
TURN_OVEN_TURN_IN = "turn_even_turn_in"
TURN_OVEN_TURN_OUT = "turn_even_turn_out"
TURN_ODD_STILL_OUT = "turn_odd_still_out"
TURN_ODD_STILL_IN = "turn_odd_still_in"
TURN_ODD_TURN_IN = "turn_odd_turn_in"
TURN_ODD_TURN_OUT = "turn_odd_turn_out"


class Action(object):
    # 存储各个阶段的数据
    phase = None

    def __init__(self, phase):
        self.phase = phase
        return

    def run(self, value, date):
        self.value = value
        self.date = date
        # 1 单边趋势
        if self.phase.turn_count == 0:
            self.trend()
        # 2 震荡状态
        # 2.1 第1次转折后的判断:
        elif self.phase.turn_count == 1:
            self.turn_first()
        # 2.2 第2次转折后的判断
        elif self.phase.turn_count == 2:
            self.turn_second()
        # 2.3 第3次转折后的判断
        elif self.phase.turn_count == 3:
            self.turn_third()
        # 2.5 第2N次(N>2)转折后，根据选择方向进行判断:
        elif self.phase.turn_count % 2 == 0:
            self.turn_oven()
        # 2.4 第2N-1(N>3)次转折后，根据选择方向进行判断:
        elif self.phase.turn_count % 2 == 1:
            self.turn_odd()
        return

    def trend(self):
        """
        单边趋势的方向选择
        """
        # 趋势延续
        if self.compare_border(TREND_STILL) or (self.phase.now == self.phase.pre):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 出现转折
        elif self.compare_border(TREND_TURN):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 提交上一阶段
            self.phase.submit()
            # 进入check阶段
            self.phase.start_check()
        return

    def turn_first(self):
        """
        第一次翻转后的方向选择
        """
        # 如果未出现转折
        if self.compare_border(TURN_FIRST_STILL):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 如果出现转折，延伸长度超出上一阶段极值，上一阶段延续
        elif self.compare_border(TURN_FIRST_TURN_OUT):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.turn_restore()
        # 如果出现转折，延伸长度未超出上一阶段极值，等待方向选择
        elif self.compare_border(TURN_FIRST_TURN_IN):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_turn()
        return

    def turn_second(self):
        """
        第二次翻转后的方向选择
        """
        # 如果未出现转折，并且超出上一阶段极值，视作上一阶段延续
        if self.compare_border(TURN_OVEN_STILL_OUT):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.trend_restore()
        # 如果未出现转折，延伸长度未超出上一阶段极值，等待方向选择
        elif self.compare_border(TURN_OVEN_STILL_IN):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 如果出现转折，延伸长度未超出第一次转折边界，等待第三次转后的方向选择
        elif self.compare_border(TURN_OVEN_TURN_IN):
            self.phase.flow(self.value, self.date, trend_flag=False)
            self.phase.merge_turn()
        # 如果出现转折，超出第一次转折边界，视作趋势翻转
        elif self.compare_border(TURN_OVEN_TURN_OUT):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 上一阶段结束
            self.phase.push()
            # check阶段结束，进入新趋势
            self.phase.enter_trend()
        return

    def turn_third(self):
        """
        第三次翻转后的方向选择
        """
        # 如果未出现转折，并且超出转折边界，视作趋势翻转
        if self.compare_border(TURN_ODD_STILL_OUT):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
            # 上一阶段结束
            self.phase.push()
        # 如果未出现转折，延伸长度未超出转折边界，等待方向选择
        elif self.compare_border(conf.TURN_ODD_STILL_IN):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 如果出现转折，延伸长度未超出上一阶段极值，上一阶段结束，进入震荡状态
        elif self.compare_border(TURN_ODD_TURN_IN):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 上一阶段结束
            self.phase.push()
            # 进入震荡状态
            self.phase.merge_turn()
            self.phase.enter_shake()
        # 如果出现转折，超出上一阶段极值，上一阶段结束，三根震荡结束，开启上一阶段同向趋势
        elif self.compare_border(TURN_ODD_TURN_OUT):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 上一阶段结束
            self.phase.push()
            # 三根震荡状态结束
            self.phase.finish_shake()
            # 开启新的趋势
            self.phase.start_trend()
        return

    def turn_oven(self):
        """
        偶数次翻转后的方向选择
        """
        # 如果未出现转折，并且超出上一阶段极值，震荡结束
        if self.compare_border(TURN_OVEN_STILL_OUT):
            self.phase.flow(self.value, self.date, trend_flag=True)
            # 偶数震荡状态结束
            self.phase.finish_shake()
            # 开启新的趋势
            self.phase.start_trend()
        # 如果未出现转折，延伸长度未超出上一阶段极值，震荡中
        elif self.compare_border(TURN_OVEN_STILL_IN):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 如果出现转折，延伸长度未超出第一次转折边界，震荡中
        elif self.compare_border(TURN_OVEN_TURN_IN):
            self.phase.flow(self.value, self.date, trend_flag=False)
            self.phase.merge_turn()
        # 如果出现转折，超出第一次转折边界，震荡结束
        elif self.compare_border(TURN_OVEN_TURN_OUT):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 偶数震荡状态结束
            self.phase.finish_shake()
            # 开启新的趋势
            self.phase.start_trend()
        return

    def turn_odd(self):
        """
        奇数次翻转后的方向选择
        震荡中的突破需要注意，偶数方向突破终点是奇数点，奇数方向突破终点是偶数点
        """
        # 如果未出现转折，并且超出转折边界，视作震荡结束
        if self.compare_border(TURN_ODD_STILL_OUT):
            self.phase.flow(self.value, self.date, trend_flag=True)
            # 奇数震荡状态结束
            self.phase.finish_shake()
            # 开启新的趋势
            self.phase.start_trend()
        # 如果未出现转折，延伸长度未超出转折边界，震荡中
        elif self.compare_border(TURN_ODD_STILL_IN):
            self.phase.flow(self.value, self.date, trend_flag=True)
            self.phase.merge_trend()
        # 如果出现转折，延伸长度未超出上一阶段极值，震荡中
        elif self.compare_border(TURN_ODD_TURN_IN):
            self.phase.flow(self.value, self.date, trend_flag=False)
            self.phase.merge_turn()
        # 如果出现转折，超出上一阶段极值，震荡结束
        elif self.compare_border(TURN_ODD_TURN_OUT):
            self.phase.flow(self.value, self.date, trend_flag=False)
            # 奇数震荡状态结束
            self.phase.finish_shake()
            # 开启新的趋势
            self.phase.start_trend()
        return

    def compare_border(self, atype):
        """
        比较不同阶段下，value与边界的情况
        """
        d = self.phase.pre_d
        pre = self.phase.pre
        border = self.phase.border
        commit_max = self.phase.commit.max
        commit_min = self.phase.commit.min
        value = self.value

        up_switch = {
            TREND_STILL: value > pre,
            TREND_TURN: value < pre,
            TURN_FIRST_STILL: value > pre,
            TURN_FIRST_TURN_OUT: value < commit_min,
            TURN_FIRST_TURN_IN: commit_min <= value <= border,
            TURN_OVEN_STILL_OUT: value > commit_max,
            TURN_OVEN_STILL_IN: pre <= value <= commit_max,
            TURN_OVEN_TURN_IN: border <= value < pre,
            TURN_OVEN_TURN_OUT: value < border,
            TURN_ODD_STILL_OUT: value > border,
            TURN_ODD_STILL_IN: pre <= value <= border,
            TURN_ODD_TURN_IN: commit_min <= value < pre,
            TURN_ODD_TURN_OUT: value < commit_min
        }
        down_switch = {
            TREND_STILL: value < pre,
            TREND_TURN: value > pre,
            TURN_FIRST_STILL: value < pre,
            TURN_FIRST_TURN_OUT: value > commit_max,
            TURN_FIRST_TURN_IN: commit_max >= value >= border,
            TURN_OVEN_STILL_OUT: value < commit_min,
            TURN_OVEN_STILL_IN: pre >= value >= commit_min,
            TURN_OVEN_TURN_IN: border >= value > pre,
            TURN_OVEN_TURN_OUT: value > border,
            TURN_ODD_STILL_OUT: value < border,
            TURN_ODD_STILL_IN: pre >= value >= border,
            TURN_ODD_TURN_IN: commit_max >= value > pre,
            TURN_ODD_TURN_OUT: commit_max > border
        }
        if d == phase.DIRECTION_UP:
            return up_switch[atype]
        elif d == phase.DIRECTION_DOWN:
            return down_switch[atype]
        else:
            raise Exception("边界获取异常")
