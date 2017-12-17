from controller import grafana
# import h5py
from library import conf

code_list = ["002273"]
omit_list = ["000", "001", "600", "601", "602", "603", "300", "900"]
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
# arrange.code_detail(code_list, start_date)
# arrange.all_macd_trend(code_list, start_date)
# arrange.code_classify(code_list, classify_list)


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
# index.all_share(omit_list, init_flag)
# index.all_index(init_flag)
# index.all_classify(classify_list, init_flag)


# strategy
# screen.daily(omit_list)
# f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
# for code in code_list:
#     screen._daily_code(f, code)
# f.close()


# wrap
# wrap.filter_share(code_list, start_date)
# wrap.all_index()
# wrap.all_classify()

# grafana
# grafana.basic_detail()
# grafana.classify_detail(classify_list)
# grafana.index_detail()
# grafana.share_detail(code_list)
# grafana.daily_filter("2017-12-16")
# grafana.code_classify()
