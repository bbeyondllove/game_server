package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"game_server/core/base"
	"io"
	"sort"
	"strings"
)

//生成签名值
func GenTonken(dataMap map[string]interface{}) string {
	srcData := ""

	tmpSlice := make([]string, 0)
	for k := range dataMap {
		if k != "sign" {
			tmpSlice = append(tmpSlice, k)
		}
	}
	sort.Strings(tmpSlice)

	for _, v := range tmpSlice {
		value := Strval(dataMap[v])
		if value == "" {
			continue
		}
		srcData = srcData + v + "=" + value + "&"
	}
	srcData = srcData + "key=" + base.Setting.Base.SecurityKey
	//fmt.Printf("srcData=%+v\n", srcData)
	genToken := strings.ToUpper(HmacSha256(srcData, base.Setting.Base.SecurityKey))

	return genToken
}

//HmacSha256算法
func HmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

//校验token参数
func CheckEequest(token string, param []string) bool {
	if param == nil {
		return false
	}
	var str string
	for key, item := range param {
		if len(param) >= 1 && (key+1) >= len(param) {
			str += item
		} else {
			str += item + "|"
		}
	}
	//fmt.Println(SHA256(str))
	if token != SHA256(str) {
		return false
	}
	return true
}

/**
 * 对字符串进行SHA256哈希
 * @param data string 要加密的字符串
 */
func SHA256(data string) string {
	t := sha256.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
