import h5py
from tsSource import classify, share, basic
from library import conf, count, error, console


def get_classify():
    """
    获取类别数据
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        conf.HDF5_CLASSIFY_HOT,
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY
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
    # 初始化相关文件
    f_classify = h5py.File(conf.HDF5_CLASSIFY_HOT, 'a')
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # 初始化error记录
    error.init_batch(conf.HDF5_ERROR_SHARE_GET)
    # 获取对应类别
    group_list = f_classify[conf.HDF5_CLASSIFY_INDUSTRY].keys()
    for classify_name in group_list:
        console.write_head(classify_name)
        for row in f_classify[conf.HDF5_CLASSIFY_INDUSTRY + '/' + classify_name]:
            code = row[0].astype(str)
            # 按编码前3位分组，如果group不存在，则创建新分组
            code_prefix = code[0:3]
            code_group_path = '/' + code_prefix + '/' + code
            if f.get(code_group_path) is None:
                f.create_group(code_group_path)
            share.get_share_data(code, f[code_group_path])
        # 记录错误内容
        error.write_batch()
        # 输出获取情况
        count.show_result()
        console.write_tail(classify_name)
    f.close()
    f_classify.close()
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
