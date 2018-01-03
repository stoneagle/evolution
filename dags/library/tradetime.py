import calendar
from library import conf
from datetime import datetime, timedelta, date
from dateutil import parser
import math
import time


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
        "S": "%Y-%m-%d %H:%M:%S",
        conf.KTYPE_DAY: "%Y-%m-%d",
        conf.KTYPE_THIRTY: "%Y-%m-%d %H:%M",
        conf.KTYPE_FIVE: "%Y-%m-%d %H:%M",
    }
    transfer_date = datetime.strptime(date_str, time_switcher.get(itype, 'error'))
    return transfer_date.strftime(time_switcher.get(otype, 'error'))


def get_today():
    today = datetime.now(conf.TZ)
    date = str(today.year) + "-" + str(today.month) + "-" + str(today.day)
    return date


def get_preday(pre_num):
    today = datetime.now(conf.TZ)
    preday = today + timedelta(days=-pre_num)
    return preday.strftime("%Y-%m-%d")


def get_remain_second(ttype):
    today = datetime.now(conf.TZ)
    if ttype == "5":
        remain_minute = 5 - today.minute % 5
        remain_second = remain_minute * 60 - today.second
    elif ttype == "1":
        remain_second = 60 - today.second

    return remain_second


def get_ashare_remain_second(date_str):
    transfer_date = datetime.strptime(date_str, "%Y-%m-%d %H:%M:%S")
    minute = 60 - transfer_date.minute
    if (transfer_date.hour == 9 and transfer_date.minute >= 30):
        remain_second = minute * 60 + 60 * 60 * 3.5
    elif transfer_date.hour == 10:
        remain_second = minute * 60 + 60 * 60 * 2.5
    elif (transfer_date.hour == 11 and transfer_date.minute <= 30):
        remain_second = minute * 60 + 60 * 60 * 2
    elif transfer_date.hour == 13:
        remain_second = minute * 60 + 60 * 60
    elif transfer_date.hour == 14:
        remain_second = minute * 60
    elif transfer_date.hour == 15:
        remain_second = 0
    return remain_second


def get_unixtime(date_str=None):
    if date_str is None:
        dtime = datetime.now(conf.TZ)
    else:
        dtime = datetime.strptime(date_str, "%Y-%m-%d %H:%M:%S")
    return time.mktime(dtime.timetuple())


def transfer_unixtime(unixtime, ktype):
    time_switcher = {
        "M": "%Y-%m",
        "W": "%Y-%W",
        "D": "%Y-%m-%d",
        "30": "%Y-%m-%d %H:%M:%S",
        "5": "%Y-%m-%d %H:%M:%S",
        "1": "%Y-%m-%d %H:%M:%S",
    }
    return datetime.fromtimestamp(unixtime).strftime(time_switcher.get(ktype, '%Y-%m-%d %H:%M:%S'))


def get_date_by_barnum(barnum, ktype, local=False):
    """
    根据barnum的数量，获取对应时间点
    """
    if local is False:
        d = datetime.now(conf.TZ)
    else:
        d = datetime.now(conf.TZ_UTC)
    switcher = {
        "M": 30 * 24 * 60,
        "W": 7 * 24 * 60,
        "D": 24 * 60,
        "H": 60,
        "60": 60,
        "30": 30,
        "5": 5,
        "1": 1,
    }
    diff = barnum * switcher.get(ktype, 'error')
    diff_datetime = datetime(d.year, d.month, d.day, d.hour, d.minute, d.second) - timedelta(minutes=diff)
    return diff_datetime.strftime("%Y-%m-%d %H:%M:%S")


def get_barnum_by_date(date_str, ktype, local=False):
    """
    根据barnum的数量，获取对应时间点
    """
    ntime = datetime.now()
    dtime = datetime.strptime(date_str, "%Y-%m-%d %H:%M:%S")

    if local is True:
        ntime = ntime.replace(tzinfo=conf.TZ)
        dtime = dtime.replace(tzinfo=conf.TZ)
    else:
        ntime = ntime.replace(tzinfo=conf.TZ_UTC)
        dtime = dtime.replace(tzinfo=conf.TZ_UTC)

    switcher = {
        conf.BINSIZE_ONE_DAY: 60 * 24,
        conf.BINSIZE_FOUR_HOUR: 60 * 4,
        conf.BINSIZE_ONE_HOUR: 60,
        conf.BINSIZE_THIRTY_MINUTE: 30,
        conf.BINSIZE_FIVE_MINUTE: 5,
        conf.BINSIZE_ONE_MINUTE: 1,
    }
    diff_seconds = (ntime - dtime).total_seconds()
    barnum = round(diff_seconds / 60 / switcher.get(ktype, 'error'), 0) + 1
    return barnum


def fix_merge_barnum(bin_size, date_str=None, local=False):
    """
    获取聚合的修正barnum
    """
    if date_str is None:
        dtime = datetime.now()
    else:
        dtime = datetime.strptime(date_str, "%Y-%m-%d %H:%M:%S")

    if local is True:
        dtime = dtime.replace(tzinfo=conf.TZ)
    else:
        dtime = dtime.replace(tzinfo=conf.TZ_UTC)

    if bin_size == "30m":
        diff = dtime.minute % 30
        barnum = math.ceil(diff / 5)
    elif bin_size == "4h":
        diff = dtime.hour % 4
        barnum = diff
    return barnum


def transfer_iso_datetime(date_str, itype):
    time_switcher = {
        "M": "%Y-%m",
        "W": "%Y-%W",
        "D": "%Y-%m-%d",
        "M": "%Y-%m-%d %H:%M:%S",
    }
    date = parser.parse(date_str).astimezone(conf.TZ)
    return date.strftime(time_switcher.get(itype, 'error'))


def get_iso_datetime(date_str, itype):
    time_switcher = {
        "M": "%Y-%m",
        "W": "%Y-%W",
        "D": "%Y-%m-%d",
        "S": "%Y-%m-%d %H:%M:%S",
    }
    dtime = datetime.strptime(date_str, time_switcher.get(itype, 'error'))
    return dtime.isoformat()


def check_pull_time(date_str, ktype, local=True):
    time_switcher = {
        conf.KTYPE_DAY: "%Y-%m-%d",
        conf.KTYPE_THIRTY: "%Y-%m-%d %H:%M:%S",
        conf.KTYPE_FIVE: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_ONE_MINUTE: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_FIVE_MINUTE: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_THIRTY_MINUTE: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_ONE_HOUR: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_FOUR_HOUR: "%Y-%m-%d %H:%M:%S",
        conf.BINSIZE_ONE_DAY: "%Y-%m-%d %H:%M:%S",
    }
    ntime = datetime.now()
    dtime = datetime.strptime(date_str, time_switcher.get(ktype, 'error'))

    if local is True:
        ntime = ntime.replace(tzinfo=conf.TZ)
        dtime = dtime.replace(tzinfo=conf.TZ)
    else:
        ntime = ntime.replace(tzinfo=conf.TZ_UTC)
        dtime = dtime.replace(tzinfo=conf.TZ_UTC)

    switcher = {
        conf.KTYPE_ONE: 1,
        conf.KTYPE_FIVE: 5,
        conf.KTYPE_THIRTY: 30,
        conf.KTYPE_DAY: 60 * 24,
        conf.BINSIZE_ONE_MINUTE: 1,
        conf.BINSIZE_FIVE_MINUTE: 5,
        conf.BINSIZE_THIRTY_MINUTE: 30,
        conf.BINSIZE_FOUR_HOUR: 60 * 4,
        conf.BINSIZE_ONE_DAY: 60 * 24,
    }
    diff = (ntime - dtime).total_seconds()
    span = switcher.get(ktype, 'error')
    if span == 'error':
        raise Exception(ktype + "不合法")
    if diff >= span * 60:
        ret = True
    else:
        ret = False
    return ret


def move_delta(date_str, mtype, num):
    dtime = datetime.strptime(date_str, "%Y-%m-%d %H:%M:%S")
    move_switcher = {
        "D": dtime + timedelta(days=num),
        "M": dtime + timedelta(minutes=num),
        "S": dtime + timedelta(seconds=num),
    }
    return date.strftime(move_switcher.get(mtype, 'error'), "%Y-%m-%d %H:%M:%S")


LAST_MONTH_LAST_DAY = get_last_day_of_last_month()
LAST_WEEK_LAST_DAY = get_last_day_of_last_week()
LAST_TRADE_DAY = get_last_of_trade_day()
