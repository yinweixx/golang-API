package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"e.coding.net/anyun-cloud-api-gateway/common"
	"e.coding.net/anyun-cloud-api-gateway/event"
	"e.coding.net/anyun-cloud-api-gateway/pool"
	"github.com/coreos/etcd/clientv3"
	_ "github.com/go-sql-driver/mysql"
	nats "github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
)

const (
	//ETCD_KEY_PREFIX_FQDN_SERVICE ETCD服务FQDN配置KEY前缀
	ETCD_KEY_PREFIX_FQDN_SERVICE = "/config/global/fqdn/service/"

	FIRST_API          = "first_api"
	SECOND_DETAIL      = "second_detail"
	THIRD_VERSION      = "third_version"
	FORTH_METHOD       = "forth_method"
	FIVE_PARAMS        = "five_params"
	SIXTH_BODY         = "sixth_body"
	SEVENTH_PRODUCT    = "seventh_product"
	REDIS_TIME_TO_LIVE = 600
	MYSQLMAXCONNS      = 10
)

// NewQueryClient -- 创建云平台资源查询客户端
func NewQueryClient(ctx context.Context, connInfo *event.ConnInfo) (*QueryClient, error) {
	//TODO: 创建云平台资源查询客户端
	//* 1. 创建etcd客户端连接
	//* 2. 创建dns查询客户端
	//* 3.创建redis客户端连接
	etcdClient := make(chan *clientv3.Client)
	// mysqlClient := make(chan *sql.DB)
	managerNatsClient := make(chan *nats.Conn)
	serviceNatsClient := make(chan *nats.Conn)
	// go func() {
	// 	client, err := sql.Open("mysql", connInfo.MysqlClient) //"user:password@/dbname"
	// 	client.SetMaxIdleConns(MYSQLMAXCONNS)
	// 	mysqlClient <- client
	// 	if err != nil {
	// 		log.WithFields(log.Fields{
	// 			"error": err.Error(),
	// 		}).Error("mysql client error")
	// 	}
	// }()
	go func() {
		client, err := getEtcdClient(connInfo.EtcdClient)
		etcdClient <- client
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("etcd client error")
		}
	}()
	go func() {
		client, err := nats.Connect(connInfo.NatsManageClient) //"nats://message-business.service.consul:4222"
		managerNatsClient <- client
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("managerNats Client error")
		}
	}()
	go func() {
		client, err := nats.Connect(connInfo.NatsServiceClient) //"nats://message-business.service.dc-anyuncloud.consul:4222"
		serviceNatsClient <- client
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("serviceNats client error")
		}
	}()
	return &QueryClient{
		// MySQLClient:       <-mysqlClient,
		EtcdClient:        <-etcdClient,
		ParentContext:     ctx,
		RedisClient:       pool.GetRedisPool(connInfo.RedisClient).Get(),
		ManageNatsClient:  <-managerNatsClient,
		ServiceNatsClient: <-serviceNatsClient,
	}, nil
}

func (_this *QueryClient) CheckFromEtcd(mess string) (string, error) {
	val, err := logRecord(_this, mess)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
		return "", err
	}
	if val != nil {
		return val.(string), nil
	}
	val, err = queryEtcd(_this, mess)
	if err != nil {
		return "", err
	}
	pool.Set_TTL(_this.RedisClient, mess, val.(string), REDIS_TIME_TO_LIVE)
	return val.(string), nil
}

//CheckQueryURL 检查请求地址
func (_this *QueryClient) CheckQueryURL(request *RequestDetail) ([]byte, string, error) {
	//TODO: 根据URI查询API的元数据信息
	//* 0. 缓存查询
	//* 1. 解析URI生成URI的hashid
	//* 2. 根据hashid去etcd查找endpoint定义的键目录
	// _self.etcdClient.Get(_self.ParentContext)
	// m := make(map[string][]string)
	// for k, param := range request.Params {
	// 	var s []string
	// 	for _, p := range param {
	// 		s = append(s, p)
	// 	}
	// 	m[k] = s
	// }
	r, _ := json.Marshal(&request)
	md5param := md5Param(r)
	val, err := logRecord(_this, md5param)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
		return nil, "", err
	}
	if val != nil && len(val.(string)) != 0 {
		return nil, val.(string), nil
	}
	mess, err := _this.ManageNatsClient.Request("API_CONTROLLER_CHANNEL_TEST1", []byte(r), 1000*time.Millisecond)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("illegar params")
		return nil, "", err
	}
	pool.Set_TTL(_this.RedisClient, md5param, string(mess.Data), REDIS_TIME_TO_LIVE)
	return nil, string(mess.Data), nil

	//检查uri地址是否合法
	// api, version, application, err := common.CheckURL(request.URL)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"error": err.Error(),
	// 	}).Error("illegar params")
	// 	return nil, "", err
	// }
	// params := request.Values
	// direction := FIRST_API
	// var uuid string

	// for {
	// 	switch direction {
	// 	case FIRST_API:
	// 		val, err := logRecord(_this, api)
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		pool.Set_TTL(_this.RedisClient, api, val.(string), REDIS_TIME_TO_LIVE)
	// 		uuid = val.(string)
	// 		direction = SECOND_DETAIL
	// case SECOND_DETAIL:
	// uri := "/meta/api_meta/" + uuid + "/desc"
	// val, err := logRecord(_this, uri)
	// if err != nil {
	// 	return nil, "", err
	// }
	// pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// direction = THIRD_VERSION
	// 	case THIRD_VERSION:
	// 		uri := "/meta/api_meta/" + uuid + "/versions"
	// 		val, err := logRecord(_this, uri)
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		if !strings.Contains(val.(string), version) {
	// 			return nil, "", errors.New("illager version")
	// 		}
	// 		pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// 		direction = FORTH_METHOD
	// 	case FORTH_METHOD:
	// 		uri := "/meta/api_meta/" + uuid + "/" + version + "/method"
	// 		val, err := logRecord(_this, uri)
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		if !strings.Contains(strings.ToUpper(val.(string)), strings.ToUpper(request.Method)) {
	// 			return nil, "", errors.New("FORTH_METHOD illager method")
	// 		}
	// 		pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// 		direction = FIVE_PARAMS
	// 	case FIVE_PARAMS:
	// 		if params == nil {
	// 			return nil, "", errors.New("illager param")
	// 		}
	// 		for k, param := range params {
	// 			for i, p := range param {
	// 				uri := "/meta/api_meta/" + uuid + "/" + version + "/request/" + request.Method + "/params/" + k + ":" + strconv.Itoa(i)
	// 				val, err := logRecord(_this, uri)
	// 				if err != nil {
	// 					return nil, "", err
	// 				}
	// 				if val.(string) != p {
	// 					return nil, "", errors.New("FIVE_METHOD illager param")
	// 				}
	// 				pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// 			}
	// 		}
	// 		switchMethod(request.Method, &direction)
	// 	case SIXTH_BODY:
	// 		uri := "/meta/api_meta/" + uuid + "/" + version + "/request/POST/body"
	// 		val, err := logRecord(_this, uri)
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		if !strings.Contains(strings.ToUpper(val.(string)), strings.ToUpper(request.ContentType)) {
	// 			return nil, "", errors.New("SIXTH_BODY illager param")
	// 		}
	// 		pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// 		direction = SEVENTH_PRODUCT
	// 	case SEVENTH_PRODUCT:
	// 		uri := "/meta/api_meta/" + uuid + "/" + version + "/request/POST/product"
	// 		val, err := logRecord(_this, uri)
	// 		if err != nil {
	// 			return nil, "", err
	// 		}
	// 		if !strings.Contains(strings.ToUpper(val.(string)), strings.ToUpper(request.AcceptType)) {
	// 			return nil, "", errors.New("SEVENTH_BODY illager param")
	// 		}
	// 		pool.Set_TTL(_this.RedisClient, uri, val.(string), REDIS_TIME_TO_LIVE)
	// 		header := &MessageHeader{
	// 			Version:     version,
	// 			Application: application,
	// 			Time:        time.Now().Unix(),
	// 			Type:        "req",
	// 		}
	// 		requestMessage := &RequestMessage{
	// 			MessageHeader: *header,
	// 			Business:      "service",
	// 			Content:       params,
	// 		}
	// 		data, _ := json.Marshal(requestMessage)
	// 		return data, val.(string), nil
	// 	default:
	// 		return nil, "", errors.New("illager param")
	// 	}
	// }
}

//QueryNatsDatas 查询nats获取数据
func (_this *QueryClient) QueryNatsDatas(val *NatsRequestDetail) (interface{}, error) {
	msg, err := _this.ServiceNatsClient.Request(val.channel, val.data, 10*time.Millisecond)
	if err != nil {
		log.WithFields(log.Fields{
			"prefix": "discovery.QueryNatsDatas",
			"error":  err.Error(),
		}).Error("QueryNatsDatas error")
	}
	fmt.Println(msg)
	return nil, nil
}

//QueryAPIMetaDatas -- 根据URI查询API的元数������信息
//* ����������������������������������要缓存
func (_this *QueryClient) QueryAPIMetaDatas(request *RequestDetail) (*EndpointAPIMetaData, error) {
	uri := request.URL
	serviceID := common.GetMD5Hash(uri)
	val, err := logRecord(_this, serviceID)
	if err != nil {
		return nil, err
	}
	fmt.Println(val)
	return nil, nil
}

//QueryAllAPIMetaDatas -- 获取所有的API的元数�����
func (_this *QueryClient) QueryAllAPIMetaDatas() (*[]EndpointAPIMetaData, error) {
	//TODO: 根据URI查询API的元数据信息
	//* 1.��etcd查找endpoint目录下所有的子目录并遍历所有API元数�����
	return nil, nil
}

//QueryServiceSRVRecords -- ����询服务的DNS记录
//* 服务的端口信息
//* 服务的IP
func (_this *QueryClient) QueryServiceSRVRecords(serviceFQDN string) ([]string, error) {
	//TODO: �������询服务的DNS记录
	return nil, nil
}

//QueryServiceCacheRecoeds -- 查询缓存的服务DNS记录
func (_this *QueryClient) QueryServiceCacheRecoeds(serviceFQDN string) ([]string, error) {
	//TODO: 查询缓存的服务DNS记录
	return nil, nil
}
