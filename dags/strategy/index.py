import numpy as np
INDEX_MA_BORDER = "ma_border"


def mean_position(index_df, value_column):
    # 计算close价格所处均线位置
    tmp_mean_dict = dict()
    for i in range(0, 241):
        tmp_mean_dict[i] = index_df[value_column].rolling(window=i).mean()

    for index, row in index_df.iterrows():
        # 先确认均线是上行还是下行
        ma5_price = tmp_mean_dict[5].loc[index]
        ma10_price = tmp_mean_dict[10].loc[index]
        if np.isnan(ma5_price) or np.isnan(ma10_price):
            index_df.loc[index, INDEX_MA_BORDER] = np.NaN
            continue

        if ma5_price >= ma10_price:
            up_flag = True
        else:
            up_flag = False

        border = 0
        # 1. 如果下行，从小于找起直到出现大于;如果上行，从大于找起直到出现小于
        for i in [30, 60, 240]:
            compare_price = tmp_mean_dict[i].loc[index]
            if compare_price is np.NaN:
                break

            if (ma5_price > compare_price and up_flag is False) or (ma5_price < compare_price and up_flag is True):
                border = i
                break

        # 2. 根据边界均线，按比例调整，找出离close最贴近的均线
        if border != 0:
            border_price = tmp_mean_dict[border].loc[index]
            index_df.loc[index, INDEX_MA_BORDER] = border
            min_diff = 0
            min_diff_ma_num = 0
            if (up_flag is False and row[value_column] > border_price) or (up_flag is True and row[value_column] < border_price):
                for i in range(0, border - 5):
                    tmp_mean = tmp_mean_dict[border - i]
                    dynamic_price = tmp_mean.loc[index, ]

                    # 寻找close最贴近的均线
                    if min_diff == 0 or abs(row[value_column] - dynamic_price) <= min_diff:
                        min_diff = abs(row[value_column] - dynamic_price)
                        min_diff_ma_num = border - i
                index_df.loc[index, INDEX_MA_BORDER] = min_diff_ma_num
        else:
            index_df.loc[index, INDEX_MA_BORDER] = 5
    index_df = index_df.drop(value_column, axis=1)
    return index_df
