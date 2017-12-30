from quota.util import action
import h5py
from library import conf, tool


def test():
    # code_list = ["000725", "600519"]
    code_list = ["000725"]
    ktype_list = ["5"]
    # 估算股票转折比例
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    for code in code_list:
        for ktype in ktype_list:
            code_prefix = code[0:3]
            df = tool.df_from_dataset(f[code_prefix][code], ktype, None)
            df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            trend_df = action.arrange_trend(df.tail(100), 0.1)
            length = len(trend_df)
            for index, row in trend_df.iterrows():
                if index != length - 1:
                    next_row = trend_df.iloc[index + 1]
                    if next_row[action.INDEX_STATUS] == action.STATUS_SHAKE and row[action.INDEX_STATUS] != action.STATUS_SHAKE:
                        start_row = trend_df.iloc[index - row[action.INDEX_TREND_COUNT]]
                        detail = "end_date:%s, trend_count:%d, status:%s, macd_range:%f, macd_diff:%f, diff_per:%d"
                        macd_range = row["macd"] - start_row["macd"]
                        macd_diff = next_row["macd"] - row["macd"]
                        print(detail % (
                            next_row['date'],
                            row[action.INDEX_TREND_COUNT],
                            row[action.INDEX_STATUS],
                            macd_range,
                            macd_diff,
                            round(macd_diff * 100 / macd_range, 0)
                        ))
    f.close()
    return
