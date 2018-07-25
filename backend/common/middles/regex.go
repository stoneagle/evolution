package middles

import (
	"evolution/backend/common/resp"
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
			regexRule = `[0-9]*`
		default:
			resp.ErrorBusiness(ctx, resp.ErrorParams, "regex type illegal", nil)
		}

		r, err := regexp.Compile(regexRule)
		if err != nil {
			resp.ErrorBusiness(ctx, resp.ErrorParams, "regex params error", err)
		}

		number := ctx.Param(numberName)
		if r.MatchString(number) != true {
			resp.ErrorBusiness(ctx, resp.ErrorParams, numberName+" params error", nil)
		} else {
			result, _ := strconv.Atoi(number)
			ctx.Set(numberName, result)
		}
	}
}
