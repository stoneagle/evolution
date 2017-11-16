import numpy as np


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
            z[k] = df[k].str.encode('latin').astype('S')
        else:
            z[k] = v[:, i]
    return z, dtype


def create_df_dataset(f, dset_name, df):
    df_array, df_type = df_to_sarray(df)
    f.create_dataset(dset_name, (len(df_array), 1), data=df_array, maxshape=(None, None), dtype=df_type)
    return


def append_df_dataset(f, dset_name, df):
    df_array, df_type = df_to_sarray(df)
    f[dset_name].resize(f[dset_name].shape[0] + len(df_array), axis=0)
    row_num = 0
    for row in df_array:
        f[dset_name][-len(df_array) + row_num, :] = row
        row_num += 1
    return
