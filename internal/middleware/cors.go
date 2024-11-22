package middleware

import (
	"crawler/pkg/config"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Cors(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(cfg.Server.AllowedOrigins, ","))
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.Server.AllowedMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.Server.AllowedHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
