from pytz import timezone
TZ = timezone('Asia/Shanghai')

HDF5_FILE_ROOT = '/home/wuzhongyang/database/hdf5'
HDF5_FILE_CLASSIFY = HDF5_FILE_ROOT + '/classify.h5'
HDF5_FILE_SHARE = HDF5_FILE_ROOT + '/share.h5'
HDF5_FILE_BASIC = HDF5_FILE_ROOT + '/basic.h5'
HDF5_FILE_ERROR = HDF5_FILE_ROOT + '/error.h5'

HDF5_ERROR_SHARE_GET = 'share_get'

HDF5_CLASSIFY_INDUSTRY = 'industry'
HDF5_CLASSIFY_CONCEPT = 'concept'
HDF5_CLASSIFY_HOT = 'hot'
HDF5_BASIC_DETAIL = 'detail'
