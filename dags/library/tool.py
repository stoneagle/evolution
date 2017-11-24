import numpy as np
import h5py
from library import conf
import pandas as pd


def df_to_sarray(df):
    def make_col_type(col_type, col):
        if 'numpy.object_' in str(col_type.type):
            maxlens = col.dropna().str.len()
            if maxlens.any():
                maxlen = maxlens.max().astype(int)
                col_type = ('S%s' % maxlen, 1)
            else:
                col_type = 'f2'
        return col.name, col_type
    v = df.values
    types = df.dtypes
    numpy_struct_types = [make_col_type(types[col], df.loc[:, col]) for col in df.columns]
    dtype = np.dtype(numpy_struct_types)
    z = np.zeros(v.shape[0], dtype)
    for (i, k) in enumerate(z.dtype.names):
        if dtype[i].str.startswith('|S'):
            z[k] = df[k].str.encode('utf-8').astype('S')
        else:
            z[k] = v[:, i]
    return z, dtype


def create_df_dataset(f, dset_name, df):
    if f.get(dset_name) is not None:
        del f[dset_name]

    df_array, df_type = df_to_sarray(df)
    # 使用shape，第二维度不能限制为1，否则会无法转换会dataFrame
    f.create_dataset(dset_name, (len(df_array), ), data=df_array, maxshape=(None, ), dtype=df_type)
    return


def append_df_dataset(f, dset_name, df):
    if f.get(dset_name) is not None:
        df_array, df_type = df_to_sarray(df)
        f[dset_name].resize(f[dset_name].shape[0] + len(df_array), axis=0)
        row_num = 0
        for row in df_array:
            f[dset_name][-len(df_array) + row_num] = row
            row_num += 1
    else:
        create_df_dataset(f, dset_name, df)
    return


def merge_df_dataset(f, dset_name, df):
    if f.get(dset_name) is not None:
        old_df = f[dset_name][:]
        df_array, df_type = df_to_sarray(df)
        for row in df_array:
            if row not in old_df:
                f[dset_name].resize(f[dset_name].shape[0] + 1, axis=0)
                f[dset_name][-1] = row
    else:
        create_df_dataset(f, dset_name, df)
    return


def delete_dataset(f, dset_name):
    if f.get(dset_name) is not None:
        del f[dset_name]
    return


def df_from_dataset(f, dset_name, columns):
    if f.get(dset_name) is not None:
        if columns is None:
            ret = pd.DataFrame(f[dset_name][:])
        else:
            ret = pd.DataFrame(f[dset_name][:], columns=columns)
    else:
        ret = None
    return ret


def op_attr_by_codelist(operator, code_list, attr_name, attr_value):
    f_share = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f_share.get(code_group_path) is not None:
            if operator == conf.HDF5_OPERATE_ADD:
                f_share[code_group_path].attrs[attr_name] = attr_value
            elif operator == conf.HDF5_OPERATE_DEL:
                if f_share[code_group_path].attrs.get(attr_name) is not None:
                    del f_share[code_group_path].attrs[attr_name]
    f_share.close()
    return


def init_empty_df(columns):
    ret = pd.DataFrame(columns=columns)
    return ret
