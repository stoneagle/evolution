import tushare as ts
import yaml
# import pandas as pd
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
share_fields = 'symbol,name,market,list_status,list_date,delist_date'
gplb = pro.stock_basic(list_status='L', fields=share_fields)
gplb.to_sql('tushare_share', cn, index=True, if_exists='append')
