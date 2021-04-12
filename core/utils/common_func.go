package utils

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io"
)

//公共方法

//是否为密码规定字符
func IsHaveInvalidPwdChar(strChar string) bool {
	bReuslt := false
	ilen := len(strChar)

	for i := 0; i < ilen; i++ {
		if (strChar[i] < '0' || strChar[i] > '9') && (strChar[i] < 'A' || strChar[i] > 'Z') && (strChar[i] < 'a' || strChar[i] > 'z') {
			bReuslt = true
			break
		}
	}
	return bReuslt
}

/**
 * 对字符串进行MD5哈希
 * @param data string 要加密的字符串
 */
func MD5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
