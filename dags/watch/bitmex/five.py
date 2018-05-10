from strategy.macd import reverseAndLever as ral
from library import tradetime, conf, console
from source.bitmex import account, order
import threading
import os
import sys
monitor_timer = None
wallet_amount = None
strategy_dict = dict()
exec_level = conf.BINSIZE_FIVE_MINUTE


def exec(code_list, backtest, rewrite):
    global strategy_dict
    global wallet_amount
    global current_position
    global exec_level
    factor_macd_range = 0.1

    # 初始化钱包数额
    wallet_detail = account.wallet()
    wallet_amount = wallet_detail['amount']

    for code in code_list:
        console.write_msg("初始化%s数据，进行%s级别操作" % (code, exec_level))
        obj = ral.strategy(code, conf.STYPE_BITMEX, backtest, rewrite, exec_level, factor_macd_range)
        # 初始化数据
        obj.prepare()
        # print(obj.small.tail(80)[["date", "close", "macd", "trend_count", "phase_status"]])
        # 初始检查
        # obj.check_all()
        strategy_dict[code] = obj
    if backtest is False:
        monitor()
    return


def monitor():
    global monitor_timer
    global wallet_amount
    global strategy_dict
    global exec_level
    # 仓位倍率
    factor_multi = 0.01
    factor_price_range = 5
    for code in strategy_dict:
        try:
            obj = strategy_dict[code]

            # 更新数据并检查
            obj.update()
            ret = obj.check_new()
            console.write_msg("时间%s，bar的操作结果:%s" % (obj.small.iloc[-1]["date"], ret))
            if ret is not False:
                obj.output()
                # 1.结算上次交易
                # 如果仓位小于0，则卖空需要赎回，如果仓位大于0，则做多需要卖出
                current_position = account.position(10)
                current = current_position["current"]
                if current != 0:
                    if current > 0:
                        side = order.SIDE_SELL
                    elif current < 0:
                        side = order.SIDE_BUY
                    result = order.create_contract(code, side, abs(current), order.TYPE_MARKET)
                    if result["price"]:
                        msg = "合约%s兑现, 操作%s, 份额%f，价格%f, 下单成功"
                        console.write_msg(msg % (code, side, abs(current), result["price"]))

                # 2.更新钱包数额(暂时不更新钱包)
                # wallet_detail = account.wallet()
                # wallet_amount = wallet_detail['amount']

                # 3.判断交易类别，进行做空或做多
                # 方案A，取两者的平均值，基本很难触发(已舍弃)
                # 方案B，取最新价格的偏差，目前固定为5(低级别不适用)
                trend_price = obj.small.iloc[-1]["close"]
                amount = wallet_amount * factor_multi
                if ret == conf.BUY_SIDE:
                    price = trend_price - factor_price_range
                    result = order.create_simple(code, order.SIDE_BUY, amount, order.TYPE_LIMIT, price)
                elif ret == conf.SELL_SIDE:
                    price = trend_price + factor_price_range
                    result = order.create_simple(code, order.SIDE_SELL, amount, order.TYPE_LIMIT, price)
                else:
                    raise Exception("交易类别异常")
                if result["price"]:
                    msg = "合约%s, 操作%s, 份额%f，价格%f, 下单成功"
                    console.write_msg(msg % (code, ret, amount, price))
                else:
                    print(result)

                # 4.标记过期时间
                result = order.cancel_all_after(30000)
                if result["cancelTime"]:
                    console.write_msg("预计过期时间:%s" % (result["cancelTime"]))
                else:
                    print(result)
        except Exception as er:
            print(obj.small.tail(10))
            print(obj.phase.tail(10))
            print(str(er))
            exc_type, exc_obj, exc_tb = sys.exc_info()
            fname = os.path.split(exc_tb.tb_frame.f_code.co_filename)[1]
            print(exc_type, fname, exc_tb.tb_lineno)
            sys.exit(2)
    if exec_level == conf.BINSIZE_FIVE_MINUTE:
        time_span = "5"
    elif exec_level == conf.BINSIZE_ONE_MINUTE:
        time_span = "1"
    else:
        raise Exception("级别异常")
    remain_second = tradetime.get_remain_second(time_span)
    monitor_timer = threading.Timer(remain_second + 2, monitor)
    monitor_timer.start()
    return
