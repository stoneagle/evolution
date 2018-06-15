package middles

import (
	"quant/backend/common"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RegexType int

const (
	RegexNumber RegexType = iota
)

func Regex(numberName string, rtype RegexType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var regexRule string
		switch rtype {
		case RegexNumber:
			regexRule = `[0-9]`
		default:
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "regex type illegal", nil)
		}

		r, err := regexp.Compile(regexRule)
		if err != nil {
			common.ResponseErrorBusiness(ctx, common.ErrorParams, "regex params error", err)
		}

		number := ctx.Param(numberName)
		if r.MatchString(number) != true {
			common.ResponseErrorBusiness(ctx, common.ErrorParams, numberName+" params error", nil)
		} else {
			result, _ := strconv.Atoi(number)
			ctx.Set(numberName, result)
		}
	}
}
