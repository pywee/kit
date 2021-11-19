package args

// 2021.10.12
// 接收参数逻辑
// 将 POST、GET 参数的数据反射到对应的 interface{} 结构体内

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ahmek/kit/funcs"
	"github.com/ahmek/kit/types"
)

type (
	// 是否检查 get 参数中的 ts 时间与服务器的间隔
	Config struct {
		CheckArgTs float64
	}
	RetJSON struct {
		Code int16       `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}
	RouteCBK func(*HTTPContext) (interface{}, error)
	ErrCBK   func(*HTTPContext) error
)

var cfg *Config

func SetConfig(c *Config) {
	cfg = c
}

type GetAPIArg struct {
	Ts int64 `json:"ts"`
}

// 解析 http post 请求数据
func (c *HTTPContext) PostArg(arg interface{}) error {
	var (
		err error
		buf []byte
	)
	if buf, err = io.ReadAll(c.r.Body); err != nil {
		return err
	}
	if len(buf) > 0 {
		if err = json.Unmarshal(buf, arg); err != nil {
			if strErr := err.Error(); strings.Contains(strErr, "cannot unmarshal") {
				// 此处代码是有效的，当输出类似以下的报错时，屏蔽掉结构体名称，仅输出：
				// json: cannot unmarshal string into Go struct field amount of type uint8
				if idx := strings.Index(strErr, "Go struct field "); idx != -1 {
					strErr = strErr[:idx+16]
				}
				if idx := strings.Index(err.Error(), "."); idx != -1 {
					strErr += err.Error()[idx+1:]
				}
				err = errors.New(strErr)
				// err = errors.New("参数值类型错误")
			}
			return err
		}
	}

	rvf := reflect.ValueOf(arg).Elem()
	numField := rvf.NumField()
	for i := 0; i < numField; i++ {
		field := rvf.Field(i)
		if field.CanSet() {
			if kind := field.Kind().String(); kind == "string" {
				field.SetString(strings.Trim(field.String(), " "))
			} else if funcs.IsKindInt(kind) {
				field.SetInt(field.Int())
			} else if funcs.IsKindUint(kind) {
				field.SetUint(field.Uint())
			}
		}
		// fmt.Println(rvf.Field(i).Kind(), rvf.Field(i))
	}

	if err = StructCheck(arg); err != nil {
		return err
	}
	return nil
}

// 解析 http get 请求数据
func (c *HTTPContext) GetArg(arg interface{}) error {
	var (
		rv       = reflect.ValueOf(arg)
		rt       = rv.Type()
		cQuery   = c.GetURL().Query()
		numField = rt.Elem().NumField()
	)
	for i := 0; i < numField; i++ {
		rField := rt.Elem().Field(i)
		tag := rField.Tag.Get("json")
		tag = strings.Replace(tag, ",omitempty", "", 1)
		tagValue := strings.TrimSpace(cQuery.Get(tag))
		if len(tagValue) == 0 {
			continue
		}

		// 反射get数据到结构体
		this := rv.Elem().Field(i)
		switch rField.Type.String() {
		case "int64", "int", "int32", "int8", "int16":
			n, err := strconv.ParseInt(tagValue, 10, 64)
			if err != nil {
				return err
			}
			this.SetInt(n)
		case "string":
			this.SetString(tagValue)
		case "byte", "uint64", "uint8", "uint16", "uint32", "uint", "rune":
			n, err := strconv.ParseUint(tagValue, 10, 64)
			if err != nil {
				return err
			}
			this.SetUint(n)
		case "bool":
			n, err := strconv.ParseBool(tagValue)
			if err != nil {
				return err
			}
			this.SetBool(n)
		case "float64", "float32":
			n, err := strconv.ParseFloat(tagValue, 64)
			if err != nil {
				return err
			}
			this.SetFloat(n)
		case "time.Duration":
			n, err := strconv.ParseInt(tagValue, 10, 64)
			if err != nil {
				return err
			}
			this.Set(reflect.ValueOf(time.Duration(n)))
		}
	}

	// 判断必选参数等
	if err := StructCheck(arg); err != nil {
		return err
	}
	return nil
}

// 设置必填项 required:"请输入id"
// 设置检查手机号 phone:"true" phoneErr:"手机号不合法"
// 设置检查邮箱 email:"true" email:"邮箱格式不合法"
// 设置检查最小长度 minLen:"10" minLenErr:"姓名不得小于10个字"
// 设置检查最大长度 maxLen:"12" maxLenErr:"姓名长度不得大于12个字"
// 设置只允许的整型值 int:"1,3" intErr:"分类id只能输入1或3的数字"

// 示例:
// type TestF struct {
// 	Id    int    `json:"id" required:"请输入id"`
// 	Cid   int    `json:"cid" int:"1,3" intErr:"分类id只能输入1或3的数字"`
// 	Phone string `json:"phone" required:"请输入手机号" phone:"true" phoneErr:"手机号不合法"`
// 	Name  string `json:"name" minLen:"10" minLenErr:"姓名不得小于10个字" maxLen:"12" maxLenErr:"姓名长度不得大于12个字"`
// }

// StructCheck 检查结构体是否为空等 用于 post 数据检查
func StructCheck(arg interface{}) error {
	rv := reflect.ValueOf(arg)
	tv := rv.Type()
	for i := 0; i < tv.Elem().NumField(); i++ {
		rvf := rv.Elem().Field(i)
		tvf := tv.Elem().Field(i)
		kind := rvf.Kind().String()
		if required := tvf.Tag.Get("required"); len(required) > 0 {
			if kind == "string" {
				if value := rvf.String(); len(value) == 0 {
					return errors.New(required)
				}
			} else if funcs.IsKindInt(kind) {
				if value := rvf.Int(); value == 0 {
					return errors.New(required)
				}
			} else if funcs.IsKindFloat(kind) {
				if value := rvf.Float(); value == 0 {
					return errors.New(required)
				}
			}
		}

		// 只针对 string 类型
		// 其他类型暂时不提供支持
		if rvf.Kind().String() == "string" && len(rvf.String()) > 0 {
			if enum := tvf.Tag.Get("enum"); len(enum) > 0 {
				enum = strings.TrimSpace(enum)
				arr := strings.Split(enum, ",")
				ok := false
				for _, check := range arr {
					if strings.EqualFold(check, rvf.String()) {
						ok = true
					}
				}
				if !ok {
					return types.ErrorParams
				}
			}
		}

		if kind == "string" {
			tag := tvf.Tag
			if p := tag.Get("phone"); p == "true" {
				if e := tag.Get("phoneErr"); len(e) > 0 {
					phone := rvf.String()
					if len(phone) != 11 {
						return errors.New(e)
					}
					if !isMobile(phone) {
						return errors.New(e)
					}
				}
			}
			if n := tag.Get("minLen"); len(n) > 0 {
				num, _ := strconv.Atoi(n)
				if strings.Count(rvf.String(), "")-1 < num {
					return errors.New(tag.Get("minLenErr"))
				}
			}
			if n := tag.Get("maxLen"); len(n) > 0 {
				num, _ := strconv.Atoi(n)
				if strings.Count(rvf.String(), "")-1 > num {
					return errors.New(tag.Get("maxLenErr"))
				}
			}
		} else if funcs.IsKindInt(kind) {
			if p := tvf.Tag.Get("int"); len(p) > 0 {
				i := rvf.Int()
				p = strings.Replace(p, " ", "", -1)
				arr := strings.Split(p, ",")
				ok := false
				for _, v := range arr {
					if t, _ := strconv.ParseInt(v, 10, 64); t == i {
						ok = true
						break
					}
				}
				if !ok {
					if errstr := tvf.Tag.Get("intErr"); len(errstr) > 0 {
						return errors.New(errstr)
					}
					return errors.New("参数错误")
				}
			}
		} else if funcs.IsKindUint(kind) {
			if p := tvf.Tag.Get("int"); len(p) > 0 {
				ok := false
				i := rvf.Uint()
				p = strings.Replace(p, " ", "", -1)
				arr := strings.Split(p, ",")

				for _, v := range arr {
					if t, _ := strconv.ParseUint(v, 10, 64); t == i {
						ok = true
						break
					}
				}
				if !ok {
					if errstr := tvf.Tag.Get("intErr"); len(errstr) > 0 {
						return errors.New(errstr)
					}
					return errors.New("参数错误")
				}
			}
		}
	}
	return nil
}

func isMobile(mobile string) bool {
	if result, err := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, mobile); err == nil {
		return result
	}
	return false
}
