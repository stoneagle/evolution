import tushare as ts
import time
from tsSource import cons
from library import tool, count, conf, error, console
from datetime import datetime, timedelta


# 闭市
GET_DETAIL_CLOSE = "close"
# 获取异常
GET_DETAIL_OTHER = "other"


def get_detail(f):
    # 按日间隔获取
    start_date = datetime.strptime("2016-08-09", "%Y-%m-%d")
    # 获取历史错误数据
    history = error.get_file()
    close_history = list()
    if history is not None:
        history["type"] = history["type"].str.decode("utf-8")
        history["date"] = history["date"].str.decode("utf-8")
        close_history = history[history["type"] == "close"]["date"].values

    while start_date <= datetime.now():
        try:
            start_date_str = datetime.strftime(start_date, "%Y-%m-%d")
            # 如果是周六日，已获取，或者闭盘的日子则跳过
            if start_date.weekday() < 5 and start_date_str not in close_history and f.get(start_date_str) is None:
                df = ts.get_stock_basics(start_date_str)
                time.sleep(cons.REQUEST_BLANK)
                if df is not None and df.empty is not True:
                    df = df.drop("name", axis=1)
                    df = df.drop("area", axis=1)
                    df = df.drop("industry", axis=1)
                    tool.create_df_dataset(f, start_date_str, df.reset_index())
                    count.inc_by_index(conf.HDF5_COUNT_GET)
                    console.write_exec()
            else:
                count.inc_by_index(conf.HDF5_COUNT_PASS)
        except Exception as er:
            time.sleep(cons.REQUEST_BLANK)
            if str(er) != "HTTP Error 404: Not Found":
                error.add_row([GET_DETAIL_OTHER, start_date_str])
                print(str(er))
            else:
                error.add_row([GET_DETAIL_CLOSE, start_date_str])
        start_date = start_date + timedelta(days=1)
    return


def arrange_detail():
    # 按个股整理至share文件下
    return


def get_achievement():
    # 获取个股的财报
    return


def get_quit(f):
    # 获取终止上市和暂定上市的股票列表
    if f.get(conf.HDF5_BASIC_QUIT_TERMINATE) is not None:
        del f[conf.HDF5_BASIC_QUIT_TERMINATE]
    df = ts.get_terminated()
    df = df.drop("name", axis=1)
    tool.create_df_dataset(f, conf.HDF5_BASIC_QUIT_TERMINATE, df)

    if f.get(conf.HDF5_BASIC_QUIT_SUSPEND) is not None:
        del f[conf.HDF5_BASIC_QUIT_SUSPEND]
    df = ts.get_suspended()
    df = df.drop("name", axis=1)
    tool.create_df_dataset(f, conf.HDF5_BASIC_QUIT_SUSPEND, df)
    return
