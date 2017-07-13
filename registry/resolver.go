package registry

import (
	"errors"
	"fmt"
	"strings"

	etcd3 "github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc/naming"
)

// resolver is the implementaion of grpc.naming.Resolver
type resolver struct {
	serviceName string // service name to resolve
}

// NewResolver return resolver with service name
func NewResolver(serviceName string) *resolver {
	return &resolver{serviceName: serviceName}
}

// Resolve to resolve the service from etcd, etcdAddrs is the dial address of etcd
// etcdAddrs example: "http://127.0.0.1:2379,http://127.0.0.1:12379,http://127.0.0.1:22379"
func (re *resolver) Resolve(etcdAddrs string) (naming.Watcher, error) {
	if re.serviceName == "" {
		return nil, errors.New("grpclb: no service name provided")
	}

	// generate etcd client
	client, err := etcd3.New(etcd3.Config{
		Endpoints: strings.Split(etcdAddrs, ","),
	})
	if err != nil {
		return nil, fmt.Errorf("grpclb: creat etcd3 client failed: %s", err.Error())
	}

	// Return watcher
	return &watcher{re: re, client: *client}, nil
}
