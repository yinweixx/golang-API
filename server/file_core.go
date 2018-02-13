package server

import (
	"regexp"

	"e.coding.net/anyun-cloud-api-gateway/discovery"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//FileMiddleware -- 网关只匹配URI以"/FILE/"开头的调用为应用的API调用
func (_this *AnyunCloudGateway) FileMiddleware() gin.HandlerFunc {
	BuiltInAPIRegexp := "^/file/"
	return func(c *gin.Context) {
		url := c.Request.URL.Path
		if matched, _ := regexp.MatchString(BuiltInAPIRegexp, url); matched {
			log.WithFields(log.Fields{
				"prefix":         "server.FILEMiddleware",
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
		}
	}
}
