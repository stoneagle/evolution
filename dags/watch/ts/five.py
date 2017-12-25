import time
import tushare as ts
from library import conf, console, tradetime
from quota import macd
from quota.util import action
import threading

timer = None
code_dict = dict()
INIT_FLAG = "init"
CHECK_FLAG = "check"
DF_START_NUM = 26
DF_SPLIT_NUM = 48 + 48


def exec(code_list):
    global code_dict
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
        # TODO 获取个股股东人数、解禁、股东出货等情况
        code_dict[code] = dict()
        code_dict[code][INIT_FLAG] = True
        code_dict[code][CHECK_FLAG] = False
    five_minutes()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def five_minutes():
    """
    个股5min监听
    """
    global code_dict
    # 监听仓位中的股票，发出5min和30min的出入信号
    for code, flag_dict in code_dict.items():
        ktype = "5"
        # TODO，应该是提前两个交易日，而不是自然日
        last_datetime = tradetime.get_preday(5)
        trend_df = get_trend(code, ktype, last_datetime)

        # 当天第一次开启监控时，读取数据并初始化
        if flag_dict[INIT_FLAG] is True:
            for i in range(DF_START_NUM, len(trend_df) + 1):
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
    """
    检查macd趋势
    """
    # TODO (应该放在筛选中处理)添加中枢(空间)的判断，估算涨跌空间
    # TODO (应该放在筛选中处理)成交量对股票波动幅度的影响
    global code_dict

    now = trend_df.iloc[-1]
    pre = trend_df.iloc[-2]
    start = trend_df.iloc[-2 - pre[action.INDEX_TREND_COUNT]]
    phase_range = pre["macd"] - start["macd"]
    price_range = pre["close"] - start["close"]

    # 背离分析，如果macd是下降趋势，但是价格上涨，则属于背离；反之同理
    reverse_flag = False
    if (phase_range > 0 and price_range < 0) or (phase_range < 0 and price_range > 0):
        reverse_flag = True

    if code_dict[code][CHECK_FLAG] is False:
        # 判断第一次出现的波动
        # 如果该段趋势大于1(太少没有做T必要)，则判断波动情况
        if abs(phase_range) >= 0.01 and now[action.INDEX_STATUS] == action.STATUS_SHAKE:
            # 检查macd波动幅度
            macd_diff = (now["macd"] - pre["macd"])
            macd_range = action.FACTOR_MACD_RANGE * 1.5
            if abs(macd_diff) > macd_range:
                # 如果macd波动值超出范围，视为转折
                output(code, now["date"], pre["status"], pre["trend_count"], macd_diff, reverse_flag)
            else:
                # 如果macd波动值在范围以内，则视为波动，观察后续走势
                code_dict[code][CHECK_FLAG] = pre["date"]
    else:
        # 判断波动过程
        border = trend_df[trend_df["date"] == code_dict[code][CHECK_FLAG]]
        if now["status"] != action.STATUS_SHAKE:
            # 波动结束趋势逆转
            if now["status"] != border["status"].values[0]:
                raw_status = border["status"].values[0]
                trend_count = border["trend_count"].values[0]
                macd_diff = now["macd"] - border["macd"].values[0]
                output(code, now["date"], raw_status, trend_count, macd_diff, reverse_flag)
            code_dict[code][CHECK_FLAG] = False
        else:
            macd_diff = now["macd"] - border["macd"].values[0]
            if abs(macd_diff) > action.FACTOR_MACD_RANGE * 2:
                raw_status = border[action.INDEX_STATUS].values[0]
                trend_count = border[action.INDEX_TREND_COUNT].values[0]
                output(code, now["date"], raw_status, trend_count, macd_diff, reverse_flag)
                code_dict[code][CHECK_FLAG] = False
    return


def output(code, date, raw_status, trend_count, macd_diff, reverse_flag):
    """
    输出检查结果
    """
    # TODO，能否估算当前5min对应的30min情况，将未来因子考虑进去
    # 考虑30min对5min的影响，30min上升对5min卖点，30min下降对5min买点，相反方向会产生压制
    thirty_status, thirty_macd_diff, thirty_trend_count = get_relate(code, "30", date)

    # 根据股票类型，获取对应大盘状况
    code_prefix = code[0:1]
    if code_prefix == "0":
        index = "sz"
    elif code_prefix == "6":
        index = "sh"
    index_status, index_macd_diff, index_trend_count = get_relate(index, "30", date)

    # 仓位计算
    positions = 1
    if raw_status == action.STATUS_UP:
        if reverse_flag is True:
            trade_type = "背离买点"
        else:
            trade_type = "正常卖点"
        if thirty_status == action.STATUS_DOWN:
            positions += 1
        if index_status == action.STATUS_DOWN:
            positions += 1
    else:
        if reverse_flag is True:
            trade_type = "背离卖点"
        else:
            trade_type = "正常买点"
        if thirty_status == action.STATUS_UP:
            positions += 1
        if index_status == action.STATUS_UP:
            positions += 1
    console.write_msg("【%s, %s, %s】" % (code, date, trade_type))
    console.write_msg("个股5min，趋势%s，连续%d次，macd差值%f" % (raw_status, trend_count, macd_diff))
    console.write_msg("个股30min，趋势%s，连续%d次，macd差值%f" % (thirty_status, thirty_trend_count, thirty_macd_diff))
    console.write_msg("大盘30min，趋势%s，连续%d次，macd差值%f" % (index_status, index_trend_count, index_macd_diff))
    # 根据trend_count和当前时间，判断买卖时机(背离情况要计算两次买卖点trend_count的差值)
    if reverse_flag is False:
        remain_seconds = tradetime.get_trade_day_remain_second(date, "S")
        upper_estimate = (trend_count + 1) * 5
        lower_estimate = (trend_count - 1) * 5
        if (upper_estimate * 60) <= remain_seconds:
            trade_opportunity = "T+0"
        else:
            trade_opportunity = "T+1"
        console.write_msg("剩余时间%d分钟，下个交易点预估需要%d-%d分钟，模式%s" % (round(remain_seconds / 60, 0), lower_estimate, upper_estimate, trade_opportunity))
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
    """
    获取趋势
    """
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
    ].tail(DF_SPLIT_NUM)
    return trend_df
