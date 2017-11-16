from __future__ import print_function
import airflow
from airflow.operators.python_operator import PythonOperator
from airflow.models import DAG

from datetime import timedelta
import h5py
from tsSource import classify, share
from library import conf, console

args = {
    'owner': 'wuzhongyang',
    'start_date': airflow.utils.dates.days_ago(1)
}

dag = DAG(
    dag_id='TS-AShare',
    default_args=args,
    description='从tushare源获取A股数据',
    schedule_interval=timedelta(days=1, hours=1)
)


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
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    # 获取对应类别
    group_list = f_classify[conf.HDF5_CLASSIFY_INDUSTRY].keys()
    console.write_head(conf.HDF5_CLASSIFY_INDUSTRY)
    for classify_name in group_list:
        for row in f_classify[conf.HDF5_CLASSIFY_INDUSTRY + '/' + classify_name]:
            code = row[0].astype(str)
            # 按编码前3位分组，如果group不存在，则创建新分组
            code_prefix = code[0:3]
            code_group_path = '/' + code_prefix + '/' + code
            if f.get(code_group_path) is None:
                f.create_group(code_group_path)
            code_group_f = f[code_group_path]
            share.get_share_data(code, code_group_f)
    console.write_tail(conf.HDF5_CLASSIFY_INDUSTRY)
    f.close()
    f_classify.close()
    return


# 先获取类别
first = PythonOperator(
    task_id='get_classify',
    python_callable=get_classify,
    dag=dag)

# 基于类别获取share数据
second = PythonOperator(
    task_id='get_share',
    python_callable=get_share,
    # op_kwargs={'random_base': float(i) / 10},
    dag=dag)
second.set_upstream(first)
