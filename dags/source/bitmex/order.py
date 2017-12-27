from library import conf, tool, bitmexClient, tradetime
SIDE_BUY = "Buy"
SIDE_SELL = "Sell"
TYPE_MARKET = "Market"
TYPE_LIMIT = "Limit"
TYPE_STOP = "Stop"
TYPE_STOP_LIMIT = "StopLimit"
TYPE_MARKET_IF_TOUCH = "MarketIfTouched"
TYPE_LIMIT_IF_TOUCH = "LimitIfTouched"
TYPE_MARKET_WITH_LEFT_OVER_AS_LIMIT = "MarketWithLeftOverAsLimit"
TYPE_PEGGED = "Pegged"


def book(symbol, depth):
    client = bitmexClient.Client(conf.BITMEX_URL_ORDERBOOK)
    params = {
        "symbol": symbol,
        "depth": depth
    }
    data_json = client.get(params)
    df = tool.init_empty_df(None)
    for one in data_json:
        row_dict = dict()
        row_dict['price'] = one['price']
        row_dict['side'] = one['side']
        row_dict['size'] = one['size']
        df = df.append(row_dict, ignore_index=True)
    return df


def list(symbol, count, start, filter_dict=None):
    """
    交易历史列表
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER)
    params = {
        "symbol": symbol,
        "count": count,
        "start": start,
        "reverse": True
    }
    if filter_dict is not None:
        params['filter'] = filter_dict
    data_json = client.get(params)
    df = tool.init_empty_df(None)
    for one in data_json:
        row_dict = dict()
        row_dict['date'] = tradetime.get_iso_datetime(one['timestamp'], "M")
        row_dict['symbol'] = one['symbol']
        row_dict['side'] = one['side']
        row_dict['type'] = one['ordType']
        row_dict['status'] = one['ordStatus']
        row_dict['price'] = one['price']
        row_dict['order'] = one['simpleOrderQty']
        row_dict['cum'] = one['simpleCumQty']
        row_dict['id'] = one['orderID']
        df = df.append(row_dict, ignore_index=True)
    return df


def create(symbol, side, simple_qty, otype=TYPE_LIMIT, price=None):
    """
    创建订单
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER)
    content = {
        "symbol": symbol,
        "side": side,
        "simpleOrderQty": simple_qty,
        "ordType": otype,
    }
    if otype == TYPE_LIMIT or otype == TYPE_STOP_LIMIT or otype == TYPE_LIMIT_IF_TOUCH:
        content["price"] = price
    return client.post(content)


def cancel(oid):
    """
    取消订单
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER)
    content = {
        "orderID": oid,
    }
    return client.delete(content)


def cancel_all(symbol):
    """
    取消全部订单
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER_CANCEL_ALL)
    content = {
        "symbol": symbol,
    }
    return client.delete(content)


def amend(oid, simple_qty, price=None):
    """
    修改订单
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER)
    content = {
        "orderID": oid,
        "simpleOrderQty": simple_qty,
    }
    if price is not None:
        content["price"] = price
    return client.put(content)


def cancel_all_after(timeout):
    """
    标记删除
    timeout: ms
    """
    client = bitmexClient.Client(conf.BITMEX_URL_ORDER_CANCEL_ALL_AFTER)
    content = {
        "timeout": timeout,
    }
    return client.post(content)
