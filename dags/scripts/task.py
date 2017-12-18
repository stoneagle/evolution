import airflow
from airflow.models import DAG
from airflow.operators.python_operator import PythonOperator

from datetime import timedelta
from controller import obtain
# from controller import obtain, arrange, index
# from library import conf

args = {
    'owner': 'wuzhongyang',
    'start_date': airflow.utils.dates.days_ago(1)
}

dag = DAG(
    dag_id='TS-AShare',
    default_args=args,
    description='从tushare源获取A股数据',
    schedule_interval=timedelta(days=1)
)


def get_basic():
    """
    获取基本信息
    """
    # 获取基本面
    obtain.get_basic()
    # 获取ipo股票列表
    obtain.get_ipo()
    # 获取融资融券
    obtain.get_margin()
    # 获取限售股解禁数据
    obtain.get_xsg()
    return


def get_share():
    """
    获取股票数据
    """
    # 获取类别
    obtain.get_classify()
    # 根据类别获取所有股票
    obtain.get_all_share(False)
    # 根据类别获取所有指数
    obtain.get_index_share()
    # 获取暂停上市股票列表
    obtain.get_quit()
    # 获取风险警示板股票列表
    obtain.get_st()
    return


# def arrange_all():
#     """
#     处理并聚合已获取数据
#     """
#     # 标记股票是否退市
#     arrange.operate_quit(conf.HDF5_OPERATE_ADD)
#     # 标记股票是否st
#     arrange.operate_st(conf.HDF5_OPERATE_ADD)
#     # 将basic的detail，聚合至share对应code
#     arrange.arrange_detail(None)
#     # 按照分类code，获取分类的均值
#     arrange.arrange_all_classify_detail(True, None)
#     # 聚合xsg数据
#     arrange.arrange_xsg()
#     # 聚合ipo数据
#     arrange.arrange_ipo()
#     # 聚合sh融资融券
#     arrange.arrange_margins("sh")
#     # 聚合sz融资融券
#     arrange.arrange_margins("sz")
#     return


# def index_share():
#     """
#     清洗并获取指数
#     """
#     # 获取所有股票的均值、macd等
#     index.all_index(True, None)
#     # code_list = []
#     # # 整理对应股票的macd趋势
#     # arrange.arrange_all_macd_trend(code_list, None)
#     # # 整理对应股票的缠论k线
#     # arrange.arrange_all_wrap(code_list, None)
#     return


basic_act = PythonOperator(
    task_id='get_basic',
    provide_context=True,
    python_callable=get_basic,
    dag=dag)


share_act = PythonOperator(
    task_id='get_share',
    provide_context=True,
    python_callable=get_share,
    dag=dag)


# arrange_act = PythonOperator(
#     task_id='arrange_all',
#     provide_context=True,
#     python_callable=arrange_all,
#     dag=dag)


# index_act = PythonOperator(
#     task_id='index_share',
#     python_callable=index_share,
#     dag=dag)

share_act.set_upstream(basic_act)
# arrange_act.set_upstream(share_act)
# index_act.set_upstream(arrange_act)
