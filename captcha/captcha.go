package captcha

import (
	"bytes"
	"github.com/dchest/captcha"
	"github.com/miaocansky/go-tool/util/file"
	"os"
	"strconv"
	"time"
)

type CaptchaTool struct {
	Path string
}

//
//  NewCaptchaTool
//  @Description: 使用基础Tool
//  @return *CaptchaTool
//
func NewCaptchaTool() *CaptchaTool {
	return &CaptchaTool{}

}

//
//  NewPathCaptchaTool
//  @Description: 使用指定图片存储位置
//  @param path
//  @return *CaptchaTool
//
func NewPathCaptchaTool(path string) *CaptchaTool {
	return &CaptchaTool{
		Path: path,
	}

}

//
//  GetCaptcha
//  @Description:生成验证码图片或者音频文件
//  @param captchaType
//  @return string
//  @return string
//  @return error
//
func (captchaTool *CaptchaTool) GetCaptcha(captchaType string) (string, string, error) {
	path, err := captchaTool.GetPath()
	if err != nil {
		return "", "", err
	}
	captchaId := captcha.New()
	ext := ".png"
	if captchaType == "audio" {
		ext = ".wav"
	}
	unix := time.Now().Unix()
	itoa := strconv.FormatInt(unix, 10)

	path = path + "captchaId_" + itoa + ext
	f, _ := os.Create(path)

	d := captcha.RandomDigits(6)

	switch captchaType {
	case "img":
		image := captcha.NewImage(captchaId, d, captcha.StdWidth, captcha.StdHeight)
		image.WriteTo(f)
		return captchaId, path, nil
	case "audio":
		audio := captcha.NewAudio(captchaId, d, "zh")
		audio.WriteTo(f)
		return captchaId, path, nil

	default:
		return captchaId, "", captcha.ErrNotFound

	}

}

//
//  GetCaptchaId
//  @Description: 获取验证码 captchaId
//  @return string
//
func (captchaTool *CaptchaTool) GetCaptchaId() string {
	captchaId := captcha.New()
	return captchaId
}

//
//  GetImageByte
//  @Description: 输出验证码图片或者视频
//  @param captchaId
//  @param ext
//  @param lang
//  @param width
//  @param height
//  @return []byte
//  @return error
//
func (captchaTool *CaptchaTool) GetImageByte(captchaId, ext, lang string, width, height int) ([]byte, error) {
	var content bytes.Buffer
	captcha.Reload(captchaId)
	switch ext {
	case "img":
		captcha.WriteImage(&content, captchaId, captcha.StdWidth, captcha.StdHeight)
	case "audio":
		captcha.WriteAudio(&content, captchaId, lang)
	default:
		return nil, captcha.ErrNotFound
	}
	return content.Bytes(), nil
}

func (captchaTool *CaptchaTool) GetPath() (string, error) {
	DirPath, err := file.GetPath()
	if err != nil {
		return "", err
	}
	path := "upload/captcha/"
	uploadParth := DirPath + "/" + path
	if captchaTool.Path != "" && len(captchaTool.Path) > 0 {
		uploadParth = captchaTool.Path
	} else {

	}
	exists, _ := file.PathExists(uploadParth)
	if !exists {
		//  创建目录
		_ = file.MakeDir(uploadParth)
	}

	return uploadParth, nil
}

//
//  CheckedCaptcha
//  @Description: 验证验证码
//  @param captchaId
//  @param val
//  @return bool
//
func (captchaTool *CaptchaTool) CheckedCaptcha(captchaId, val string) bool {
	return captcha.VerifyString(captchaId, val)
}
