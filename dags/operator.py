import airflow
from airflow.models import DAG
from airflow.operators.python_operator import PythonOperator

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


def get_all():
    getData.get_classify()
    getData.get_basic()
    getData.get_share()
    return


# 先获取类别
first = PythonOperator(
    task_id='get_all',
    python_callable=get_all,
    dag=dag)
