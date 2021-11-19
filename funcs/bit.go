package funcs

import "fmt"

// AddArrToBit 将数组放入bigmap
func AddArrToBit(arr []int64) int64 {
	var i int64
	for _, v := range arr {
		i = AddBit(i, v)
	}
	return i
}

// AddBit 向位图中新增一个数字
func AddBit(bit, num int64) int64 {
	return bit | (0x01 << (num % 64))
}

// GetBit .
func GetBit(bit, num int64) bool {
	return (bit & (0x01 << (num % 64))) != 0
}

// RangeBit 找出当前整型中所有的数字
func RangeBit(bit int64) []int64 {
	arr := make([]int64, 0, 10)
	var i int64
	for i = 0; i < 64; i++ {
		if bit&(0x01<<(i%64)) != 0 {
			arr = append(arr, i)
		}
	}
	return arr
}

// RangeBitSum 取 bit 数字的总和
func RangeBitSum(bit int64) int64 {
	var (
		i   int64
		sum int64
	)
	for i = 0; i < 64; i++ {
		if bit&(0x01<<(i%64)) != 0 {
			sum += i
		}
	}
	return sum
}

// RangeBitStr 找出当前整型中所有的数字
func RangeBitStr(bit int) string {
	str := ""
	for i := 0; i < 64; i++ {
		if bit&(0x01<<(i%64)) != 0 {
			str += fmt.Sprintf("%d ", i)
		}
	}
	return str
}
