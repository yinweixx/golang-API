package server

import (
	"e.coding.net/anyun-cloud-api-gateway/discovery"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//APIMiddleware -- 服务查询调用中间件
//! 网关只匹配URI以"/api/"开头的调用为应用的API调用
func (_this *AnyunCloudGateway) APIMiddleware() gin.HandlerFunc {
	// BuiltInAPIRegexp := "^/api/"
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		log.WithFields(log.Fields{
			"prefix":         "server.APIMiddleware",
			"request-path":   url,
			"request-method": c.Request.Method,
			"client-ip":      c.ClientIP(),
			"params":         c.Request.URL.Query(),
			"content-type":   c.ContentType(),
			"content-length": c.Request.ContentLength,
			"accept-type":    c.Request.Header.Get("Accept"),
		}).Info("request info")
		// form, _ := c.MultipartForm()
		// for _, headers := range form.File {
		// 	for _, head := range headers {
		// 		fmt.Println(head)
		// 	}
		// }
		val, err := doCheck(&discovery.RequestDetail{
			URL:         url,
			Method:      c.Request.Method,
			Params:      c.Request.URL.Query(),
			ContentType: c.ContentType(),
			AcceptType:  c.Request.Header.Get("Accept"),
		})

		if err != nil {
			log.WithFields(log.Fields{
				"prefix": "server.APIMiddleware.docheck",
			}).Error("error")
			c.JSON(SYSTEMERROR, gin.H{
				"message": "error",
			})
		} else {
			c.JSON(OK, gin.H{
				"message": val,
			})
		}
		// if matched, _ := regexp.MatchString(BuiltInAPIRegexp, url); matched {
		// 	log.WithFields(log.Fields{
		// 		"prefix":         "server.APIMiddleware",
		// 		"request-path":   url,
		// 		"request-method": c.Request.Method,
		// 		"client-ip":      c.ClientIP(),
		// 		"params":         c.Request.URL.Query(),
		// 		"content-type":   c.ContentType(),
		// 		"content-length": c.Request.ContentLength,
		// 		"accept-type":    c.Request.Header.Get("Accept"),
		// 	}).Info("request info")
		// 	// form, _ := c.MultipartForm()
		// 	// for _, headers := range form.File {
		// 	// 	for _, head := range headers {
		// 	// 		fmt.Println(head)
		// 	// 	}
		// 	// }
		// 	val, err := doCheck(&discovery.RequestDetail{
		// 		URL:         url,
		// 		Method:      c.Request.Method,
		// 		Params:      c.Request.URL.Query(),
		// 		ContentType: c.ContentType(),
		// 		AcceptType:  c.Request.Header.Get("Accept"),
		// 	})

		// 	if err != nil {
		// 		log.WithFields(log.Fields{
		// 			"prefix": "server.APIMiddleware.docheck",
		// 		}).Error("error")
		// 		c.JSON(SYSTEMERROR, gin.H{
		// 			"message": "error",
		// 		})
		// 	} else {
		// 		c.JSON(OK, gin.H{
		// 			"message": val,
		// 		})
		// 	}
		// 	//TODO: 服务查询，调用
		// 	//* 1.根据服务资源ID查询服务缓存，确定是否有服务实例在运行
		// 	//* 2.查询服务DNS SRV，如果查找到，将SRV的A记录存如缓存
		// 	//* 3.通知容器管理平台，等待服务容器实例创建完毕并且加入DNS SRV记录
		// 	//* 4.回步骤2
		// }
	}
}
