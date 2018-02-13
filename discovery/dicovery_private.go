package discovery

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

func getEtcdClient(etcd []string) (*clientv3.Client, error) {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd,
		DialTimeout: 6 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return etcdClient, nil
}

func logRecord(_this *QueryClient, uri string) (interface{}, error) {
	val, err := doCheck(_this, uri)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func switchMethod(method string, direction *string) {
	switch method {
	case "GET":
		*direction = SEVENTH_PRODUCT
	case "OPTION":
		*direction = SEVENTH_PRODUCT
	case "HEAD":
		*direction = SEVENTH_PRODUCT
	case "POST":
		*direction = SIXTH_BODY
	case "PUT":
		*direction = SIXTH_BODY
	case "DELETE":
		*direction = SIXTH_BODY
	default:
		*direction = "error"
	}
}

func doCheck(_this *QueryClient, str string) (string, error) {
	value, err := redis.String(_this.RedisClient.Do("GET", str))
	log.WithFields(log.Fields{
		"prefix": "discovery.doCheck",
		"value":  value,
		"err":    err,
	}).Info("redisClient working")
	if value != "" {
		return value, nil
	}
	return "", nil
}

func queryEtcd(_this *QueryClient, mess string) (string, error) {
	val, err := _this.EtcdClient.Get(_this.ParentContext, mess)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("etcdClient error")
		return "", err
	}
	if len(val.Kvs) == 0 {
		return "", nil
	}
	return string(val.Kvs[0].Value), nil
}

func BuildMessage(request *RequestDetail) (data []byte) {
	header := &MessageHeader{
		Application: "api-gateway",
		Time:        time.Now().Unix(),
		Type:        "req",
		Version:     "1.0.0",
	}

	// var s1 []string
	// s1 = append(s1, "05862547438137427641")
	// var s2 []string
	// s2 = append(s2, "param2-value1", "param2-value2")
	// m := make(map[string][]string)
	// m["gbid"] = s1
	// m["key2"] = s2

	requestMessage := &RequestMessage{
		MessageHeader: *header,
		Business:      "service",
		Content:       request.Params,
	}
	data, _ = json.Marshal(requestMessage)
	return
}

func md5Param(param []byte) string {
	hasher := md5.New()
	hasher.Write(param)
	return hex.EncodeToString(hasher.Sum(nil))
}
