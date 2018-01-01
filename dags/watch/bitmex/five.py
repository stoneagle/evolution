from strategy.macd import fiveMinuteWithRelate as fmwr
from library import tradetime
import threading
timer = None
strategy_dict = dict()


def exec(code_list, backtest, rewrite):
    """
    根据仓位情况，选择级别
    5min决定多空方向，1min的趋势进行交易，背离决定杠杆
    例如5min单段趋势向下，即将出现1min买点，则选择做多
    1min第一个买点杠杆1x，上涨后出售，如果dea、dif未上穿零轴，则5min单段趋势继续
    当第二个背离买点出现，则杠杆增加至2x，上涨后出售，背离杠杆最多允许增加至3x
    正常出售后，当1min出现卖点背离，则5min单段趋势准备逆转，方向由做多转为做空，1min买卖同理
    """
    factor_macd_range = 0.1
    global strategy_dict

    for code in code_list:
        obj = fmwr.strategy(code, fmwr.STYPE_BITMEX, backtest, rewrite, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # 初始检查
        obj.check_all()
        strategy_dict[code] = obj
    monitor()
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
