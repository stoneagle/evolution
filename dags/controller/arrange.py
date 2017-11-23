import h5py
from library import conf, console, tool


def arrange_detail():
    # 将basic的detail内容，按个股整理至share文件下
    return


def operate_quit(action_type):
    # 将退市quit内容，转换成标签添加在对应code下
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        action_type,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_QUIT
    )
    path = '/' + conf.HDF5_BASIC_QUIT
    if f.get(path) is None:
        return

    quit_list = [
        conf.HDF5_BASIC_QUIT_TERMINATE,
        conf.HDF5_BASIC_QUIT_SUSPEND,
    ]
    for qtype in quit_list:
        # 如果文件不存在，则退出
        quit_df = tool.df_from_dataset(f[path], qtype, None)
        if quit_df is not None and quit_df.empty is not True:
            quit_df["code"] = quit_df["code"].str.decode("utf-8")
            # 将退市内容，转换成标签添加在对应code下
            tool.op_attr_by_codelist(action_type, quit_df["code"].values, conf.HDF5_BASIC_QUIT, True)
    console.write_tail()
    f.close()
    return


def operate_st(action_type):
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        action_type,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_ST
    )
    path = '/' + conf.HDF5_BASIC_ST
    # 如果文件不存在，则退出
    if f.get(path) is None:
        return
    st_df = tool.df_from_dataset(f[path], conf.HDF5_BASIC_ST, None)
    if st_df is not None and st_df.empty is not True:
        st_df["code"] = st_df["code"].str.decode("utf-8")
        # 将st内容，转换成标签添加在对应code下
        tool.op_attr_by_codelist(action_type, st_df["code"].values, conf.HDF5_BASIC_ST, True)
    console.write_tail()
    f.close()
    return


def arrange_classify_mean():
    # 整理分类均线
    return
