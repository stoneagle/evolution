package middles

import (
	"evolution/backend/common/resp"
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
	resp.ExceptionServer(ctx, err.(error))
}
