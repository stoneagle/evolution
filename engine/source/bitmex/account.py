from library import conf, tool, tradetime, bitmexClient
WALLET_COLS = ['status', 'address', 'amount', 'fee', 'date', 'balance']


def wallet_history(count, start):
    """
    钱包历史列表
    """
    client = bitmexClient.Client(conf.BITMEX_URL_WALLET_HISTORY)
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
        row_dict['date'] = tradetime.transfer_iso_datetime(one['timestamp'], "M")
        df = df.append(row_dict, ignore_index=True)
    return df


def wallet(currency=conf.BITMEX_CURRENCY_XBT):
    """
    钱包状态
    """
    client = bitmexClient.Client(conf.BITMEX_URL_WALLET)
    params = {
        "currency": currency,
    }
    data_json = client.get(params)
    row_dict = dict()
    row_dict['date'] = tradetime.transfer_iso_datetime(data_json['timestamp'], "M")
    row_dict['amount'] = round(data_json['amount'] / 1000 / 1000 / 100, 4)
    return row_dict


def position(count):
    """
    获取仓位列表
    """
    # TODO 优化成支持多个symbol
    client = bitmexClient.Client(conf.BITMEX_URL_POSITION)
    params = {
        "count": count,
    }
    data_json = client.get(params)
    row_dict = dict()
    row_dict['date'] = tradetime.transfer_iso_datetime(data_json[0]['openingTimestamp'], "M")
    row_dict['symbol'] = data_json[0]['symbol']
    row_dict['is_open'] = data_json[0]['isOpen']
    row_dict['price'] = data_json[0]['avgEntryPrice']
    row_dict['current'] = data_json[0]['currentQty']
    return row_dict
