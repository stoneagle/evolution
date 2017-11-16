from library import conf, tool
import pandas as pd
import h5py


batch = list()


def write_df_batch(g_name, d_name, data_frame):
    f = h5py.File(conf.HDF5_FILE_ERROR)
    group_path = '/' + g_name
    if f.get(group_path) is None:
        f.create_group(group_path)
    group_f = f[group_path]

    if group_f.get(d_name) is None:
        tool.create_df_dataset(group_f, d_name, data_frame.reset_index())
    else:
        tool.append_df_dataset(group_f, d_name, data_frame.reset_index())
    return


def add_row(row):
    return
