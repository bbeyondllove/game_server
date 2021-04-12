// 调用本地生活接口基础结构体.
package local_life

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/proto"
	"sort"
	"strings"
	"time"
)

// ResponseMsg 返回消息格式结构体.
type ResponseMsg struct {
	// Code 状态码.
	Code int `json:"code"`
	// Msg 状态消息.
	Msg string `json:"msg"`
	// Data 返回的数据.
	Data map[string]interface{} `json:"data"`
}

// Base Base struct.
type Base struct {
	resMsg *ResponseMsg
}

// Request 向本地生活服务发起接口调用.
func (b *Base) Request(url string, requestMethod string, requestData map[string]interface{}) (map[string]interface{}, error) {
	logger.Infof("Request api[%v], args[%v]:", url, requestData)
	var requestUrl string
	if base.Setting.Life.Debug {
		requestUrl = base.Setting.Life.LocalIp + ":" + base.Setting.Life.LocalPort + url
	} else {
		requestUrl = base.Setting.Life.HttpUrl + url
	}
	result := make(map[string]interface{})
	var data, tmpRequestArgs []byte
	var err error

	// 拼接sign
	requestData["sign"] = strings.ToUpper(b.generateSign(requestData, base.Setting.Life.Secret))

	done := make(chan struct{})
	go func() {
		defer close(done)
		switch requestMethod {
		case "POST":
			tmpRequestArgs, err = json.Marshal(requestData)
			if err != nil {
				logger.Errorf("http post request with json error, err=", err.Error())
				return
			}
			data, err = utils.HttpPost(requestUrl, string(tmpRequestArgs), proto.JSON)
		case "GET":
			data, err = utils.HttpGet(requestUrl, requestData)
		case "DELETE":
			data, err = utils.HttpDelete(requestUrl, requestData)
		default:
			data, err = utils.HttpGet(requestUrl, requestData)
		}
		if err != nil {
			logger.Errorf("http request url[%v] error, err=", url, err.Error())
			return
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(base.Setting.Life.Timeout)*time.Second)
	select {
	// 超时处理.
	case <-ctx.Done():
		logger.Errorf("http request url[%v] timeout", url)
		return nil, errors.New("http request url[" + url + "] timeout")
		// 处理返回的数据.
	case <-done:
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(data, &result)
		if err != nil {
			logger.Errorf("json-Unmarshal url err: %v", err)
			return nil, err
		}

		if code, _ := result["code"]; int(code.(float64)) != LocalLifeSuccess {
			logger.Errorf("http request url[%v] failed, code: %v", url, code)
			return nil, errors.New("http request url[" + url + "] statusCode is not 200")
		}
		// 这里兼容有些接口返回值是没有(data字段)数据的.
		if resultData, ok := result["data"]; ok && resultData != nil {
			return result["data"].(map[string]interface{}), err
		}
		return nil, err
	}
}

// CheckArgs 检查公共参数.
// 目前公共参数为token.
func (b *Base) checkPubArgs(args map[string]interface{}) (err error, statusCode int) {
	mustArgs := b.mustPublicArgs()
	for _, v := range mustArgs {
		switch v {
		case "token":
			token, ok := args[v]
			if !ok || token == "" {
				return errors.New(statusCodeMessage[LocalLifeRequireToken]), LocalLifeRequireToken
			}
			flag, _, _ := utils.GetUserByToken(token.(string))
			if !flag {
				return errors.New(statusCodeMessage[LocalLifeTokenError]), LocalLifeTokenError
			}
		// 检查剩余其它参数.
		default:
			value, ok := args[v]
			if !ok || value == "" {
				return errors.New(fmt.Sprintf("missing public [%v] parameters", value)), LocalLifeRequireArgs
			}
		}
	}
	return nil, LocalLifeVerifyOk
}

// ParseArgs 解析请求参数.
func (b *Base) parseArgsToMap(requestMsg *utils.Packet) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(requestMsg.Bytes(), &m)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		return nil, errors.New(statusCodeMessage[LocalLifeArgumentError])
	}
	return m, nil
}

// mustPublicArgs 必要公共参数.
func (b *Base) mustPublicArgs() []string {
	return []string{
		"token", // 登入base系统后返回的token.
	}
}

// generateResponseMessage 生成返回返回给客户端格式消息.
func (b *Base) generateResponseMessage(opCode uint16, message interface{}) *utils.Packet {
	rsp := &utils.Packet{}
	rsp.Initialize(opCode)
	rsp.WriteData(message)
	return rsp
}

// beforeRequest 请求本地生活接口之前检查.
func (b *Base) beforeRequest(requestMsg *utils.Packet, opoRespCmd uint16, businessArgs []string) (map[string]interface{}, *utils.Packet) {
	// 解析参数.
	requestArgs, err := b.parseArgsToMap(requestMsg)
	if err != nil {
		b.resMsg.Code = LocalLifeArgumentError
		b.resMsg.Msg = statusCodeMessage[LocalLifeArgumentError]
		return requestArgs, b.generateResponseMessage(opoRespCmd, b.resMsg)
	}

	// 检查公共参数.
	err, statusCode := b.checkPubArgs(requestArgs)
	if err != nil {
		b.resMsg.Code = statusCode
		b.resMsg.Msg = err.Error()
		return requestArgs, b.generateResponseMessage(opoRespCmd, b.resMsg)
	}

	// 检查各个接口必要业务参数.
	for _, k := range businessArgs {
		if v, ok := requestArgs[k]; !ok || v == "" {
			b.resMsg.Code = LocalLifeRequireArgs
			b.resMsg.Msg = fmt.Sprintf(statusCodeMessage[LocalLifeRequireArgs], k)
			return requestArgs, b.generateResponseMessage(opoRespCmd, b.resMsg)
		}
	}

	return requestArgs, nil
}

// generateSign 生成签名.
// 采用md5方式.
func (b *Base) generateSign(data map[string]interface{}, secret string) string {
	var sign string
	var key []string

	// 按照参数数组的key升序排序
	for k := range data {
		key = append(key, k)
	}
	sort.Strings(key)

	// 生成签名
	for _, k := range key {
		// sign参数不参于签名
		if k == "sign" {
			continue
		}
		sign += k + "=" + utils.Strval(data[k]) + "&"
	}
	md5HashInBytes := md5.Sum([]byte(sign + "secret=" + secret))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])

	return md5HashInString
}

// NewBase 实例化Base结构体.
func NewBase() *Base {
	return &Base{resMsg: &ResponseMsg{
		Code: 0,
		Msg:  "",
		Data: nil,
	}}
}
