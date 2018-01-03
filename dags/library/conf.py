from pytz import timezone
import os
from library import basic
basic.import_env()
TZ = timezone('Asia/Shanghai')
TZ_UTC = timezone('Utc')

# 请求间隔
REQUEST_BLANK = 1

RUNMODE = os.environ.get("RUNMODE")
if RUNMODE == "dev":
    ASHARE_FILE_ROOT = '/home/wuzhongyang/database/ashare'
    FUTURE_FILE_ROOT = '/home/wuzhongyang/database/future'
else:
    ASHARE_FILE_ROOT = '/tmp/hdf5'

HDF5_FILE_CLASSIFY = ASHARE_FILE_ROOT + '/classify.h5'
HDF5_FILE_SHARE = ASHARE_FILE_ROOT + '/share.h5'
HDF5_FILE_INDEX = ASHARE_FILE_ROOT + '/index.h5'
HDF5_FILE_BASIC = ASHARE_FILE_ROOT + '/basic.h5'
HDF5_FILE_ERROR = ASHARE_FILE_ROOT + '/error.h5'
HDF5_FILE_FUNDAMENTAL = ASHARE_FILE_ROOT + '/fundamental.h5'
HDF5_FILE_SCREEN = ASHARE_FILE_ROOT + '/screen.h5'
HDF5_FILE_OTHER = ASHARE_FILE_ROOT + '/other.h5'

HDF5_FILE_BITMEX = FUTURE_FILE_ROOT + '/bitmex.h5'

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

HDF5_CLASSIFY_REFRESH_DAYS_BLANK = 7
HDF5_CLASSIFY_NAME_ATTR = 'name'
HDF5_CLASSIFY_REFRESH_ATTR = 'datetime'

HDF5_CLASSIFY_DS_CODE = 'codelist'
HDF5_CLASSIFY_DS_DETAIL = 'detail'

HDF5_BASIC_DETAIL = 'detail'
HDF5_BASIC_ST = 'st'
HDF5_BASIC_QUIT = 'quit'
HDF5_BASIC_QUIT_SUSPEND = 'suspend'
HDF5_BASIC_QUIT_TERMINATE = 'terminate'

HDF5_SHARE_DETAIL = "share"
HDF5_SHARE_KTYPE = ["M", "W", "D", "30", "5"]
HDF5_SHARE_WRAP_KTYPE = ["D", "30"]
HDF5_SHARE_DATE_INDEX = "date"
HDF5_SHARE_COLUMN = ["open", "high", "close", "low", "volume", "turnover"]

HDF5_INDEX_COLUMN = ["dif", "dea", "macd"]
HDF5_INDEX_DETAIL = "index"
HDF5_INDEX_MACD_TREND = "macd_trend"
HDF5_INDEX_WRAP = "wrap"
HDF5_INDEX_PHASE = "phase"

HDF5_RESOURCE_TUSHARE = "tushare"
HDF5_RESOURCE_BITMEX = "bitmex"
HDF5_RESOURCE_DATASET = "dataset"
HDF5_RESOURCE_GROUP = "group"
HDF5_RESOURCE_ATTR = "attr"

HDF5_OTHER_CODE_CLASSIFY = "code_classify"

HDF5_OPERATE_GET = "get"
HDF5_OPERATE_ARRANGE = "arrange"
HDF5_OPERATE_INDEX = "index"
HDF5_OPERATE_TACTICS = "tactics"
HDF5_OPERATE_SCREEN = "screen"
HDF5_OPERATE_WRAP = "wrap"
HDF5_OPERATE_PUSH = "push"
HDF5_OPERATE_ADD = "add"
HDF5_OPERATE_DEL = "del"

HDF5_FUNDAMENTAL_XSG = "xsg"
HDF5_FUNDAMENTAL_IPO = "ipo"
HDF5_FUNDAMENTAL_SH_MARGINS = "shm"
HDF5_FUNDAMENTAL_SZ_MARGINS = "szm"
HDF5_FUNDAMENTAL_XSG_DETAIL = "xsg_detail"
HDF5_FUNDAMENTAL_IPO_DETAIL = "ipo_detail"
HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL = "shm_detail"
HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL = "szm_detail"

HDF5_TACTICS_MACD = "macd"
HDF5_TACTICS_WRAP_K = "wrap_k"
HDF5_TACTICS_WRAP_K_BOX = "wrap_k_box"

INFLUXDB_USER = "wuzhongyang"
INFLUXDB_PASSWORD = "a1b2c3d4E"
INFLUXDB_DBNAME = "tushare"
INFLUXDB_HOST = "localhost"
INFLUXDB_PORT = 18086
INFLUXDB_PROTOCOL_JSON = "json"

SCREEN_SHARE_FILTER = "share_filter"
SCREEN_SHARE_GRADE = "share_grade"

MEASUREMENT_SHARE = "share"
MEASUREMENT_SHARE_WRAP = "share_wrap"
MEASUREMENT_SHARE_BASIC = "share_basic"
MEASUREMENT_BASIC = "basic"
MEASUREMENT_SCREEN = "screen"
MEASUREMENT_INDEX = "index"
MEASUREMENT_INDEX_WRAP = "index_wrap"
MEASUREMENT_CLASSIFY = "classify"
MEASUREMENT_CLASSIFY_WRAP = "classify_wrap"
MEASUREMENT_CODE_CLASSIFY = "code_classify"
MEASUREMENT_FILTER_SHARE = "filter_share"
MEASUREMENT_FILTER_SHARE_GRADE = "filter_share_grade"

WEIXIN_BOT_CACHE_PATH = "/home/wuzhongyang/www/airflow/tmp"

STRATEGY_TREND_AND_REVERSE = "trend_and_reverse"

BITMEX_HOST = "https://www.bitmex.com"
BITMEX_CURRENCY_XBT = "XBt"
BITMEX_URL_HISTORY = "/api/udf/history"
BITMEX_URL_TRADE_BUCKETED = "/api/v1/trade/bucketed"
BITMEX_URL_ORDERBOOK = "/api/v1/orderBook/L2"
BITMEX_URL_WALLET = "/api/v1/user/wallet"
BITMEX_URL_POSITION = "/api/v1/position"
BITMEX_URL_WALLET_HISTORY = "/api/v1/user/walletHistory"
BITMEX_URL_ORDER = "/api/v1/order"
BITMEX_URL_ORDER_CANCEL_ALL = "/api/v1/order/all"
BITMEX_URL_ORDER_CANCEL_ALL_AFTER = "/api/v1/order/cancelAllAfter"
BITMEX_XBTUSD = "XBTUSD"
BITMEX_BXBT = ".BXBT"

BITMEX_APIKEY = os.environ.get("APIKEY")
BITMEX_APISECRET = os.environ.get("APISECRET")

KTYPE_MONTH = "M"
KTYPE_WEEK = "W"
KTYPE_DAY = "D"
KTYPE_THIRTY = "30"
KTYPE_FIVE = "5"
KTYPE_ONE = "1"

BINSIZE_ONE_MINUTE = "1m"
BINSIZE_FIVE_MINUTE = "5m"
BINSIZE_THIRTY_MINUTE = "30m"
BINSIZE_ONE_HOUR = "1h"
BINSIZE_FOUR_HOUR = "4h"
BINSIZE_ONE_DAY = "1d"

STYPE_ASHARE = "ashare"
STYPE_BITMEX = "bitmex"

BUY_SIDE = "buy"
SELL_SIDE = "sell"
