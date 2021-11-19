package funcs

import (
	"encoding/json"
	"reflect"
	"strings"
)

// CopyStruct 复制结构体
func CopyStruct(src, dst interface{}) error {
	bufA, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bufA, dst)
}

// GetField 从一个结构体中查找指定数据 返回第一个被找到的数据
func GetField(values []*reflect.Value, fieldName string) *reflect.Value {
	fields := strings.Split(fieldName, ".")
	for _, value := range values {
		if value == nil || !value.IsValid() {
			continue
		}
		for i := 0; i < len(fields); i++ {
			field := fields[i]
			// fmt.Println(values, field)
			if value.Type().String() == "*model."+Case2CamelS(field) {
				if i+1 < len(fields) {
					if fdata := value.Elem().FieldByName(Case2CamelS(fields[i+1])); fdata.IsValid() {
						return &fdata
					}
				}
				return nil
			}
		}
	}
	return nil
}
