import h5py
from library import conf
from tsSource import share 


f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
f = h5py.File(conf.HDF5_FILE_SHARE, 'a')

# 获取对应类别
group_list = f_classify[conf.HDF5_CLASSIFY_INDUSTRY].keys()
for classify_name in group_list:
    for row in f_classify[conf.HDF5_CLASSIFY_INDUSTRY + '/' + classify_name]:
        code = row[0].astype('str')
        # 如果group不存在，则穿件新分组
        if f.get('/' + code) is None:
            f.create_group('/' + code)
        code_group = f['/' + code]
        share.get_share_data(code, code_group)

f.close()
f_classify.close()
