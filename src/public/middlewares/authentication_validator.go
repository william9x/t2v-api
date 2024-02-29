package middlewares

import (
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/response"
)

func Authenticate(port ports.AuthenticationPort) gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		shouldValidate := c.GetHeader("Disable-Authorization")
		if shouldValidate == "true" {
			c.Next()
			return
		}

		token := c.GetHeader("Authorization")
		if token == "" {
			response.WriteError(c.Writer, exception.New(401, "missing access token"))
			return
		}
		_, err := port.Authenticate(c, token)
		if err != nil {
			response.WriteError(c.Writer, exception.New(403, "invalid access token"))
			return
		}
		c.Next()
	}
}
