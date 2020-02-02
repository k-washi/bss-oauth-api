package server

import (
	"context"
	"fmt"

	"github.com/k-washi/bss-utils/logger"
	"google.golang.org/grpc/peer"

	pb "github.com/k-washi/bss-oauth-api/access_token_proto"
	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"
	ge "github.com/k-washi/bss-oauth-api/src/grpc_errors"

	"github.com/k-washi/bss-oauth-api/src/service"
)

type accessTokenServiceServer struct {
	service service.Service
	//mu      sync.RWMutex
}

//NewAccessTokenServiceServer service injection for db selecting
/*
dbをインジェクション
atHandler := http.NewAccessTokenHandler(
	access_token.NewService(rest.NewRestUsersRepository(), db.NewRepository())
)
*/
func NewAccessTokenServiceServer(service service.Service) *accessTokenServiceServer {
	return &accessTokenServiceServer{
		service: service,
	}
}

func (h *accessTokenServiceServer) Check(ctx context.Context, r *pb.AccessTokenCheck) (*pb.AccessTokenResponse, error) {
	//h.mu.Lock()
	//defer h.mu.Unlock()

	at := accesstoken.AccessToken{
		UserID:      r.AccessToken.GetUserId(),
		AccessToken: r.AccessToken.GetAccessToken(),
	}

	logger.Info(fmt.Sprintf("Access token check: UserID: %s, IP: %s", at.UserID, getIP(ctx)))

	upAt, err := h.service.Check(at)
	if err != nil {
		logger.Error(fmt.Sprintf("Access token check error: UserID: %s", at.UserID), err)
		return nil, err
	}

	return &pb.AccessTokenResponse{
		AccessToken: &pb.AccessToken{
			UserId:      upAt.UserID,
			AccessToken: upAt.AccessToken,
		},
		Status: 0,
		Msg:    "Check access token OK!",
	}, nil
}

func (h *accessTokenServiceServer) Create(ctx context.Context, r *pb.AccessTokenCreate) (*pb.AccessTokenResponse, error) {
	//h.mu.Lock()
	//defer h.mu.Unlock()
	userID := r.UserId.GetUserId()
	logger.Info(fmt.Sprintf("Access token create: UserID: %s, IP: %s", userID, getIP(ctx)))

	at, err := h.service.Create(userID)

	if err != nil {
		logger.Error(fmt.Sprintf("Access token create error: UserID: %s", at.UserID), err)
		return nil, err
	}

	return &pb.AccessTokenResponse{
		AccessToken: &pb.AccessToken{
			UserId:      userID,
			AccessToken: at.AccessToken,
		},
		Status: 0,
		Msg:    "Create access token OK!",
	}, nil

}

func (h *accessTokenServiceServer) Get(ctx context.Context, r *pb.AccessTokenGet) (*pb.AccessTokenResponse, error) {
	return nil, ge.UnImplemented("Get method dose not implement")
}

func getIP(ctx context.Context) string {
	if pr, ok := peer.FromContext(ctx); ok {
		return pr.Addr.String()
	}
	return ""
}
