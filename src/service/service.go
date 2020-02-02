package service

import (
	ge "github.com/k-washi/bss-oauth-api/src/grpc_errors"
	"github.com/k-washi/bss-utils/logger"

	"github.com/k-washi/bss-oauth-api/src/db"
	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"
)

//Service service interface
type Service interface {
	Check(accesstoken.AccessToken) (*accesstoken.AccessToken, error)
	Create(string) (*accesstoken.AccessToken, error)
}

type service struct {
	dbRepo db.DbRepositry
}

//NewService injection db repository
func NewService(dbRepo db.DbRepositry) Service {
	return &service{
		dbRepo: dbRepo,
	}
}

//Check user idに対するアクセストークが一致かつ、有効期限以内実行
func (s *service) Check(at accesstoken.AccessToken) (*accesstoken.AccessToken, error) {
	if err := accesstoken.ValidUserID(at.AccessToken); err != nil {
		logger.Debug("Invalid access token")
		return nil, ge.InvalidArgument("Invalid access token")
	}
	if err := accesstoken.ValidUserID(at.UserID); err != nil {
		logger.Debug("Invalid user id")
		return nil, ge.InvalidArgument("Invalid user id")
	}

	dbAccessToken, err := s.dbRepo.GetByID(at.UserID)
	if err != nil {
		return nil, err
	}

	//AccessTokenの一致をチェック
	if dbAccessToken.AccessToken != at.AccessToken {
		logger.Debug("Access token is mismatch")
		return nil, ge.InvalidArgument("Access token is mismatch")
	}

	//AccessTokenの有効期限が切れている(頻繁な更新を避けるため、以下の処理を行う)
	if dbAccessToken.IsExpired() {
		//更新期間を超えているか
		if dbAccessToken.IsDelete() {
			//削除
			if err := s.dbRepo.Delete(at.UserID); err != nil {
				return nil, err
			}
			logger.Debug("Access token is expired, and delete access token")
			return nil, ge.InvalidArgument("Access token is expired")
		}
		//更新(時間のみ)
		if err := s.dbRepo.UpdateExpirationTime(at.UserID); err != nil {
			return nil, err
		}

	}
	return dbAccessToken, nil
}

func (s *service) Create(userID string) (*accesstoken.AccessToken, error) {
	if err := accesstoken.ValidUserID(userID); err != nil {
		logger.Debug("Invalid access token")
		return nil, ge.InvalidArgument("Invalid access token")
	}

	at := accesstoken.GetNewAccessToken(userID)
	at.Generate()

	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	dbAccessToken, err := s.dbRepo.GetByID(at.UserID)
	if err != nil {
		return nil, err
	}

	return dbAccessToken, nil

}
