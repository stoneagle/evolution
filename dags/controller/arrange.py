import h5py
from library import conf, console, tool


def arrange_detail(start_date):
    """
    将basic的detail内容，按个股整理至share文件下
    """
    # 获取basic所有日期的detail，并遍历读取详细信息
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    f_share = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_DETAIL
    )
    path = '/' + conf.HDF5_BASIC_DETAIL
    if f.get(path) is None:
        return

    code_basic_dict = dict()
    for date in f[path]:
        if start_date is not None and date < start_date:
            continue
        df = tool.df_from_dataset(f[path], date, None)
        df["code"] = df["code"].str.decode("utf-8")
        df = df.set_index("code")
        for code in df.index:
            if code == "600035":
                if code not in code_basic_dict:
                    code_basic_dict[code] = tool.init_empty_df(df.columns)
                code_basic_dict[code].loc[date] = df.loc[code, :]
    for code, code_df in code_basic_dict.items():
        code_df.index.name = conf.HDF5_SHARE_DATE_INDEX
        code_df = code_df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])

        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f_share.get(code_group_path) is None:
            continue

        if start_date is None:
            tool.delete_dataset(f_share[code_group_path], conf.HDF5_BASIC_DETAIL)
        tool.merge_df_dataset(f_share[code_group_path], conf.HDF5_BASIC_DETAIL, code_df)
    console.write_tail()
    f_share.close()
    f.close()
    # 从最初时间开始，按照code聚合detail
    return


def operate_quit(action_type):
    """
    将退市quit内容，转换成标签添加在对应code下
    """
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
    """
    将st内容，转换成标签添加在对应code下
    """
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


def arrange_all_classify_detail(gem_flag, start_date):
    """
    遍历所有分类，聚合所有code获取分类均值
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f_classify[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_ARRANGE,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )
            if classify_name != "chgn_700997":
                continue

            if f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                continue

            for ktype in conf.HDF5_SHARE_KTYPE:
                mean_df = arrange_one_classify_detail(f, f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE), gem_flag, ktype, start_date)
                # 如果start_date为空，则重置该数据
                if start_date is None:
                    tool.delete_dataset(f_classify[ctype][classify_name], conf.HDF5_CLASSIFY_DS_DETAIL)
                tool.merge_df_dataset(f_classify[ctype][classify_name], conf.HDF5_CLASSIFY_DS_DETAIL, mean_df)
            console.write_tail()
    f_classify.close()
    f.close()
    return


def arrange_one_classify_detail(f, code_list, gem_flag, ktype, start_date):
    """
    根据单个分类，聚合所有code获取分类均值
    """
    # 初始化一个全日期的空DataFrame，并初始化一列作为统计个数
    init_df = tool.init_empty_df(conf.HDF5_SHARE_COLUMN)

    # 按照列表顺序，获取code并逐一添加至初始化DF，并递增该日期的个数
    for row in code_list:
        code = row[0].astype(str)
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

        if f[code_group_path].get(ktype) is None:
            continue

        add_df = tool.df_from_dataset(f[code_group_path], ktype, None)
        add_df[conf.HDF5_SHARE_DATE_INDEX] = add_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        add_df["num"] = 1
        add_df = add_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
        init_df = init_df.add(add_df, fill_value=0)
    # 总数除以数量，得到平均值
    init_df = init_df.div(init_df.num, axis=0)
    init_df = init_df.drop("num", axis=1)
    if start_date is not None:
        init_df = init_df.ix[start_date:]
    init_df = init_df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
    return init_df
