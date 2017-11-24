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


def arrange_all_classify_detail(gem_flag):
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    # 获取classify列表
    for ctype in classify_list:
        console.write_head(
            conf.HDF5_OPERATE_ARRANGE,
            conf.HDF5_RESOURCE_TUSHARE,
            ctype
        )
        for classify_name in f_classify[ctype]:
            if f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                continue
            for row in f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE):
                code = row[0].astype(str)
                arrange_code_detail(f, code, gem_flag)
        console.write_tail()
    f_classify.close()
    f.close()
    return


def arrange_code_detail(f, code, gem_flag):
    code_prefix = code[0:3]
    # 判断是否跳过创业板
    if gem_flag is True and code_prefix == "300":
        return
    code_group_path = '/' + code_prefix + '/' + code
    if f.get(code_group_path) is None:
        # TODO 为空时立刻获取数据
        print(code)
    else:
        # 忽略停牌、退市、无法获取的情况
        if f[code_group_path].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
            return

        if f[code_group_path].attrs.get(conf.HDF5_BASIC_ST) is not None:
            return
    # 初始化一个全日期的空DataFrame，并初始化一列作为统计个数
    # 按照各周线顺序，列表顺序，日期顺序，获取code并逐一添加至初始化DF，并递增该日期的个数
    # 总数除以数量，得到平均值
    return
