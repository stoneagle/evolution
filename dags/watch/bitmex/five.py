from strategy.macd import reverseAndLever as ral
from strategy.common import phase
from library import tradetime, conf, console
from source.bitmex import account, order
import threading
monitor_timer = None
wallet_amount = None
current_position = None
strategy_dict = dict()


def exec(code_list, backtest, rewrite):
    global strategy_dict
    global wallet_amount
    global current_position
    factor_macd_range = 0.1

    # 初始化钱包数额
    wallet_detail = account.wallet()
    wallet_amount = wallet_detail['amount']
    current_position = account.position(10)

    for code in code_list:
        obj = ral.strategy(code, conf.STYPE_BITMEX, backtest, rewrite, conf.BINSIZE_ONE_MINUTE, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # print(obj.small.tail(80)[["date", "close", "macd", "trend_count", "phase_status"]])
        # 初始检查
        # obj.check_all()
        # 初始化历史交易清单
        strategy_dict[code] = obj
    if backtest is False:
        monitor()
    return


def monitor():
    global current_position
    global monitor_timer
    global wallet_amount
    global strategy_dict
    # 仓位倍率
    factor_multi = 0.05
    for code in strategy_dict:
        obj = strategy_dict[code]

        # 更新数据
        if obj.backtest is False:
            obj.update()

        ret = obj.check_new()
        if ret is not False:
            obj.output()

            # 1.结算上次交易
            # 如果仓位小于0，则卖空需要赎回，如果仓位大于0，则做多需要卖出
            current = current_position["current"]
            if current != 0:
                if current > 0:
                    side = order.SIDE_SELL
                elif current < 0:
                    side = order.SIDE_BUY
                result = order.create_contract(code, side, abs(current), order.TYPE_MARKET)
                if len(result) > 0:
                    msg = "合约%s兑现, 操作%s, 份额%f，价格%f, 下单成功"
                    console.write_msg(msg % (code, side, abs(current), result["price"]))

            # 2.更新钱包数额
            wallet_detail = account.wallet()
            wallet_amount = wallet_detail['amount']

            # 3.判断交易类别，进行做空或做多
            phase_price = obj.phase.iloc[-1][phase.PRICE_START]
            trend_price = obj.small.iloc[-1]["close"]
            price = round((phase_price + trend_price) / 2, 1)
            amount = wallet_amount * factor_multi
            if ret == conf.BUY_SIDE:
                result = order.create_simple(code, order.SIDE_BUY, amount, order.TYPE_LIMIT, price)
            elif ret == conf.SELL_SIDE:
                result = order.create_simple(code, order.SIDE_SELL, amount, order.TYPE_LIMIT, price)
            else:
                raise Exception("交易类别异常")
            if len(result) > 0:
                msg = "合约%s, 操作%s, 份额%f，价格%f, 下单成功"
                console.write_msg(msg % (code, ret, amount, price))

            # 4.标记过期时间
            result = order.cancel_all_after(60000)
            if len(result) > 0:
                console.write_msg("预计过期时间:%s" % (result["cancelTime"]))
    remain_second = tradetime.get_remain_second("1")
    monitor_timer = threading.Timer(remain_second + 10, monitor)
    monitor_timer.start()
    return
