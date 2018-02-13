package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/discovery"
	"e.coding.net/anyun-cloud-api-gateway/event"
	"e.coding.net/anyun-cloud-api-gateway/newConn"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//Statistics -- 网关统计信息
//* 运行时间统计
//* 访问统计信息
//* API部署情况统计
type Statistics struct {
	Timer struct {
		Uptime         int64 //当前API网关节点的启动时间
		LastUpdateTime int64 //当前API网关最后一次更改配置的重启时间
	} //网关时间指标
	Counter struct {
		AccessTotalCount int            `json:"total"`   //系统启动后的连接总量
		ProcessCount     int            `json:"process"` //当前处理器个数
		APIAccessCount   map[string]int `json:"access"`  //根据URL分类的API访问总量
	}
	API struct {
		DeploySuccessCount        int            //已成功部署的API个数
		DeployErrorCount          int            //未成功部署的API个数
		LastAPIDeployTime         int64          //API最后一次部署的UNIX时间
		TopTimeAPI                map[string]int //最长调用时间的API信息
		AuthenticationFailedCount map[string]int //认证失败的API调用信息
	}
}

//SetUptime -- 设置API网关进程第一次启动的时间
func (_this *Statistics) SetUptime() {
	if _this.Timer.Uptime == 0 { //启动时间只能设置一次
		return
	}
	now := time.Now().Unix()
	_this.Timer.Uptime = now
	_this.Timer.LastUpdateTime = now
}

//UpdateConfigTime -- 更新最后一次API网关配置修改的时间
func (_this *Statistics) UpdateConfigTime() {
	now := time.Now().Unix()
	_this.Timer.LastUpdateTime = now
}

//APIStatisticsMiddleware -- 统计中间件
func (_this *Statistics) APIStatisticsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO: 需要检测API是否在缓存里面存在
		//* 不存在的API不需要作为API调用统计
		//* 不在的API也不算作总的API调用个数里面
		//* 不存在的API记录在另外的统计
		_this.Counter.AccessTotalCount++
		_this.Counter.ProcessCount++
		if _this.Counter.APIAccessCount == nil {
			_this.Counter.APIAccessCount = map[string]int{}
		}
		if _, exist := _this.Counter.APIAccessCount[c.Request.URL.Path]; !exist {
			_this.Counter.APIAccessCount[c.Request.URL.Path] = 1
		} else {
			_this.Counter.APIAccessCount[c.Request.URL.Path]++
		}
		// log.WithFields(log.Fields{
		// 	"prefix":  "server.APIStatisticsMiddleware",
		// 	"api-url": c.Request.URL,
		// }).Debug("waiting for other resource")
		c.Next() //! 等待API处理器完成
		_this.Counter.ProcessCount--
		log.WithFields(log.Fields{
			"prefix":            "server.APIStatisticsMiddleware",
			"api-total-count":   _this.Counter.AccessTotalCount,
			"api-process-count": _this.Counter.ProcessCount,
			"api-url":           c.Request.URL,
			"api-current-count": _this.Counter.APIAccessCount[c.Request.URL.Path],
		}).Debug("api counter")
	}
}

//Initialization -- 初始化内置管理API
//* API访问统计
func (_this *AnyunCloudGateway) Initialization() {
	gm1 := _this.Engine.Group("/gateway/v1/")
	{
		gm1.GET(_this.statistics())
	}
}

//statistics -- API访问统计服务
//* 返回API调用次数统计信息
func (_this *AnyunCloudGateway) statistics() (string, func(c *gin.Context)) {
	return "/statistics", func(c *gin.Context) {
		c.JSON(http.StatusOK, _this.Statistics.Counter)
	}
}

func doCheck(requestDetail *discovery.RequestDetail, info ...*event.ConnInfo) (string, error) {
	var aInfo *event.ConnInfo = nil
	for _, fo := range info {
		aInfo = fo
	}
	queryClient, err := newConn.GetConn(context.Background(), aInfo)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "queryClient error",
		}).Error("API gateway doCheck error", err)
		return "", nil
	}
	// 查询API对应的ID
	_, mess, err := queryClient.CheckQueryURL(requestDetail)
	if err != nil {
		return "", err
	}
	var p discovery.ReponseDetail
	json.Unmarshal([]byte(mess), &p)
	log.Info(p.URL + "的请求返回数据为" + p.Uuid)
	if len(p.Uuid) == 0 {
		return "无服务", nil
	}
	//通过ID检查nats管道
	val, err := queryClient.CheckFromEtcd(p.Uuid)
	if val == "" {
		return "服务未绑定", errors.New("服务未绑定")
	}
	result, err := queryClient.ServiceNatsClient.Request(p.Uuid, discovery.BuildMessage(requestDetail), 1000*time.Millisecond)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error")
		return "", err
	}
	return string(result.Data), nil
}

func ginReturn(requestDetail *discovery.RequestDetail) (string, error) {
	val, err := doCheck(requestDetail)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "server.APIMiddleware.docheck",
		}).Error("error")
		return string(SYSTEMERROR), errors.New("system error")
	}
	return val, nil
}
