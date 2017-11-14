import h5py
from tsSource import classify
from library import conf

f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
if f.get(conf.HDF5_CLASSIFY_INDUSTRY) is None:
    f.create_group(conf.HDF5_CLASSIFY_INDUSTRY)
classify.get_industry_classified(f[conf.HDF5_CLASSIFY_INDUSTRY])


if f.get(conf.HDF5_CLASSIFY_CONCEPT) is None:
    f.create_group(conf.HDF5_CLASSIFY_CONCEPT)
classify.get_concept_classified(f[conf.HDF5_CLASSIFY_CONCEPT])


if f.get(conf.HDF5_CLASSIFY_HOT) is None:
    f.create_group(conf.HDF5_CLASSIFY_HOT)
classify.get_hot_classified(f[conf.HDF5_CLASSIFY_HOT])


f.close()
