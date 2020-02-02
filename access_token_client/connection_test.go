package access_token_client

import (
	"fmt"
	"testing"

	"github.com/k-washi/bss-oauth-api/src/app"

	"github.com/k-washi/bss-utils/uniquegenerator"

	"github.com/k-washi/bss-utils/logger"

	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"

	"github.com/tj/assert"
)

func TestConnectionTimeOutError(t *testing.T) {

	at := accesstoken.AccessToken{
		AccessToken: "donadonaTest",
		UserID:      "absss",
	}
	logger.Log.Info("Connection error test start")
	err := Check(at)
	assert.Error(t, err)
	logger.Log.Info("Connection error test finish")
}

func TestAtCreate(t *testing.T) {

	go func() {
		app.StartApplication()
	}()

	n := 1000
	for range make([]int, n) {

		userID := uniquegenerator.Get()

		at, err := Create(userID)
		assert.Nil(t, err)
		logger.Log.Info(fmt.Sprintf("%s : %s", at.UserID, at.AccessToken))

		err = Check(accesstoken.AccessToken{
			AccessToken: at.AccessToken,
			UserID:      at.UserID,
		})
		assert.Nil(t, err)

	}

}
