from pytz import timezone
TZ = timezone('Asia/Shanghai')

HDF5_FILE_ROOT = '/tmp/hdf5'
HDF5_FILE_CLASSIFY = HDF5_FILE_ROOT + '/classify.h5'
HDF5_FILE_SHARE = HDF5_FILE_ROOT + '/tmp/hdf5/share.h5'
HDF5_FILE_ERROR = HDF5_FILE_ROOT + '/tmp/hdf5/error.h5'

HDF5_CLASSIFY_INDUSTRY = '/industry'
HDF5_CLASSIFY_CONCEPT = '/concept'
HDF5_CLASSIFY_HOT = '/hot'
