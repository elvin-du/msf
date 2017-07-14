package registry

import (
	"fmt"
	"strings"
	"time"

	etcd3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
)

// Prefix should start and end with no slash
var Prefix = "push"
var client *etcd3.Client
var serviceKey string

var stopSignal = make(chan bool, 1)

type ConfigOption struct {
	Interval       int64 //Internal MUST less than TTL.uint: sec
	TTL            int64
	DialTimeout    time.Duration
	RequestTimeout time.Duration
}

var defaultConfig = &ConfigOption{
	Interval:       5 * 60,
	TTL:            6 * 60,
	DialTimeout:    time.Second * 5,
	RequestTimeout: time.Second * 5,
}

func Register(name, host string, port int, etcdAddrs string) error {
	return RegisterOpt(name, host, port, etcdAddrs, defaultConfig)
}

func RegisterOpt(name, host string, port int, etcdAddrs string, cfg *ConfigOption) error {
	if cfg.Interval >= cfg.TTL {
		err := fmt.Errorf("registry: register interval must less than ttl,but got interval:%d,ttl:%d", cfg.Interval, cfg.TTL)
		grpclog.Errorln(err)
		return err
	}

	serviceValue := fmt.Sprintf("%s:%d", host, port)
	serviceKey = fmt.Sprintf("/%s/%s/%s", Prefix, name, serviceValue)

	// get endpoints for register dial address
	var err error
	client, err = etcd3.New(etcd3.Config{
		Endpoints:   strings.Split(etcdAddrs, ","),
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		grpclog.Errorf("registry: create etcd3 client failed: %v", err)
		return fmt.Errorf("registry: create etcd3 client failed: %v", err)
	}
	//	defer client.Close()

	go func() {
		// invoke self-register with ticker
		ticker := time.NewTicker(time.Duration(cfg.Interval))
		for {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.RequestTimeout)
			defer cancel()
			// minimum lease TTL is ttl-second
			resp, err := client.Grant(ctx, cfg.TTL)
			if nil != err {
				grpclog.Errorln(err)
				continue
			}

			// should get first, if not exist, set it
			ctx, cancel = context.WithTimeout(context.Background(), cfg.RequestTimeout)
			defer cancel()
			_, err = client.Get(ctx, serviceKey)
			if err != nil {
				if err == rpctypes.ErrKeyNotFound {
					ctx, cancel = context.WithTimeout(context.Background(), cfg.RequestTimeout)
					defer cancel()
					if _, err := client.Put(ctx, serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
						grpclog.Errorf("registry: set service '%s' with ttl to etcd3 failed: %s", name, err.Error())
					}
				} else {
					grpclog.Printf("registry: service '%s' connect to etcd3 failed: %s", name, err.Error())
				}
			} else {
				// refresh set to true for not notifying the watcher
				ctx, cancel = context.WithTimeout(context.Background(), cfg.RequestTimeout)
				defer cancel()
				if _, err := client.Put(ctx, serviceKey, serviceValue, etcd3.WithLease(resp.ID)); err != nil {
					grpclog.Errorf("registry: refresh service '%s' with ttl to etcd3 failed: %s", name, err.Error())
				}
			}

			select {
			case <-stopSignal:
				return
			case <-ticker.C:
			}
		}
	}()

	return nil
}

// UnRegister delete registered service from etcd
func UnRegister() error {
	defer client.Close()
	stopSignal <- true
	stopSignal = make(chan bool, 1) // just a hack to avoid multi UnRegister deadlock
	if _, err := client.Delete(context.Background(), serviceKey); err != nil {
		grpclog.Errorf("registry: deregister '%s' failed: %s", serviceKey, err.Error())
		return err
	}

	grpclog.Printf("registry: deregister '%s' ok.", serviceKey)
	return nil
}
