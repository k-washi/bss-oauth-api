package access_token_client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/k-washi/bss-oauth-api/access_token_proto"
	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"
	"github.com/k-washi/bss-utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	grpcTimeOutSeccnds = 1
)

var (
	host    string
	port    string
	address string
	tlsTF   bool
	caFile  string
	opts    []grpc.DialOption
)

func init() {

	host = os.Getenv("GRPC_AT_HOST")
	port := os.Getenv("GRPC_ACT_PORT")
	if port == "" {
		port = "30001"
	}
	address = fmt.Sprintf("%s:%s", host, port)

	var err error

	tlsTF, err = strconv.ParseBool(os.Getenv("GRPC_AT_TLS"))
	if err != nil {
		tlsTF = false
	}

	if tlsTF {
		caFile = os.Getenv("GRPC_AT_CA")
		if caFile == "" {
			caFile = "./ca.cert"
		}

		b, _ := ioutil.ReadFile(caFile)
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			panic("credentials: failed to append certificates")
		}
		//creds, err := credentials.NewClientTLSFromFile(caFile)
		config := &tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            cp,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	//opts = append(opts, grpc.WithUnaryInterceptor(grpc_zap.UnaryClientInterceptor(logger.Log)))
	//opts = append(opts, grpc.WithBlock()) //接続確立のためブロック

}

func connection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("GRPC Client connection error: %s", err.Error()))
		return nil, errors.New("GRPC Client connection error")
	}
	return conn, nil
}

//Check grpc check method
func Check(at accesstoken.AccessToken) error {
	conn, err := connection()
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Second*grpcTimeOutSeccnds)
	defer ctxCancel()

	client := pb.NewAccessTokenServiceClient(conn)

	grpcAt := &pb.AccessTokenCheck{
		AccessToken: &pb.AccessToken{
			UserId:      at.UserID,
			AccessToken: at.AccessToken,
		},
	}
	_, err = client.Check(ctx, grpcAt)
	if err != nil {
		grpcErrorCheck(err)
		return err
	}
	logger.Log.Debug(fmt.Sprintf("gRPC AT Check successs:%s", at.UserID))
	return nil

}

//Create grpc create method
func Create(userID string) (*accesstoken.AccessToken, error) {
	conn, err := connection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, ctxCanncel := context.WithTimeout(context.Background(), time.Second*grpcTimeOutSeccnds)
	defer ctxCanncel()

	client := pb.NewAccessTokenServiceClient(conn)

	grpcAt := &pb.AccessTokenCreate{
		UserId: &pb.UserID{
			UserId: userID,
		},
	}

	reply, err := client.Create(ctx, grpcAt)
	if err != nil {
		grpcErrorCheck(err)
		return nil, err
	}
	logger.Log.Debug(fmt.Sprintf("gRPC AT Create successs:%s", userID))

	return &accesstoken.AccessToken{
		UserID:      userID,
		AccessToken: reply.AccessToken.GetAccessToken(),
	}, nil

}

func grpcErrorCheck(err error) {
	respErr, ok := status.FromError(err)
	if ok {
		if respErr.Code() == codes.InvalidArgument {
			logger.Log.Error(fmt.Sprintf("gRPC AT invalid argument: %s", err.Error()))
			//invalid argument
			//this error is possible unauthorized access by user
		} else if respErr.Code() == codes.DeadlineExceeded {
			logger.Log.Error(fmt.Sprintf("gRPC AT Time Out: %s", err.Error()))
		} else if respErr.Code() == codes.Internal {
			logger.Log.Error(fmt.Sprintf("gRPC Internal error: %s", err.Error()))
		} else {
			logger.Log.Error(fmt.Sprintf("gRPC AT Unkown error: %s", err.Error()))
		}
	} else {
		logger.Log.Error(fmt.Sprintf("gRPC Check error: %s", err.Error()))
	}
}
