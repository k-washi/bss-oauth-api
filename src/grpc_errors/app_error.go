package grpc_errors

import (
	"google.golang.org/grpc/status"
)

const (
	accessTokenMismatch = "Access token is mismatch"
	accessTokenExpired  = "Access token is expired"

	codeBadAccessToken = 401
)

//AccessTokenMismatch mismatch arrcess token
func AccessTokenMismatch() error {
	return status.Errorf(codeBadAccessToken, accessTokenMismatch)
}

//AccessTokenExpired access token is expired
func AccessTokenExpired() error {
	return status.Errorf(codeBadAccessToken, accessTokenExpired)
}
