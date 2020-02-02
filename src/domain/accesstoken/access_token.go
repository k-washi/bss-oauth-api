package accesstoken

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/k-washi/bss-utils/uniquegenerator"

	"github.com/k-washi/bss-utils/encryption"
)

const (
	expirationTime       = 24
	expirationUpdateTime = 24
)

//AccessToken for session
type AccessToken struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
}

//Validate accesstoken validation
func (at *AccessToken) Validate() error {
	at.UserID = strings.TrimSpace(at.UserID)
	at.AccessToken = strings.TrimSpace(at.AccessToken)

	if err := ValidAccessToken(at.AccessToken); err != nil {
		return errors.New("Invalid access token id")
	}
	if err := ValidUserID(at.UserID); err != nil {
		return err
	}
	if at.Expires <= 0 {
		return errors.New("Invalid expiration time")
	}

	return nil
}

//GetNewAccessToken get AccessToken
func GetNewAccessToken(userID string) AccessToken {
	return AccessToken{
		UserID:  userID,
		Expires: time.Now().UTC().Add(expirationTime * time.Hour).Unix(),
	}
}

//Generate generate access token
func (at *AccessToken) Generate() {
	at.AccessToken = encryption.GetMD5(at.UserID, strconv.FormatInt(at.Expires, 10))
}

//UpdateExpirationTime update expiration time
func (at *AccessToken) UpdateExpirationTime() {
	at.Expires = time.Now().UTC().Add(expirationTime * time.Hour).Unix()
}

//IsExpired check access token expired: 期限切れを超えている => true
func (at *AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

//IsDelete check access token delete time: 期限切れ+更新時間を超えている => true
func (at *AccessToken) IsDelete() bool {
	return time.Unix(at.Expires, 0).Add(expirationUpdateTime * time.Hour).Before(time.Now().UTC())

}

//ValidUserID validate userid by  md5
func ValidUserID(userID string) error {
	if err := uniquegenerator.Validation(userID); err != nil {
		return err
	}
	return nil
}

//ValidAccessToken validate accesstoken
func ValidAccessToken(at string) error {
	if err := encryption.Validation(at); err != nil {
		return err
	}
	return nil
}
