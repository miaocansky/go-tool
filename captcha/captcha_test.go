package captcha

import (
	"testing"
)

func TestGetCaptcha(t *testing.T) {
	tool := NewCaptchaTool()
	captchaId, path, _ := tool.GetCaptcha("img")
	t.Log(captchaId)
	t.Log(path)
}
func TestCaptchaTool_GetCaptchaId(t *testing.T) {
	tool := NewCaptchaTool()
	captchaId := tool.GetCaptchaId()
	t.Log(captchaId)
}
func TestCaptchaTool_GetImageByte(t *testing.T) {
	tool := NewCaptchaTool()
	captchaId := tool.GetCaptchaId()
	ext := "img"
	//if ext == "img" {
	//	c.Writer.Header().Set("Content-Type", "image/png")
	//} else if ext == "audio" {
	//	c.Writer.Header().Set("Content-Type", "audio/x-wav")
	//}
	//if download == "1" {
	//	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	//
	//}
	imageByte, _ := tool.GetImageByte(captchaId, ext, "zh", 240, 80)
	t.Log(imageByte)
	//c.Writer.Write(imageByte)

}

//
//func TestGetCaptchaVideo(t *testing.T) {
//	path := GetPath()
//	path = path + "11.wav"
//	f, _ := os.Create(path)
//	captchaId := captcha.New()
//	t.Log(captchaId)
//	d := captcha.RandomDigits(6)
//
//	audio := captcha.NewAudio(captchaId, d, "zh")
//	audio.WriteTo(f)
//	t.Log(d)
//}
//
