package client

import (
	"msf/log"
	"msf/registry"
	"time"

	//	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	DialTimeout time.Duration = time.Second * 20
)

func Get(serviceName, etcdAddr string) (*grpc.ClientConn, error) {
	r := registry.NewResolver(serviceName)
	b := grpc.RoundRobin(r)

	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	defer cancel()
	conn, err := grpc.DialContext(
		ctx,
		etcdAddr,
		grpc.WithInsecure(),
		grpc.WithBalancer(b),
		grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()),
		grpc.WithUnaryInterceptor(grpc_prometheus.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return conn, nil
}
