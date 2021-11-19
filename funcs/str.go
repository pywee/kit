package funcs

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// CheckStringType 取得字符串类型
// 0-无法识别 1-整型 2-浮点型 3-字符串 4-布尔型(true) 5-布尔型(false) 6-变量
func CheckStringType(src string, isStr bool) (int8, string) {
	l := len(src)
	if l == 0 {
		return 0, ""
	}

	if l >= 4 {
		if src == "true" {
			return 4, "true"
		}
		if src == "false" {
			return 5, "false"
		}
	}

	// 数值判断
	if i := IsDigit(src); i == 1 {
		return 1, src
	} else if i == 2 {
		return 2, src
	}

	// 变量判断
	if ok, err := regexp.MatchString(`^[a-zA-Z\_\.]{1}[0-9a-zA-Z\.\_]+$`, src); err != nil {
		return 0, ""
	} else if ok {
		return 6, src
	}

	// 字符串判断
	if str, f, err := CheckAndGetString(src, isStr); err != nil {
		return 0, err.Error()
	} else if f {
		return 3, str
	}

	// 如果未能识别 当字符串处理
	return 3, src
}

// IsDigit 判断字符串是否为数字字符串
// f: 0-其他类型，如字符串 1-整型 2-浮点
func IsDigit(str string) int8 {
	var f int8 = 1
	for _, x := range str {
		if x == '.' {
			f = 2
			continue
		}
		if !unicode.IsDigit(x) {
			return 0
		}
	}
	return f
}

func SplitStrToInt(str string) []int64 {
	var arrInt = make([]int64, 0, 10)
	str = strings.TrimRight(str, ",")
	arr := strings.Split(str, ",")
	for _, v := range arr {
		m, _ := strconv.ParseInt(v, 10, 64)
		arrInt = append(arrInt, m)
	}
	return arrInt
}

// CheckAndGetString 检查并返回字符串
// str 如果等于 true ，表示这个是从变量中解析出来的数据
func CheckAndGetString(src string, str bool) (string, bool, error) {
	if slen := len(src); slen > 0 {
		if src[0] == '"' && src[slen-1] == '"' {
			if strings.Contains(src, "\"") || strings.Contains(src, "\\") {
				yc := strings.Count(src, "\"")
				xc := strings.Count(src, "\\")
				if c := yc + xc; c%2 != 0 {
					return "", false, errors.New("字符串格式不正确")
				}
				if src[0] == '"' && src[slen-1] == '"' {
					src = src[1 : slen-1]
				}
				src = strings.Replace(src, "\\\"", "\"", -1)
				src = strings.Replace(src, "\\\\", "\\", -1)
			}
			return src, true, nil
		}
	}

	// 深入判断是否为字符串
	if str {
		for _, v := range src {
			if (v < 40 || v > 58) && v != 37 { // 没有将运算符 ^ 计算在内
				return src, true, nil
			}
		}
	}

	return src, false, nil
}

func IsWordEn(s rune) bool {
	return s >= 65 && s <= 90
}
