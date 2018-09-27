from library import conf, tool
import random
import pandas as pd
import h5py


batch_dict = dict()
# TODO 检查index是否存在
index = ''


def set_index(current_index):
    global index
    index = current_index
    return


def write_batch():
    global index
    if index in batch_dict and batch_dict[index].empty is not True:
        f = h5py.File(conf.HDF5_FILE_ERROR)
        if f.get(index) is None:
            tool.create_df_dataset(f, index, batch_dict[index])
        else:
            tool.append_df_dataset(f, index, batch_dict[index])
        batch_dict[index].drop(batch_dict[index].index, inplace=True)
    return


def merge_batch():
    global index
    if index in batch_dict and batch_dict[index].empty is not True:
        f = h5py.File(conf.HDF5_FILE_ERROR)
        if f.get(index) is None:
            tool.create_df_dataset(f, index, batch_dict[index])
        else:
            tool.merge_df_dataset(f, index, batch_dict[index])
        batch_dict[index].drop(batch_dict[index].index, inplace=True)
    return


def delete_batch(d_name):
    f = h5py.File(conf.HDF5_FILE_ERROR)
    if f.get(d_name) is not None:
        del f[d_name]
    return


def init_batch(current_index):
    set_index(current_index)
    global index
    if index not in batch_dict:
        columns = conf.HDF5_ERROR_COLUMN_MAP.get(index, ["error"])
        batch_dict[index] = pd.DataFrame([], columns=columns)
    return


def get_batch():
    global index
    if index not in batch_dict:
        ret = pd.DataFrame([], columns=['A'])
    else:
        ret = batch_dict[index]
    return ret


def add_row(list_row):
    global index
    default = ['A', 'B', 'C', 'D', 'E', 'F', 'G']
    if index not in batch_dict:
        batch_dict[index] = pd.DataFrame([list_row], columns=random.sample(default, len(list_row)))
    else:
        df = pd.DataFrame([list_row], columns=batch_dict[index].columns)
        batch_dict[index] = batch_dict[index].append(df)
    return


def get_file():
    global index
    f = h5py.File(conf.HDF5_FILE_ERROR)
    if f.get(index) is not None:
        columns = conf.HDF5_ERROR_COLUMN_MAP.get(index, ["error"])
        ret = tool.df_from_dataset(f, index, columns)
    else:
        ret = None
    return ret
