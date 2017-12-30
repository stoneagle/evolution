from strategy.macd import fiveMinuteWithRelate as fmwr
from library import tradetime
import threading
timer = None
strategy_dict = dict()


def exec(code_list, backtest=False):
    factor_macd_range = 0.1
    global strategy_dict

    # TODO 获取个股股东人数、解禁、股东出货等情况
    for code in code_list:
        obj = fmwr.strategy(code, fmwr.STYPE_TUSHARE, factor_macd_range, not backtest)
        strategy_dict[code] = obj
        # 初始化数据
        obj.prepare()
        # 初始检查
        obj.check_all()
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
            obj.output()

    remain_second = tradetime.get_remain_second("5")
    timer = threading.Timer(remain_second + 15, monitor)
    timer.start()
    return
