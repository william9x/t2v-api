package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

func LogExtraHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		sessID := c.GetHeader("X-Session-ID")
		agent := c.GetHeader("X-Agent")
		appVer := c.GetHeader("X-App-Version")

		log.Infof("Extra headers: X-Request-ID: %s X-Session-ID: %s X-Agent: %s X-App-Version: %s", reqID, sessID, agent, appVer)
		c.Next()
	}
}
