package email

import (
	"gopkg.in/gomail.v2"
)

//  邮件消息类型 html 与文本
const EMAIL_CONTENT_TYPE_HTML string = "text/html"
const EMAIL_CONTENT_TYPE_TEXT string = "text/plain"

//
//  ConEmail
//  @Description: ConEmail struct
//
type ConEmail struct {
	con       *gomail.Dialer
	EmailType string
	Host      string
	port      int
	UserName  string
	Password  string
}

//
//  Msg
//  @Description: 消息体
//
type Msg struct {
	HeaderTo    string
	HeaderCc    string
	HeaderTitle string
	Content     string
	AttachFile  string
}

//
//  NewGoEmail
//  @Description: 构建email
//  @param host 服务地址
//  @param port 服务端口
//  @param userName 用户名
//  @param password 密码
//  @param emailType 邮件类型
//  @return *ConEmail
//
func NewGoEmail(host string, port int, userName string, password string, emailType string) *ConEmail {
	//"smtp.qq.com"
	//"587"
	//"1535474309@qq.com"
	//"rugptsndrbomfifd"
	dialer := gomail.NewDialer(host, port, userName, password)
	return &ConEmail{
		con:       dialer,
		EmailType: emailType,
		Host:      host,
		port:      port,
		UserName:  userName,
		Password:  password,
	}
}

//
//  Send
//  @Description: 发送邮件
//  @receiver em
//  @param msg 邮件内容
//  @return error
//
func (em *ConEmail) Send(msg Msg) error {
	message := gomail.NewMessage()
	message.SetHeader("From", em.UserName)
	message.SetHeader("To", msg.HeaderTo)
	if msg.HeaderCc != "" {
		message.SetHeader("Cc", msg.HeaderCc)
	}

	message.SetHeader("Subject", msg.HeaderTitle)
	message.SetBody(em.EmailType, msg.Content)
	if msg.AttachFile != "" {
		message.Attach(msg.AttachFile)
	}
	// 关闭SSL协议认证
	//em.con.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return em.con.DialAndSend(message)

}
