package exjwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret []byte

func Auth(username string, uuid string) (t string, err error) {
	now := time.Now()
	// default token exp time is 86400s 60 * 60 * 24
	expSecond := 86400

	if username == "admin" || username == "root" || strings.HasSuffix(username, "bot") {
		expSecond = 86400000 // 1000d
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["uuid"] = uuid
	claims["exp"] = now.Add(time.Duration(expSecond) * time.Second).Unix()
	t, err = token.SignedString(secret)
	if err != nil {
		return "", errors.New("JWT Generate Failure")
	}
	return t, nil
}

func Parse(ts string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("token invalid")
	}
	claim := token.Claims.(jwt.MapClaims)
	return claim, nil
}
