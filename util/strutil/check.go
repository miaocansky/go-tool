package strutil

import "regexp"

//
//  CheckMobile
//  @Description: 验证手机号
//  @param phone
//  @return bool
//
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[123456789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}

// IsNumeric returns true if the given character is a numeric, otherwise false.
func IsNumeric(c byte) bool {
	return c >= '0' && c <= '9'
}
