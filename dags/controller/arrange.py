import h5py
from library import conf, console, tool, tradetime
from strategy import macd


def code_detail(code_list, start_date):
    """
    将code的basic内容，整理至share文件下
    """
    # 获取basic所有日期的detail，并遍历读取详细信息
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    f_share = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_DETAIL
    )
    path = '/' + conf.HDF5_BASIC_DETAIL
    if f.get(path) is None:
        return

    code_basic_dict = dict()
    for date in f[path]:
        if start_date is not None and date < start_date:
            console.write_msg(start_date + "起始日期大于基本数据的最大日期")
            continue
        df = tool.df_from_dataset(f[path], date, None)
        df["code"] = df["code"].str.decode("utf-8")
        df = df.set_index("code")
        for code in df.index:
            if code not in code_list:
                continue

            if code not in code_basic_dict:
                code_basic_dict[code] = tool.init_empty_df(df.columns)
            code_basic_dict[code].loc[date] = df.loc[code, :]

    for code, code_df in code_basic_dict.items():
        code_df.index.name = conf.HDF5_SHARE_DATE_INDEX
        code_df = code_df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])

        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f_share.get(code_group_path) is None:
            console.write_msg(code + "的detail文件不存在")
            continue

        if start_date is None:
            tool.delete_dataset(f_share[code_group_path], conf.HDF5_BASIC_DETAIL)
        tool.merge_df_dataset(f_share[code_group_path], conf.HDF5_BASIC_DETAIL, code_df)
        console.write_exec()
    console.write_blank()
    console.write_tail()
    f_share.close()
    f.close()
    return


def code_classify(code_list, classify_list):
    """
    按照code整理其所属的classify
    """
    f = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_OTHER_CODE_CLASSIFY
    )
    # 获取classify列表
    code_classify_df = tool.init_empty_df(["date", "code", "classify"])
    today_str = tradetime.get_today()
    for ctype in classify_list:
        for classify_name in f[ctype]:
            if f[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                console.write_msg(classify_name + "的code列表不存在")
            classify_df = tool.df_from_dataset(f[ctype][classify_name], conf.HDF5_CLASSIFY_DS_CODE, None)
            for index, row in classify_df.iterrows():
                code = row[0].astype(str)
                if code in code_list:
                    code_dict = dict()
                    code_dict["date"] = today_str
                    code_dict["code"] = code
                    code_dict["classify"] = classify_name
                    code_classify_df = code_classify_df.append(code_dict, ignore_index=True)
    console.write_tail()
    f.close()

    f_other = h5py.File(conf.HDF5_FILE_OTHER, 'a')
    tool.delete_dataset(f_other, conf.HDF5_OTHER_CODE_CLASSIFY)
    tool.merge_df_dataset(f_other, conf.HDF5_OTHER_CODE_CLASSIFY, code_classify_df)
    f_other.close()
    return


def operate_quit(action_type):
    """
    将退市quit内容，转换成标签添加在对应code下
    """
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        action_type,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_QUIT
    )
    path = '/' + conf.HDF5_BASIC_QUIT
    if f.get(path) is None:
        console.write_msg("quit的detail不存在")
        return

    quit_list = [
        conf.HDF5_BASIC_QUIT_TERMINATE,
        conf.HDF5_BASIC_QUIT_SUSPEND,
    ]
    for qtype in quit_list:
        # 如果文件不存在，则退出
        quit_df = tool.df_from_dataset(f[path], qtype, None)
        if quit_df is not None and quit_df.empty is not True:
            quit_df["code"] = quit_df["code"].str.decode("utf-8")
            # 将退市内容，转换成标签添加在对应code下
            tool.op_attr_by_codelist(action_type, quit_df["code"].values, conf.HDF5_BASIC_QUIT, True)
        else:
            console.write_msg("quit的detail数据获取失败")
    console.write_tail()
    f.close()
    return


def operate_st(action_type):
    """
    将st内容，转换成标签添加在对应code下
    """
    f = h5py.File(conf.HDF5_FILE_BASIC, 'a')
    console.write_head(
        action_type,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_BASIC_ST
    )
    path = '/' + conf.HDF5_BASIC_ST
    # 如果文件不存在，则退出
    if f.get(path) is None:
        console.write_msg("st的detail文件不存在")
        return

    st_df = tool.df_from_dataset(f[path], conf.HDF5_BASIC_ST, None)
    if st_df is not None and st_df.empty is not True:
        st_df["code"] = st_df["code"].str.decode("utf-8")
        # 将st内容，转换成标签添加在对应code下
        tool.op_attr_by_codelist(action_type, st_df["code"].values, conf.HDF5_BASIC_ST, True)
    else:
        console.write_msg("st的detail数据获取失败")
    console.write_tail()
    f.close()
    return


def all_classify_detail(classify_list, omit_list, start_date):
    """
    遍历所有分类，聚合所有code获取分类均值
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    f_classify = h5py.File(conf.HDF5_FILE_CLASSIFY, 'a')
    # 获取classify列表
    for ctype in classify_list:
        for classify_name in f_classify[ctype]:
            console.write_head(
                conf.HDF5_OPERATE_ARRANGE,
                conf.HDF5_RESOURCE_TUSHARE,
                classify_name
            )

            if f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE) is None:
                console.write_msg(classify_name + "的detail文件不存在")
                continue

            for ktype in conf.HDF5_SHARE_KTYPE:
                mean_df = one_classify_detail(f, f_classify[ctype][classify_name].get(conf.HDF5_CLASSIFY_DS_CODE), omit_list, ktype, start_date)
                ds_name = ktype
                # 如果start_date为空，则重置该数据
                if start_date is None:
                    tool.delete_dataset(f_classify[ctype][classify_name], ds_name)

                if mean_df is not None:
                    tool.merge_df_dataset(f_classify[ctype][classify_name], ds_name, mean_df)
            console.write_tail()
    f_classify.close()
    f.close()
    return


def one_classify_detail(f, code_list, omit_list, ktype, start_date):
    """
    根据单个分类，聚合所有code获取分类均值
    """
    # 初始化一个全日期的空DataFrame，并初始化一列作为统计个数
    init_df = tool.init_empty_df(conf.HDF5_SHARE_COLUMN)

    # 按照列表顺序，获取code并逐一添加至初始化DF，并递增该日期的个数
    for row in code_list:
        code = row[0].astype(str)
        code_prefix = code[0:3]
        # 判断是否跳过创业板
        if code_prefix in omit_list:
            continue

        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            console.write_msg(code + "的detail文件不存在")
            continue
        else:
            # 忽略停牌、退市、无法获取的情况
            if f[code_group_path].attrs.get(conf.HDF5_BASIC_QUIT) is not None:
                continue

            if f[code_group_path].attrs.get(conf.HDF5_BASIC_ST) is not None:
                continue

        if f[code_group_path].get(ktype) is None:
            console.write_msg(code + "的" + ktype + "文件不存在")
            continue
        add_df = tool.df_from_dataset(f[code_group_path], ktype, None)
        add_df[conf.HDF5_SHARE_DATE_INDEX] = add_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
        add_df["num"] = 1
        add_df = add_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
        init_df = init_df.add(add_df, fill_value=0)
    # 总数除以数量，得到平均值
    if len(init_df) > 0:
        init_df = init_df.div(init_df.num, axis=0)
        init_df = init_df.drop("num", axis=1)
        if start_date is not None:
            init_df = init_df.ix[start_date:]
        init_df = init_df.reset_index().sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
        init_df["volume"] = init_df["volume"].astype('float64')
        return init_df
    else:
        return None


def xsg():
    """
    聚合xsg数据
    """
    # 初始化文件
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    f_share = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_XSG
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_XSG
    xsg_sum_dict = dict()
    if f.get(path) is not None:
        for month in f[path]:
            df = tool.df_from_dataset(f[path], month, None)
            df["code"] = df["code"].str.decode("utf-8")
            df["count"] = df["count"].str.decode("utf-8")
            df[conf.HDF5_SHARE_DATE_INDEX] = df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
            for index, row in df.iterrows():
                code = row["code"]
                xsg_date_str = row[conf.HDF5_SHARE_DATE_INDEX]
                code_prefix = code[0:3]
                code_group_path = '/' + code_prefix + '/' + code
                if f_share.get(code_group_path) is None:
                    continue
                # 获取限售股解禁前一天的价格
                share_df = tool.df_from_dataset(f_share[code_group_path], "D", None)
                share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
                share_df = share_df.set_index(conf.HDF5_SHARE_DATE_INDEX)
                share_df = share_df[: xsg_date_str]
                if len(share_df) == 0:
                    continue
                close = share_df.tail(1)["close"]
                # 万股*元，统一单位为亿
                code_sum = close.values * float(row["count"]) * 10000
                sum_price = round(code_sum[0] / 10000 / 10000, 2)
                # trade_date = tradetime.get_week_of_date(xsg_date_str, "D")
                trade_date = xsg_date_str
                if trade_date in xsg_sum_dict:
                    xsg_sum_dict[trade_date] += sum_price
                else:
                    xsg_sum_dict[trade_date] = sum_price
        sum_df = tool.init_df(list(xsg_sum_dict.items()), [conf.HDF5_SHARE_DATE_INDEX, "sum"])
        if len(sum_df) > 0:
            sum_df = sum_df.sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
            tool.create_df_dataset(f, conf.HDF5_FUNDAMENTAL_XSG_DETAIL, sum_df)
    console.write_tail()
    f_share.close()
    f.close()
    return


def ipo():
    """
    聚合ipo上市数据
    """
    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    path = '/' + conf.HDF5_FUNDAMENTAL_IPO
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_FUNDAMENTAL_IPO
    )
    path = '/' + conf.HDF5_FUNDAMENTAL_IPO
    ipo_sum_dict = dict()
    if f.get(path) is not None:
        df = tool.df_from_dataset(f[path], conf.HDF5_FUNDAMENTAL_IPO, None)
        df["issue_date"] = df["issue_date"].str.decode("utf-8")
        df["ipo_date"] = df["ipo_date"].str.decode("utf-8")
        for index, row in df.iterrows():
            trade_date = row["ipo_date"]
            # 统一单位为亿元
            sum_price = round(row["funds"], 2)
            if trade_date in ipo_sum_dict:
                ipo_sum_dict[trade_date] += sum_price
            else:
                ipo_sum_dict[trade_date] = sum_price
        sum_df = tool.init_df(list(ipo_sum_dict.items()), [conf.HDF5_SHARE_DATE_INDEX, "sum"])
        if len(sum_df) > 0:
            sum_df = sum_df.sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
            tool.create_df_dataset(f[path], conf.HDF5_FUNDAMENTAL_IPO_DETAIL, sum_df)
    console.write_tail()
    f.close()
    return


def margins(mtype):
    """
    聚合融资融券数据
    """
    if mtype == "sh":
        mtype_index = conf.HDF5_FUNDAMENTAL_SH_MARGINS
        mtype_index_detail = conf.HDF5_FUNDAMENTAL_SH_MARGINS_DETAIL
    elif mtype == "sz":
        mtype_index = conf.HDF5_FUNDAMENTAL_SZ_MARGINS
        mtype_index_detail = conf.HDF5_FUNDAMENTAL_SZ_MARGINS_DETAIL
    else:
        print("mtype " + mtype + " error\r\n")
        return

    f = h5py.File(conf.HDF5_FILE_FUNDAMENTAL, 'a')
    path = '/' + mtype_index
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        mtype_index
    )
    console.write_tail()
    margin_sum_dict = dict()
    if f.get(path) is not None:
        df = tool.df_from_dataset(f[path], mtype_index, None)
        df["opDate"] = df["opDate"].str.decode("utf-8")
        for index, row in df.iterrows():
            # trade_date = tradetime.get_week_of_date(row["opDate"], "D")
            trade_date = row["opDate"]
            # 统一单位为亿元
            sum_price = round((row["rzmre"] - row["rqmcl"]) / 10000 / 10000, 2)
            if trade_date in margin_sum_dict:
                margin_sum_dict[trade_date] += sum_price
            else:
                margin_sum_dict[trade_date] = sum_price
        sum_df = tool.init_df(list(margin_sum_dict.items()), [conf.HDF5_SHARE_DATE_INDEX, "sum"])
        if len(sum_df) > 0:
            sum_df = sum_df.sort_values(by=[conf.HDF5_SHARE_DATE_INDEX])
            tool.create_df_dataset(f[mtype_index], mtype_index_detail, sum_df)
    f.close()
    return


def all_macd_trend(code_list, start_date):
    """
    整理所有股票的macd趋势数据
    """
    f = h5py.File(conf.HDF5_FILE_SHARE, 'a')
    console.write_head(
        conf.HDF5_OPERATE_ARRANGE,
        conf.HDF5_RESOURCE_TUSHARE,
        conf.HDF5_INDEX_MACD_TREND
    )
    for code in code_list:
        code_prefix = code[0:3]
        code_group_path = '/' + code_prefix + '/' + code
        if f.get(code_group_path) is None:
            continue
        # 忽略停牌、退市、无法获取的情况
        if f[code_prefix][code].attrs.get(conf.HDF5_BASIC_QUIT) is not None or f[code_prefix][code].attrs.get(conf.HDF5_BASIC_ST) is not None:
            continue
        for ktype in conf.HDF5_SHARE_KTYPE:
            trend_df = code_macd_trend(f[code_prefix][code], ktype)
            if trend_df is not None:
                ds_name = conf.HDF5_INDEX_MACD_TREND + "_" + ktype
                if f[code_prefix][code].get(ds_name) is not None:
                    tool.delete_dataset(f[code_prefix][code], ds_name)
                tool.create_df_dataset(f[code_prefix][code], ds_name, trend_df)
    console.write_tail()
    f.close()
    return


def code_macd_trend(f, ktype):
    index_ds_name = conf.HDF5_INDEX_DETAIL + "_" + ktype
    if f.get(ktype) is None:
        return None

    if f.get(index_ds_name) is None:
        return None

    share_df = tool.df_from_dataset(f, ktype, None)
    index_df = tool.df_from_dataset(f, index_ds_name, None)
    share_df = share_df.merge(index_df, left_on=conf.HDF5_SHARE_DATE_INDEX, right_on=conf.HDF5_SHARE_DATE_INDEX, how='left')
    share_df[conf.HDF5_SHARE_DATE_INDEX] = share_df[conf.HDF5_SHARE_DATE_INDEX].str.decode("utf-8")
    trend_df = macd.trend(share_df.tail(60))
    return trend_df
