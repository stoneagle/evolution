from pytz import timezone
import os

TZ = timezone('Asia/Shanghai')

RUN_MODE = os.environ.get("runmode")
if RUN_MODE == "prod":
    HDF5_FILE_ROOT = '/tmp/hdf5'
else:
    HDF5_FILE_ROOT = '/home/wuzhongyang/database/hdf5'

HDF5_FILE_CLASSIFY = HDF5_FILE_ROOT + '/classify.h5'
HDF5_FILE_SHARE = HDF5_FILE_ROOT + '/share.h5'
HDF5_FILE_BASIC = HDF5_FILE_ROOT + '/basic.h5'
HDF5_FILE_ERROR = HDF5_FILE_ROOT + '/error.h5'

HDF5_ERROR_SHARE_GET = 'share_get'
HDF5_ERROR_DETAIL_GET = 'detail_get'
HDF5_ERROR_COLUMN_MAP = {
    HDF5_ERROR_SHARE_GET: ['ktype', 'code'],
    HDF5_ERROR_DETAIL_GET: ['type', 'date'],
}

HDF5_COUNT_GET = 'get'
HDF5_COUNT_PASS = 'pass'

HDF5_CLASSIFY_INDUSTRY = 'industry'
HDF5_CLASSIFY_CONCEPT = 'concept'
HDF5_CLASSIFY_HOT = 'hot'
HDF5_BASIC_DETAIL = 'detail'
HDF5_BASIC_ST = 'st'
HDF5_BASIC_QUIT = 'quit'
HDF5_BASIC_QUIT_SUSPEND = 'suspend'
HDF5_BASIC_QUIT_TERMINATE = 'terminate'
