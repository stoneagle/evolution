import h5py
from tsSource import classify, share
from library import conf, count, error


def get_classify():
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    classify_list = [
        conf.HDF5_CLASSIFY_HOT,
        conf.HDF5_CLASSIFY_CONCEPT,
        conf.HDF5_CLASSIFY_INDUSTRY
    ]
    for ctype in classify_list:
        if f.get(ctype) is None:
            f.create_group(ctype)
        if ctype == conf.HDF5_CLASSIFY_INDUSTRY:
            # 获取工业分类
            classify.get_industry_classified(f[ctype])
        elif ctype == conf.HDF5_CLASSIFY_CONCEPT:
            # 获取概念分类
            classify.get_concept_classified(f[ctype])
        elif ctype == conf.HDF5_CLASSIFY_HOT:
            # 获取热门分类
            classify.get_hot_classified(f[ctype])
    f.close()
    return


def get_share():
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # code_list = ["000805", "600003", "600001"]
    code_list = ["000805"]
    for code in code_list:
        # 按编码前3位分组，如果group不存在，则创建新分组
        count.inc_by_index("industry")
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            f.create_group(code_group_path)
        code_group_f = f[code_group_path]
        share.get_share_data(code, code_group_f)
    count.show_result()
    f.close()
    return


# get_classify()
# get_share()

error.add_row()

# IDEA
# * 数据获取
# ** 股价、股东等基本面信息
# * 数据清洗
# ** 退市股票的标记
# * 数据筛选
# ** 分类均价等统计
