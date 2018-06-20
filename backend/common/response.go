package common

import (
	"net/http"
	"runtime"
	"strconv"

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

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: ErrorOk,
		Data: data,
		Desc: "success",
	})
}

func ResponseErrorBusiness(ctx *gin.Context, code ErrorCode, desc string, err error) {
	if err != nil {
		desc += ":" + err.Error()
		_, fn, line, _ := runtime.Caller(1)
		GetLogger().Infow("response-error-response:【" + fn + ":" + strconv.Itoa(line) + ":" + desc + "】")
	} else {
		GetLogger().Infow("response-error-response:【" + desc + "】")
	}
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Data: struct{}{},
		Desc: desc,
	})
	ctx.Abort()
}

func ResponseErrorServer(ctx *gin.Context, desc string) {
	ctx.JSON(http.StatusOK, Response{
		Code: ErrorServer,
		Data: struct{}{},
		Desc: desc,
	})
	ctx.Abort()
}
