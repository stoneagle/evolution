from library import console, conf, tool, influx, tradetime
import pandas as pd
import numpy as np
import datetime
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


def share_detail(code_list, reset_flag=False):
    """
    将股票数据推送至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            console.write_msg(code + "目录不存在")
            continue

        console.write_head(
            conf.HDF5_OPERATE_PUSH,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        # 推送原始kline
        measurement = conf.MEASUREMENT_SHARE
        _raw_kline(f[code_prefix][code], measurement, code, reset_flag)

        # 推送缠论kline
        measurement = conf.MEASUREMENT_SHARE_WRAP
        _wrap_kline(f[code_prefix][code], measurement, code, reset_flag)

        # 推送基本面数据
        measurement = conf.MEASUREMENT_SHARE_BASIC
        _basic_info(f[code_prefix][code], measurement, code, reset_flag)

        console.write_blank()
        console.write_tail()
    f.close()
    return


def share_filter(today_str=None):
    """
    推送筛选列表至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_SCREEN, 'a')
    if today_str is None:
        today_str = tradetime.get_today()
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.MEASUREMENT_FILTER_SHARE
    )

    if f[conf.SCREEN_SHARE_FILTER].get(today_str) is None:
        console.write_msg(today_str + "的筛选数据不存在")
        return
    screen_df = tool.df_from_dataset(f[conf.SCREEN_SHARE_FILTER], today_str, None)
    screen_df[conf.HDF5_SHARE_DATE_INDEX] = bytes(today_str, encoding="utf8")
    screen_df = _datetime_index(screen_df)
    screen_df = screen_df.reset_index()
    num = 1
    for index, row in screen_df.iterrows():
        screen_df.loc[index, conf.HDF5_SHARE_DATE_INDEX] = screen_df.loc[index][conf.HDF5_SHARE_DATE_INDEX] + datetime.timedelta(0, num)
        num += 1
    screen_df = screen_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    try:
        influx.write_df(screen_df, conf.MEASUREMENT_FILTER_SHARE, {"filter_date": today_str})
    except Exception as er:
        print(str(er))
    console.write_tail()
    f.close()
    return


def share_grade(today_str=None):
    """
    推送打分结果至influxdb
    """
    f = h5py.File(conf.HDF5_FILE_SCREEN, 'a')
    if today_str is None:
        today_str = tradetime.get_today()
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.MEASUREMENT_FILTER_SHARE_GRADE
    )

    if f[conf.SCREEN_SHARE_GRADE].get(today_str) is None:
        console.write_msg(today_str + "的打分数据不存在")
        return
    grade_df = tool.df_from_dataset(f[conf.SCREEN_SHARE_GRADE], today_str, None)
    grade_df[conf.HDF5_SHARE_DATE_INDEX] = bytes(today_str, encoding="utf8")
    grade_df = _datetime_index(grade_df)
    grade_df = grade_df.reset_index()
    num = 1
    for index, row in grade_df.iterrows():
        grade_df.loc[index, conf.HDF5_SHARE_DATE_INDEX] = grade_df.loc[index][conf.HDF5_SHARE_DATE_INDEX] + datetime.timedelta(0, num)
        num += 1
    grade_df = grade_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    try:
        influx.write_df(grade_df, conf.MEASUREMENT_FILTER_SHARE_GRADE, {"filter_date": today_str})
    except Exception as er:
        print(str(er))
    console.write_tail()
    f.close()
    return


def code_classify(today_str=None):
    """
    推送筛选出的股票相关分类
    """
    f = h5py.File(conf.HDF5_FILE_OTHER, 'a')
    console.write_head(
        conf.HDF5_OPERATE_PUSH,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_OTHER_CODE_CLASSIFY
    )
    if today_str is None:
        today_str = tradetime.get_today()

    if f.get(conf.HDF5_OTHER_CODE_CLASSIFY) is None:
        console.write_msg("code的分类文件不存在")
        return
    code_classify_df = tool.df_from_dataset(f, conf.HDF5_OTHER_CODE_CLASSIFY, None)
    code_classify_df[conf.HDF5_SHARE_DATE_INDEX] = bytes(today_str, encoding="utf8")
    code_classify_df = _datetime_index(code_classify_df)
    code_classify_df = code_classify_df.reset_index()
    num = 1
    for index, row in code_classify_df.iterrows():
        code_classify_df.loc[index, conf.HDF5_SHARE_DATE_INDEX] = code_classify_df.loc[index][conf.HDF5_SHARE_DATE_INDEX] + datetime.timedelta(0, num)
        num += 1
    code_classify_df = code_classify_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
    try:
        influx.write_df(code_classify_df, conf.MEASUREMENT_CODE_CLASSIFY, None)
    except Exception as er:
        print(str(er))
    console.write_tail()
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
        measurement = conf.MEASUREMENT_INDEX
        _raw_kline(f[code], measurement, code, reset_flag)

        # 推送缠论kline
        measurement = conf.MEASUREMENT_INDEX_WRAP
        _wrap_kline(f[code], measurement, code, reset_flag)

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
            measurement = conf.MEASUREMENT_CLASSIFY + "_" + ctype
            _raw_kline(f[ctype][classify_name], measurement, classify_name, reset_flag)

            # 推送缠论kline
            measurement = conf.MEASUREMENT_CLASSIFY_WRAP + "_" + ctype
            _wrap_kline(f[ctype][classify_name], measurement, classify_name, reset_flag)
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


def _raw_kline(f, measurement, code, reset_flag=False):
    """
    推送原始kline
    """
    for ktype in conf.HDF5_SHARE_KTYPE:
        ctags = {"kcode": code, "ktype": ktype}
        detail_ds_name = ktype
        index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype

        if f.get(detail_ds_name) is None:
            console.write_msg(code + "的detail数据不存在")
            continue
        if f.get(index_ds_name) is None:
            console.write_msg(code + "的index数据不存在")
            continue
        detail_df = tool.df_from_dataset(f, detail_ds_name, None)
        index_df = tool.df_from_dataset(f, index_ds_name, None)
        detail_df = detail_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='left')
        detail_df = _datetime_index(detail_df)
        last_datetime = influx.get_last_datetime(measurement, ctags)
        if last_datetime is not None and reset_flag is False:
            detail_df = detail_df.loc[detail_df.index > last_datetime]
        else:
            detail_df = detail_df.tail(DF_INIT_LIMIT)
        detail_df = detail_df.drop("ma_border", axis=1)
        if len(detail_df) > 0:
            try:
                influx.reset_df(detail_df, measurement, ctags)
                console.write_exec()
            except Exception as er:
                print(str(er))
        else:
            console.write_pass()
    return


def _wrap_kline(f, measurement, code, reset_flag=False):
    """
    推送缠论kline
    """
    for ktype in conf.HDF5_SHARE_WRAP_KTYPE:
        ctags = {"kcode": code, "ktype": ktype}
        wrap_ds_name = conf.HDF5_INDEX_WRAP + "_" + ktype
        if f.get(wrap_ds_name) is None:
            console.write_msg(code + "缠论数据不存在")
            continue

        wrap_df = tool.df_from_dataset(f, wrap_ds_name, None)
        wrap_df = _datetime_index(wrap_df)
        last_datetime = influx.get_last_datetime(measurement, ctags)
        if last_datetime is not None and reset_flag is False:
            wrap_df = wrap_df.loc[wrap_df.index > last_datetime]
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
    return


def _basic_info(f, measurement, code, reset_flag):
    ctags = {"kcode": code}
    basic_ds_name = conf.HDF5_BASIC_DETAIL
    if f.get(basic_ds_name) is None:
        console.write_msg(code + "基本面数据不存在")
        return

    basic_df = tool.df_from_dataset(f, basic_ds_name, None)
    basic_df = _datetime_index(basic_df)
    last_datetime = influx.get_last_datetime(measurement, ctags)
    if last_datetime is not None and reset_flag is False:
        basic_df = basic_df.loc[basic_df.index > last_datetime]
    else:
        basic_df = basic_df.tail(DF_INIT_LIMIT)
    if len(basic_df) > 0:
        try:
            influx.reset_df(basic_df, measurement, ctags)
            console.write_exec()
        except Exception as er:
            print(str(er))
    else:
        console.write_pass()
    return
