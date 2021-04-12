package utils

import (
	"math"
	"strings"
)

var num2char = "0123456789abcdefghijklmnopqrstuvwxyz"

const InviteCode = 100000000

func NumToBHex(num, n int64) string {
	//num = InviteCode + num
	num_str := ""
	for num != 0 {
		yu := num % n
		num_str = string(num2char[yu]) + num_str
		num = num / n
	}
	return strings.ToUpper(num_str)
}

func BHex2Num(str string, n int) int64 {
	str = strings.ToLower(str)
	v := 0.0
	length := len(str)
	for i := 0; i < length; i++ {
		s := string(str[i])
		index := strings.Index(num2char, s)
		v += float64(index) * math.Pow(float64(n), float64(length-1-i)) // 倒序
	}
	return int64(v)
}
