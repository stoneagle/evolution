import tushare as ts
from library import tradetime as ttime
from library import tool, count, error
from tsSource import cons
import time


def _add_data(code, ktype, f, end_date):
    df = ts.get_hist_data(code, ktype=ktype, pause=cons.REQUEST_BLANK, end=end_date)
    time.sleep(cons.REQUEST_BLANK)
    ret = False
    if df is not None and df.empty is not True:
        df = df[cons.SHARE_COLS]
        df = df.reset_index().sort_values(by=[cons.SHARE_DATE_INDEX])
        tool.create_df_dataset(f, ktype, df)
        ret = True
        count.inc_by_index(ktype)
    else:
        count.inc_by_index("empty")
    return ret


def _append_data(code, ktype, f, start_date, end_date):
    df = ts.get_hist_data(code, ktype=ktype, pause=cons.REQUEST_BLANK, end=end_date, start=start_date)
    time.sleep(cons.REQUEST_BLANK)
    ret = False
    if df is not None and df.empty is not True:
        df = df[cons.SHARE_COLS]
        df = df.reset_index().sort_values(by=[cons.SHARE_DATE_INDEX])
        tool.append_df_dataset(f, ktype, df)
        count.inc_by_index(ktype)
        ret = True
    else:
        count.inc_by_index("empty")
    return ret


def get_share_data(code, f):
    # 获取不同周期的数据
    for ktype in ["M", "W", "D", "30", "5"]:
        try:
            if f.get(ktype) is None:
                # 如果股票不存在，则获取17年至今数据(M取上个月月底，W取上周日)
                end_date = ttime.get_end_date(code, ktype)
                end_date = "2017-10-17"
                ret = _add_data(code, ktype, f, end_date)
            else:
                # 如果已存在，则根据存储内最后日期，获取至今
                if len(f[ktype]) == 0:
                    # 如果数据为空
                    end_date = ttime.get_end_date(code, ktype)
                    ret = _add_data(code, ktype, f, end_date)
                else:
                    # TODO 将index的0改为常量控制
                    tail_date_str = f[ktype][-1][0].astype(str)
                    start_date = ttime.get_start_date(tail_date_str, code, ktype)
                    end_date = ttime.get_end_date(code, ktype)
                    # 如果开始日期大于等于结束日期，则不需要进行处理
                    if start_date >= end_date:
                        count.inc_by_index("pass")
                        continue
                    else:
                        ret = _append_data(code, ktype, f, start_date, end_date)
            if ret is not True:
                error.add_row([ktype, code])
        except Exception as er:
            time.sleep(cons.REQUEST_BLANK)
            print(str(er))
    return
