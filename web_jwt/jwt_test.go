package webjwt

import "testing"

//
//  TestNewWebJwt
//  @Description: token生成
//  @param t
//
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

//
//  TestParseToken
//  @Description: 解析token 获取用户信息
//  @param t
//
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
