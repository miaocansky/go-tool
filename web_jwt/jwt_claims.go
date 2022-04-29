package webjwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserInfo map[string]interface{}
	jwt.StandardClaims
}
