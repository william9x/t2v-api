package middlewares

import (
	"github.com/Braly-Ltd/t2v-api-core/ports"
	"github.com/Braly-Ltd/t2v-api-public/properties"
	"github.com/coreos/go-semver/semver"
	"github.com/gin-gonic/gin"
	"github.com/golibs-starter/golib/exception"
	"github.com/golibs-starter/golib/web/context"
	"github.com/golibs-starter/golib/web/response"
	"strings"
)

var AppV2Ver = semver.New("1.1.2")

func Authenticate(port ports.AuthenticationPort, props *properties.MiddlewaresProperties) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !props.AuthenticationEnable {
			c.Next()
			return
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, response.Error(exception.New(401, "missing access token")))
			return
		}

		agent := c.GetHeader("X-Agent")
		if agent != "" && strings.ToLower(agent[:3]) == "ios" {
			agent = "ios"
		} else {
			appVer, err := semver.NewVersion(c.GetHeader("X-App-Version"))
			if err == nil && appVer.Compare(*AppV2Ver) != -1 {
				agent = "androidV2"
			} else {
				agent = "android"
			}
		}
		tokenData, err := port.Authenticate(c, agent, token)
		if err != nil {
			c.AbortWithStatusJSON(401, response.Error(exception.New(403, "invalid access token")))
			return
		}

		c.Request.Header.Set("X-User-ID", tokenData.UserID)
		context.GetOrCreateRequestAttributes(c.Request).SecurityAttributes.UserId = tokenData.UserID

		c.Next()
	}
}
