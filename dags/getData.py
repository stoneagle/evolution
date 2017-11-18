# from __future__ import print_function
# import airflow
# from airflow.operators.python_operator import PythonOperator
# from airflow.models import DAG
# from datetime import timedelta

import h5py
from tsSource import classify, share, basic
from library import conf, count, error, console

# args = {
#     'owner': 'wuzhongyang',
#     'start_date': airflow.utils.dates.days_ago(1)
# }

# dag = DAG(
#     dag_id='TS-AShare',
#     default_args=args,
#     description='从tushare源获取A股数据',
#     schedule_interval=timedelta(days=1, hours=1)
# )


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
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # 初始化error记录
    error.set_index(conf.HDF5_ERROR_SHARE_GET)
    error.init_batch(['ktype', 'code'])
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
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    path = '/' + conf.HDF5_BASIC_DETAIL
    if f.get(path) is None:
        f.create_group(path)
    basic.get_detail(f[path])
    count.show_result()
    f.close()
    return


def get_suspended():
    """
    获取退市列表
    """


get_basic()
# airflow相关执行

# # 先获取类别
# first = PythonOperator(
#     task_id='get_classify',
#     python_callable=get_classify,
#     dag=dag)

# # 基于类别获取share数据
# second = PythonOperator(
#     task_id='get_share',
#     python_callable=get_share,
#     # op_kwargs={'random_base': float(i) / 10},
#     dag=dag)
# second.set_upstream(first)
