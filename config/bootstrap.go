package config

import (
	"e.coding.net/anyun-cloud-api-gateway/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//Bootstrap -- 网关启动入口函数
func Bootstrap(ctx *cli.Context) error {
	logLevel, err := log.ParseLevel(args.LoggerLevel)
	if err != nil {
		log.WithField("prefix", "application.AgentBootstrap").Fatalf("unsuppored logger level: %s", args.LoggerLevel)
	}
	log.SetLevel(logLevel)
	formatter := new(prefixed.TextFormatter)
	log.SetFormatter(formatter)
	gateway := server.GetTestAnyunCloudGateway()
	gateway.SetUpMiddlewares()
	gateway.Start()
	gateway.Initialization()
	gateway.Join()
	return nil
}
