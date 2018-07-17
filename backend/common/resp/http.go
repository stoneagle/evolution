package resp

import (
	"evolution/backend/common/logger"
	"fmt"
	"net/http"

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
	log := logger.Get()
	log.Caller = 3
	log.Log(logger.InfoLevel, res, nil)
}

func CustomSuccess(ctx *gin.Context, res interface{}) {
	ctx.JSON(http.StatusOK, res)
	log := logger.Get()
	log.Caller = 3
	log.Log(logger.InfoLevel, res, nil)
}

func ErrorBusiness(ctx *gin.Context, code ErrorCode, desc string, err error) {
	message, ok := ErrorMessages[code]
	if !ok {
		logger.Get().Log(logger.WarnLevel, "code message not exist", nil)
		message = "unknown code message"
	}
	desc = fmt.Sprintf("%s:%s", message, desc)
	if err != nil {
		desc += "[" + err.Error() + "]"
	}
	res := Response{
		Code: code,
		Data: struct{}{},
		Desc: desc,
	}
	ctx.AbortWithStatusJSON(http.StatusOK, res)
	log := logger.Get()
	log.Caller = 3
	log.Log(logger.InfoLevel, res, err)
}

func ExceptionServer(ctx *gin.Context, err error) {
	code := ErrorServer
	message, ok := ErrorMessages[code]
	if !ok {
		logger.Get().Log(logger.WarnLevel, "code message not exist", nil)
		message = "unknown code message"
	}
	res := Response{
		Code: code,
		Data: struct{}{},
		Desc: message,
	}
	ctx.AbortWithStatusJSON(http.StatusOK, res)
	log := logger.Get()
	log.Caller = 10
	log.Log(logger.ErrorLevel, res, err)
}
