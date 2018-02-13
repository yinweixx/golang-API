package config

import (
	"e.coding.net/anyun-cloud-common/app"
	"github.com/urfave/cli"
)

//GatewayArgs -- API网关启动配置映射结构体
//
//* 日志的级别的配置
//* 分布式DNS服务器的配置
type GatewayArgs struct {
	LoggerLevel string
	DNS         string `yaml:"dns"`
}

var args GatewayArgs

//GatewayInitFlags -- 网关启动参数初始化映射
func GatewayInitFlags() *app.ApplicationArgs {
	return &app.ApplicationArgs{
		Name:    "anyun-cloud-api-gateway",
		Usage:   "Anyun Cloud Distributed API Gateway",
		Version: "1.0.0",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "logger_level",
				Usage:       "logger level",
				Value:       "debug",
				Destination: &args.LoggerLevel,
			},
			cli.StringFlag{
				Name:        "discovery_dns",
				Usage:       "discovery dns list (eg. 192.168.1.1,192.168.1.2)",
				Value:       "localdns",
				Destination: &args.DNS,
			},
		},
	}
}
