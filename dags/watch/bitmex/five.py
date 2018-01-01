# from strategy.macd import fiveMinuteWithRelate as fmwr
from strategy.macd import reverseAndLever as ral
from library import tradetime, conf
import threading
timer = None
strategy_dict = dict()


def exec(code_list, backtest, rewrite):
    global strategy_dict
    factor_macd_range = 0.1

    for code in code_list:
        obj = ral.strategy(code, conf.STYPE_BITMEX, backtest, rewrite, conf.BINSIZE_ONE_MINUTE, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # print(obj.small.head(30)[["date", "close", "macd", "trend_count", "status", "phase_status"]])
        # print(obj.phase)
        # 初始检查
        # obj.check_all()
        # obj.check_new()
        strategy_dict[code] = obj
    # monitor()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def monitor():
    global timer
    global strategy_dict
    for code in strategy_dict:
        obj = strategy_dict[code]
        ret = obj.check()
        if ret is True:
            obj.save_trade()
            # 判断交易类别，进行做空或做多
            obj.output()

    remain_second = tradetime.get_remain_second("5")
    timer = threading.Timer(remain_second + 15, monitor)
    timer.start()
    return
