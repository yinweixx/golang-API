package server

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/discovery"
	"e.coding.net/anyun-cloud-api-gateway/newConn"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

//根据一些手动默认配置测试API网关的启动
func TestAPIServerStartup(t *testing.T) {
	engine := gin.New()
	cfg := new(AnyunCloudGatewayConfig)
	cfg.HTTPS.TLSConfig = DefaultTLSConfig()
	cfg.ListenerAddr = ":9000"
	gateway := &AnyunCloudGateway{
		Engine: engine,
		Config: cfg,
	}
	gateway.Start()
	gateway.Join()
}

func getTestAnyunCloudGateway() *AnyunCloudGateway {
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

//API网关中间件测试
func TestMiddleware(T *testing.T) {
	gateway := getTestAnyunCloudGateway()
	gateway.SetUpMiddlewares()
	gateway.Start()
	vt := gateway.Engine.Group("/test")
	{
		vt.GET("/api1", func(c *gin.Context) {
			c.String(http.StatusOK, "test1")
		})
	}
	gateway.Join()
}

//测试内置的API
func TestBuildInAPI(t *testing.T) {
	gateway := getTestAnyunCloudGateway()
	gateway.SetUpMiddlewares()
	gateway.Start()
	gateway.Initialization()
	gateway.Join()
}

func TestMySQLClient(t *testing.T) {
	client, err := discovery.NewQueryClient(context.Background(), newConn.GetConnInfo())
	defer client.MySQLClient.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	err1 := client.MySQLClient.Ping()
	if err1 != nil {
		fmt.Println(err1.Error())
	}
}

func TestNatsClient(t *testing.T) {
	client, _ := discovery.NewQueryClient(context.Background(), newConn.GetConnInfo())
	mess, _ := client.ManageNatsClient.Request("API_CONTROLLER_CHANNEL", []byte("test"), 1000*time.Millisecond)
	fmt.Println(string(mess.Data))
}

func TestHTTPS(t *testing.T) {
	crt := "/Users/twitchgg/Develop/Projects/goproject/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.crt"
	key := "/Users/twitchgg/Develop/Projects/goproject/src/e.coding.net/anyun-cloud-api-gateway/ssl/server.key"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there!")
	})
	http.ListenAndServeTLS(":8081", crt, key, nil)
}
