import tushare as ts
import time
from datetime import datetime
from library import conf, tool, count, console


def get_profit():
    """
    获取高送转与分红
    """
    return


def get_xsg(f):
    """
    获取限售股解禁
    """
    for year in range(2010, datetime.today().year + 1):
        for month in range(1, 13):
            if month in range(1, 10):
                dset_name = str(year) + "0" + str(month)
            else:
                dset_name = str(year) + str(month)

            if f.get(dset_name) is not None:
                count.inc_by_index(conf.HDF5_COUNT_PASS)
                continue

            try:
                df = ts.xsg_data(year=year, month=month, pause=conf.REQUEST_BLANK)
                df = df.drop("name", axis=1)
                df = df.sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
                tool.create_df_dataset(f, dset_name, df)
                console.write_exec()
                count.inc_by_index(conf.HDF5_COUNT_GET)
            except Exception as er:
                print(str(er))
    time.sleep(conf.REQUEST_BLANK)
    return


def get_ipo(f, reset_flag=False):
    """
    获取ipo数据
    """
    df = ts.new_stocks(pause=conf.REQUEST_BLANK)
    df = df.drop("name", axis=1)
    df = df.sort_values(by=["ipo_date"])
    if reset_flag is False:
        tool.merge_df_dataset(f, conf.HDF5_FUNDAMENTAL_IPO, df)
    else:
        tool.create_df_dataset(f, conf.HDF5_FUNDAMENTAL_IPO, df)
    return


def get_sh_margins(f, reset_flag=False):
    """
    获取沪市的融资融券
    """
    df = ts.sh_margins(pause=conf.REQUEST_BLANK)
    df = df.sort_values(by=["opDate"])
    if reset_flag is False:
        tool.merge_df_dataset(f, conf.HDF5_FUNDAMENTAL_SH_MARGINS, df)
    else:
        tool.create_df_dataset(f, conf.HDF5_FUNDAMENTAL_SH_MARGINS, df)
    return


def get_sz_margins(f, reset_flag=False):
    """
    获取深市的融资融券
    """
    df = ts.sz_margins(pause=conf.REQUEST_BLANK)
    df = df.sort_values(by=["opDate"])
    if reset_flag is False:
        tool.merge_df_dataset(f, conf.HDF5_FUNDAMENTAL_SZ_MARGINS, df)
    else:
        tool.create_df_dataset(f, conf.HDF5_FUNDAMENTAL_SZ_MARGINS, df)
    return
