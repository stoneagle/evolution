import h5py
import time
import tushare as ts
from library import conf, console, tool
from strategy import macd
from strategy.util import action
import threading
# 定时器
timer = None
share_dict = dict()
init_flag = True


def tushare(code_list):
    global timer
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # 读取相关股票的数据，并存入字典
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            console.write_msg(code + "的detail文件目录不存在")
            return
        # for ktype in ["30", "5"]:
        for ktype in ["5"]:
            if f[code_group_path].get(ktype) is None:
                console.write_msg(code + "的" + ktype + "文件不存在")
                return
            if code not in share_dict:
                share_dict[code] = dict()
            share_df = tool.df_from_dataset(f[code_group_path], ktype, None)
            share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            share_df = share_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
            # share_dict[code][ktype] = share_df.head(len(share_df) - 48)
            # share_dict[code][ktype] = share_dict[code][ktype].tail(60)
            share_dict[code][ktype] = share_df.tail(60)
    f.close()
    daily()
    # bot = weixin.WXBot()
    # bot.sendHelper("lalala")
    return


def daily():
    global share_dict
    global init_flag
    # 监听仓位中的股票，发出5min和30min的出入信号
    for code, data in share_dict.items():
        for ktype, share_df in data.items():
            last_datetime = share_df.tail(1).index.values[0]
            while True:
                today_df = _get_data(code, ktype, last_datetime)
                if len(today_df) == 1:
                    console.write_msg("无法获取" + code + "-" + ktype + ":" + last_datetime + "以后的数据，休息30秒重新获取")
                    time.sleep(30)
                else:
                    break

            today_df = today_df.tail(len(today_df) - 1)
            today_df = today_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
            share_df = share_df.append(today_df)
            index_df = macd.value(share_df)
            share_df = share_df.merge(index_df, left_index=True, right_index=True, how='left')
            share_df = share_df.dropna()
            trend_df = macd.trend(share_df.reset_index())
            trend_df = trend_df[["date", "turn_count", "trend_count", "status", "phase_status", "macd", "dif", "dea"]].tail(48)

            # 当天第一次开启监控时，读取数据并初始化
            if init_flag is True:
                for i in range(2, len(trend_df)):
                    tmp_df = trend_df.head(i)
                    now = tmp_df.iloc[-1]
                    pre = tmp_df.iloc[-2]
                    _check(code, pre, now)
                init_flag = False
            else:
                now = trend_df.iloc[-1]
                pre = trend_df.iloc[-2]
                _check(code, pre, now)
    # TODO 监听筛选池的股票，发出5min和30min的出入信号

    global timer
    timer = threading.Timer(60 * 5, daily)
    timer.start()


def _check(code, pre, now):
    # TODO 添加中枢的判断
    # 判断第一次出现的shake
    if pre["trend_count"] >= 4 and now["status"] == action.STATUS_SHAKE and abs(pre["macd"]) >= 0.01:
        if pre["status"] == action.STATUS_UP:
            # dif与dea线，都在0轴以上0.01才算
            msg = "%s在%s出现5min卖点,连续趋势%d次,macd差值%f"
        elif pre["status"] == action.STATUS_DOWN:
            # dif与dea线，都在0轴以下0.01才算
            msg = "%s在%s出现5min买点,连续趋势%d次,macd差值%f"
        console.write_msg(msg % (code, now["date"], pre["trend_count"], now["macd"] - pre["macd"]))
    return


def _get_data(code, ktype, start_date):
    df = ts.get_hist_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date)
    time.sleep(conf.REQUEST_BLANK)
    if df is not None and df.empty is not True:
        df = df[['open', 'high', 'close', 'low', 'volume', 'turnover']]
        df = df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
        return df
    else:
        return None
