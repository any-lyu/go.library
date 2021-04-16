package jwt

import (
	"testing"

	"github.com/siddontang/go/log"
)

func TestTokenGenerate(t *testing.T) {
	tokenString, err := TokenGenerate("234")
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("tokenString = ", tokenString)

	token, err := TokenParse(tokenString)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("uid = ", token.UID)
}
