package access_token_client

import (
	"testing"

	"github.com/k-washi/bss-utils/logger"

	"github.com/k-washi/bss-oauth-api/src/domain/accesstoken"

	"github.com/tj/assert"
)

func TestConnectionTimeOutError(t *testing.T) {

	at := accesstoken.AccessToken{
		AccessToken: "donadonaTest",
		UserID:      "absss",
	}
	logger.Info("Connection error test start")
	err := Check(at)
	assert.Error(t, err)
	logger.Info("Connection error test finish")
}
