package newConn

import (
	"context"

	"e.coding.net/anyun-cloud-api-gateway/discovery"
	"e.coding.net/anyun-cloud-api-gateway/event"
	log "github.com/sirupsen/logrus"
)

var (
	new_conn *discovery.QueryClient
)

func NewConn(info ...*event.ConnInfo) {
	var aInfo *event.ConnInfo = nil
	var err error
	for _, fo := range info {
		aInfo = fo
	}
	new_conn, err = discovery.NewQueryClient(context.Background(), getConnInfo(aInfo))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(err.Error())
	}
}

func GetConn(ctx context.Context, info ...*event.ConnInfo) (*discovery.QueryClient, error) {
	if new_conn == nil {
		var aInfo *event.ConnInfo = nil
		var err error
		for _, fo := range info {
			aInfo = fo
		}
		new_conn, err = discovery.NewQueryClient(ctx, getConnInfo(aInfo))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(err.Error())
			return nil, err
		}
	}
	return new_conn, nil
}

func getConnInfo(info *event.ConnInfo) *event.ConnInfo {
	if info != nil {
		return info
	}
	return &event.ConnInfo{
		EtcdClient:        []string{"server.etcd.service.dc-anyuncloud.consul:2379"},
		MysqlClient:       "root:1234qwer@tcp(client.mysql.service.consul:3306)/test?charset=utf8", //http://pear.php.net/manual/en/package.database.db.intro-dsn.php
		NatsServiceClient: "nats://message-business.service.dc-anyuncloud.consul:4222",
		NatsManageClient:  "nats://message-business.service.dc-anyuncloud.consul:4222",
		RedisClient:       "redis.service.dc-anyuncloud.consul:6379",
	}
}

func GetConnInfo() *event.ConnInfo {
	return getConnInfo(nil)
}
