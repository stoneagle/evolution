import calendar
from library import conf
from datetime import datetime, timedelta, date


def get_last_day_of_last_month():
    # 获取上月最后一天
    d = datetime.now(conf.TZ)
    year = d.year
    month = d.month
    if month == 1:
        month = 12
        year -= 1
    else:
        month -= 1
    days = calendar.monthrange(year, month)[1]
    return datetime(year, month, 1, 16, 00, 00) + timedelta(days=days - 1)


def get_last_day_of_last_week():
    today = date.today()
    d = datetime.now(conf.TZ)
    return datetime(d.year, d.month, d.day, 16, 00, 00) - timedelta(days=today.weekday() + 1)


def get_last_of_trade_day():
    today = datetime.now(conf.TZ)
    # 需要考虑一月一号的情况
    if today.hour < 16:
        trade_datetime = datetime(today.year, today.month, today.day, 16, 00, 00) - timedelta(days=1)
    else:
        trade_datetime = datetime(today.year, today.month, today.day, 16, 00, 00)
    return trade_datetime


def get_end_date(code, ktype):
    # 如果h5文件中没有对应股票数据
    end_datetime_map = {
        # 上月最后一个交易日，15:30收盘后
        "M": LAST_MONTH_LAST_DAY,
        # 上周最后一个交易日，15:30收盘后
        "W": LAST_WEEK_LAST_DAY,
        # 如果当前时间小于15:30，则截止上一个交易日，否则则是当前交易日
        "D": LAST_TRADE_DAY + timedelta(days=1),
        "30": LAST_TRADE_DAY + timedelta(days=1),
        "5": LAST_TRADE_DAY + timedelta(days=1),
    }
    end_datetime = end_datetime_map.get(ktype, 'error')
    end_date = end_datetime.strftime('%Y-%m-%d')
    return end_date


def get_start_date(tail_date_str, code, ktype):
    time_switcher = {
        "M": "%Y-%m-%d",
        "W": "%Y-%m-%d",
        "D": "%Y-%m-%d",
        "30": "%Y-%m-%d %H:%M:%S",
        "5": "%Y-%m-%d %H:%M:%S",
    }
    tail_datetime = datetime.strptime(tail_date_str, time_switcher.get(ktype, 'error'))
    if tail_datetime.month == 12:
        month_start_year = tail_datetime.year + 1
        month_start_month = 1
    else:
        month_start_year = tail_datetime.year
        month_start_month = tail_datetime.month + 1

    start_datetime_map = {
        # 结尾日期该月的下个月第一天
        "M": date(day=1, month=month_start_month, year=month_start_year),
        # 结尾日期该周的下个周第一天
        "W": tail_datetime + timedelta(days=(7 - tail_datetime.weekday())),
        # 结尾日期加一个交易日
        "D": tail_datetime + timedelta(days=1),
        "30": tail_datetime + timedelta(days=1),
        "5": tail_datetime + timedelta(days=1),
    }
    start_datetime = start_datetime_map.get(ktype, 'error')
    start_date = start_datetime.strftime('%Y-%m-%d')
    return start_date


def get_week_of_date(date_str, ktype):
    time_switcher = {
        "M": "%Y-%m",
        "D": "%Y-%m-%d",
        "S": "%Y-%m-%d %H:%M:%S",
    }
    transfer_date = datetime.strptime(date_str, time_switcher.get(ktype, 'error'))
    return transfer_date.strftime("%Y-%W")


def transfer_date(date_str, itype, otype):
    time_switcher = {
        "M": "%Y-%m",
        # W模式无法转换
        "W": "%Y-%W",
        "D": "%Y-%m-%d",
        "S": "%Y-%m-%d %H:%M:%S",
    }
    transfer_date = datetime.strptime(date_str, time_switcher.get(itype, 'error'))
    return transfer_date.strftime(time_switcher.get(otype, 'error'))


LAST_MONTH_LAST_DAY = get_last_day_of_last_month()
LAST_WEEK_LAST_DAY = get_last_day_of_last_week()
LAST_TRADE_DAY = get_last_of_trade_day()
