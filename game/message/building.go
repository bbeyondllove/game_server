package message

import (
	"encoding/json"
	"errors"
	kk_core "game_server/core"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"html"
	"sync"
	"time"
)

var (
	_SleepTime = 5 * time.Minute
)
var BuildingTypes sync.Map

func GetBuildingType(smallType string) (ret model.BuildType, err error) {
	if item, ok := BuildingTypes.Load(smallType); ok {
		if info, ok := item.(model.BuildType); ok {
			return info, nil
		}
		return model.BuildType{}, errors.New("BuildType conver failed smallType=" + smallType)
	}
	return model.BuildType{}, errors.New("BuildType not found smallType=" + smallType)
}

//获取建筑简介
func (s *CSession) HandleGetBuildingDesc(requestMsg *utils.Packet) {
	logger.Debugf("in request:", requestMsg.GetBuffer())

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_BUILDING_DESC_RSP)

	responseMessage := &proto.S2CGetBuildingInfo{}
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SGetBuildingInfo{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, _, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	kk_core.PushMysql(func() {
		building, err := db_service.WorldMapIns.GetBuildDesc(s.LocateId, msg.LocationX, msg.LocationY)
		if err != nil {
			kk_core.PushWorld(func() {
				logger.Errorf("WorldMapIns.GetBuildDesc error=", err.Error())
				responseMessage.Code = errcode.ERROR_MYSQL
				responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
				rsp.WriteData(responseMessage)
				s.sendPacket(rsp)
			})
			return
		}
		kk_core.PushWorld(func() {
			building.Desc = html.UnescapeString(building.Desc) //转换html实体
			building.Desc = utils.Stripslashes(building.Desc)  //将引号去掉
			/*var buildTypeName string
			if building.SmallType != "" {
				buildType, err := GetBuildingType(building.SmallType)
				if err != nil {
					kk_core.PushWorld(func() {
						logger.Errorf("WorldMapIns.GetBuildingType error=", err.Error())
						responseMessage.Code = errcode.ERROR_MYSQL
						responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
						rsp.WriteData(responseMessage)
						s.sendPacket(rsp)
					})
					return
				}
				buildTypeName = buildType.BuildingName
			}*/

			responseMessage.Desc = building.Desc
			responseMessage.H5Url = building.H5Url
			responseMessage.WebUrl = building.WebUrl
			responseMessage.PassportAviable = building.PassportAviable
			responseMessage.ImageUrl = building.ImageUrl
			responseMessage.SmallType = building.SmallType
			responseMessage.BuildingTypeName = building.BuildingName
			responseMessage.BuildingName = building.BuildingName

			responseMessage.Code = errcode.MSG_SUCCESS
			responseMessage.Message = ""
			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
		})
	})

	return
}
