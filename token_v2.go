package goxi_v2

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v4"
)

type JWTToken struct {
}

func NewJWTTokenLogic() *JWTToken {
	return &JWTToken{}
}

// GenToken 生成jwt
func (j *JWTToken) GenToken(secretKey string, iat, seconds int64, payload string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// GenRefreshToken 生成刷新token
func (j *JWTToken) GenRefreshToken() (string, error) {
	b := make([]byte, 16) // 16字节（128位）随机数据
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
