package utils

import (
	"encoding/base64"
	"encoding/json"
	"game_server/db"
	"game_server/game/proto"
	"strings"

	"game_server/core/logger"
)

//get jwt payload
func GetPayLoad(token string) (bool, *proto.PayLoad) {
	var ret proto.PayLoad
	jwtAry := strings.Split(token, ".")
	if len(jwtAry) != 3 {
		return false, nil
	}

	decodeBytes, err := base64.RawURLEncoding.DecodeString(jwtAry[1])
	if err != nil {
		return false, nil
	}

	err = json.Unmarshal(decodeBytes, &ret)
	if err != nil {
		return false, nil
	}
	return true, &ret
}

func GetUserByToken(token string) (bool, *proto.PayLoad, map[string]string) {
	flag, payLoad := GetPayLoad(token)
	if !flag {
		logger.Errorf("token error", token)
		ShowStack()
		return false, nil, nil
	}

	//token 检验
	user_info := db.RedisMgr.HGetAll(payLoad.UserId)
	if user_info != nil && user_info["token"] != token {
		logger.Errorf("GetUserByToken user_info['token']=,requestMsg['token']=", user_info["token"], token)
		ShowStack()
		return false, nil, nil
	}
	return true, payLoad, user_info
}
