from controller import watch, obtain
from library import conf
# from source.bitmex import future, account


def test():
    # obtain
    symbol_list = [conf.BITMEX_XBTUSD, conf.BITMEX_BXBT]
    for symbol in symbol_list:
        obtain.bitmex(symbol, conf.BINSIZE_ONE_MINUTE, 7200)
        obtain.bitmex(symbol, conf.BINSIZE_FIVE_MINUTE, 1440)
        obtain.bitmex(symbol, conf.BINSIZE_THIRTY_MINUTE, 240)
        obtain.bitmex(symbol, conf.BINSIZE_FOUR_HOUR, 300)
        obtain.bitmex(symbol, conf.BINSIZE_ONE_DAY, 300)

    # future
    # result = future.history(conf.BITMEX_BXBT, conf.BINSIZE_THIRTY_MINUTE, 3)
    # result = future.history_thirty_minute(conf.BITMEX_XBTUSD, 100)
    # result = future.history_four_hour(conf.BITMEX_XBTUSD, 100)
    # print(result)

    # order
    # result = order.book(conf.BITMEX_XBTUSD, 10)
    # result = order.list(conf.BITMEX_XBTUSD, 10, 0)
    # result = order.create_simple(conf.BITMEX_XBTUSD, order.SIDE_BUY, 0.001, order.TYPE_LIMIT, 10000.00)
    # result = order.create_contract(conf.BITMEX_XBTUSD, order.SIDE_BUY, 10, order.TYPE_LIMIT, 10000.00)
    # result = order.amend_simple('1c1bfbe6-6c2c-3d4a-3271-5145d3171a28', 0.002, 8000.00)
    # result = order.amend_contract('1c1bfbe6-6c2c-3d4a-3271-5145d3171a28', 20, 8000.00)
    # result = order.cancel_all_after(10000)
    # result = order.cancel('1c1bfbe6-6c2c-3d4a-3271-5145d3171a28')
    # result = order.cancel_all(conf.BITMEX_XBTUSD)

    # account
    # account.wallet_history(100, 0)
    # result = account.wallet()
    # result = account.position(10)
    # print(result)
    return


def monitor():
    code_list = [conf.BITMEX_XBTUSD]
    # watch
    watch.bitmex(code_list)
    return
