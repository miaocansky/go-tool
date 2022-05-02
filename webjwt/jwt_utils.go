package webjwt

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type WebJwtRealize struct {
	SigningKey  string
	ExpiresTime int64
}

func NewWebJwt(signingKey string, expiresTime int64) *WebJwtRealize {
	webJwt := &WebJwtRealize{
		SigningKey:  signingKey,
		ExpiresTime: expiresTime,
	}
	return webJwt
}

/**
更加用户信息获取token
*/
func (webJwt *WebJwtRealize) GetToken(userInfo map[string]interface{}) (string, error) {
	claims := Claims{
		UserInfo: userInfo,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,               // 签名生效时间
			ExpiresAt: time.Now().Unix() + webJwt.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "web",                                  // 签名的发行者
		},
	}
	token, error := webJwt.createToken(claims)
	return token, error
}

/**
创建token
*/
func (webJwt *WebJwtRealize) createToken(claims Claims) (string, error) {
	newWithToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, error := newWithToken.SignedString(webJwt.getSigningKey())
	return token, error

}

/**
解析token
*/
func (webJwt *WebJwtRealize) ParseToken(tokenString string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return webJwt.getSigningKey(), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, err
			} else {
				return nil, err
			}
		}
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	} else {
		return nil, err
	}
	return nil, err
}

/**
获取加密串
*/
func (webJwt *WebJwtRealize) getSigningKey() []byte {
	var jwtSigningKey = []byte(webJwt.SigningKey)
	return jwtSigningKey
}
