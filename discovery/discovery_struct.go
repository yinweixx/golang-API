//Package discovery -- 服务发现的结构体定义
//* 云平台服务查询客户端
//* 元数据版本信息
//* API元数据信息
//* 服务元数信息
//* DNS SRV记录
package discovery

import (
	"context"
	"database/sql"

	"github.com/coreos/etcd/clientv3"
	"github.com/garyburd/redigo/redis"
	"github.com/nats-io/go-nats"
)

//QueryClient -- 云平台服务查询客户端
//* 通过dns查询服务运行信息
//* 通过etcd查询元数据信息
type QueryClient struct {
	DNS               string           //查询DNS
	Etcd              []string         //etcd连接字符串
	Redis             string           //redis连接字符串
	ParentContext     context.Context  //上层应用程序上下文
	MySQLClient       *sql.DB          // MySQLClient
	EtcdClient        *clientv3.Client //etcd v3客户端
	RedisClient       redis.Conn       //redis 客户端
	ManageNatsClient  *nats.Conn       //管理nats
	ServiceNatsClient *nats.Conn       //业务nats
}

//RequestDetail -- 请求数据实体
type RequestDetail struct {
	URL         string              `json:"url"`
	Method      string              `json:"method"`
	Params      map[string][]string `json:"params"`
	ContentType string              `json:"contenttype"`
	AcceptType  string              `json:"accepttype"`
}

//ReponseDetail -- 返回数据实体
type ReponseDetail struct {
	URL  string `json:"url"`
	Uuid string `json:"uuid"`
}

//NatsRequestDetail -- Nats数据请求实体
type NatsRequestDetail struct {
	channel string
	data    []byte
}

//MetaDataVersion -- 元数据版本信息结构体
type MetaDataVersion struct {
	//TODO: 元数据版本信息定义
	serviceBranch string
	tagVersion    string
	info          string
}

//EndpointAPIMetaData -- API元数据信息结构体
type EndpointAPIMetaData struct {
	//TODO: API元数据信息定义
	Service EndpointServiceMetaData
	Version MetaDataVersion
}

//EndpointServiceMetaData -- 服务元数信息据结构体
//* 如果服务支持缓存，必须设置缓存的TTL
type EndpointServiceMetaData struct {
	//TODO: 服务元数据信息定义
	projectName  string
	serviceName  string
	Version      MetaDataVersion
	SupportCache bool //服务是否支持缓存
	CacheTTL     bool //服务缓存的TTL
}

//DNSTypeSRVRecord -- DNS SRV记录
type DNSTypeSRVRecord struct {
	//TODO: DNS SRV记录定义
	SRV []struct {
		FQDN string //提供服务节点的域名信息
		IP   string //提供服务节点的IP地址信息
		Port int    //提供服务节点的端口信息
	}
}

//ServerConnectionInfo -- 服务器连接信息
type ServerConnectionInfo struct {
	Address string //URI,IP,ADDR
	Port    int    //服务端口
}

//MessageHeader -- 请求消息头
type MessageHeader struct {
	Version     string `json:"version"`
	Type        string `json:"type"`
	Application string `json:"application"`
	Time        int64  `json:"time"`
}

//RequestMessage -- 请求消息实体
type RequestMessage struct {
	MessageHeader MessageHeader `json:"header"`
	Business      string        `json:"business"`
	Content       interface{}   `json:"content"`
}

type header struct {
	Encryption  string `json:"encryption"`
	Timestamp   int64  `json:"timestamp"`
	Key         string `json:"key"`
	Partnercode int    `json:"partnercode"`
}
