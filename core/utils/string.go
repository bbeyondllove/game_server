package utils

import (
	"math/rand"
	"strconv"
	"time"
	"unicode"
)

// 字符首字母是否大写
func IsStartUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

//把字符串变成首字母小写
func toFirstLowwer(str string) string {
	if !IsStartUpper(str) {
		return str
	}

	var retStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			vv[i] += 32 // string的码表相差32位
			retStr += string(vv[i])
		} else {
			retStr += string(vv[i])
		}
	}
	return retStr

}


//
//生成随机字符串
func GetRandString(generate_num int) string {
	ret := ""
	src := "0123456789ABCDEFGHIJKLMLOPQRSTUVWXYZ"
	for i := 0; i < generate_num; i++ {
		idx := rand.Intn(len(src))
		ret += string(src[idx])

	}
	return ret
}

// Strval 获取变量的字符串值
func Strval(value interface{}) string {
	ret := ""
	if value == nil {
		return ret
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		ret = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		ret = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		ret = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		ret = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		ret = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		ret = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		ret = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		ret = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		ret = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		ret = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		ret = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		ret = strconv.FormatUint(it, 10)
	case string:
		ret = value.(string)
	case []byte:
		ret = string(value.([]byte))
	case time.Time:
		ret = Time2Str(value.(time.Time))

	default:
		ret = ""
	}

	return ret
}
