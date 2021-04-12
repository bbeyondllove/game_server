package message

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"

	plug "github.com/syyongx/php2go"
)

func CheckNoticeIsSendedTo(notice_id int, user_id string) bool {
	data, _ := db_service.UserNoticeIns.Get(notice_id, user_id)
	if data.Id != 0 {
		return true
	}
	// fmt.Println(notice_id,user_id,data)
	_, err := db_service.UserNoticeIns.Add(&model.UserNotice{
		NoticeId: notice_id,
		UserId:   user_id,
	})
	if err != nil {
		logger.Error(err)
	}
	return false
}

//检查最新升级公告
func (a *agent) HandleGetUpgradeNotice(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetUpgradeNotice in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_NOTICE_RSP)
	responseMessage := &proto.S2CNotices{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	msg := &proto.C2SNotice{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		SendPacket(a.conn, rsp)
		return
	}
	// 获取公告
	notice, err := db_service.NoticeIns.GetLastNotice()
	if err != nil {
		logger.Errorf("GetUpgradeNotice, err=", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		SendPacket(a.conn, rsp)
		return
	}
	_, PayLoad, _ := utils.GetUserByToken(a.token)

	result := []proto.NoticeInfo{}
	for _, item := range notice {
		if !(item.NoticeType == proto.NOTICE_SYSTEM_UPDATE || item.NoticeType == proto.NOFICE_FORCE) {
			continue
		}
		if item.NoticeType == proto.NOTICE_SYSTEM_UPDATE {
			// 检查是否发送给用户,发送过则不发送
			if CheckNoticeIsSendedTo(int(item.Id), PayLoad.UserId) {
				continue
			}
		} else if item.NoticeType == proto.NOFICE_FORCE {
			if len(msg.Version) > 0 {
				if plug.VersionCompare(item.Version, msg.Version, "<=") {
					continue
				}
			}
		}
		// if len(msg.Version) > 0 {
		// 	if item.NoticeType == proto.NOTICE_SYSTEM_UPDATE && len(item.Version) > 0 {
		// 		if plug.VersionCompare(item.Version, msg.Version, "<") {
		// 			continue
		// 		}
		// 	} else if item.NoticeType == proto.NOFICE_FORCE {
		// 		if plug.VersionCompare(item.Version, msg.Version, "<=") {
		// 			continue
		// 		}
		// 	}
		// }

		value := proto.NoticeInfo{}
		err = utils.CopyFields(&value, item)
		if err != nil {
			logger.Error(err)
		}
		result = append(result, value)
	}
	responseMessage.Notice = result
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = ""
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleGetUpgradeNotice end")
	SendPacket(a.conn, rsp)
}
