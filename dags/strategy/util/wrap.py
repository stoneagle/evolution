from library import tool
INDEX_DATE = "date"
INDEX_VALUE = "value"
TURN_LIMIT = 4


class Wrap(object):
    # 存储数据集合
    raw_df = None
    wrap_df = None
    df = tool.init_empty_df([INDEX_DATE, INDEX_VALUE])

    def __init__(self, df):
        self.raw_df = df
        self.wrap_df = tool.init_empty_df(self.raw_df.columns)
        return

    def merge_line(self, high_column, low_column, merge_price):
        for index, row in self.raw_df.iterrows():
            if index == 0:
                self.wrap_df = self.wrap_df.append(row)
                continue

            length = len(self.wrap_df)
            pre_row = self.wrap_df.iloc[-1]
            # 递归向前合并?
            if (pre_row[high_column] <= row[high_column] and pre_row[low_column] >= row[low_column]) or \
               (pre_row[high_column] >= row[high_column] and pre_row[low_column] <= row[low_column]):
                # 上升过程中，N-1被N包含，两K线最高点当高点，低点中较高者当低点
                # 下降过程中，N-1包含N，两K线最低点当低点，高点中较低者当高点
                self.wrap_df.loc[length - 1:length, high_column] = row[high_column]
                # TODO 合并后，向前检查直到无法合并
                for i in range(1, (len(self.wrap_df) - 1)):
                    for_pre_row = self.wrap_df.iloc[-i - 1]
                    for_row = self.wrap_df.iloc[-i]
                    if (for_pre_row[high_column] <= for_row[high_column] and for_pre_row[low_column] >= for_row[low_column]) or \
                       (for_pre_row[high_column] >= for_row[high_column] and for_pre_row[low_column] <= for_row[low_column]):
                        self.wrap_df.loc[length - 1 - i:length - i, high_column] = for_row[high_column]
                    else:
                        self.wrap_df = self.wrap_df.head(len(self.wrap_df) - i)
                        break
            else:
                self.wrap_df = self.wrap_df.append(row)
                self.wrap_df = self.wrap_df.reset_index(drop=True)
        if merge_price is True:
            return self._merge_price(high_column, low_column)
        else:
            return self.wrap_df

    def _merge_price(self, high_column, low_column):
        # 将open、close、high、low聚合成一个数值
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
        return self.df
