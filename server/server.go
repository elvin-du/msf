package server

import (
	"msf/grpclb"
	"msf/log"
	"msf/util"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
)

type RPCServer struct {
	ServiceName    string
	ServiceVersion string
	Meta           map[string]string
	ETCDAddrs      string

	ip         string
	port       int
	grpcServer *grpc.Server
}

func NewRPCServer(serviceName, serviceVer, etcdAddrs string, meta map[string]string) *RPCServer {
	grpc_logrus.ReplaceGrpcLogger(log.LogrusEntry)
	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor),
			),
			grpc_logrus.UnaryServerInterceptor(log.LogrusEntry),
		))

	return &RPCServer{
		ServiceName:    serviceName,
		ServiceVersion: serviceVer,
		Meta:           meta,
		grpcServer:     srv,
		ETCDAddrs:      etcdAddrs,
	}
}

func (s *RPCServer) Server() *grpc.Server {
	return s.grpcServer
}

func (s *RPCServer) Run() {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	s.port = listener.Addr().(*net.TCPAddr).Port

	//TODO
	err = grpclb.Register(s.ServiceName, util.InternalIP, s.port, s.ETCDAddrs, time.Second*10, 15)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		stop := <-ch
		log.LogrusEntry.Errorf("receive signal '%v'", stop)
		grpclb.UnRegister()
		os.Exit(1)
	}()

	//	pb.RegisterGateServer(srv, &Gate{})

	err = s.grpcServer.Serve(listener)
	if err != nil {
		//		log.Fatalf("failed to listen: %v", err)
	}
}
