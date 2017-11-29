import h5py
from tsSource import classify, share, basic, fundamental
from library import conf, count, error, console
from controller import arrange


def get_classify():
    """
    获取类别数据
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    for ctype in classify_list:
        console.write_head(
            conf.HDF5_OPERATE_GET,
            conf.HDF5_RESOURCE_TUSHARE,
            ctype
        )
        cpath = '/' + ctype
        if f.get(cpath) is None:
            f.create_group(cpath)

        if ctype == conf.HDF5_CLASSIFY_INDUSTRY:
            # 获取工业分类
            classify.get_industry_classified(f[cpath])
        elif ctype == conf.HDF5_CLASSIFY_CONCEPT:
            # 获取概念分类
            classify.get_concept_classified(f[cpath])
        elif ctype == conf.HDF5_CLASSIFY_HOT:
            # 获取热门分类
            classify.get_hot_classified(f[cpath])
        count.show_result()
        console.write_tail()
    f.close()
    return


def get_code_share(f, code, gem_flag):
    # 按编码前3位分组，如果group不存在，则创建新分组
    code_prefix = code[0:3]

    # 判断是否跳过创业板
    if gem_flag is True and code_prefix == "300":
        return

    code_group_path = '/' + code_prefix + '/' + code
    if f.get(code_group_path) is None:
        f.create_group(code_group_path)
    else:
        # 忽略停牌、退市、无法获取的情况
        if f[code_group_path].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
            return

        if f[code_group_path].attrs.get(conf.HDF5_BASIC_ST) is not None:
            return

    # 获取不同周期的数据
    for ktype in conf.HDF5_SHARE_KTYPE:
        share.get_share_data(code, f[code_group_path], ktype, share.SHARE_TYPE)
    return


def get_index_share():
    index_list = ["sh", "sz", "hs300", "sz50", "zxb", "cyb"]
    # 初始化相关文件
    f = h5py.File(conf.HDF5_FILE_INDEX, 'a')
    for code in index_list:
        console.write_head(
            conf.HDF5_OPERATE_GET,
            conf.HDF5_RESOURCE_TUSHARE,
            code
        )
        code_group_path = '/' + code
        if f.get(code_group_path) is None:
            f.create_group(code_group_path)
        # 获取不同周期的数据
        for ktype in conf.HDF5_SHARE_KTYPE:
            share.get_share_data(code, f[code_group_path], ktype, share.INDEX_TYPE)
        console.write_tail()
    f.close()
    return


def get_all_share(gem_flag):
    """
    根据类别获取所有share数据
    """
    classify_list = [
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    # 初始化相关文件
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    for ctype in classify_list:
        # 初始化error记录
        error.init_batch(conf.HDF5_ERROR_SHARE_GET)

        # 获取概念下具体类别名称
        for classify_name in f_classify[ctype]:
            # 如果不存在code list，则跳过
            if f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                continue

            console.write_head(
                conf.HDF5_OPERATE_GET,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )
            # 获取单个分类(code的list)下的share数据
            for row in f_classify[ctype][classify_name][conf.HDF5_CLASSIFY_DS_CODE]:
                code = row[0].astype(str)
                # TODO, 筛选错误数据
                # history = error.get_file()
                # error_history = list()
                # if history is not None:
                #     history["ktype"] = history["ktype"].str.decode("utf-8")
                #     history["code"] = history["code"].str.decode("utf-8")
                #     error_history = history.values
                get_code_share(f, code, gem_flag)
            # 记录错误内容
            error.write_batch()
            # 输出获取情况
            count.show_result()
            console.write_tail()
    f_classify.close()
    f.close()
    return


def get_basic():
    """
    获取基本面数据
    """
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    # 初始化error记录
    error.init_batch(conf.HDF5_ERROR_DETAIL_GET)
    # 初始化打点记录
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_DETAIL
    )
    path = '/' + conf.HDF5_BASIC_DETAIL
    if f.get(path) is None:
        f.create_group(path)
    basic.get_detail(f[path])

    # 记录错误内容
    error.merge_batch()
    # 展示打点结果
    count.show_result()
    console.write_tail()
    f.close()
    return


def get_quit():
    """
    获取退市、终止上市列表
    """
    # 删除以往添加的标签
    arrange.operate_quit(conf.HDF5_OPERATE_DEL)

    # 获取最新的数据
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_QUIT
    )
    path = '/' + conf.HDF5_BASIC_QUIT
    if f.get(path) is None:
        f.create_group(path)
    basic.get_quit(f[path])
    console.write_tail()
    f.close()

    # 添加标签
    arrange.operate_quit(conf.HDF5_OPERATE_ADD)
    return


def get_st():
    """
    获取风险警示板列表
    """
    # 删除以往添加的标签
    arrange.operate_st(conf.HDF5_OPERATE_DEL)

    # 获取最新的数据
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_ST
    )
    path = '/' + conf.HDF5_BASIC_ST
    if f.get(path) is None:
        f.create_group(path)
    basic.get_st(f[path])
    console.write_tail()
    f.close()

    # 添加标签
    arrange.operate_st(conf.HDF5_OPERATE_ADD)
    return


def get_xsg():
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    # 获取限售股解禁数据
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_XSG
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_XSG
    if f.get(path) is None:
        f.create_group(path)
    fundamental.get_xsg(f[path])
    count.show_result()
    console.write_tail()
    f.close()
    return


def get_ipo(reset_flag=False):
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    # 获取ipo数据
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_IPO
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_IPO
    if f.get(path) is None:
        f.create_group(path)
    fundamental.get_ipo(f[path], reset_flag)
    count.show_result()
    console.write_tail()
    f.close()
    return


def get_margin(reset_flag=False):
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    # 获取沪市融资融券
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SH_MARGINS
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_SH_MARGINS
    if f.get(path) is None:
        f.create_group(path)
    fundamental.get_sh_margins(f[path], reset_flag)
    count.show_result()
    console.write_tail()

    # 获取深市融资融券
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_SZ_MARGINS
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_SZ_MARGINS
    if f.get(path) is None:
        f.create_group(path)
    fundamental.get_sz_margins(f[path], reset_flag)
    count.show_result()
    console.write_tail()
    f.close()
    return
