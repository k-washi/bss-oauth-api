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
	tls      bool
	keyFile  string
	certFile string
)

func init() {
	host = os.Getenv("GRPC_AT_HOST")
	port := os.Getenv("GRPC_ACT_PORT")
	if port == "" {
		port = "30001"
	}
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

}

func StartApplication() {
	grpcHandler := server.NewAccessTokenServiceServer(service.NewService(db.NewRepository()))
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		logger.Error("Failed to listen", err)
		panic(err)
	}

	var opts []grpc.ServerOption
	if tls {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Error("Failed to generate credentials", err)
			panic(err)
		}
		opts = []grpc.ServerOption{
			grpc.Creds(creds),
		}
	}

	//Keepalive
	/*
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             2 * time.Second,
			PermitWithoutStream: true,
		}))
	*/

	grpcServer := grpc.NewServer(opts...)
	logger.Info("NewServer...")
	pb.RegisterAccessTokenServiceServer(grpcServer, grpcHandler)
	logger.Debug("Set grpc Handler")
	grpcServer.Serve(lis)
	logger.Info("Start session application")

}
