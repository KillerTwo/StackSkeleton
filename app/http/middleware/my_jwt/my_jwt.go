package my_jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"goskeleton/app/global/my_errors"
	"time"
)

// CreateMyJWT 使用工厂创建一个 JWT 结构体
func CreateMyJWT(signKey string) *JwtSign {
	if len(signKey) <= 0 {
		signKey = "goskeleton"
	}
	return &JwtSign{
		[]byte(signKey),
	}
}

// JwtSign 定义一个 JWT验签 结构体
type JwtSign struct {
	SigningKey []byte
}

// CreateToken 生成一个token
func (j *JwtSign) CreateToken(claims CustomClaims) (string, error) {
	// 生成jwt格式的header、claims 部分
	tokenPartA := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 继续添加秘钥值，生成最后一部分
	return tokenPartA.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JwtSign) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New(my_errors.ErrorsTokenMalFormed)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New(my_errors.ErrorsTokenExpired)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New(my_errors.ErrorsTokenNotActiveYet)
			} else {
				return nil, errors.New(my_errors.ErrorsTokenInvalid)
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New(my_errors.ErrorsTokenInvalid)
}

// RefreshToken 更新token
func (j *JwtSign) RefreshToken(tokenString string, extraAddSeconds int64) (string, error) {

	if CustomClaims, err := j.ParseToken(tokenString); err == nil {
		CustomClaims.ExpiresAt = time.Now().Unix() + extraAddSeconds
		return j.CreateToken(*CustomClaims)
	} else {
		return "", err
	}
}
