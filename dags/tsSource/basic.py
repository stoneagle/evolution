import tushare as ts
import time
from tsSource import cons
from library import tool, count
from datetime import datetime, timedelta


def get_detail(f):
    # 按周间隔获取
    start_date = datetime.strptime("2016-08-12", "%Y-%m-%d")
    while start_date <= datetime.now():
        try:
            time.sleep(cons.REQUEST_BLANK)
            start_date_str = datetime.strftime(start_date, "%Y-%m-%d")
            start_date = start_date + timedelta(days=7)
            if f.get(start_date_str) is not None:
                count.inc_by_index("pass")
                continue
            df = ts.get_stock_basics(start_date_str)
            if df is not None and df.empty is not True:
                df = df.drop("name", axis=1)
                df = df.drop("area", axis=1)
                df = df.drop("industry", axis=1)
                tool.create_df_dataset(f, start_date_str, df.reset_index())
                count.inc_by_index("get")
        except Exception as er:
            # TODO 如果无法获取数据，则后退一天，两次后未获取则忽略
            print(str(er))
            continue
    return


def arrange_detail():
    # 按个股整理至share文件下
    return


def get_achievement():
    return
