package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/any-lyu/go.library/errors"
	"strings"
	"time"
)

// TokenGenerateClaims jwt token add uid
type TokenGenerateClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

var hmacSampleSecret = []byte("124")

// TokenGenerate jwt generate token
func TokenGenerate(uid string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenGenerateClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(120) * time.Hour).Unix(),
			Issuer:    "system",
		},
	})
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(hmacSampleSecret)
}

// TokenParse jwt parse token
func TokenParse(tokenString string) (token *TokenGenerateClaims, err error) {
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return nil, errors.Wrapf(errors.ErrToken, "TokenParse:%s", "token invalid")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token = new(TokenGenerateClaims)
	claims, err := jwt.ParseWithClaims(tokenString, token, func(token *jwt.Token) (i interface{}, e error) {
		return hmacSampleSecret, nil
	})
	if err != nil {
		return nil, errors.Wrapf(errors.ErrToken, "TokenParse:%s", err.Error())
	}
	if token, ok := claims.Claims.(*TokenGenerateClaims); ok && claims.Valid {
		return token, nil
	}
	if !claims.Valid {
		return nil, errors.Wrapf(errors.ErrToken, "TokenParse:%s", "valid")
	}
	return nil, errors.Wrapf(errors.ErrToken, "TokenParse")
}
