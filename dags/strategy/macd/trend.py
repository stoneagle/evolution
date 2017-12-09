from strategy.util import phase


def get_macd_phases(trend_df, share_df):
    """
    获取macd的阶段情况
    按照up、down与shake，记录各个阶段的红柱、绿柱面积
    记录up、down趋势下，从最低点/最高点到0周阶段中，价格的变化与macd是否一致，价格波动幅度以及波动比例
    """
    result = phase.Phase(trend_df, share_df).merge(phase.PTYPE_MACD)
    return result
