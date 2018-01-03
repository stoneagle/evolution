from strategy.macd import reverseAndLever as ral
from library import tradetime, conf
from source.bitmex import account
import threading
monitor_timer = None
updater_timer = None
strategy_dict = dict()
wallet_amount = None


def exec(code_list, backtest, rewrite):
    global strategy_dict
    global wallet_amount
    factor_macd_range = 0.1

    # 初始化钱包数额
    wallet_detail = account.wallet()
    wallet_amount = wallet_detail['amount']

    for code in code_list:
        obj = ral.strategy(code, conf.STYPE_BITMEX, backtest, rewrite, conf.BINSIZE_ONE_MINUTE, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # print(obj.small.tail(80)[["date", "close", "macd", "trend_count", "phase_status"]])
        # 初始检查
        obj.check_all()
        # obj.check_new()
        strategy_dict[code] = obj
    if backtest is False:
        monitor()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def monitor():
    global monitor_timer
    global strategy_dict
    for code in strategy_dict:
        obj = strategy_dict[code]
        # 更新数据
        if obj.backtest is False:
            obj.update()

        ret = obj.check_new()
        if ret is not False:
            obj.output()
            # 判断交易类别，进行做空或做多
            if ret == conf.BUY_SIDE:
                print(ret)
            elif ret == conf.SELL_SIDE:
                print(ret)
            # TODO 结算上次交易
            # TODO 开始新的交易，并标记过期时间
    remain_second = tradetime.get_remain_second("5")
    monitor_timer = threading.Timer(remain_second + 10, monitor)
    monitor_timer.start()
    return


def updater():
    global updater_timer
    global wallet_amount
    # 更新钱包数额
    wallet_detail = account.wallet()
    wallet_amount = wallet_detail['amount']

    # 检查交割单
    # 如果存在未完成的新交易，则删除
    # 如果存在未成交的旧交易，市价成交

    updater_timer = threading.Timer(30, updater)
    updater_timer.start()
    return
