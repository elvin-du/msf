package server

import (
	"msf/log"
	myprome "msf/prometheus"
	"msf/registry"
	"msf/util"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	//	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type RPCServer struct {
	ServiceName string
	ETCDAddrs   string

	grpcServer *grpc.Server
	config     *Config
}

type Config struct {
	AuthFunc     grpc_auth.AuthFunc
	PromhttpAddr string
}

func defaultAuthFunc(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

//    var interceptor grpc.UnaryServerInterceptor
func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	//	err = auth(ctx)
	//	if err != nil {
	//		return
	//	}
	// 继续处理请求
	start := time.Now()
	resp, err = handler(ctx, req)
	if nil != err {
		log.Errorln(err)
	}
	log.Infof("spent time:%+v", time.Now().Sub(start))

	log.Debugf("%+v\n", req)
	log.Debugf("%+v\n", *info)

	//	myprome.LogPrometheusCount(&myprome.CountLables{
	//		Caller:info.FullMethod
	//	})
	return
}

func NewRPCServer(serviceName, etcdAddrs string, cfg *Config) *RPCServer {
	grpc_logrus.ReplaceGrpcLogger(log.LogrusEntry)
	if nil == cfg.AuthFunc {
		cfg.AuthFunc = defaultAuthFunc
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
			),
			grpc_logrus.UnaryServerInterceptor(log.LogrusEntry),
			grpc_auth.UnaryServerInterceptor(cfg.AuthFunc),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_opentracing.UnaryServerInterceptor(),
			Interceptor,
		))

	return &RPCServer{
		ServiceName: serviceName,
		grpcServer:  srv,
		ETCDAddrs:   etcdAddrs,
		config:      cfg,
	}
}

//Must first call this method
func (s *RPCServer) Init(f func() error) {
	myprome.Init()

	if err := f(); nil != err {
		log.Fatalln(err)
	}
}

func (s *RPCServer) Run() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	log.Infof("%s rpc listening on %s:%d", s.ServiceName, util.InternalIP, port)

	err = registry.Register(s.ServiceName, util.InternalIP, port, s.ETCDAddrs)
	if err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		stop := <-ch
		log.Errorf("receive signal '%v'", stop)
		registry.UnRegister()
		os.Exit(1)
	}()

	//	pb.RegisterGateServer(srv, &Gate{})

	go func() {
		//After all your registrations, make sure all of the Prometheus metrics are initialized.
		grpc_prometheus.Register(s.Server())
		// Register Prometheus metrics handler.
		http.Handle("/metrics", promhttp.Handler())
		if err = http.ListenAndServe(s.config.PromhttpAddr, nil); nil != err {
			log.Fatalln(err)
		}
	}()

	err = s.grpcServer.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *RPCServer) Server() *grpc.Server {
	return s.grpcServer
}
