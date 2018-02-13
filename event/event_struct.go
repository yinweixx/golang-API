package event

//Event -- 平台事件结构体
type Event struct {
	From     string //事件来源URI
	Describe string //事件的详细信息,JSON
}

type ConnInfo struct{
	MysqlClient string
	EtcdClient []string
	RedisClient string
	NatsServiceClient string
	NatsManageClient string
}