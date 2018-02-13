package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/newConn"
	"e.coding.net/anyun-cloud-api-gateway/pool"
	"github.com/gin-gonic/gin"
	"github.com/kabukky/httpscerts"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//NewDefauleAnyunCloudGateway -- 创建默认配置的API网关
func NewDefauleAnyunCloudGateway() (Gateway, error) {
	//TODO: 创建API网关
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	cfg := new(AnyunCloudGatewayConfig)
	cfg.HTTPS.TLSConfig = DefaultTLSConfig()
	cfg.ListenerAddr = ":9000"
	gateway := &AnyunCloudGateway{
		Engine: engine,
		Config: cfg,
	}
	gateway.SetUpMiddlewares()
	return gateway, nil
}

//DefaultTLSConfig -- 获取默认的HTTPS TLS配置
func DefaultTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
}

//Start --  启动API网关
func (_this *AnyunCloudGateway) Start() error {
	//TODO: 启动API网关
	_this.Server = &http.Server{
		Addr:      _this.Config.ListenerAddr,
		Handler:   _this.Engine,
		TLSConfig: &tls.Config{},
	}
	crt := "/home/yw/goworkspace/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.crt"
	key := "/home/yw/goworkspace/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.key"
	if err := httpscerts.Check(crt, key); err != nil {
		log.WithFields(log.Fields{
			"prefix": "server.Start",
		}).Error("server certificate check error")
		return err
	}
	go func() {
		if err := _this.Server.ListenAndServeTLS(crt, key); err != nil {
			log.WithFields(log.Fields{
				"prefix": "server.Start",
			}).Fatal("API gateway start error")
		}
	}()
	b := make(chan bool)
	go func() {
		b <- false
		pool.NewRedisPool("redis.service.dc-anyuncloud.consul:6379")
	}()
	<-b
	c := make(chan bool)
	go func() {
		newConn.NewConn()
		c <- false
	}()
	<-c
	return nil
}

func GetTestAnyunCloudGateway() *AnyunCloudGateway {
	gin.SetMode(gin.ReleaseMode)
	log.SetLevel(log.DebugLevel)
	formatter := new(prefixed.TextFormatter)
	log.SetFormatter(formatter)
	engine := gin.New()
	cfg := new(AnyunCloudGatewayConfig)
	cfg.HTTPS.TLSConfig = DefaultTLSConfig()
	cfg.ListenerAddr = ":9000"
	gateway := &AnyunCloudGateway{
		Engine: engine,
		Config: cfg,
	}
	return gateway
}

//Join -- 维持API网关进程
//* 如果有其它任务后台执行，那么就不需要网关执行join操作
//* 如果不执行join操作，那么网关需要手动清零资源
//* 网关启动后添加系统interrupt信号，接收到interrupt信号后会先关闭服务并且清理资源
func (_this *AnyunCloudGateway) Join() error {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.WithFields(log.Fields{
		"prefix": "server.Start",
	}).Warn("recive system interrupt signal")
	var cancel context.CancelFunc
	_this.ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := _this.Server.Shutdown(_this.ctx); err != nil {
		log.WithFields(log.Fields{
			"prefix": "server.Start",
		}).Fatal("API gateway stop error", err)
	}
	log.WithFields(log.Fields{
		"prefix": "server.Start",
	}).Info("API gateway exist")
	return nil
}

//Stop -- 停止API网关
//* 网关的终止不会清理网关资源
//TODO: 需要添加资源的清理
func (_this *AnyunCloudGateway) Stop() error {
	if err := _this.Server.Shutdown(_this.ctx); err != nil {
		log.WithFields(log.Fields{
			"prefix": "server.Stop",
		}).Error("API gateway stop error", err)
		return err
	}
	return nil
}

//Restart -- 重启API网关
//! 网关被重启的条件
//* 1.配置变更
//* 2.某些条件下的故障恢复
func (_this *AnyunCloudGateway) Restart() error {
	//TODO: 重启API网关
	if err := _this.Stop(); err != nil {
		return err
	}
	return _this.Start()
}

//SetUpMiddlewares -- 获取所有的API网关中间件
//* 1.添加统计中间件
//* 2.添加API的调用中间件
func (_this *AnyunCloudGateway) SetUpMiddlewares() {
	_this.Statistics = new(Statistics)
	_this.Statistics.SetUptime()                                                            //设置网关启动时间
	_this.Statistics.UpdateConfigTime()                                                     //设置最后网关配置时间
	_this.Engine.Use(_this.Statistics.APIStatisticsMiddleware())                            //添加API调用统计中间件
	_this.Engine.Use(_this.APIMiddleware(), _this.BYTEMiddleware(), _this.FileMiddleware()) //添加API中间件
	// _this.Engine.Use(_this.APIMiddleware())
}
