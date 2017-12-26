from library import conf, tool, auth, tradetime
import json
try:
    from urllib.request import urlopen, Request
except ImportError:
    from urllib2 import urlopen, Request
WALLET_COLS = ['status', 'address', 'amount', 'fee', 'date', 'balance']


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
