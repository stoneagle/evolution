from library import tool
INDEX_DATE = "date"


class Wrap(object):
    # 存储数据集合
    raw_df = None
    wrap_df = None

    def __init__(self, df):
        self.df = df
        self.wrap_df = tool.init_empty_df(self.df.columns)
        return

    def merge(self, high_column, low_column):
        for index, row in self.df.iterrows():
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
        return self.wrap_df
