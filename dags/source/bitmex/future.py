from library import conf, tool, tradetime, bitmexClient
from datetime import datetime
SYMBOL_COLS = ['date', 'open', 'high', 'close', 'low', 'volume']
BINSIZE_ONE_MINUTE = "1m"
BINSIZE_FIVE_MINUTE = "5m"
BINSIZE_THIRTY_MINUTE = "30m"
BINSIZE_ONE_HOUR = "1h"
BINSIZE_ONE_DAY = "1d"


def history(symbol, bin_size, count, start_time=None, end_time=None):
    """
    获取交易历史
    """
    merge_flag = False
    client = bitmexClient.Client(conf.BITMEX_URL_TRADE_BUCKETED)
    # 如果是30min，基于5min数据进行聚合
    if bin_size == BINSIZE_THIRTY_MINUTE:
        bin_size = BINSIZE_FIVE_MINUTE
        count = count * 6
        merge_flag = True
    params = {
        "symbol": symbol,
        "binSize": bin_size,
        "count": count,
        "partial": False,
        "reverse": True,
    }
    if start_time is not None:
        params['startTime'] = start_time
    if end_time is not None:
        params['endTime'] = end_time
    data_json = client.get(params)
    df = tool.init_empty_df(None)
    for one in data_json:
        row_dict = dict()
        row_dict['date'] = tradetime.get_iso_datetime(one['timestamp'], "M")
        row_dict['open'] = one['open']
        row_dict['low'] = one['low']
        row_dict['high'] = one['high']
        row_dict['close'] = one['close']
        row_dict['volume'] = one['volume']
        # row_dict['turnover'] = one['turnover']
        df = df.append(row_dict, ignore_index=True)
    df["volume"] = df["volume"].astype('float64')
    df = df.reindex(index=df.index[::-1]).reset_index(drop=True)

    # 如果是30min，基于5min数据基于聚合
    merge_df = tool.init_empty_df(None)
    if merge_flag is True:
        for index, row in df.iterrows():
            datetime_obj = datetime.strptime(row[conf.HDF5_SHARE_DATE_INDEX], "%Y-%m-%d %H:%M:%S")
            if datetime_obj.minute % 30 == 0:
                if index < 5:
                    continue
                else:
                    row_dict = dict()
                    row_dict['date'] = row['date']
                    row_dict['close'] = row['close']
                    row_dict['high'] = row['high']
                    row_dict['low'] = row['low']
                    row_dict['volume'] = row['volume']
                    for i in range(0, 6):
                        one = df.iloc[index - 5 + i]
                        if i == 0:
                            row_dict['open'] = one['open']
                        row_dict['volume'] += one['volume']
                        row_dict['high'] = max(one['high'], row_dict['high'])
                        row_dict['low'] = min(one['low'], row_dict['low'])
                    merge_df = merge_df.append(row_dict, ignore_index=True)
            else:
                continue
        ret = merge_df
    else:
        ret = df
    return ret


def _get_raw_data(symbol, ktype, start, end):
    """
    根据爬虫地址获取(弃用)
    """
    client = bitmexClient.Client(conf.BITMEX_URL_HISTORY)
    params = {
        "symbol": symbol,
        "resolution": ktype,
        "from": start,
        "to": end,
    }
    data_json = client.get(params)
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
