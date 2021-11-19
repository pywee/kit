package types

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrorParams = ErrJSON(1003, "不支持的参数值，请检查")
)

type JSON struct {
	Code int16       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ErrJSON(code int16, msg string, data ...interface{}) error {
	j := &JSON{Code: code, Msg: msg}
	if len(data) > 0 {
		j.Data = data[0]
	}
	buf, _ := json.Marshal(j)
	return errors.New(string(buf))
}

// WithError 在 {code:-1, msg:""} 的原有基本上为msg字段新增error msg
func WithError(err, msg error) error {
	var (
		ret JSON
		str = err.Error()
	)
	json.Unmarshal([]byte(str), &ret)
	if ret.Code > 0 && msg != nil {
		ret.Msg += "; reason: " + msg.Error()
		b, _ := json.Marshal(ret)
		return errors.New(string(b))
	}
	return nil
}

func GetErrorJSON(err error, data interface{}) []byte {
	if err == nil {
		b, _ := json.Marshal(&JSON{
			Msg:  "success",
			Data: data,
		})
		return b
	}

	var (
		ret JSON
		str = err.Error()
	)
	if strings.Contains(str, `"code":`) {
		if strings.Contains(str, `"data":`) {
			json.Unmarshal([]byte(str), &ret)
			ret.Data = data
			b, _ := json.Marshal(ret)
			return b
		}
		return []byte(str)
	}

	b, _ := json.Marshal(&JSON{Code: -1, Msg: err.Error()})
	return b
}
