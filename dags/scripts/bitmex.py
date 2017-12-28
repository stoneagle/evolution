from controller import obtain
from library import conf
from source.bitmex import future


def test():
    # obtain
    # symbol_list = [conf.BITMEX_XBTUSD, conf.BITMEX_BXBT]
    # for symbol in symbol_list:
    #     obtain.bitmex(symbol, future.BINSIZE_ONE_DAY, 300)
    #     obtain.bitmex(symbol, future.BINSIZE_FIVE_MINUTE, 300)
    #     obtain.bitmex(symbol, future.BINSIZE_FOUR_HOUR, 300)
    #     obtain.bitmex(symbol, future.BINSIZE_THIRTY_MINUTE, 300)

    # future
    # result = future.history(conf.BITMEX_BXBT, future.BINSIZE_THIRTY_MINUTE, 3)
    # result = future.history_thirty_minute(conf.BITMEX_XBTUSD, 100)
    # result = future.history_four_hour(conf.BITMEX_XBTUSD, 100)
    # print(result)

    # order
    # result = order.book(conf.BITMEX_XBTUSD, 10)
    # result = order.list(conf.BITMEX_XBTUSD, 10, 0)
    # result = order.create(conf.BITMEX_XBTUSD, order.SIDE_BUY, 0.001, order.TYPE_LIMIT, 10000.00)
    # result = order.amend('1c1bfbe6-6c2c-3d4a-3271-5145d3171a28', 0.002, 8000.00)
    # result = order.cancel_all_after(10000)
    # result = order.cancel('1c1bfbe6-6c2c-3d4a-3271-5145d3171a28')
    # result = order.cancel_all(conf.BITMEX_XBTUSD)

    # account
    # account.wallet_history(100, 0)
    # result = account.wallet()
    # print(result)
    return
