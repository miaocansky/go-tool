package strutil

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"strings"
)

//
//  Md5
//  @Description: md5 加密
//  @param str
//  @return string
//
func Md5(str string) string {
	bytes := []byte(str)
	h := md5.New()
	h.Write(bytes)
	return hex.EncodeToString(h.Sum(nil))
}

//
//  BasBase64Encodee64
//  @Description: base64加密
//  @param str
//  @return string
//

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//
//  Base64Decode
//  @Description: base64解密
//  @param str
//  @return string
//
func Base64Decode(str string) string {
	decodeString, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return ""
	}
	return string(decodeString)
}

//
//  URLEncode
//  @Description: 加密url
//  @param s
//  @return string
//
func URLEncode(s string) string {
	if pos := strings.IndexRune(s, '?'); pos > -1 { // escape query data
		return s[0:pos+1] + url.QueryEscape(s[pos+1:])
	}

	return s
}

//
//  URLDecode
//  @Description: url解密
//  @param s
//  @return string
//
func URLDecode(s string) string {
	if pos := strings.IndexRune(s, '?'); pos > -1 { // un-escape query data
		qy, err := url.QueryUnescape(s[pos+1:])
		if err == nil {
			return s[0:pos+1] + qy
		}
	}

	return s
}
