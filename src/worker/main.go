package main

import (
	"github.com/Braly-Ltd/t2v-api-worker/bootstrap"
	"go.uber.org/fx"
)

// @title Voice Changer API Public
// @version 1.0.0
func main() {
	fx.New(bootstrap.All()).Run()
}
