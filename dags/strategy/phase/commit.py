class Commit(object):
    # 所有phase的数据集合
    phases = dict()
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

    def submit(self, phase_obj, end_date):
        self.end_date = end_date
        self.start_date = phase_obj.start_date
        self.max = phase_obj.max
        self.min = phase_obj.min
        self.sum = phase_obj.sum
        self.d = phase_obj.d
        return

    def reset(self):
        self.end_date = None
        self.start_date = None
        self.max = 0
        self.min = 0
        self.sum = 0
        self.d = None
        return

    def push(self):
        phase = dict()
        phase["d"] = self.d
        phase["max"] = self.max
        phase["min"] = self.min
        phase["sum"] = self.sum
        phase["start_date"] = self.start_date
        phase["end_date"] = self.end_date
        self.phases[self.start_date] = phase
        self.reset()
        return

    def result(self):
        return self.phases
