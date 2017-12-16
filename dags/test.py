from controller import index 
# import h5py
from library import conf

code_list = ["002273"]
omit_list = ["300", "900"]
start_date = None
init_flag = True
classify_list = [
    # conf.HDF5_CLASSIFY_INDUSTRY,
    # conf.HDF5_CLASSIFY_CONCEPT,
    conf.HDF5_CLASSIFY_HOT,
]


# obtain
# obtain.classify_detail(classify_list)
# obtain.quit()
# obtain.st()
# obtain.index_share()
# obtain.all_share(True)
# obtain.basic_environment(start_date)
# obtain.xsg()
# obtain.ipo()
# obtain.margin()
# f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
# for code in code_list:
#     obtain.code_share(f, code, True)
# f.close()


# arrange
# arrange.operate_quit(conf.HDF5_OPERATE_ADD)
# arrange.operate_st(conf.HDF5_OPERATE_ADD)
# arrange.xsg()
# arrange.ipo()
# arrange.margins("sh")
# arrange.margins("sz")
# arrange.all_classify_detail(classify_list, omit_list, start_date)
# arrange.share_detail(start_date, omit_list)
# arrange.all_macd_trend(code_list, start_date)
# arrange.all_wrap(code_list, start_date)


# index
# f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
# for code in code_list:
#     for ktype in conf.HDF5_SHARE_KTYPE:
#         code_prefix = code[0:3]
#         df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
#         index_df = index.one_df(df, True)
#         tool.delete_dataset(f[code_prefix][code], conf.HDF5_INDEX_DETAIL + "_" + ktype)
#         tool.merge_df_dataset(f[code_prefix][code], conf.HDF5_INDEX_DETAIL + "_" + ktype, index_df.reset_index())
# f.close()
index.all_share(omit_list, init_flag)
# index.all_index(init_flag)
# index.all_classify(classify_list, init_flag)


# strategy
# screen.daily(code_list, True)


# grafana
# grafana.classify_detail(classify_list)
# grafana.basic_detail()
# grafana.index_detail()

"""
TODO
[]index的macd与均线，share不计算均线所处位置
"""


def idea():
    """
    数据还原:
        [] 删除所有清洗生成的数据
    数据获取:
    数据清洗:
        [x] 投资基本面，平均每周的ipo、限售股、融资融券数据
        [x] 价格所处均线位置
        [x] 概念均价
        [x] 股票股东人数变化
        [x] 个股K线图，向缠论转换
        [x] 个股MACD趋势
        [] 成交量?
    策略思路:
        基本面
            [x]市场抽水:限售股解禁、ipo上市、融券
            [x]市场填水:融资，沪港通
        比价
            []横向:各概念均价所处均线位置
            []纵向:个股季度财报，行业内龙头比对
        出入信号
            []出入时机:macd的背离(三类买卖点),成交量
            []波动空间:缠论K线的箱体，背离级别与均线位置
            []持续时间:股东人数
        管理
            []入场风格:左侧/右侧交易
            []资金管理
    数据展示:
        [x] 将量化数据传导至influxdb
        [] 梳理需要展示的内容，以及脚本执行顺序
        [] 在grafana上进行配置
    """
