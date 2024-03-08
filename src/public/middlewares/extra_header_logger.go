package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/log"
)

func LogExtraHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		agent := c.GetHeader("X-Agent")
		log.Infof("Extra headers: X-Agent: %s", agent)
		c.Next()
	}
}
