from library import influx
from library import console, conf, tool
import pandas as pd
import numpy as np
import h5py


def push_classify_detail():
    """
    将分类推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        # conf.HDF5_CLASSIFY_CONCEPT,
        # conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    # 获取classify列表
    for ctype in classify_list:
        # for classify_name in f[ctype]:
        classify_name = "chgn_700997"
        console.write_head(
            conf.HDF5_OPERATE_PUSH,
            conf.HDF5_RESOURCE_TUSHARE,
            classify_name
        )

        if f[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_DETAIL) is not None:
            df = tool.df_from_dataset(f[ctype][classify_name], conf.HDF5_CLASSIFY_DS_DETAIL, None)
            df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
            df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
            df = df.replace(np.inf, 0)
            influx.write_df(df, "demo")
        console.write_tail()
    f.close()

    # 获取classify列表
    # influx.init_db()
    return
