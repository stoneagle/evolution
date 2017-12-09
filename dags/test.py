from controller import strategy
# import h5py
# from library import conf, tool

code_list = ["002466"]

# obtain
# obtain.get_classify()
# obtain.get_basic()
# obtain.get_quit()
# obtain.get_st()
# obtain.get_all_share(True)
# obtain.get_xsg()
# obtain.get_ipo()
# obtain.get_margin()
# obtain.get_index_share()
# f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
# for code in code_list:
#     obtain.get_code_share(f, code, True)
# f.close()

# index
# f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
# for code in code_list:
#     for ktype in conf.HDF5_SHARE_KTYPE:
#         code_prefix = code[0:3]
#         index_df = index.one_index(f, code, ktype, None)
#         tool.delete_dataset(f[code_prefix][code], conf.HDF5_INDEX_DETAIL + "_" + ktype)
#         tool.merge_df_dataset(f[code_prefix][code], conf.HDF5_INDEX_DETAIL + "_" + ktype, index_df.reset_index())
# f.close()
# index.all_index(True, None)

# arrange
# arrange.arrange_all_classify_detail(True, None)
# arrange.arrange_detail(None)
# arrange.arrange_xsg()
# arrange.arrange_ipo()
# arrange.arrange_all_macd_trend(code_list, None)
# arrange.arrange_all_wrap(code_list, None)
# arrange.arrange_margins("sh")
# arrange.arrange_margins("sz")

# grafana
# grafana.push_classify_detail()

# strategy
strategy.macd_divergence(code_list, "30", None)


def idea():
    """
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
            []市场抽水:限售股解禁、ipo上市、融券
            []市场填水:融资，沪港通
        比价
            []横向:各概念均价所处均线位置
            []纵向:个股季度财报，行业内龙头比对
        出入信号
            []出入时机:macd的背离(三类买卖点),成交量
            []波动空间:缠论K线的箱体
            []持续时间:股东人数
        管理
            入场风格:左侧/右侧交易
            资金管理
    数据展示:
        [x] 将量化数据传导至influxdb
        [] 梳理需要展示的内容，以及脚本执行顺序
        [] 在grafana上进行配置
    """
