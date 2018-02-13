package server

import (
	"net"
	"net/http"
	"time"

	"github.com/coreos/etcd/client"
)

//Config --
type Config struct {
	EndPoints []string
	Transport CancelableTransport
}

//DefaultTransport --
var DefaultTransport CancelableTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	Dial: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 10 * time.Second,
}

//CancelableTransport --
type CancelableTransport interface {
	http.RoundTripper
	CancelRequest(req *http.Request)
}

func (cfg *Config) getTransport() CancelableTransport {
	if cfg.Transport == nil {
		return DefaultTransport
	}
	return cfg.Transport
}

func clientEtcd(s string, cg *Config) (*client.Response, error) {
	// client, err := clientv3.New(clientv3.Config{
	// 	Endpoints:   []string{etcdEndpoints},
	// 	DialTimeout: 5 * time.Second,
	// })
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"prefix": "server.clientEtcd",
	// 		"error":  err.Error(),
	// 	}).Error("clientEtcd work error")
	// 	return nil, err
	// }
	// kApi := clientv3.NewKeysAPI(c)
	// v, err := kApi.Get(context.Background(), s)
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"prefix": "server.clientEtcd_kApi.get",
	// 		"error":  err.Error(),
	// 	}).Error("kApi.Get work error")
	// 	return nil, err
	// }
	return nil, nil
}
