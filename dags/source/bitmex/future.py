from library import conf, tool, tradetime
import json
try:
    from urllib.request import urlopen, Request
except ImportError:
    from urllib2 import urlopen, Request
SYMBOL_COLS = ['date', 'open', 'high', 'close', 'low', 'volume']


def latest(symbol, ktype, start_date):
    start = tradetime.get_unixtime(start_date, ktype)
    end = tradetime.get_unixtime()
    df = _get_raw_data(symbol, ktype, start, end)
    df["volume"] = df["volume"].astype('float64')
    # 忽略最后一条未确定的bar
    return df.head(len(df) - 1)


def trade_history():
    """
    获取交易历史
    """
    return


def _get_raw_data(symbol, ktype, start, end):
    url = conf.BITMEX_HOST + conf.BITMEX_HISTORY_URL
    headers = auth.bitmex_header("GET", conf.BITMEX_HISTORY_URL % (symbol, ktype, start, end), '')
    request = Request(url % (symbol, ktype, start, end), headers=headers)
    data_str = urlopen(request, timeout=10).read()
    data_str = data_str.decode('GBK')
    data_json = json.loads(data_str)
    df = tool.init_empty_df(SYMBOL_COLS)
    if data_json['s'] == 'ok':
        open_arr = data_json['o']
        close_arr = data_json['c']
        high_arr = data_json['h']
        low_arr = data_json['l']
        time_arr = data_json['t']
        volume_arr = data_json['v']

        for index in range(0, len(time_arr)):
            row_dict = dict()
            row_dict['date'] = tradetime.transfer_unixtime(time_arr[index], ktype)
            row_dict['open'] = open_arr[index]
            row_dict['close'] = close_arr[index]
            row_dict['high'] = high_arr[index]
            row_dict['low'] = low_arr[index]
            row_dict['volume'] = volume_arr[index]
            df = df.append(row_dict, ignore_index=True)
    return df
