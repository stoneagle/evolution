package rpc

import (
	"context"
	"encoding/json"
	"net"
	"quant/backend/rpc/engine"
	"quant/backend/rpc/models"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Rpc struct {
	Host   string
	Port   string
	Client *engine.EngineServiceClient
}

func NewRpc(host, port string) *Rpc {
	// golang访问python，需要使用TBuffered，并且python服务端配置为127.0.0.1
	// 因为python使用tcpv6
	ret := Rpc{
		Host: host,
		Port: port,
	}
	return &ret
}

func (r *Rpc) client() (client *engine.EngineServiceClient, transport *thrift.TSocket, err error) {
	transportFactory := thrift.NewTBufferedTransportFactory(1024)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transport, err = thrift.NewTSocket(net.JoinHostPort(r.Host, r.Port))
	if err != nil {
		panic(err)
	}
	useTransport, err := transportFactory.GetTransport(transport)
	if err != nil {
		panic(err)
	}
	client = engine.NewEngineServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		panic(err)
	}
	return
}

func (r *Rpc) GetType(atype engine.AssetType) (ret map[string]map[string][]string, err error) {
	ret = make(map[string]map[string][]string)
	client, transport, err := r.client()
	if err != nil {
		return
	}
	defer transport.Close()
	d, err := client.GetType(context.Background(), atype)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(d.Data), &ret)
	if err != nil {
		return
	}
	return
}

func (r *Rpc) GetClassify(atype engine.AssetType, ctype, source, sub string) (ret []models.Classify, err error) {
	ret = []models.Classify{}
	client, transport, err := r.client()
	if err != nil {
		return
	}
	defer transport.Close()
	d, err := client.GetClassify(context.Background(), atype, ctype, source, sub)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(d.Data), &ret)
	if err != nil {
		return
	}
	return
}
