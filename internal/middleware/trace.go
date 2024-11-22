package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 生成请求追踪ID
		traceID := uuid.New().String()
		c.Set("trace_id", traceID)
		// 添加到响应头
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}
