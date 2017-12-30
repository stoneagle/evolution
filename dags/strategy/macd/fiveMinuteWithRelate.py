from quota.util import action
from quota import macd
from library import tool, conf, console, tradetime
import tushare as ts
import h5py
import time
DF_ONE = "one"
DF_FIVE = "five"
DF_THIRTY = "thirty"
DF_BIG = "big"
DF_INDEX = "index"
TRADE_TYPE = "type"
TRADE_POSITIONS = "positions"
TRADE_STATUS = "_status"
TRADE_MACD_PHASE = "_macd_phase"
TRADE_MACD_DIFF = "_macd_diff"
TRADE_TREND_COUNT = "_trend_count"
STYPE_TUSHARE = "tushare"
STYPE_BITMEX = "bitmex"
DF_FILE_NUM = 128
DF_START_NUM = 26


class strategy(object):
    """
    5min级别的macd趋势震荡策略
    核心:
    1. 状态判断(震荡因子):5min的macd-shake幅度
    2. 仓位控制(比例因子):大方向(例如a股日线,bitmex的4h线)、30min级别、指数的30min级别
    3. 交易时机(未来因子):1min的背离
    """
    # 用于存储code名称
    code = None
    # 数据来源类别
    stype = None
    # check时是否更新数据源
    update = None
    # 5min是否处于shake状态
    five_shake = False
    # 5min的macd趋势是否处于背离
    five_trend_reverse = False
    # 1min的macd阶段是否存在背离
    one_phase_reverse = False

    # 存储1m数据
    one = None
    # 存储5m数据
    five = None
    # 存储30m数据
    thirty = None
    # 存储大方向(例如a股日线、bitmex的4h线)数据
    big = None
    # 存储相关指数的30m数据
    index = None
    # 存储处理得到的买卖点
    trade = tool.init_empty_df(None)
    # 震荡因子，shake幅度百分比，默认10%
    factor_macd_range = 0.1

    def __init__(self, code, stype, factor_macd_range=None, update=True):
        self.code = code
        self.stype = stype
        self.update = update
        if factor_macd_range is not None:
            self.factor_macd_range = factor_macd_range
        return

    def prepare(self):
        """
        初始化数据
        """
        if self.stype == STYPE_TUSHARE:
            ktype_dict = {
                DF_FIVE: conf.KTYPE_FIVE,
                DF_THIRTY: conf.KTYPE_THIRTY,
                DF_BIG: conf.KTYPE_DAY,
                DF_INDEX: conf.KTYPE_THIRTY,
            }
            code_prefix = self.code[0:1]
            if code_prefix == "0":
                index = "sz"
            elif code_prefix == "6":
                index = "sh"
        elif self.stype == STYPE_BITMEX:
            ktype_dict = {
                DF_FIVE: conf.BINSIZE_FIVE_MINUTE,
                DF_THIRTY: conf.BINSIZE_THIRTY_MINUTE,
                DF_BIG: conf.BINSIZE_FOUR_HOUR,
                DF_INDEX: conf.BINSIZE_THIRTY_MINUTE,
            }
            index = conf.BITMEX_BXBT
        else:
            raise Exception("stype未配置或配置错误")

        for key in ktype_dict:
            ktype = ktype_dict[key]
            if key == DF_INDEX:
                setattr(self, key, self.get_trend_from_file(ktype, index))
            else:
                setattr(self, key, self.get_trend_from_file(ktype))

        if self.update is True:
            self.update()
        return

    def update(self):
        if self.stype == STYPE_TUSHARE:
            ktype_dict = {
                DF_FIVE: conf.KTYPE_FIVE,
                DF_THIRTY: conf.KTYPE_THIRTY,
                DF_BIG: conf.KTYPE_DAY,
                DF_INDEX: conf.KTYPE_THIRTY,
            }
            code_prefix = self.code[0:1]
            if code_prefix == "0":
                index = "sz"
            elif code_prefix == "6":
                index = "sh"
        elif self.stype == STYPE_BITMEX:
            ktype_dict = {
                DF_FIVE: conf.BINSIZE_FIVE_MINUTE,
                DF_THIRTY: conf.BINSIZE_THIRTY_MINUTE,
                DF_BIG: conf.BINSIZE_FOUR_HOUR,
                DF_INDEX: conf.BINSIZE_THIRTY_MINUTE,
            }
            index = conf.BITMEX_BXBT
        else:
            raise Exception("stype未配置或配置错误")

        for key in ktype_dict:
            # 以df最后的时间为准
            last_date = getattr(self, key).iloc[-1][conf.HDF5_SHARE_DATE_INDEX]
            # 根据时间节点判断是否拉取
            ktype = ktype_dict[key]
            pull_flag = tradetime.check_pull_time(last_date, ktype)
            if pull_flag is False:
                continue
            if key == DF_INDEX:
                new_df = self.get_trend_from_remote(ktype, last_date, index)
            else:
                new_df = self.get_trend_from_remote(ktype, last_date)
            setattr(self, key, new_df)
            # TODO 将读取的数据，写回文件
        return

    def check_all(self):
        """
        遍历5min-macd所有趋势，罗列历史买卖信号
        """
        trend_df = self.five
        for i in range(3, len(trend_df) + 1):
            tmp_df = trend_df.head(i)
            self.five = tmp_df
            result = self.check_new()
            if result is True:
                self.save_trade()
                self.output()
        self.five = trend_df
        return

    def check_new(self):
        """
        检查5min-macd最新趋势，是否存在买卖信号
        """
        if self.update is True:
            self.update()
        start, shake_before, pre, now = self._get_phase_row()
        if start is None:
            return False

        # 当前macd趋势的macd波动幅度
        phase_range = abs(shake_before["macd"] - start["macd"])
        # 当前macd趋势的price波动幅度
        price_range = pre["close"] - start["close"]

        # 背离分析，如果macd是下降趋势，但是价格上涨，则属于背离；反之同理
        self.five_trend_reverse = False
        if (phase_range > 0 and price_range < 0) or (phase_range < 0 and price_range > 0):
            self.five_trend_reverse = True

        ret = False
        if self.five_shake is False:
            # 判断第一次出现的波动，如果该段趋势太小则没有做T必要
            # TODO phase_range的大小判断，避免判断太小的range趋势(目前用的是旧方案，trend数量)
            if shake_before[action.INDEX_TREND_COUNT] >= 4 and now[action.INDEX_STATUS] == action.STATUS_SHAKE:
                # 检查macd波动幅度
                macd_diff = abs(now["macd"] - pre["macd"])
                if macd_diff > self.factor_macd_range * phase_range:
                    # 如果macd波动值超出范围，视为转折
                    ret = True
                else:
                    # 如果macd波动值在范围以内，则视为波动，观察后续走势
                    self.five_shake = True
        else:
            # 判断波动过程
            if now["status"] != action.STATUS_SHAKE:
                # 波动结束趋势逆转
                self.five_shake = False
                if now["status"] != shake_before[action.INDEX_STATUS]:
                    ret = True
            else:
                macd_diff = abs(now["macd"] - shake_before["macd"])
                # 波动超出边界，震荡结束
                if macd_diff > self.factor_macd_range * phase_range:
                    self.five_shake = False
                    ret = True
        return ret

    def reverse(self):
        """
        定位1min-phase背离
        """
        return

    def save_trade(self):
        """
        存储买卖点，并输出结果
        """
        # 获取当前5min非shake的状态
        trade_dict = dict()
        start, shake_before, pre, now = self._get_phase_row()
        phase_range = abs(shake_before["macd"] - start["macd"])

        if pre[action.INDEX_STATUS] == action.STATUS_SHAKE:
            shake_pre_status = shake_before[action.INDEX_STATUS]
            macd_diff = abs(now["macd"] - shake_before["macd"])
        else:
            shake_pre_status = pre[action.INDEX_STATUS]
            macd_diff = abs(now["macd"] - pre["macd"])

        # 判断交易点性质
        if shake_pre_status == action.STATUS_UP:
            if self.five_trend_reverse is True:
                trade_type = "背离买点"
            else:
                trade_type = "正常卖点"
        else:
            if self.five_trend_reverse is True:
                trade_type = "背离卖点"
            else:
                trade_type = "正常买点"
        trade_dict = dict()
        trade_dict[conf.HDF5_SHARE_DATE_INDEX] = now[conf.HDF5_SHARE_DATE_INDEX]
        trade_dict[TRADE_TYPE] = trade_type
        trade_dict[DF_FIVE + TRADE_STATUS] = shake_pre_status
        trade_dict[DF_FIVE + TRADE_MACD_PHASE] = phase_range
        trade_dict[DF_FIVE + TRADE_MACD_DIFF] = round(macd_diff * 100 / phase_range, 0)
        trade_dict[DF_FIVE + TRADE_TREND_COUNT] = shake_before[action.INDEX_TREND_COUNT]

        # 获取仓位与状态
        positions = 0
        dtype_dict = {
            DF_THIRTY: "个股30min",
            DF_INDEX: "指数30min",
            DF_BIG: "个股大趋势",
        }
        for dtype in dtype_dict:
            status, trend_count, macd_diff, phase_range = self.get_relate_status(dtype)
            positions = self.count_positions(shake_pre_status, status, positions)
            trade_dict[dtype + TRADE_STATUS] = status
            trade_dict[dtype + TRADE_MACD_PHASE] = phase_range
            trade_dict[dtype + TRADE_MACD_DIFF] = round(macd_diff * 100 / phase_range, 0)
            trade_dict[dtype + TRADE_TREND_COUNT] = trend_count
        trade_dict[TRADE_POSITIONS] = positions
        # 更新trade列表
        self.trade = self.trade.append(trade_dict, ignore_index=True)
        return

    def output(self):
        now = self.trade.iloc[-1]
        console.write_msg("【%s, %s, %s】" % (self.code, now[conf.HDF5_SHARE_DATE_INDEX], now["type"]))
        msg = "个股5min，趋势%s，连续%d次，macd趋势%f, macd差值%d%%"
        console.write_msg(msg % (
            now[DF_FIVE + TRADE_STATUS],
            now[DF_FIVE + TRADE_TREND_COUNT],
            now[DF_FIVE + TRADE_MACD_PHASE],
            now[DF_FIVE + TRADE_MACD_DIFF]
        ))
        dtype_dict = {
            DF_THIRTY: "个股30min",
            DF_INDEX: "指数30min",
            DF_BIG: "个股大趋势",
        }
        for dtype in dtype_dict:
            msg = "%s，趋势%s，连续%d次，macd趋势%f, macd差值%d%%"
            console.write_msg(msg % (
                dtype_dict[dtype],
                now[dtype + TRADE_STATUS],
                now[dtype + TRADE_TREND_COUNT],
                now[dtype + TRADE_MACD_PHASE],
                now[dtype + TRADE_MACD_DIFF]
            ))

        # TODO 背离情况要计算两次买卖点trend_count的差值
        trend_count = now[DF_FIVE + TRADE_TREND_COUNT]
        upper_estimate = (trend_count + 1) * 5
        lower_estimate = (trend_count - 1) * 5
        if self.stype == STYPE_TUSHARE:
            remain_seconds = tradetime.get_ashare_remain_second(now[conf.HDF5_SHARE_DATE_INDEX])
            remain_minutes = round(remain_seconds / 60, 0)
            msg = "剩余时间%d分钟，下个交易点预估需要%d-%d分钟，模式%s"
            if (upper_estimate * 60) <= remain_seconds:
                trade_opportunity = "T+0"
            else:
                trade_opportunity = "T+1"
            console.write_msg(msg % (remain_minutes, lower_estimate, upper_estimate, trade_opportunity))
        elif self.stype == STYPE_BITMEX:
            msg = "下个交易点预估需要%d-%d分钟"
            console.write_msg(msg % (lower_estimate, upper_estimate))

        console.write_msg("建议仓位：%d/6" % (now[TRADE_POSITIONS]))
        console.write_blank()
        return

    def count_positions(self, pre_status, status, positions):
        """
        仓位估算
        考虑相关趋势的影响，相反方向会产生压制，例如relate上升对5min卖点，relate下降对5min买点
        """
        # 检查big 30min index的趋势,计算仓位
        if pre_status == action.STATUS_UP:
            if status == action.STATUS_DOWN:
                positions += 2
            elif status == action.STATUS_UP:
                positions += 0
            else:
                positions += 1
        else:
            if status == action.STATUS_UP:
                positions += 2
            elif status == action.STATUS_DOWN:
                positions += 0
            else:
                positions += 1
        return positions

    def get_relate_status(self, dtype):
        """
        获取index, big, 30min的最新状态
        """
        now_date = self.five.iloc[-1][conf.HDF5_SHARE_DATE_INDEX]
        if dtype == DF_BIG:
            trend_df = self.big
        elif dtype == DF_THIRTY:
            trend_df = self.thirty
        elif dtype == DF_INDEX:
            trend_df = self.index
        trend_df = trend_df[trend_df[conf.HDF5_SHARE_DATE_INDEX] <= now_date]

        start, shake_before, pre, now = self._get_phase_row(trend_df)
        # 如果最新状态是震荡，则追溯至非震荡状态的时间点，非震荡则直接获取状态
        if now[action.INDEX_STATUS] == action.STATUS_SHAKE:
            status = shake_before[action.INDEX_STATUS] + "-" + now[action.INDEX_STATUS]
            macd_diff = abs(now["macd"] - shake_before["macd"])
            phase_range = abs(shake_before["macd"] - start["macd"])
            trend_count = shake_before[action.INDEX_TREND_COUNT]
        else:
            macd_diff = abs(now["macd"] - pre["macd"])
            phase_range = abs(shake_before["macd"] - start["macd"])
            status = now[action.INDEX_STATUS]
            trend_count = now[action.INDEX_TREND_COUNT]
        return status, trend_count, macd_diff, phase_range

    def get_trend_from_remote(self, ktype, start_date, code=None):
        """
        从远端获取最新数据，并计算趋势
        """
        if code is None:
            code = self.code

        # TODO (重要)，支持ip池并发获取，要不然多code的高频获取过于缓慢
        while True:
            time.sleep(conf.REQUEST_BLANK)
            if self.stype == STYPE_TUSHARE:
                df = ts.get_k_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date)
            if df is not None and df.empty is not True:
                break
            else:
                console.write_msg("无法获取" + code + "-" + ktype + ":" + start_date + "以后的数据，休息30秒重新获取")
                time.sleep(30)
        return macd.value_and_trend(df, self.factor_macd_range)

    def get_trend_from_file(self, ktype, code=None):
        """
        从文件获取历史数据，并计算趋势
        """
        if code is None:
            code = self.code
        df = None
        if self.stype == STYPE_TUSHARE:
            if code.isdigit():
                f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
                code_prefix = self.code[0:3]
                path = '/' + code_prefix + '/' + code
            else:
                f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
                path = '/' + code

            if f.get(path) is not None:
                df = tool.df_from_dataset(f[path], ktype, None)
                df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
                df = df.tail(DF_FILE_NUM)
            else:
                raise Exception(code + "的" + ktype + "文件数据不存在")
            f.close
        elif self.stype == STYPE_BITMEX:
            f = h5py.File(conf.HDF5_FILE_FUTURE, 'a')
        else:
            raise Exception("数据源不存在或未配置")
        return macd.value_and_trend(df, self.factor_macd_range)

    def _get_phase_row(self, trend_df=None):
        """
        获取5min的不同时机的row
        """
        if trend_df is None:
            trend_df = self.five
        # 最新bar
        now = trend_df.iloc[-1]
        # 次新bar
        pre = trend_df.iloc[-2]
        # 当前macd趋势开始震荡前的bar
        trend_no_shake_df = trend_df[trend_df[action.INDEX_STATUS] != action.STATUS_SHAKE]
        shake_before = trend_no_shake_df.iloc[-1]
        # 当前macd趋势的开始bar的数据
        if abs(1 + shake_before[action.INDEX_TREND_COUNT]) < len(trend_no_shake_df):
            start = trend_no_shake_df.iloc[-1 - shake_before[action.INDEX_TREND_COUNT]]
        else:
            start = None
        return start, shake_before, pre, now

    def _valid_exist(self):
        if self.five is None:
            raise Exception("5min的数据不存在")
        if self.thirty is None:
            raise Exception("30min的数据不存在")
        if self.big is None:
            raise Exception("big的数据不存在")
        if self.index is None:
            raise Exception("index的数据不存在")
        return
