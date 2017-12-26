from library import conf, tool, auth
import json
try:
    from urllib.request import urlopen, Request
except ImportError:
    from urllib2 import urlopen, Request
ORDER_COLS = ['price', 'side', 'size']


def book(symbol, depth):
    url = conf.BITMEX_HOST + conf.BITMEX_ORDERBOOK_URL
    headers = auth.bitmex_header("GET", conf.BITMEX_ORDERBOOK_URL % (symbol, depth), '')
    request = Request(url % (symbol, depth), headers=headers)
    data_str = urlopen(request, timeout=10).read()
    data_str = data_str.decode('GBK')
    data_json = json.loads(data_str)
    df = tool.init_empty_df(ORDER_COLS)
    for one in data_json:
        row_dict = dict()
        row_dict['price'] = one['price']
        row_dict['side'] = one['side']
        row_dict['size'] = one['size']
        df = df.append(row_dict, ignore_index=True)
    return df
