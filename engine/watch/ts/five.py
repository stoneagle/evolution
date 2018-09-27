from strategy.macd import fiveMinuteWithRelate as fmwr
from library import tradetime, conf
import threading
timer = None
strategy_dict = dict()


def exec(code_list, backtest, rewrite):
    factor_macd_range = 0.1
    global strategy_dict

    # TODO 获取个股股东人数、解禁、股东出货等情况
    # TODO 5min与1min的背离结合
    # TODO 部分交易点，连续次数存在失真情况，例如天齐锂业的18-1-4
    for code in code_list:
        obj = fmwr.strategy(code, conf.STYPE_ASHARE, backtest, rewrite, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # 初始检查
        obj.check_all()
        # print(obj.five[["date", "macd", "status", "close", "trend_count"]])
        strategy_dict[code] = obj
    if backtest is False:
        monitor()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def monitor():
    global timer
    global strategy_dict
    for code in strategy_dict:
        obj = strategy_dict[code]
        # 更新数据
        if obj.backtest is False:
            obj.update()

        # 检查最新数据
        ret = obj.check_new()
        if ret is True:
            obj.save_trade()
            obj.output()
    remain_second = tradetime.get_remain_second("5")
    timer = threading.Timer(remain_second + 15, monitor)
    timer.start()
    return
