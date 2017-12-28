from library import conf, tool, tradetime, bitmexClient
from datetime import datetime
SYMBOL_COLS = ['date', 'open', 'high', 'close', 'low', 'volume']
BINSIZE_ONE_MINUTE = "1m"
BINSIZE_FIVE_MINUTE = "5m"
BINSIZE_ONE_HOUR = "1h"
BINSIZE_ONE_DAY = "1d"
BINSIZE_THIRTY_MINUTE = "30m"
BINSIZE_FOUR_HOUR = "4h"
BUCKET_LIMIT = 10000
TYPE_API = "api"
TYPE_SPIDER = "spider"
SPIDER_BINSIZE_DICT = {
    BINSIZE_ONE_MINUTE: "1",
    BINSIZE_FIVE_MINUTE: "5",
    BINSIZE_THIRTY_MINUTE: "5",
    BINSIZE_ONE_HOUR: "60",
    BINSIZE_FOUR_HOUR: "60",
    BINSIZE_ONE_DAY: "D",
}


def history(symbol, bin_size, count, gtype=TYPE_SPIDER):
    """
    获取交易历史
    """
    if count >= BUCKET_LIMIT:
        return None

    if gtype == TYPE_SPIDER:
        spider_ktype = SPIDER_BINSIZE_DICT[bin_size]
        sdate = tradetime.get_date_by_barnum(count, spider_ktype)
        start = tradetime.get_unixtime(sdate)
        end = tradetime.get_unixtime()
        df = _get_data_by_spider(symbol, spider_ktype, start, end)
    elif gtype == TYPE_API:
        df = _get_data_by_api(symbol, bin_size, count)
    else:
        df = None
    return df.head(len(df) - 1)


def history_merge(symbol, bin_size, count, gtype=TYPE_SPIDER):
    """
    获取聚合类交易历史
    """
    if bin_size == BINSIZE_FOUR_HOUR:
        df = _get_merge_data(symbol, bin_size, count, 4, 'H', gtype)
        df = df.head(len(df) - 1)
        merge_df = _get_merge_df(df, 4, 'H')
    elif bin_size == BINSIZE_THIRTY_MINUTE:
        df = _get_merge_data(symbol, bin_size, count, 6, '5', gtype)
        df = df.head(len(df) - 1)
        merge_df = _get_merge_df(df, 6, '30')
    return merge_df


def _get_merge_data(symbol, bin_size, count, multi, ktype, gtype):
    """
    按倍数放大获取基础数据，并进行聚合
    """
    multi_count = count * multi
    if count >= BUCKET_LIMIT or multi_count >= BUCKET_LIMIT:
        return None

    # 获取时需要根据聚合周期，进行偏移
    if gtype == TYPE_SPIDER:
        spider_ktype = SPIDER_BINSIZE_DICT[bin_size]
        fix_barnum = tradetime.fix_merge_barnum(bin_size)
        sdate = tradetime.get_date_by_barnum(multi_count + fix_barnum, spider_ktype)
        start = tradetime.get_unixtime(sdate)
        end = tradetime.get_unixtime()
        df = _get_data_by_spider(symbol, spider_ktype, start, end)
    elif gtype == TYPE_API:
        sdate = tradetime.get_date_by_barnum(multi_count, ktype, True)
        sdate_iso = tradetime.get_iso_datetime(sdate, 'S')
        df = _get_data_by_api(symbol, bin_size, multi_count, sdate_iso)
    else:
        df = None
    return df


def _get_merge_df(df, multi, ktype):
    """
    根据获取的数据，进行聚合(需要注意聚合数据的起点)
    """
    merge_df = tool.init_empty_df(None)
    for index, row in df.iterrows():
        datetime_obj = datetime.strptime(row[conf.HDF5_SHARE_DATE_INDEX], "%Y-%m-%d %H:%M:%S")
        if (ktype == '30' and datetime_obj.minute % 30 == 0) or (ktype == 'H' and datetime_obj.hour % 4 == 0):
            if (index + multi - 1) < len(df):
                row_dict = dict()
                row_dict['date'] = row['date']
                row_dict['open'] = row['open']
                row_dict['close'] = df.iloc[index + multi - 1]['close']
                row_dict['high'] = row['high']
                row_dict['low'] = row['low']
                row_dict['volume'] = row['volume']
                for i in range(1, multi):
                    one = df.iloc[index + i]
                    row_dict['volume'] += one['volume']
                    row_dict['high'] = max(one['high'], row_dict['high'])
                    row_dict['low'] = min(one['low'], row_dict['low'])
                merge_df = merge_df.append(row_dict, ignore_index=True)
        else:
            continue
    return merge_df


def _transfer_json_to_df(data_json):
    """
    将返回的json转化为df
    """
    df = tool.init_empty_df(None)
    for one in data_json:
        row_dict = dict()
        row_dict['date'] = tradetime.transfer_iso_datetime(one['timestamp'], "M")
        row_dict['open'] = one['open']
        row_dict['low'] = one['low']
        row_dict['high'] = one['high']
        row_dict['close'] = one['close']
        row_dict['volume'] = one['volume']
        # row_dict['turnover'] = one['turnover']
        df = df.append(row_dict, ignore_index=True)
    if len(df) > 0:
        df["volume"] = df["volume"].astype('float64')
        df = df.reindex(index=df.index[::-1]).reset_index(drop=True)
    return df


def _get_data_by_api(symbol, bin_size, count, start=None, end=None):
    client = bitmexClient.Client(conf.BITMEX_URL_TRADE_BUCKETED)
    params = {
        "symbol": symbol,
        "binSize": bin_size,
        "count": count,
        "partial": False,
        "reverse": True,
    }
    if start is not None:
        params['startTime'] = start
    if end is not None:
        params['endTime'] = end
    data_json = client.get(params)
    df = _transfer_json_to_df(data_json)
    return df


def _get_data_by_spider(symbol, ktype, start, end):
    """
    根据爬虫地址获取
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
            row_dict['open'] = float(open_arr[index])
            row_dict['close'] = float(close_arr[index])
            row_dict['high'] = float(high_arr[index])
            row_dict['low'] = float(low_arr[index])
            row_dict['volume'] = volume_arr[index]
            df = df.append(row_dict, ignore_index=True)
    if len(df) > 0:
        df["volume"] = df["volume"].astype('float64')
    return df
