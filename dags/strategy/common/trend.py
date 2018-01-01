from quota import macd
from strategy.common import phase
from source.bitmex import future
from library import conf, tool, tradetime
from source.ts import share as tss
import time
import tushare as ts
import h5py
DF_MACD_MIN_NUM = 34


def get_from_file(ktype, stype, code, factor_macd_range, df_file_num=48, direct_turn=False):
    """
    从文件获取历史数据，并计算趋势
    """
    df = None
    df_file_num = df_file_num + DF_MACD_MIN_NUM
    if stype == conf.STYPE_BITMEX:
        f = h5py.File(conf.HDF5_FILE_FUTURE, 'a')
        path = '/' + conf.HDF5_RESOURCE_BITMEX + '/' + code
        if f.get(path) is not None:
            df = tool.df_from_dataset(f[path], ktype, None)
            df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            df = df.tail(df_file_num)
        else:
            raise Exception(code + "的" + ktype + "文件数据不存在")
        f.close()
    elif stype == conf.STYPE_ASHARE:
        if code.isdigit():
            f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
            code_prefix = code[0:3]
            path = '/' + code_prefix + '/' + code
        else:
            f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
            path = '/' + code

        if f.get(path) is not None:
            df = tool.df_from_dataset(f[path], ktype, None)
            df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            df = df.tail(df_file_num)
        else:
            raise Exception(code + "的" + ktype + "文件数据不存在")
        f.close
    else:
        raise Exception("数据源不存在或未配置")
    return macd.value_and_trend(df, factor_macd_range, direct_turn)


def get_from_remote(ktype, stype, start_date, code, rewrite):
    """
    从远端获取最新数据
    """
    if stype == conf.STYPE_BITMEX:
        count = tradetime.get_barnum_by_date(start_date, ktype)
        if ktype in [
            conf.BINSIZE_ONE_DAY,
            conf.BINSIZE_ONE_HOUR,
            conf.BINSIZE_FIVE_MINUTE,
            conf.BINSIZE_ONE_MINUTE,
        ]:
            df = future.history(code, ktype, count)
        elif ktype in [
            conf.BINSIZE_THIRTY_MINUTE,
            conf.BINSIZE_FOUR_HOUR
        ]:
            df = future.history_merge(code, ktype, count)

        if rewrite:
            # 将读取的remote数据写回文件
            f = h5py.File(conf.HDF5_FILE_FUTURE, 'a')
            path = '/' + conf.HDF5_RESOURCE_BITMEX + '/' + code
            if f.get(path) is None:
                f.create(path)
            tool.merge_df_dataset(f[path], ktype, df)
            f.close()
    elif stype == conf.STYPE_ASHARE:
        # TODO (重要)，支持ip池并发获取，要不然多code的高频获取过于缓慢
        if rewrite:
            # backtest时，将读取的hist数据写回文件
            df = ts.get_hist_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date)
            if code.isdigit():
                f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
                code_prefix = code[0:3]
                path = '/' + code_prefix + '/' + code
                df = df[tss.SHARE_COLS]
            else:
                f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
                path = '/' + code
                df = df[tss.INDEX_COLS]
            df = df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
            tool.merge_df_dataset(f[path], ktype, df)
            f.close
        else:
            # get_k_data无法获取换手率，但是在线监听时，get_k_data存在时间延迟
            df = ts.get_k_data(code, ktype=ktype, pause=conf.REQUEST_BLANK, start=start_date)
        time.sleep(conf.REQUEST_BLANK)
    else:
        raise Exception("数据源不存在或未配置")

    if df is None and df.empty is True:
        raise Exception("无法获取" + code + "-" + ktype + ":" + start_date + "以后的数据，休息30秒重新获取")
    return df


def append_and_macd(old_df, new_df, last_date, factor_macd_range, direct_turn=False):
    """
    拼接原始与已计算macd的数据
    """
    new_df = new_df[new_df[conf.HDF5_SHARE_DATE_INDEX] > last_date][["date", "close"]]
    old_df = old_df[["date", "close"]]
    old_df = old_df.append(new_df).drop_duplicates(conf.HDF5_SHARE_DATE_INDEX)
    new_trend_df = macd.value_and_trend(old_df, factor_macd_range, direct_turn).tail(len(new_df) + DF_MACD_MIN_NUM)
    old_df = old_df.append(new_trend_df).drop_duplicates(conf.HDF5_SHARE_DATE_INDEX)
    return old_df


def check_reverse(now_date, check_df):
    """
    检查trend是否存在背离
    """
    ret = False
    check_df = check_df[check_df[conf.HDF5_SHARE_DATE_INDEX] <= now_date]
    start, now = phase.now(check_df)
    if start is None:
        return ret
    macd_diff = now["macd"] - start["macd"]
    price_diff = now["close"] - start["close"]
    if (macd_diff < 0 and price_diff > 0) or (macd_diff > 0 and price_diff < 0):
        ret = True
    return ret
