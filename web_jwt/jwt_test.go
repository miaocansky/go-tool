package webjwt

import "testing"

func TestNewWebJwt(t *testing.T) {

	m := make(map[string]interface{})
	m["username"] = "test"
	m["id"] = 111
	token, err := GetWebJwt().GetToken(m)
	t.Log(token)
	t.Log(err)
}
func GetWebJwt() *WebJwt {
	signingKey := "aaaa"
	var expiresTime int64 = 3600
	jwt := NewWebJwt(signingKey, expiresTime)
	return jwt

}

func TestParseToken(t *testing.T) {
	m := make(map[string]interface{})
	m["username"] = "test"
	m["id"] = 111
	token, err := GetWebJwt().GetToken(m)
	t.Log(token)
	t.Log(err)
	parseToken, err := GetWebJwt().ParseToken(token)
	t.Log(parseToken)
	t.Log(err)

}
