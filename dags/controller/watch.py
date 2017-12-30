from watch.ts import five as tsfive
from library import basic
import os
basic.import_env()


def tushare(code_list):
    backtest = os.environ.get("BACKTEST")
    tsfive.exec(code_list, backtest)
    return
