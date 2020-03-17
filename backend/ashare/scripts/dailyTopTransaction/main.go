package main

import (
	"evolution/backend/ashare/models"
	"evolution/backend/common/config"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	// 获取目标数据库配置
	yamlFile, err := ioutil.ReadFile("../../../config/.config.yaml")
	if err != nil {
		panic(err)
	}
	conf := &config.Conf{}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		panic(err)
	}
	dbConfig := conf.Ashare.Database
	source := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Target)
	destEng, err := xorm.NewEngine(dbConfig.Type, source)
	if err != nil {
		panic(err)
	}
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	destEng.TZLocation = location
	destEng.StoreEngine("InnoDB")
	destEng.Charset("utf8")

	// 获取 excel 内容
	inFile := "/Users/wuzhongyang/Downloads/top.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	importDailyTopTransaction(destEng, xlFile)
}

func importDailyTopTransaction(des *xorm.Engine, xlFile *xlsx.File) {
	// 遍历sheet页读取
	for _, sheet := range xlFile.Sheets {
		var currentFundType int
		switch sheet.Name {
		case "国家队":
			currentFundType = models.TopFundTypeCountry
		case "派系":
			currentFundType = models.TopFundTypeSite
		case "营业厅":
			currentFundType = models.TopFundTypeHall
		case "牛散":
			currentFundType = models.TopFundTypePerson
		default:
			panic(fmt.Sprintf("%s 资金所属类别异常 \r\n", sheet.Name))
		}
		// 遍历行读取
		for rowNumber, row := range sheet.Rows {
			if rowNumber == 0 {
				continue
			}
			dateRaw := row.Cells[0].String()
			date, err := time.Parse("2006-01-02", convertToFormatDay(dateRaw))
			var cstSh, _ = time.LoadLocation("Asia/Shanghai")
			date = date.Add(time.Duration(-8) * time.Hour).In(cstSh)
			checkErrorPanic(fmt.Sprintf("%s date transfer", dateRaw), err)
			print(currentFundType)
			// TODO 导入所有股票名称和代码
			// topFundName := row.Cells[1].String()
			// topFund := checkAndSaveTopFund(des, topFundName, currentFundType)
		}
	}
	fmt.Println("\nimport success")
}

func checkErrorPanic(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s:%v\r\n", msg, err.Error()))
	}
}

func checkAndSaveTopFund(des *xorm.Engine, name string, fundType int) models.TopFund {
	topFund := models.TopFund{}
	topFund.Name = name
	topFund.Type = fundType
	has, err := des.Where("name = ?", topFund.Name).And("type = ?", topFund.Type).Get(&topFund)
	checkErrorPanic(fmt.Sprintf("top fund %s has check", topFund.Name), err)
	if !has {
		_, err = des.Insert(&topFund)
		checkErrorPanic(fmt.Sprintf("top fund %s insert", topFund.Name), err)
		fmt.Printf("top fund %s insert success\r\n", topFund.Name)
	}
	// return share
	return topFund
}

func checkAndSaveShare(des *xorm.Engine, name, code string) models.Share {
	share := models.Share{}
	share.Name = name
	share.Code = code
	switch code[0:1] {
	case "6":
		share.IndexType = models.ShareIndexTypeShanghai
	case "0":
		share.IndexType = models.ShareIndexTypeShenzhen
	case "3":
		share.IndexType = models.ShareIndexTypeChuangye
	}
	has, err := des.Where("name = ?", share.Name).And("code = ?", share.Code).Get(&share)
	checkErrorPanic(fmt.Sprintf("share %s has check", share.Name), err)
	if !has {
		_, err = des.Insert(&share)
		checkErrorPanic(fmt.Sprintf("share %s insert", share.Name), err)
		fmt.Printf("share %s insert success\r\n", share.Name)
	}
	return share
}

func convertToFormatDay(excelDaysString string) string {
	// 2006-01-02 距离 1900-01-01的天数
	baseDiffDay := 38719 //在网上工具计算的天数需要加2天，什么原因没弄清楚
	curDiffDay := excelDaysString
	b, _ := strconv.Atoi(curDiffDay)
	// 获取excel的日期距离2006-01-02的天数
	realDiffDay := b - baseDiffDay
	//fmt.Println("realDiffDay:",realDiffDay)
	// 距离2006-01-02 秒数
	realDiffSecond := realDiffDay * 24 * 3600
	//fmt.Println("realDiffSecond:",realDiffSecond)
	// 2006-01-02 15:04:05距离1970-01-01 08:00:00的秒数 网上工具可查出
	baseOriginSecond := 1136185445
	resultTime := time.Unix(int64(baseOriginSecond+realDiffSecond), 0).Format("2006-01-02")
	return resultTime
}
