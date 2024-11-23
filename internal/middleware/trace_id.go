package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	TraceIDHeader = "X-Request-ID"
	TraceIDKey    = "trace_id"
)

func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先从请求头获取
		traceID := c.GetHeader(TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		c.Set(TraceIDKey, traceID)
		c.Header(TraceIDHeader, traceID)
		c.Next()
	}
}
