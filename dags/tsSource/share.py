import tushare as ts
from library import tradetime as ttime
from tsSource import cons
import h5py
import time


def _add_data(code, ktype, f, end_date):
    df = ts.get_hist_data(code, ktype=ktype, pause=1, end=end_date)
    if df.empty is not True:
        df = df[cons.SHARE_COLS]
        df.reset_index(level=0, inplace=True)
        df = df.sort_values(by=['date'])
        df['date'] = df['date'].astype('S')
        f.create_dataset(ktype, (len(df), 7), data=df, maxshape=(None, None), dtype=h5py.special_dtype(vlen=str))
    return


def _append_data(code, ktype, f, start_date, end_date):
    # 如果开始日期大于等于结束日期，则不需要进行处理
    if start_date >= end_date:
        continue
    df = ts.get_hist_data(code, ktype=ktype, pause=1, end=end_date, start=start_date)
    if df.empty is not True:
        df = df[cons.SHARE_COLS]
        df.reset_index(level=0, inplace=True)
        df = df.sort_values(by=['date'])
        df['date'] = df['date'].astype('S')
        if f.get('tmp') is not None:
            del f['tmp']
        tmp_set = f.create_dataset('tmp', (len(df), 7), data=df, maxshape=(len(df), 7), dtype=h5py.special_dtype(vlen=str))
        f.resize(f.shape[0] + len(tmp_set), axis=0)
        f[-len(tmp_set):] = tmp_set
        del f['tmp']
    return


def get_share_data(code, f):
    # 获取不同周期的数据
    for ktype in ['M', 'W', 'D', '30', '5']:
        if f.get(ktype) is None:
            # 如果股票不存在，则获取17年至今数据(M取上个月月底，W取上周日)
            end_date = ttime.get_end_date(code, ktype)
            _add_data(code, ktype, f, end_date)
        else:
            # 如果股票已存在，则根据存储内最后日期，获取至今
            share_set = f[ktype]
            if len(share_set) == 0:
                # 如果数据为空
                end_date = ttime.get_end_date(code, ktype)
                _add_data(code, ktype, f, end_date)
            else:
                tail_date_str = share_set[-1, 0]
                start_date = ttime.get_start_date(tail_date_str, code, ktype)
                end_date = ttime.get_end_date(code, ktype)
                _append_data(code, ktype, f, start_date, end_date)
        time.sleep(1)
    return
