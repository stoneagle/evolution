from controller import obtain, arrange, index, grafana
from library import conf


def get_basic(classify_list):
    """
    获取基本信息
    """
    # 获取基本面
    obtain.basic_environment()
    # 获取ipo股票列表
    obtain.ipo()
    # 获取融资融券
    obtain.margin()
    # 获取限售股解禁数据
    obtain.xsg()
    # 获取类别
    obtain.classify_detail(classify_list)
    return


def get_share():
    """
    获取股票数据
    """
    # 根据类别获取所有股票
    obtain.all_share(False)
    # 根据类别获取所有指数
    obtain.index_share()
    # 获取暂停上市股票列表
    obtain.quit()
    # 获取风险警示板股票列表
    obtain.st()
    return


def arrange_all(classify_list, start_date):
    """
    处理并聚合已获取数据
    """
    # 标记股票是否退市
    arrange.operate_quit(conf.HDF5_OPERATE_ADD)
    # 标记股票是否st
    arrange.operate_st(conf.HDF5_OPERATE_ADD)
    # 聚合xsg数据
    arrange.xsg()
    # 聚合ipo数据
    arrange.ipo()
    # 聚合sh融资融券
    arrange.margins("sh")
    # 聚合sz融资融券
    arrange.margins("sz")
    # 按照分类code，获取分类的均值
    arrange.all_classify_detail(classify_list, True, start_date)
    # # 将basic的detail，聚合至share对应code
    # arrange.share_detail(start_date)
    return


def index_exec(classify_list, start_date):
    """
    清洗并获取指数
    """
    # 整理classify分类的均值、macd等
    index.all_classify(classify_list, start_date)
    # 整理index的均值、macd等
    index.all_index(start_date)
    # 获取所有股票的均值、macd等
    index.all_share(True, None)
    return


def strategy_share():
    """
    策略选股
    """
    # 整理对应股票的macd趋势
    # code_list = []
    # arrange.all_macd_trend(code_list, None)
    # 整理对应股票的缠论k线
    # arrange.all_wrap(code_list, None)
    return


def grafana_push(classify_list):
    """
    推送相关数据至influxdb
    """
    # 推送ipo、xsg、shm、szm等数据
    grafana.basic_detail()
    # 推送classify数据
    grafana.classify_detail(classify_list)
    # 推送index数据
    grafana.index_detail()
    return


start_date = None
classify_list = [
    # conf.HDF5_CLASSIFY_INDUSTRY,
    # conf.HDF5_CLASSIFY_CONCEPT,
    conf.HDF5_CLASSIFY_HOT,
]


get_basic(classify_list)
get_share()
arrange_all(classify_list, start_date)
index_exec(classify_list, start_date)
# strategy_share()
grafana_push(classify_list)
