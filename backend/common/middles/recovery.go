package middles

import (
	"errors"
	"evolution/backend/common/resp"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func Recovery(f func(c *gin.Context, err interface{})) gin.HandlerFunc {
	return RecoveryWithWriter(f, gin.DefaultErrorWriter)
}

func RecoveryWithWriter(f func(c *gin.Context, err interface{}), out io.Writer) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				f(c, err)
			}
		}()
		c.Next()
	}
}

func RecoveryHandler(ctx *gin.Context, err interface{}) {
	exception := errors.New(fmt.Sprintf("%v", err))
	resp.ExceptionServer(ctx, exception)
	return
}
