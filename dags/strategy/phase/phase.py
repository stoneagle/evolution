from strategy.phase import commit
# 第一次反转后单边趋势突破需要次数
FIRST_TREND_LIMIT = 2
# 阶段方向
DIRECTION_UP = "up"
DIRECTION_DOWN = "down"
DIRECTION_CHECK = "check"
DIRECTION_SHAKE = "shake"


class Phase(object):
    # 当前时间的值
    now = 0
    # 当前时间的方向
    now_d = None
    # 当前时间
    now_date = None
    # 上一时间的值
    pre = 0
    # 上一时间的方向
    pre_d = None
    # 上一时间
    pre_date = None
    # MACD阶段方向
    d = None
    # 阶段最大值
    max = 0
    # 阶段最小值
    min = 0
    # 阶段macd总值
    sum = 0
    # 阶段开始日期
    start_date = None
    # 阶段结束日期
    end_date = None
    # 趋势持续次数
    count = 0
    # 第一次转折的边界
    border = 0
    # 翻转开始后的转折次数
    turn_count = 0
    # 待提交区
    commit = None

    def __init__(self, first, second, start_date):
        """
        根据最初两个时间点的数据，初始化阶段
        """
        self.commit = commit.Commit()
        self.pre = second
        self.sum = abs(first) + abs(second)
        self.max = max(first, second)
        self.min = min(first, second)
        self.count = 2
        self.start_date = start_date
        diff = first - second
        if diff > 0 or (diff == 0 and second >= 0):
            direction = DIRECTION_UP
        elif diff < 0 or (diff == 0 and second < 0):
            direction = DIRECTION_DOWN
        self.now_d = direction
        self.pre_d = direction
        self.d = direction
        return

    def reset(self):
        """
        重置所有数据
        """
        self.pre = 0
        self.pre_d = None
        self.pre_date = None
        self.now = 0
        self.now_d = None
        self.now_date = None
        self.sum = 0
        self.max = 0
        self.min = 0
        self.count = 0
        self.start_date = None
        self.end_date = None
        self.d = None
        return

    def flow(self, value, date, trend_flag):
        """
        now与pre值的流动
        value、date、d向前推移
        """
        pre_d = self.pre_d
        self.pre = self.now
        self.pre_d = self.now_d
        self.pre_date = self.now_date
        if trend_flag is True:
            now_d = pre_d
        else:
            now_d = self.get_op_d(pre_d)
        self.now = value
        self.now_d = now_d
        self.now_date = date
        return

    def start_trend(self):
        """
        进入新的趋势
        """
        self._start(self.now, self.now_date)
        new_d = self.get_op_d(None)
        self.now_d = new_d
        self.d = new_d
        return

    def start_check(self):
        """
        进入check等待趋势出现
        """
        self._start(self.now, self.now_date)
        new_d = self.get_op_d(None)
        self.border = self.now
        self.now_d = new_d
        self.d = DIRECTION_CHECK
        return

    def _start(self, value, start_date):
        self.prefix = value
        self.sum = abs(value)
        self.max = value
        self.min = value
        self.start_date = start_date
        return

    def enter_shake(self):
        """
        check阶段结束，进入shake状态
        """
        if self.d == DIRECTION_CHECK:
            self.d = DIRECTION_SHAKE
        return

    def enter_trend(self):
        """
        check阶段结束，进入trend状态
        """
        if self.d == DIRECTION_CHECK:
            self.d = self.now_d
        return

    def finish_shake(self):
        """
        shake阶段结束
        """
        self.enter_shake()
        self.submit()
        self.push()
        return

    def merge_trend(self):
        """
        合并趋势数据
        """
        self._merge()
        self.now_d = self.pre_d
        self.count += 1
        if self.d == DIRECTION_CHECK and self.turn_count == 1:
            # 第一次转折边界延伸
            self.border = self.now
            # 如果第一次反转后，单边延续次数超出界限，则视为新趋势启动
            if self.count >= FIRST_TREND_LIMIT:
                self.d = self.pre_d
                self.now_d = self.pre_d
                self.push()
        return

    def merge_turn(self):
        """
        合并转折数据
        """
        self._merge()
        self.now_d = self.get_op_d(None)
        self.turn_count += 1
        return

    def _merge(self):
        self.pre = self.now
        self.sum += abs(self.now)
        self.max = max(self.max, self.now)
        self.min = min(self.min, self.now)
        return

    def submit(self):
        """
        将phase数据提交到commit区域
        """
        self.commit.submit(self, end_date=self.now_date)
        self.reset()
        return

    def trend_restore(self):
        """
        趋势突破，将commit区域数据恢复至phase
        """
        self.merge_trend()
        self._resotre()
        return

    def turn_restore(self):
        """
        转折突破，将commit区域数据恢复至phase
        """
        self.merge_turn()
        self._resotre()
        return

    def _restore(self):
        self.d = self.commit.d
        self.now_d = self.commit.d
        self.sum += self.commit.sum
        self.max = max(self.max, self.commit.max)
        self.min = min(self.min, self.commit.min)
        self.commit.reset()
        return

    def push(self):
        """
        将commit区域推送至结果集合
        """
        self.commit.push()
        return

    def result(self):
        """
        获取结果集合
        """
        self.submit()
        self.push()
        return self.commit.result()

    def get_op_d(self, d):
        """
        获取相反方向
        """
        if d is None:
            d = self.pre_d
        if d == DIRECTION_UP:
            return DIRECTION_DOWN
        elif d == DIRECTION_DOWN:
            return DIRECTION_UP
