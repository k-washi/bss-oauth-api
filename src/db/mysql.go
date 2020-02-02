package db

import (
	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"
)

//DbRepositry  mysql db repository
type DbRepositry interface {
	GetByID(string) (*accesstoken.AccessToken, error)
	Create(accesstoken.AccessToken) error
	Delete(string) error
	UpdateExpirationTime(accesstoken.AccessToken) error
}

//NewRepository get db repository
func NewRepository(dbName string) DbRepositry {
	if dbName == "mysql" {
		return &dbRepositry{}
	}
	return &testdbRepositry{}
}

type dbRepositry struct{}

//GetByID get accesstoken by id
func (r *dbRepositry) GetByID(userID string) (*accesstoken.AccessToken, error) {
	var result accesstoken.AccessToken
	return &result, nil
}

//Create
func (r *dbRepositry) Create(at accesstoken.AccessToken) error {
	return nil
}

//Delete
func (r *dbRepositry) Delete(userID string) error {
	return nil
}

//UpdateExpirationTime
func (r *dbRepositry) UpdateExpirationTime(at accesstoken.AccessToken) error {
	return nil
}
