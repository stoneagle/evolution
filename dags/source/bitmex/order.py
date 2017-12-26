from library import conf, tool, bitmexClient
ORDER_COLS = ['price', 'side', 'size']


def book(symbol, depth):
    client = bitmexClient.Client(conf.BITMEX_ORDERBOOK_URL)
    params = {
        "symbol": symbol,
        "depth": depth
    }
    data_json = client.get(params)
    df = tool.init_empty_df(ORDER_COLS)
    for one in data_json:
        row_dict = dict()
        row_dict['price'] = one['price']
        row_dict['side'] = one['side']
        row_dict['size'] = one['size']
        df = df.append(row_dict, ignore_index=True)
    return df


def cancel():
    """
    取消订单
    """
    return


def cancel_all():
    """
    取消全部订单
    """
    return


def create():
    """
    创建订单
    """
    return


def amend():
    """
    修改订单
    """
    return


def cancel_all_after():
    """
    标记删除
    """
    return
