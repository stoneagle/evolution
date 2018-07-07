package resp

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"

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

func ErrorBusiness(ctx *gin.Context, code ErrorCode, desc string, err error) {
	res := Response{
		Code: code,
		Data: struct{}{},
		Desc: desc,
	}
	FormatResponseLog(res)
	FormatErrorLog(err)
	ctx.AbortWithStatusJSON(http.StatusOK, res)
}

func FormatResponseLog(response Response) {
	logResponse, _ := json.Marshal(response)
	if string(logResponse) != "" {
		logger.Get().Infow("response:【" + string(logResponse) + "】")
	}
}

func FormatErrorLog(err error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(2)
		logger.Get().Infow("response-error:【" + fn + ":" + strconv.Itoa(line) + ":" + err.Error() + "】")
	}
}
