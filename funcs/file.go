package funcs

import (
	"io/ioutil"
	"os"
)

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

// WriteFileString 写入数据到文件
func WriteFileString(fileName, data string) error {
	return WriteFile(fileName, []byte(data))
}

// WriteFile 写入数据到文件
func WriteFile(fileName string, buf []byte) error {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(buf)
	return err
}
