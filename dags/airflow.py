from __future__ import print_function
from airflow.operators.python_operator import PythonOperator
from airflow.models import DAG
from datetime import timedelta
from controller import getData

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

# airflow相关执行

# 先获取类别
first = PythonOperator(
    task_id='get_classify',
    python_callable=get_data.get_classify,
    dag=dag)

# 基于类别获取share数据
second = PythonOperator(
    task_id='get_basic',
    python_callable=get_data.get_basic,
    dag=dag)

# 基于类别获取share数据
third = PythonOperator(
    task_id='get_share',
    python_callable=get_data.get_share,
    # op_kwargs={'random_base': float(i) / 10},
    dag=dag)

second.set_upstream(first)
third.set_upstream(second)
