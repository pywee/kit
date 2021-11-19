package funcs

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// ToString reflect.Value 数据类型转换
func ToString(v reflect.Value, strFormat bool) string {
	kind := v.Kind().String()
	if kind == "bool" {
		if v.Bool() {
			return "true"
		}
		return "false"
	}
	if kind == "string" {
		if strFormat {
			return formatString(v.String())
		}
		return v.String()
	}
	if kind == "slice" {
		return "[]slice"
	}
	if IsKindInt(kind) {
		return strconv.FormatInt(v.Int(), 10)
	}
	if IsKindUint(kind) {
		return strconv.FormatUint(v.Uint(), 10)
	}
	if IsKindFloat(kind) {
		return fmt.Sprintf("%.f", v.Float())
	}
	return ""
}

// IsInt 是否为整型基础类型
// 包含有符号、无符号整型、浮点型
func IsInt(k string) bool {
	return IsKindInt(k) || IsKindFloat(k) || IsKindUint(k)
}

// IsKindInt 判断当前 reflect 类型是否为整型
func IsKindInt(k string) bool {
	return k == "int64" || k == "int32" || k == "int" || k == "int8" || k == "int16"
}

// IsKindUint 判断当前 reflect 类型是否为非负数整型
func IsKindUint(k string) bool {
	return k == "uint64" || k == "uint8" || k == "byte" || k == "uint32" || k == "uint16" || k == "uint" || k == "rune"
}

// IsKindFloat 判断当前 reflect 类型是否为整型
func IsKindFloat(k string) bool {
	return k == "float64" || k == "float" || k == "float32"
}

// Name2Case 驼峰转下划线
func Name2Case(str string) string {
	m := make([]rune, 0, 10)
	for k, v := range str {
		if k == 0 && IsWordEn(v) {
			v += 32
		} else if IsWordEn(v) {
			m = append(m, '_')
			v += 32
		}
		m = append(m, v)
	}
	return string(m)
}

// Case2Camel .
func Case2Camel(name []byte) []byte {
	name = bytes.Replace(name, []byte{95}, []byte{32}, -1)
	name = bytes.Title(name)
	return bytes.Replace(name, []byte{32}, nil, -1)
}

// Case2CamelS 下划线转驼峰
func Case2CamelS(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func Camel2Case(name string) string {
	str := ""
	for i := 0; i < len(name); i++ {
		k := name[i]
		if k == '_' || i == 0 {
			str += strings.ToUpper(string(k))
		} else {
			str += string(k)
		}
	}
	return str
}

// LoadFile 业务相关 载入内容文件数据
func LoadFile(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return src, nil
}

func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// IsDigit2 判断字符串是否为数字字符串
func IsDigit2(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

// GetFiles 获取指定目录下所有文件
func GetFiles(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + strings.ToLower(fi.Name())
			s = append(s, fullName)
		}
	}
	return s, nil
}

func ToMd5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func StringInArray(arr []string, n string) bool {
	for _, v := range arr {
		if n == v {
			return true
		}
	}
	return false
}

// 生成大写随机字符串
// A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// formatString 处理字符串数据 如果没有双引号，则需要加入双引号
func formatString(str string) string {
	if len(str) > 1 && str[0] == '"' && str[1] == '"' {
		return str
	}
	str = strings.Replace(str, "\"", "\\\"", -1)
	return `"` + str + `"`
}

// RangeRand 生成区间[m, n]的安全随机数
func RangeRand(min, max int) int {
	if min > max {
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max) + min
}
