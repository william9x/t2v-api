package middlewares

import (
	"github.com/coreos/go-semver/semver"
	"strings"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	appVer, err := semver.NewVersion("1.1.5")
	agent := "ios/16.7.2"
	if agent != "" && strings.ToLower(agent[:3]) == "ios" {
		if err == nil && appVer.Compare(*IOSV2Ver) != -1 {
			agent = "iosV2"
		} else {
			agent = "ios"
		}
	} else {
		if err == nil && appVer.Compare(*AndroidV2Ver) != -1 {
			agent = "androidV2"
		} else {
			agent = "android"
		}
	}
	println(agent)
}
