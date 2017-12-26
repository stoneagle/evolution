from library import conf, tool, tradetime, bitmexClient
WALLET_COLS = ['status', 'address', 'amount', 'fee', 'date', 'balance']
ORDER_COLS = ['symbol', 'date', 'price', 'side', 'type', 'value', 'order_amount', 'fill_amount']


def wallet_history(count, start):
    """
    钱包历史列表
    """
    client = bitmexClient.Client(conf.BITMEX_WALLET_HISTORY_URL)
    params = {
        "count": count,
        "start": start,
        "reverse": True
    }
    data_json = client.get(params)
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
    print(df)
    return df


def order_history(symbol, count, start):
    """
    交易历史列表
    """
    client = bitmexClient.Client(conf.BITMEX_ORDER_LIST_URL)
    params = {
        "symbol": symbol,
        "count": count,
        "start": start,
        "reverse": True
    }
    data_json = client.get(params)
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
