from library import console, conf, tool, influx
import pandas as pd
import numpy as np
import h5py


def basic_detail():
    """
    将基本数据推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    # 推送xsg
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_XSG
    )
    if f.get(conf.HDF5_FUNDAMENTAL_XSG_DETAIL) is not None:
        df = tool.df_from_dataset(f, conf.HDF5_FUNDAMENTAL_XSG_DETAIL, None)
        df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
        df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
        df = df.replace(np.inf, 0)
        influx.reset_df(df, conf.MEASUREMENT_BASIC, {"btype": "xsg"})
    console.write_tail()

    # 推送ipo
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_IPO
    )
    if f.get(conf.HDF5_FUNDAMENTAL_IPO) and f[conf.HDF5_FUNDAMENTAL_IPO].get(conf.HDF5_FUNDAMENTAL_IPO_DETAIL) is not None:
        df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_IPO], conf.HDF5_FUNDAMENTAL_IPO_DETAIL, None)
        df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
        df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
        df = df.replace(np.inf, 0)
        influx.reset_df(df, conf.MEASUREMENT_BASIC, {"btype": "ipo"})
    console.write_tail()

    # 推送shm融资融券
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SH_MARGINS
    )
    console.write_tail()
    if f.get(conf.HDF5_FUNDAMENTAL_SH_MARGINS) and f[conf.HDF5_FUNDAMENTAL_SH_MARGINS].get(conf.HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL) is not None:
        df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_SH_MARGINS], conf.HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL, None)
        df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
        df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
        df = df.replace(np.inf, 0)
        influx.reset_df(df, conf.MEASUREMENT_BASIC, {"btype": "shm"})

    # 推送shz融资融券
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SZ_MARGINS
    )
    if f.get(conf.HDF5_FUNDAMENTAL_SZ_MARGINS) and f[conf.HDF5_FUNDAMENTAL_SZ_MARGINS].get(conf.HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL) is not None:
        df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_SZ_MARGINS], conf.HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL, None)
        df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
        df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
        df = df.replace(np.inf, 0)
        influx.reset_df(df, conf.MEASUREMENT_BASIC, {"btype": "szm"})
    console.write_tail()
    f.close()
    return


def share_detail(code_list):
    """
    将股票详情推送至influxdb
    """
    return


def classify_detail():
    """
    将分类推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        # conf.HDF5_CLASSIFY_INDUSTRY,
        # conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_HOT,
    ]
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_PUSH,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )

            for ktype in conf.HDF5_SHARE_KTYPE:
                detail_ds_name = conf.HDF5_CLASSIFY_DS_DETAIL + "_" + ktype
                if f[ctype][classify_name].get(detail_ds_name) is None:
                    continue

                index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if f[ctype][classify_name].get(index_ds_name) is None:
                    continue

                detail_df = tool.df_from_dataset(f[ctype][classify_name], detail_ds_name, None)
                index_df = tool.df_from_dataset(f[ctype][classify_name], index_ds_name, None)
                detail_df = detail_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')

                detail_df[conf.HDF5_SHARE_DATE_INDEX] = detail_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
                detail_df.index = pd.to_datetime(detail_df[conf.HDF5_SHARE_DATE_INDEX])
                detail_df = detail_df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
                detail_df = detail_df.replace(np.inf, 0)
                detail_df = detail_df.replace(np.nan, 0)
                influx.write_df(detail_df, conf.MEASUREMENT_CLASSIFY + "-" + ctype, {"ktype": ktype, "classify": classify_name})
            console.write_tail()
    f.close()
    return
