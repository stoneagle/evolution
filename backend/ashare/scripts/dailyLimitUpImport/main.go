package main

import (
	"evolution/backend/ashare/models"
	"evolution/backend/common/config"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/tealeg/xlsx"
	yaml "gopkg.in/yaml.v2"
)

var (
	date                       time.Time
	currentConcept             models.Concept
	currentDailyConceptLimitUp models.DailyConceptLimitUp
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
	// inFile := "./import.xlsx"
	inFile := "/Users/wuzhongyang/Downloads/import.xlsx"
	xlFile, err := xlsx.OpenFile(inFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	importDailyLimitUp(destEng, xlFile)
}

func importDailyLimitUp(des *xorm.Engine, xlFile *xlsx.File) {
	// 遍历sheet页读取
	var err error
	for _, sheet := range xlFile.Sheets {
		// 获取第一行第一列的日期
		dateRaw := sheet.Rows[1].Cells[0].String()
		date, err = time.Parse("2006-01-02", convertToFormatDay(dateRaw))
		var cstSh, _ = time.LoadLocation("Asia/Shanghai")
		date = date.Add(time.Duration(-8) * time.Hour).In(cstSh)
		checkErrorPanic(fmt.Sprintf("%s date transfer", dateRaw), err)

		// 遍历行读取
		for rowNumber, row := range sheet.Rows {
			// 跳过首行
			if rowNumber == 0 {
				continue
			}

			conceptRaw := row.Cells[1].String()
			if conceptRaw != "" {
				// 记录上阶段汇总数据
				if currentConcept.Name != "" {
					_, err = des.Insert(&currentDailyConceptLimitUp)
					checkErrorPanic(fmt.Sprintf("daily concept limit up %s insert", currentConcept.Name), err)
					fmt.Printf("daily concept limit up %s insert success\r\n", currentConcept.Name)
				}
				// 判断概念是否存在
				currentConcept = checkAndSaveConcept(des, conceptRaw)
				currentDailyConceptLimitUp = models.DailyConceptLimitUp{}
				currentDailyConceptLimitUp.ConceptId = currentConcept.Id
				currentDailyConceptLimitUp.Date = date
				currentDailyConceptLimitUp.LimitUpCount = 0
				currentDailyConceptLimitUp.LimitCompetitionSum = 0
				currentDailyConceptLimitUp.LimitContinuationSum = 0
				currentDailyConceptLimitUp.LimitContinuationMax = 1
				currentDailyConceptLimitUp.CirculateSum = 0
				currentDailyConceptLimitUp.HighTransactionSum = 0
				currentDailyConceptLimitUp.TransactionSum, err = row.Cells[2].Int()
				checkErrorPanic(fmt.Sprintf("daily concept limit up %s transactionSum", currentConcept.Name), err)
			}

			// 判断个股是否存在
			shareName := row.Cells[3].String()
			shareCode := row.Cells[4].String()
			share := checkAndSaveShare(des, shareName, shareCode)
			shareLimitUp := models.DailyLimitUp{}

			// 插入个股涨停数据
			has, err := des.Where("share_id = ?", share.Id).And("date = ?", date).Get(&shareLimitUp)
			checkErrorPanic(fmt.Sprintf("share limit up %s has check", share.Name), err)
			if !has {
				shareLimitUp.ShareId = share.Id
				shareLimitUp.Date = date
				// 计算涨停时间
				rawLimitTime, err := row.Cells[5].GetTime(false)
				checkErrorPanic(fmt.Sprintf("%s limit time", "test"), err)
				hour, minute, second := rawLimitTime.Clock()
				s := time.Duration(second) * time.Second
				m := time.Duration(minute) * time.Minute
				h := time.Duration(hour) * time.Hour
				shareLimitUp.LimitTime = date.Add(h).Add(m).Add(s)
				if hour == 9 && minute == 25 {
					shareLimitUp.LimitCompetition = 1
				}
				shareLimitUp.LimitContinuation, err = row.Cells[6].Int()
				checkErrorPanic(fmt.Sprintf("%s continuation", share.Name), err)
				shareLimitUp.HighTransaction, err = row.Cells[7].Int()
				checkErrorPanic(fmt.Sprintf("%s high transaction", share.Name), err)
				shareLimitUp.Circulate, err = row.Cells[8].Int()
				checkErrorPanic(fmt.Sprintf("%s circulate", share.Name), err)
				shareLimitUp.MainConceptId = currentConcept.Id
				AdditionalConceptName := row.Cells[9].String()
				if AdditionalConceptName != "" {
					shareLimitUp.AdditionalConceptId = checkAndSaveConcept(des, AdditionalConceptName).Id
				} else {
					shareLimitUp.AdditionalConceptId = 0
				}
				_, err = des.Insert(&shareLimitUp)
				checkErrorPanic(fmt.Sprintf("share limit up %s insert", share.Name), err)
				fmt.Printf("share limit up %s insert success\r\n", share.Name)
				fmt.Printf("%s\n", shareLimitUp.Date)

				currentDailyConceptLimitUp.LimitUpCount += 1
				if shareLimitUp.LimitCompetition > 0 {
					currentDailyConceptLimitUp.LimitCompetitionSum += 1
				}
				if shareLimitUp.LimitContinuation > 1 {
					currentDailyConceptLimitUp.LimitContinuationSum += 1
				}
				if shareLimitUp.LimitContinuation > currentDailyConceptLimitUp.LimitContinuationMax {
					currentDailyConceptLimitUp.LimitContinuationMax = shareLimitUp.LimitContinuation
				}
				currentDailyConceptLimitUp.CirculateSum += shareLimitUp.Circulate
				currentDailyConceptLimitUp.HighTransactionSum += shareLimitUp.HighTransaction
			}
		}

		// 汇总最后的 concept 信息
		_, err = des.Insert(&currentDailyConceptLimitUp)
		checkErrorPanic(fmt.Sprintf("daily concept limit up %s insert", currentConcept.Name), err)
		fmt.Printf("daily concept limit up %s insert success\r\n", currentConcept.Name)
	}
	fmt.Println("\nimport success")
}

func checkErrorPanic(msg string, err error) {
	if err != nil {
		panic(fmt.Sprintf("%s:%v\r\n", msg, err.Error()))
	}
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

// 检查并存储概念
func checkAndSaveConcept(des *xorm.Engine, name string) models.Concept {
	currentConcept := models.Concept{}
	currentConcept.Name = name
	has, err := des.Where("name = ?", currentConcept.Name).Get(&currentConcept)
	checkErrorPanic(fmt.Sprintf("concept %s check", currentConcept.Name), err)
	if !has {
		_, err = des.Insert(&currentConcept)
		checkErrorPanic(fmt.Sprintf("concept %s insert", currentConcept.Name), err)
		fmt.Printf("concept %s insert success\r\n", currentConcept.Name)
	}
	return currentConcept
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
