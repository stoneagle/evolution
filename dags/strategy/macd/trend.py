from strategy.phase import phase, action


def get_macd_trends(index_df):
    """
    获取某个时间段的macd分布情况
    """
    index_df = index_df.reset_index()
    first = index_df.iloc[1]["macd"]
    second = index_df.iloc[2]["macd"]
    start_date = index_df.iloc[1]["date"]
    phase_obj = phase.Phase(first, second, start_date)
    action_machine = action.Action(phase_obj)

    for index, row in index_df[3:].iterrows():
        action_machine.run(row["macd"], row["date"])
    return action_machine.phase.result()


def get_price_relate(trend_df):
    """
    根据阶段时间范围，获取以下数值
    记录各个红柱、绿柱面积
    记录各个红柱、绿柱面积范围的价格变化
    比较前后两段红柱/绿柱的价格高点/低点，判断趋势方向(上行、下行、震荡)
    比较两段绿柱，价格的波动幅度
    """
    return
