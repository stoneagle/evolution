import h5py
from source.ts import classify, share, basic, fundamental
from source.bitmex import future
from library import conf, count, error, console, tool
from controller import arrange
from quota import macd


def classify_detail(classify_list):
    """
    获取类别数据
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
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


def code_share(f, code):
    """
    获取对应code的股票数据
    """
    # 按编码前3位分组，如果group不存在，则创建新分组
    code_prefix = code[0:3]

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


def index_share():
    """
    获取对应指数数据
    """
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


def all_share(omit_list):
    """
    根据类别获取所有share数据
    """
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_RESOURCE_TUSHARE
    )
    # 初始化相关文件
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    code_list = []
    classify_list = [
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY,
        conf.HDF5_CLASSIFY_HOT,
    ]
    for ctype in classify_list:
        # 获取概念下具体类别名称
        for classify_name in f_classify[ctype]:
            # 如果不存在code list，则跳过
            if f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                continue
            # 获取单个分类(code的list)下的share数据
            for row in f_classify[ctype][classify_name][conf.HDF5_CLASSIFY_DS_CODE]:
                code = row[0].astype(str)
                if code not in code_list:
                    code_list.append(code)
    f_classify.close()

    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # 初始化error记录
    error.init_batch(conf.HDF5_ERROR_SHARE_GET)
    for code in code_list:
        code_prefix = code[0:3]
        if code_prefix in omit_list:
            continue
        # TODO, 筛选错误数据
        # history = error.get_file()
        # error_history = list()
        # if history is not None:
        #     history["ktype"] = history["ktype"].str.decode("utf-8")
        #     history["code"] = history["code"].str.decode("utf-8")
        #     error_history = history.values
        code_share(f, code)
    # 记录错误内容
    error.write_batch()
    # 输出获取情况
    count.show_result()
    f.close()
    console.write_tail()
    return


def basic_environment(start_date):
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
    basic.get_detail(f[path], start_date)

    # 记录错误内容
    error.merge_batch()
    # 展示打点结果
    count.show_result()
    console.write_tail()
    f.close()
    return


def quit():
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


def st():
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


def xsg():
    """
    获取限售股解禁数据
    """
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
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


def ipo(reset_flag=False):
    """
    获取ipo数据
    """
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
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


def margin(reset_flag=False):
    """
    获取沪市融资融券
    """
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
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


def bitmex(symbol, bin_size, count):
    """
    bitmex期货数据
    """
    console.write_head(
        conf.HDF5_OPERATE_GET,
        conf.HDF5_RESOURCE_BITMEX,
        symbol + '-' + bin_size
    )
    f = h5py.File(conf.HDF5_FILE_FUTURE, 'a')
    if f.get(conf.HDF5_RESOURCE_BITMEX) is None:
        f.create_group(conf.HDF5_RESOURCE_BITMEX)
    if f[conf.HDF5_RESOURCE_BITMEX].get(symbol) is None:
        f[conf.HDF5_RESOURCE_BITMEX].create_group(symbol)

    # 处理D,5,1数据
    if bin_size in [
        future.BINSIZE_ONE_DAY,
        future.BINSIZE_ONE_HOUR,
        future.BINSIZE_FIVE_MINUTE,
        future.BINSIZE_ONE_MINUTE,
    ]:
        df = future.history(symbol, bin_size, count)
    # 处理30m数据
    elif bin_size in [
        future.BINSIZE_THIRTY_MINUTE,
        future.BINSIZE_FOUR_HOUR
    ]:
        df = future.history_merge(symbol, bin_size, count)

    if df is not None and df.empty is not True:
        df = df.set_index(conf.HDF5_SHARE_DATE_INDEX)
        index_df = macd.value(df)
        df = df.merge(index_df, left_index=True, right_index=True, how='left')
        df = df.dropna().reset_index()
        tool.merge_df_dataset(f[conf.HDF5_RESOURCE_BITMEX][symbol], bin_size, df)

    f.close()
    console.write_tail()
    return
