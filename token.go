package goxi_v2

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type TokenLogic struct {
	Salt string
}

func NewTokenLogic(salt string) *TokenLogic {
	return &TokenLogic{
		Salt: salt,
	}
}

type CustomClaims struct {
	Uid string `json:"uid"`
	jwt.RegisteredClaims
}

// GenerateToken 生成jwt
func (l *TokenLogic) GenerateToken(uid string) (tokenStr string, err error) {
	claim := CustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * time.Duration(180))), // 过期时间180天
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                          // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                          // 生效时间
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法
	tokenStr, err = token.SignedString([]byte(l.Salt))
	return tokenStr, err
}

func (l *TokenLogic) Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(l.Salt), nil
	}
}

func (l *TokenLogic) AnalyzeToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, l.Secret())
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token not active yet")
			} else {
				return nil, errors.New("couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("couldn't handle this token")
}
