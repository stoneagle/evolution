from library import tool, conf, tradetime, console
from quota.util import action
from strategy.common import trend, phase
DF_SMALL = "small"
DF_MEDIUM = "medium"
DF_BIG = "big"
TREND_REVERSE = "trend_reverse"
PHASE_COLUMNS = [conf.HDF5_SHARE_DATE_INDEX, phase.MACD_DIFF, phase.PRICE_START, phase.PRICE_END, phase.DIF_END, phase.COUNT, action.INDEX_PHASE_STATUS]
BITMEX_LEVEL_DICT = {
    conf.BINSIZE_ONE_MINUTE: conf.BINSIZE_FIVE_MINUTE,
    conf.BINSIZE_FIVE_MINUTE: conf.BINSIZE_THIRTY_MINUTE,
    conf.BINSIZE_THIRTY_MINUTE: conf.BINSIZE_FOUR_HOUR,
    conf.BINSIZE_FOUR_HOUR: conf.BINSIZE_ONE_DAY,
}


class strategy(object):
    """
    大-中-小，三重级别递归策略
    1. 大与中决定多空方向，小决定交易节点
    2. 小级别的phase背离出现次数，决定杠杆倍率，并推动中级别方向转换
    3. 大或中出现trend背离，子级别方向转换压制

    例如父级别趋势向下，1min出现买点，则起步选择做多
    子级别第一个买点杠杆1x，上涨后macd转折出售，如果dea与dif未上传零轴，5min方向不变
    第二个买点出现，正常买点杠杆不变，背离买点杠杆增加至2x，背离杠杆最多增加至3x
    期间子级别的相反趋势，可以用0.5x杠杆进行交易
    当子级别出现卖点背离时，父级别趋势逆转，方向由做多转为做空
    """
    # 用于存储code名称
    code = None
    # 数据来源类别
    stype = None
    # 回测标签，回测时不更新数据源
    backtest = None
    # 是否将最新数据回写至文件
    rewrite = None
    # 父级别方向
    medium_side = None
    # macd波动因子
    factor_macd_range = None

    # 大级别对应时间节点
    big_level = None
    # 大级别trend背离标签
    big_trend_reverse = False
    # 存储大级别数据
    big = None
    # 中级别对应时间节点
    medium_level = None
    # 中级别trend背离标签
    medium_trend_reverse = False
    # 存储中级别数据
    medium = None
    # 小级对应时间节点
    small_level = None
    # 小级别phase背离标签
    small_phase_reverse = False
    # 存储小级别数据
    small = None
    # 存储小级别趋势数据
    phase = tool.init_empty_df(PHASE_COLUMNS)

    def __init__(self, code, stype, backtest, rewrite, small_level, factor_macd_range=0.1):
        self.code = code
        self.stype = stype
        self.backtest = backtest
        self.rewrite = rewrite
        self.small_level = small_level
        self.factor_macd_range = factor_macd_range
        if stype == conf.STYPE_BITMEX:
            self.medium_level = BITMEX_LEVEL_DICT[self.small_level]
            self.big_level = BITMEX_LEVEL_DICT[self.medium_level]
        return

    def prepare(self):
        """
        初始化数据
        """
        update_dict = {
            DF_SMALL: self.small_level,
            DF_MEDIUM: self.medium_level,
            DF_BIG: self.big_level,
        }
        num_dict = {
            self.small_level: 180,
            self.medium_level: 48,
            self.big_level: 48,
        }
        for key in update_dict:
            ktype = update_dict[key]
            file_num = num_dict[ktype]
            direct_turn = False
            if ktype in [conf.BINSIZE_ONE_MINUTE]:
                direct_turn = False
            df = trend.get_from_file(ktype, conf.STYPE_BITMEX, self.code, self.factor_macd_range, file_num, direct_turn)
            setattr(self, key, df)
        if self.backtest is False:
            self.update()
        return

    def update(self):
        update_dict = {
            DF_SMALL: self.small_level,
            DF_MEDIUM: self.medium_level,
            DF_BIG: self.big_level,
        }
        for key in update_dict:
            # 以df最后的时间为准
            last_date = getattr(self, key).iloc[-1][conf.HDF5_SHARE_DATE_INDEX]
            ktype = update_dict[key]
            pull_flag = tradetime.check_pull_time(last_date, ktype)
            if pull_flag is False:
                continue
            new_df = trend.get_from_remote(ktype, conf.STYPE_BITMEX, last_date, self.code, self.rewrite)

            # 更新macd趋势列表
            trend_df = getattr(self, key)
            df_length = len(trend_df)
            direct_turn = False
            if ktype in [conf.BINSIZE_ONE_MINUTE]:
                direct_turn = False
            trend_df = trend.append_and_macd(trend_df, new_df, last_date, self.factor_macd_range, direct_turn)
            setattr(self, key, trend_df.tail(df_length).reset_index(drop=True))
        return

    def check_all(self):
        """
        遍历子级别macd趋势，并罗列买卖信号
        """
        trend_df = self.small
        for i in range(3, len(trend_df) + 1):
            self.small = trend_df.head(i)
            result = self.check_new()
            if result is not False:
                self.output()
        self.small = trend_df
        return

    def check_new(self):
        """
        检查最新子级别macd趋势
        """
        # 更新趋势聚合
        self.merge_phase()

        now = self.small.iloc[-1]
        now_date = now[conf.HDF5_SHARE_DATE_INDEX]

        # 检查medium与big是否存在trend背离
        check_dict = [DF_MEDIUM, DF_BIG]
        for df_name in check_dict:
            check_df = getattr(self, df_name)
            result = trend.check_reverse(now_date, check_df)
            setattr(self, df_name + "_" + TREND_REVERSE, result)

        # TODO 检查small的phase背离

        # 获取最新的phase状态
        phase_now = self.phase.iloc[-1]
        phase_status = phase_now[action.INDEX_PHASE_STATUS]
        phase_date = phase_now[conf.HDF5_SHARE_DATE_INDEX]

        ret = False
        if phase_date == now_date:
            ret = self._get_side(phase_status)
        return ret

    def merge_phase(self):
        self.phase = phase.latest_dict(self.small, self.phase)
        return

    def output(self):
        """
        打印结果
        """
        trend_now = self.small.iloc[-1]
        phase_now = self.phase.iloc[-1]
        # phase_pre = self.phase.iloc[-2]
        side = self._get_side(phase_now[action.INDEX_PHASE_STATUS])
        console.write_msg("【%s, %s, %s】" % (self.code, trend_now[conf.HDF5_SHARE_DATE_INDEX], side))
        if len(self.phase) > 1:
            phase_pre = self.phase.iloc[-2]
            msg = "上阶段，macd差值%f，price差值%d，连续%d次，dif位置%f"
            console.write_msg(msg % (
                phase_pre[phase.MACD_DIFF],
                phase_pre[phase.PRICE_END] - phase_pre[phase.PRICE_START],
                phase_pre[phase.COUNT],
                phase_pre[phase.DIF_END]))
        msg = "建议价格%d - %d"
        console.write_msg(msg % (phase_now[phase.PRICE_START], trend_now["close"]))
        msg = "medium级别背离情况：%s"
        console.write_msg(msg % (self.medium_trend_reverse))
        msg = "big级别背离情况：%s"
        console.write_msg(msg % (self.big_trend_reverse))
        console.write_blank()
        return

    def _get_side(self, phase_status):
        if phase_status == action.STATUS_UP:
            ret = conf.BUY_SIDE
        else:
            ret = conf.SELL_SIDE
        return ret
