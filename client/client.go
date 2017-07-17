package client

import (
	//	"msf/log"
	//	"msf/registry"
	//	registry "grpclb/etcdv3"
	"time"

	//	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	//	"github.com/grpc-ecosystem/go-grpc-prometheus"
	//	"github.com/coreos/etcd/clientv3"
	"fmt"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	DialTimeout time.Duration = time.Second * 20
)

func Get(serviceName, etcdAddr string) (*grpc.ClientConn, error) {
	//	cli, err := clientv3.NewFromURL("http://localhost:2379")
	//	if nil != err {
	//		return nil, err
	//	}
	cli, err := etcd3.New(etcd3.Config{
		Endpoints:   strings.Split("http://127.0.0.1:2379", ","),
		DialTimeout: DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}



	r := &etcdnaming.GRPCResolver{
		Client: cli,
	}
	//	r := registry.NewResolver(serviceName)
	b := grpc.RoundRobin(r)

	resp, err := cli.Get(context.Background(), "/push",etcd3.WithPrefix())
	fmt.Println(resp, err)

	//	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	//	defer cancel()
//	ctx := context.Background()
	//TODO
	conn, err := grpc.Dial(
//		ctx,
		"/push",
		grpc.WithInsecure(),
		grpc.WithBalancer(b),
		grpc.WithBlock(),
		grpc.WithTimeout(DialTimeout),
		//		grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()),
		//		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		//		log.Fatalln(err)
		panic(err)
	}

	return conn, nil
}
