package main

import (
	"github.com/coreos/pkg/capnslog"
	"immortality-demo/config"
	"immortality-demo/pkg/app"
	"immortality-demo/pkg/gredis"
	"immortality-demo/service/workerpool"
)

var ulog = capnslog.NewPackageLogger("immortality", "main")

func init() {
	config.LoadConfig()
	config.SetupLogging()
	gredis.Setup()
}

// @title immortality API
// @version 1.0
// @description This is API of immortality.
// @license.name Apache 2.0
func main() {
	agent := app.NewImmortalityAgent()
	agent.Start()
	workerpool.NewWorkPool(10)
}
