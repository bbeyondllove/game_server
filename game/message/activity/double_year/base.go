package double_year

import (
	"encoding/json"
	"errors"
	"game_server/core/logger"
	"game_server/core/utils"
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

type Base struct {
	ResMsg *ResponseMsg
}

// generateResponseMessage 生成返回返回给客户端格式消息.
func (b *Base) GenerateResponseMessage(opCode uint16, message interface{}) *utils.Packet {
	rsp := &utils.Packet{}
	rsp.Initialize(opCode)
	rsp.WriteData(message)
	return rsp
}

// ParseArgsToMap 解析请求参数.
func (b *Base) ParseArgsToMap(requestMsg *utils.Packet) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(requestMsg.Bytes(), &m)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		return nil, errors.New(StatusCodeMessage[ActiveDoubleYearArgsError])
	}

	// 解析token.
	token, ok := m["token"]
	if !ok || token.(string) == "" {
		logger.Errorf("require token")
		return nil, errors.New(StatusCodeMessage[ActiveDoubleYearArgsError])
	}

	flag, _, userInfo := utils.GetUserByToken(token.(string))
	if !flag {
		logger.Errorf("token is illegal")
		return nil, errors.New(StatusCodeMessage[ActiveDoubleYearIllegalToken])
	}

	userId, ok := userInfo["user_id"]
	if !ok || userId == "" {
		logger.Errorf("token is illegal")
		return nil, errors.New(StatusCodeMessage[ActiveDoubleYearIllegalToken])
	}
	m["userId"] = userId

	return m, nil
}

// NewBase 实例化Base结构体.
func NewBase() *Base {
	return &Base{ResMsg: &ResponseMsg{}}
}
