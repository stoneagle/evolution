from strategy.phase import phase


def get_macd_phases(index_df, trend_df):
    """
    获取macd的阶段情况
    记录各个红柱、绿柱面积
    记录各个红柱、绿柱面积范围的价格变化
    比较前后两段红柱/绿柱的价格高点/低点，判断趋势方向(上行、下行、震荡)
    比较两段绿柱，价格的波动幅度
    """
    phase_obj = phase.Phase()
    return phase_obj
