from strategy.macd import trendAndReverse
from library import conf


def daily(stype, omit_list, today_str=None):
    if stype == conf.STRATEGY_TREND_AND_REVERSE:
        trendAndReverse.all_exec(omit_list)
        trendAndReverse.mark_grade(today_str)
    return
