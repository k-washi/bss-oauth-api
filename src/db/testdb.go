package db

import (
	"fmt"

	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"
	ge "github.com/k-washi/bss-oauth-api/src/grpc_errors"
	"github.com/k-washi/bss-utils/logger"
)

type testdbRepositry struct{}

var (
	testDB = map[string]accesstoken.AccessToken{}
)

//GetByID get accesstoken by id
func (r *testdbRepositry) GetByID(userID string) (*accesstoken.AccessToken, error) {
	if at, ok := testDB[userID]; ok {
		return &at, nil
	}

	logger.Log.Error(fmt.Sprintf("Can not get access toekn: %s", userID))
	return nil, ge.InternalError("Can not get access token")
}

//Create
func (r *testdbRepositry) Create(at accesstoken.AccessToken) error {
	testDB[at.UserID] = at
	return nil
}

//Delete
func (r *testdbRepositry) Delete(userID string) error {
	delete(testDB, userID)
	return nil
}

//UpdateExpirationTime
func (r *testdbRepositry) UpdateExpirationTime(at accesstoken.AccessToken) error {
	testDB[at.UserID] = at
	return nil
}
