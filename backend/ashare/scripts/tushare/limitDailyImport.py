import tushare as ts
import yaml
import time
from sqlalchemy import create_engine


with open("config.yaml", "r") as ymlfile:
    cfg = yaml.load(ymlfile)
host = 'mysql+pymysql://' \
    + cfg['mysql']['user'] + ':' + cfg['mysql']['passwd'] + '@' \
    + cfg['mysql']['host'] + ':3306/' + cfg['mysql']['db'] + '?charset=utf8'
cn = create_engine(host)
token = cfg['tushare']['token']
ts.set_token(token)
pro = ts.pro_api()
date = time.strftime("%Y%m%d", time.localtime())
gplb = pro.limit_list(trade_date=date, limit_type='U')
gplb.to_sql('tushare_limit_daily', cn, index=False, if_exists='append')
