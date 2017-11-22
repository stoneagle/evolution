import h5py
from tsSource import classify, share, basic
from library import conf, count, error, console


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
        console.write_head(ctype)
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
        console.write_tail(ctype)
    f.close()
    return


def get_share():
    """
    根据类别获取share数据
    """
    classify_list = [
        # conf.HDF5_CLASSIFY_CONCEPT,
        # conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    for ctype in classify_list:
        # 初始化相关文件
        # 初始化error记录
        error.init_batch(conf.HDF5_ERROR_SHARE_GET)

        # TODO, 获取错误列表
        # history = error.get_file()
        # error_history = list()
        # if history is not None:
        #     history["ktype"] = history["ktype"].str.decode("utf-8")
        #     history["code"] = history["code"].str.decode("utf-8")
        #     error_history = history.values

        # 获取对应类别
        group_list = f_classify[ctype].keys()
        for classify_name in group_list:
            console.write_head(classify_name)
            for row in f_classify[ctype + '/' + classify_name]:
                code = row[0].astype(str)
                # 按编码前3位分组，如果group不存在，则创建新分组
                code_prefix = code[0:3]
                code_group_path = '/' + code_prefix + '/' + code
                if f.get(code_group_path) is None:
                    f.create_group(code_group_path)

                # 忽略停牌、退市、无法获取的情况
                # 获取不同周期的数据
                for ktype in ["M", "W", "D", "30", "5"]:
                    share.get_share_data(code, f[code_group_path], ktype)
            # 记录错误内容
            error.write_batch()
            # 输出获取情况
            count.show_result()
            console.write_tail(classify_name)
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
    console.write_head(conf.HDF5_BASIC_DETAIL)
    path = '/' + conf.HDF5_BASIC_DETAIL
    if f.get(path) is None:
        f.create_group(path)
    basic.get_detail(f[path])

    # 记录错误内容
    error.merge_batch()
    # 展示打点结果
    count.show_result()
    console.write_tail(conf.HDF5_BASIC_DETAIL)
    f.close()
    return


def get_quit():
    """
    获取退市、终止上市列表
    """
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(conf.HDF5_BASIC_QUIT)
    path = '/' + conf.HDF5_BASIC_QUIT
    if f.get(path) is None:
        f.create_group(path)
    basic.get_quit(f[path])
    console.write_tail(conf.HDF5_BASIC_QUIT)
    f.close()
    return


def idea():
    """
    数据获取:
        指数
    数据清洗:
        处理退市、已过期的股票，避免二次获取
    量化策略结果:
        概念均价
        股票股东人数变化
    """