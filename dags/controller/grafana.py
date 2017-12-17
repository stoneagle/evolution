from library import console, conf, tool, influx, tradetime
import pandas as pd
import numpy as np
import h5py

DF_INIT_LIMIT = 500


def basic_detail():
    """
    聚合xsg、ipo、shm、szm等数据，推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    # 获取xsg
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_XSG
    )
    if f.get(conf.HDF5_FUNDAMENTAL_XSG_DETAIL) is not None:
        xsg_df = tool.df_from_dataset(f, conf.HDF5_FUNDAMENTAL_XSG_DETAIL, None)
        xsg_df = _datetime_index(xsg_df)
        influx.reset_df(xsg_df, conf.MEASUREMENT_BASIC, {"btype": "xsg"})
    console.write_tail()

    # 获取ipo
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_IPO
    )
    if f.get(conf.HDF5_FUNDAMENTAL_IPO) and f[conf.HDF5_FUNDAMENTAL_IPO].get(conf.HDF5_FUNDAMENTAL_IPO_DETAIL) is not None:
        ipo_df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_IPO], conf.HDF5_FUNDAMENTAL_IPO_DETAIL, None)
        ipo_df = _datetime_index(ipo_df)
        influx.reset_df(ipo_df, conf.MEASUREMENT_BASIC, {"btype": "ipo"})
    console.write_tail()

    # 获取sh融资融券
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SH_MARGINS
    )
    console.write_tail()
    if f.get(conf.HDF5_FUNDAMENTAL_SH_MARGINS) and f[conf.HDF5_FUNDAMENTAL_SH_MARGINS].get(conf.HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL) is not None:
        shm_df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_SH_MARGINS], conf.HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL, None)
        shm_df = _datetime_index(shm_df)
        influx.reset_df(shm_df, conf.MEASUREMENT_BASIC, {"btype": "shm"})

    # 获取sz融资融券
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SZ_MARGINS
    )
    if f.get(conf.HDF5_FUNDAMENTAL_SZ_MARGINS) and f[conf.HDF5_FUNDAMENTAL_SZ_MARGINS].get(conf.HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL) is not None:
        shz_df = tool.df_from_dataset(f[conf.HDF5_FUNDAMENTAL_SZ_MARGINS], conf.HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL, None)
        shz_df = _datetime_index(shz_df)
        influx.reset_df(shz_df, conf.MEASUREMENT_BASIC, {"btype": "szm"})
    console.write_tail()
    f.close()
    return


def share_detail(code_list):
    """
    将股票数据推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code in f:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            continue

        console.write_head(
            conf.HDF5_OPERATE_PUSH,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )

        # 推送原始kline
        for ktype in conf.HDF5_SHARE_KTYPE:
            if f[code_prefix][code].get(ktype) is None:
                continue
            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if f[code_prefix][code].get(index_ds_name) is None:
                continue
            detail_df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
            index_df = tool.df_from_dataset(f[code_prefix][code], index_ds_name, None)
            detail_df = detail_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')
            detail_df = _datetime_index(detail_df)
            influx.reset_df(detail_df, conf.MEASUREMENT_SHARE, {"code": code, "ktype": ktype})

        # 推送缠论kline
        for ktype in ["D", "30"]:
            wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
            if f[code_prefix][code].get(wrap_ds_name) is None:
                continue
            wrap_df = tool.df_from_dataset(f[code], wrap_ds_name, None)
            wrap_df = _datetime_index(wrap_df)
            influx.reset_df(wrap_df, conf.MEASUREMENT_SHARE_WRAP, {"code": code, "ktype": ktype})
        console.write_tail()
    f.close()
    return


def daily_filter():
    """
    推送每日筛选列表至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_SCREEN, 'a')
    today_str = tradetime.get_today()
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        today_str
    )
    if f.get(today_str) is None:
        return
    screen_df = tool.df_from_dataset(f, today_str, None)
    screen_df = _datetime_index(screen_df)
    influx.reset_df(screen_df, conf.MEASUREMENT_FILTER_SHARE, None)
    f.close()
    return


def code_classify():
    """
    推送筛选出的股票相关分类
    """
    f = h5py.File(conf.HDF5_FILE_OTHER, 'a')
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_OTHER_CODE_CLASSIFY
    )
    if f.get(conf.HDF5_FILE_OTHER) is None:
        return
    code_classify_df = tool.df_from_dataset(f, conf.HDF5_FILE_OTHER, None)
    code_classify_df = _datetime_index(code_classify_df)
    influx.reset_df(code_classify_df, conf.MEASUREMENT_CODE_CLASSIFY, None)
    f.close()
    return


def index_detail(reset_flag=False):
    """
    将指数数据推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
    for code in f:
        console.write_head(
            conf.HDF5_OPERATE_PUSH,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        # 推送原始kline
        for ktype in conf.HDF5_SHARE_KTYPE:
            ctags = {"itype": code, "ktype": ktype}
            measurement = conf.MEASUREMENT_INDEX

            if f[code].get(ktype) is None:
                continue
            index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
            if f[code].get(index_ds_name) is None:
                continue

            detail_df = tool.df_from_dataset(f[code], ktype, None)
            index_df = tool.df_from_dataset(f[code], index_ds_name, None)
            detail_df = detail_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')
            detail_df = _datetime_index(detail_df)

            datetime = influx.get_last_datetime(measurement, ctags)
            if datetime is not None and reset_flag is False:
                detail_df = detail_df.loc[detail_df.index > datetime]
            else:
                detail_df = detail_df.tail(DF_INIT_LIMIT)
            detail_df = detail_df.drop("ma_border", axis=1)
            if len(detail_df) > 0:
                try:
                    # influx.reset_df(detail_df, measurement, ctags)
                    console.write_exec()
                except Exception as er:
                    print(str(er))
            else:
                console.write_pass()

        # 推送缠论kline
        for ktype in ["D", "30"]:
            ctags = {"itype": code, "ktype": ktype}
            measurement = conf.MEASUREMENT_INDEX_WRAP

            wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
            if f[code].get(wrap_ds_name) is None:
                console.write_msg(wrap_ds_name + "的wrap数据不存在")
                continue
            wrap_df = tool.df_from_dataset(f[code], wrap_ds_name, None)
            wrap_df = _datetime_index(wrap_df)
            datetime = influx.get_last_datetime(measurement, ctags)
            if datetime is not None and reset_flag is False:
                wrap_df = wrap_df.loc[wrap_df.index > datetime]
            else:
                wrap_df = wrap_df.tail(DF_INIT_LIMIT)
            if len(wrap_df) > 0:
                try:
                    influx.reset_df(wrap_df, measurement, ctags)
                    console.write_exec()
                except Exception as er:
                    print(str(er))
            else:
                console.write_pass()
        console.write_blank()
        console.write_tail()
    f.close()
    return


def classify_detail(classify_list, reset_flag=False):
    """
    将分类推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    # 获取classify列表
    for ctype in classify_list:
        console.write_head(
            conf.HDF5_OPERATE_PUSH,
            conf.HDF5_RESOURCE_TUSHARE,
            ctype
        )
        for classify_name in f[ctype]:

            # 推送原始kline
            for ktype in conf.HDF5_SHARE_KTYPE:
                ctags = {"ktype": ktype, "classify": classify_name}
                measurement = conf.MEASUREMENT_CLASSIFY + "_" + ctype
                detail_ds_name = conf.HDF5_CLASSIFY_DS_DETAIL + "_" + ktype
                if f[ctype][classify_name].get(detail_ds_name) is None:
                    continue

                index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
                if f[ctype][classify_name].get(index_ds_name) is None:
                    continue

                detail_df = tool.df_from_dataset(f[ctype][classify_name], detail_ds_name, None)
                index_df = tool.df_from_dataset(f[ctype][classify_name], index_ds_name, None)
                detail_df = detail_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='outer')
                detail_df = _datetime_index(detail_df)
                datetime = influx.get_last_datetime(measurement, ctags)
                if datetime is not None and reset_flag is False:
                    detail_df = detail_df.loc[detail_df.index > datetime]
                else:
                    detail_df = detail_df.tail(DF_INIT_LIMIT)
                if len(detail_df) > 0:
                    try:
                        influx.reset_df(detail_df, measurement, ctags)
                        console.write_exec()
                    except Exception as er:
                        print(str(er))
                else:
                    console.write_pass()

            # 推送缠论kline
            for ktype in ["D", "30"]:
                ctags = {"ktype": ktype, "classify": classify_name}
                measurement = conf.MEASUREMENT_CLASSIFY_WRAP + "_" + ctype
                wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
                if f[ctype][classify_name].get(wrap_ds_name) is None:
                    continue
                wrap_df = tool.df_from_dataset(f[ctype][classify_name], wrap_ds_name, None)
                wrap_df = _datetime_index(wrap_df)
                datetime = influx.get_last_datetime(measurement, ctags)
                if datetime is not None and reset_flag is False:
                    wrap_df = wrap_df.loc[wrap_df.index > datetime]
                else:
                    detail_df = detail_df.tail(DF_INIT_LIMIT)
                if len(wrap_df) > 0:
                    try:
                        influx.reset_df(wrap_df, measurement, ctags)
                        console.write_exec()
                    except Exception as er:
                        print(str(er))
                else:
                    console.write_pass()
        console.write_blank()
        console.write_tail()
    f.close()
    return


def _datetime_index(df):
    df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    df.index = pd.to_datetime(df[conf.HDF5_SHARE_DATE_INDEX])
    df = df.drop(conf.HDF5_SHARE_DATE_INDEX, axis=1)
    df = df.replace(np.inf, 0)
    df = df.replace(-np.inf, 0)
    df = df.replace(np.nan, 0)
    return df
