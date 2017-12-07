from library import tool
from strategy.util import action
INDEX_DATE = "date"
INDEX_VALUE = "value"


class Wrap(object):
    # 存储数据集合
    raw_df = None
    wrap_df = None
    df = tool.init_empty_df([INDEX_DATE, INDEX_VALUE])

    def __init__(self, df):
        self.raw_df = df
        self.wrap_df = tool.init_empty_df(self.raw_df.columns)
        return

    def merge(self, high_column, low_column):
        for index, row in self.raw_df.iterrows():
            if index == 0:
                self.wrap_df = self.wrap_df.append(row)
                continue

            length = len(self.wrap_df)
            pre_row = self.wrap_df.iloc[-1]
            if pre_row[high_column] <= row[high_column] and pre_row[low_column] >= row[low_column]:
                # 上升过程中，N-1被N包含，两K线最高点当高点，低点中较高者当低点
                self.wrap_df.loc[length - 1:length, high_column] = row[high_column]
            elif pre_row[high_column] >= row[high_column] and pre_row[low_column] <= row[low_column]:
                # 下降过程中，N-1包含N，两K线最低点当低点，高点中较低者当高点
                self.wrap_df.loc[length - 1:length, high_column] = row[high_column]
            else:
                self.wrap_df = self.wrap_df.append(row)

        for index, row in self.wrap_df.iterrows():
            one = dict()
            one[INDEX_DATE] = row[INDEX_DATE]
            if index == 0:
                one[INDEX_VALUE] = row[high_column]
                continue
            pre_row = self.wrap_df.iloc[-1]
            if pre_row[low_column] < row[low_column] and pre_row[high_column] < row[high_column]:
                # 如果上涨，取high作为value
                one[INDEX_VALUE] = row[high_column]
            else:
                # 如果下跌，取low作为value
                one[INDEX_VALUE] = row[low_column]
            self.df = self.df.append(one, ignore_index=True)
        return action.Action().run(self.df, INDEX_DATE, INDEX_VALUE)
