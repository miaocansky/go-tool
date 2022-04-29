package email

import "testing"

func TestConEmail_Send(t *testing.T) {

	//"smtp.qq.com"
	//"587"
	//"1535474309@qq.com"
	//"rugptsndrbomfifd"
	email := NewGoEmail("smtp.qq.com", 587, "1535474309@qq.com", "rugptsndrbomfifd", EMAIL_CONTENT_TYPE_HTML)

	msg := Msg{
		HeaderTo:    "miaocansky@163.com",
		HeaderTitle: "测试",
		Content:     "<h1>来自Test</h1>",
	}
	err := email.Send(msg)
	t.Log(err)

}
