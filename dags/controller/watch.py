import time
import tushare as ts
from library import conf, console, tradetime
from strategy import macd
from strategy.util import action
import threading
# 定时器
timer = None
code_dict = dict()
INIT_FLAG = "init"
CHECK_FLAG = "check"


def tushare(code_list):
    global timer
    global code_dict
    # 读取相关股票的数据，并存入字典
    for code in code_list:
        # 检查股票日线状态，只有趋势上升才监听
        ktype = "D"
        last_datetime = tradetime.get_preday(60)
        trend_df = get_trend(code, ktype, last_datetime)
        tail_row = trend_df.iloc[-1]
        if tail_row[action.INDEX_STATUS] == action.STATUS_DOWN:
            console.write_msg("%s日线当前状态%s，趋势状态%s，忽略监听" % (code, tail_row[action.INDEX_STATUS]))
        else:
            for i in range(1, len(trend_df)):
                row = trend_df.iloc[-i]
                if row[action.INDEX_STATUS] == action.STATUS_SHAKE:
                    continue
                else:
                    break
            console.write_msg("%s日线当前状态%s，趋势状态%s" % (code, tail_row[action.INDEX_STATUS], row[action.INDEX_STATUS]))
        code_dict[code] = dict()
        code_dict[code][INIT_FLAG] = True
        code_dict[code][CHECK_FLAG] = False
        # TODO 监听筛选池的股票，发出5min和30min的出入信号
    five_minutes()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def five_minutes():
    global code_dict
    # 监听仓位中的股票，发出5min和30min的出入信号
    for code, flag_dict in code_dict.items():
        ktype = "5"
        last_datetime = tradetime.get_preday(2)
        trend_df = get_trend(code, ktype, last_datetime)

        # 当天第一次开启监控时，读取数据并初始化
        if flag_dict[INIT_FLAG] is True:
            for i in range(2, len(trend_df) + 1):
                tmp_df = trend_df.head(i)
                check(code, tmp_df)
            code_dict[code][INIT_FLAG] = False
        else:
            check(code, trend_df)

    global timer
    remain_second = tradetime.get_remain_second("5")
    timer = threading.Timer(remain_second + 30, five_minutes)
    timer.start()


def check(code, trend_df):
    # TODO 添加中枢(空间)的判断
    # TODO 修复phase_status的回溯数量
    # TODO 震荡空间只跟trend_count相关，一般7-8根有2个点左右，时间滞后性有1个点成本
    # TODO 成交量对股票波动幅度的影响
    # TODO 结合背离的分析
    global code_dict
    if code_dict[code][CHECK_FLAG] is False:
        # 判断第一次出现的波动
        now = trend_df.iloc[-1]
        pre = trend_df.iloc[-2]
        if pre[action.INDEX_TREND_COUNT] >= 5 and now[action.INDEX_STATUS] == action.STATUS_SHAKE:
            # 检查macd波动幅度
            if abs(now["macd"] - pre["macd"]) > action.FACTOR_MACD_RANGE:
                # 如果macd波动值超出范围，视为转折
                macd_diff = now["macd"] - pre["macd"]
                output(code, now["date"], pre["status"], pre["trend_count"], macd_diff)
            else:
                # 如果macd波动值在范围以内，则视为波动，观察后续走势
                code_dict[code][CHECK_FLAG] = pre["date"]
    else:
        # 判断波动过程
        border = trend_df[trend_df["date"] == code_dict[code][CHECK_FLAG]]
        now = trend_df.iloc[-1]
        pre = trend_df.iloc[-2]
        if now["status"] != action.STATUS_SHAKE:
            # 波动结束趋势逆转
            if now["status"] != border["status"].values[0]:
                raw_status = border["status"].values[0]
                trend_count = border["trend_count"]
                macd_diff = now["macd"] - border["macd"].values[0]
                output(code, now["date"], raw_status, trend_count, macd_diff)
            code_dict[code][CHECK_FLAG] = False
        else:
            macd_diff = now["macd"] - border["macd"].values[0]
            if abs(macd_diff) > 0.015:
                raw_status = border[action.INDEX_STATUS].values[0]
                trend_count = border[action.INDEX_TREND_COUNT]
                output(code, now["date"], raw_status, trend_count, macd_diff)
                code_dict[code][CHECK_FLAG] = False
    return


def output(code, date, raw_status, trend_count, macd_diff):
    # 考虑30min对5min的影响，30min上升对5min卖点，30min下降对5min买点，相反方向会产生压制
    # TODO，能否估算当前5min对应的30min情况，将未来因子考虑进去
    thirty_status, thirty_macd_diff, thirty_trend_count = get_relate(code, "30", date)

    # 根据股票类型，获取对应大盘状况
    code_prefix = code[0:1]
    if code_prefix == "0":
        index = "sz"
    elif code_prefix == "6":
        index = "sh"
    index_status, index_macd_diff, index_trend_count = get_relate(index, "5", date)

    positions = 3
    if raw_status == action.STATUS_UP:
        trade_type = "卖点"
        if thirty_status == action.STATUS_UP:
            positions -= 1
        if index_status == action.STATUS_UP:
            positions -= 1
    else:
        trade_type = "买点"
        if thirty_status == action.STATUS_DOWN:
            positions -= 1
        if index_status == action.STATUS_DOWN:
            positions -= 1
    console.write_msg("【%s, %s, %s】" % (code, date, trade_type))
    console.write_msg("个股5min，趋势%s，连续%d次，macd差值%f" % (raw_status, trend_count, macd_diff))
    console.write_msg("个股30min，趋势%s，连续%d次，macd差值%f" % (thirty_status, thirty_trend_count, thirty_macd_diff))
    console.write_msg("大盘5min，趋势%s，连续%d次，macd差值%f" % (index_status, index_trend_count, index_macd_diff))
    console.write_msg("建议仓位：%d/3" % (positions))
    return


def get_relate(code, ktype, end=None):
    """
    获取相关资讯
    """
    if ktype == "5":
        last_datetime = tradetime.get_preday(2)
    elif ktype == "30":
        last_datetime = tradetime.get_preday(16)
    elif ktype == "D":
        last_datetime = tradetime.get_preday(128)

    if end is not None:
        trend_df = get_trend(code, ktype, last_datetime, end)
    else:
        trend_df = get_trend(code, ktype, last_datetime)

    tail_row = trend_df.iloc[-1]
    if tail_row[action.INDEX_STATUS] == action.STATUS_SHAKE:
        for i in range(2, len(trend_df)):
            tmp_row = trend_df.iloc[-i]
            if tmp_row[action.INDEX_STATUS] != action.STATUS_SHAKE:
                status = tmp_row[action.INDEX_STATUS] + "-" + tail_row[action.INDEX_STATUS]
                macd_diff = tail_row["macd"] - tmp_row["macd"]
                trend_count = i - 1
                break
    else:
        tmp_row = trend_df.iloc[-2]
        macd_diff = tail_row["macd"] - tmp_row["macd"]
        status = tail_row[action.INDEX_STATUS]
        trend_count = tail_row[action.INDEX_TREND_COUNT]
    return status, macd_diff, trend_count


def get_trend(code, ktype, start_date, end_date=None):
    # TODO (重要)，支持ip池并发获取，要不然无法支持多code的高频获取
    while True:
        if end_date is None:
            df = ts.get_hist_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date)
        else:
            df = ts.get_hist_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date, end=end_date)
        time.sleep(conf.REQUEST_BLANK)
        if df is not None and df.empty is not True:
            df = df[['open', 'high', 'close', 'low', 'volume']]
            df = df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
            df = df.set_index(conf.HDF5_SHARE_DATE_INDEX)
            index_df = macd.value(df)
            df = df.merge(index_df, left_index=True, right_index=True, how='left')
            df = df.dropna()
            trend_df = macd.trend(df.reset_index())
            break
        else:
            console.write_msg("无法获取" + code + "-" + ktype + ":" + start_date + "以后的数据，休息30秒重新获取")
            time.sleep(30)
    trend_df = trend_df[
        [
            conf.HDF5_SHARE_DATE_INDEX,
            action.INDEX_TURN_COUNT,
            action.INDEX_TREND_COUNT,
            action.INDEX_STATUS,
            action.INDEX_PHASE_STATUS,
            "close", "macd", "dif", "dea"
        ]
    ].tail(48)
    return trend_df
