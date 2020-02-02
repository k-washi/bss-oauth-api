package app

import (
	"fmt"
	"net"
	"os"
	"strconv"

	pb "github.com/k-washi/bss-oauth-api/access_token_proto"
	"github.com/k-washi/bss-oauth-api/src/db"
	"github.com/k-washi/bss-oauth-api/src/server"
	"github.com/k-washi/bss-oauth-api/src/service"
	"github.com/k-washi/bss-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	host     string
	port     string
	address  string
	tls      bool
	keyFile  string
	certFile string
	database string
)

func init() {
	host = os.Getenv("GRPC_AT_HOST")
	port := os.Getenv("GRPC_AT_PORT")
	if port == "" {
		port = "30001"
	}
	address = fmt.Sprintf("%s:%s", host, port)

	//tls setting
	var err error
	tls, err = strconv.ParseBool(os.Getenv("GRPC_AT_TLS"))
	if err != nil {
		tls = false
	}
	if tls {
		keyFile = os.Getenv("GRPC_AT_KEY")
		if keyFile == "" {
			keyFile = "./grpc_at_server1.key"
		}
		certFile = os.Getenv("GRPC_AT_CERT")
		if certFile == "" {
			certFile = "./grpc_at_server1.pem"
		}
	}

	//database
	database = os.Getenv("AT_DATABASE")

}

//StartApplication start grpc server for access token managoing
func StartApplication() {
	grpcHandler := server.NewAccessTokenServiceServer(service.NewService(db.NewRepository(database)))
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Failed to listen: %s", err.Error()))
		panic(err)
	}

	var opts []grpc.ServerOption
	if tls {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Failed to generate credentials: %s", err.Error()))
			panic(err)
		}
		opts = []grpc.ServerOption{
			grpc.Creds(creds),
		}
	}

	//opts = append(opts, grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(logger.Log)))

	//Keepalive
	/*
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             2 * time.Second,
			PermitWithoutStream: true,
		}))
	*/

	grpcServer := grpc.NewServer(opts...)
	logger.Log.Info(fmt.Sprintf("NewServer... %s", address))
	pb.RegisterAccessTokenServiceServer(grpcServer, grpcHandler)
	logger.Log.Debug("Set grpc Handler")
	grpcServer.Serve(lis)
	logger.Log.Info("Start session application")

}
