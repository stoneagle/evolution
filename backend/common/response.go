package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Redirect(ctx *gin.Context, uri string) {
	ctx.Redirect(http.StatusFound, uri)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": ErrorOk,
		"data": data,
		"desc": "success",
	})
}

func ResponseErrorBusiness(ctx *gin.Context, code ErrorCode, desc string, err error) {
	if err != nil {
		desc += ":" + err.Error()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": struct{}{},
		"desc": desc,
	})
	ctx.Abort()
}

func ResponseErrorServer(ctx *gin.Context, desc string) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": ErrorServer,
		"data": struct{}{},
		"desc": desc,
	})
	ctx.Abort()
}
