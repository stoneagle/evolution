package resp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"evolution/backend/common/logger"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code ErrorCode   `json:"code"`
	Data interface{} `json:"data"`
	Desc string      `json:"desc"`
}

func Redirect(ctx *gin.Context, uri string) {
	ctx.Redirect(http.StatusFound, uri)
}

func Success(ctx *gin.Context, data interface{}) {
	res := Response{
		Code: ErrorOk,
		Data: data,
		Desc: "success",
	}
	ctx.JSON(http.StatusOK, res)
	FormatResponseLog(res)
}

func CustomSuccess(ctx *gin.Context, res interface{}) {
	ctx.JSON(http.StatusOK, res)
	FormatResponseLog(res)
}

func ErrorBusiness(ctx *gin.Context, code ErrorCode, desc string, err error) {
	if err != nil {
		desc += ":" + err.Error()
	}
	res := Response{
		Code: code,
		Data: struct{}{},
		Desc: desc,
	}
	FormatResponseLog(res)
	FormatErrorLog(err)
	ctx.AbortWithStatusJSON(http.StatusOK, res)
}

func FormatResponseLog(response interface{}) {
	logResponse, _ := json.Marshal(response)
	if string(logResponse) != "" {
		logger.Get().Infow("response:【" + string(logResponse) + "】")
	}
}

func FormatErrorLog(err error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(3)
		logger.Get().Infow(fmt.Sprintf("response-error:【%s:%d:%v】", fn, line, err))
	}
}
