from library import conf, tool, auth, tradetime
import json
try:
    from urllib.request import urlopen, Request
except ImportError:
    from urllib2 import urlopen, Request
WALLET_COLS = ['status', 'address', 'order_amount', 'fill_amount', 'fee', 'date', 'balance']
ORDER_COLS = ['symbol', 'date', 'price', 'side', 'type', 'value', 'amount']


def wallet_history(count, start):
    """
    钱包历史列表
    """
    url = conf.BITMEX_HOST + conf.BITMEX_WALLET_HISTORY_URL
    headers = auth.bitmex_header("GET", conf.BITMEX_WALLET_HISTORY_URL % (count, start), '')
    request = Request(url % (count, start), headers=headers)
    data_str = urlopen(request, timeout=10).read()
    data_str = data_str.decode('GBK')
    data_json = json.loads(data_str)
    df = tool.init_empty_df(WALLET_COLS)
    for one in data_json:
        row_dict = dict()
        row_dict['status'] = one['transactStatus']
        row_dict['address'] = one['address']
        row_dict['amount'] = one['amount']
        row_dict['fee'] = one['fee']
        row_dict['balance'] = one['walletBalance']
        row_dict['date'] = tradetime.get_iso_datetime(one['timestamp'], "M")
        df = df.append(row_dict, ignore_index=True)
    return df


def order_history(symbol, count, start):
    """
    交易历史列表
    """
    url = conf.BITMEX_HOST + conf.BITMEX_ORDER_LIST_URL
    headers = auth.bitmex_header("GET", conf.BITMEX_ORDER_LIST_URL % (symbol, count, start), '')
    request = Request(url % (symbol, count, start), headers=headers)
    data_str = urlopen(request, timeout=10).read()
    data_str = data_str.decode('GBK')
    data_json = json.loads(data_str)
    df = tool.init_empty_df(ORDER_COLS)
    for one in data_json:
        row_dict = dict()
        row_dict['symbol'] = one['symbol']
        row_dict['date'] = tradetime.get_iso_datetime(one['timestamp'], "M")
        row_dict['price'] = one['price']
        row_dict['side'] = one['side']
        row_dict['type'] = one['ordType']
        row_dict['value'] = one['simpleCumQty']
        row_dict['order_amount'] = one['orderQty']
        row_dict['fill_amount'] = one['cumQty']
        df = df.append(row_dict, ignore_index=True)
    return df
