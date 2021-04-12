package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"time"
)

type UserInvitation struct {
}

const USER_INVITER = "user:inviter:"         // + 用户ID
const USER_INVIT_LIST = "user:inviter:list:" //+用户ID

func NewUserInvitation() *UserInvitation {
	return &UserInvitation{}
}

/**
绑定邀请码
*/
func (this *UserInvitation) BindingInviter(userId string, inviteCode string) (bool, error) {
	if inviteCode == "" {
		return false, errors.New("gameCode is null")
	}
	var baseLogic BaseLogic
	//检查绑定关系是否已存在
	key := USER_INVITER + userId
	result := db.RedisMgr.Get(key)
	if result == "1" {
		return false, errors.New("Binding relationship already exists")
	}
	// 获取用户信息
	user, err := db_service.UserIns.GetDataByInvite(inviteCode)
	if err != nil {
		return false, err
	}
	// 没有找到数据表示不存在
	if user.Id == "" || user.Id == "0" {
		return false, errors.New("gameCode not found")
	}

	baseInvit, err := baseLogic.BaseInvitation(userId) //  去base系统获取用户邀请信息
	fmt.Println(string(baseInvit))
	if err != nil {
		logger.Errorf("Invitation error", err.Error())
		return false, err
	}
	resp := &proto.S2C_HTTP{}
	err = json.Unmarshal(baseInvit, resp)
	if err != nil {
		logger.Errorf("BaseInvitation Invitation to json", err.Error())
		return false, err
	}

	//已绑定的数据添加 redis
	if resp.Code == 0 {
		fmt.Println(db.RedisMgr.Incr(key))
		return false, errors.New("Binding relationship already exists")
	}

	// 不存在记录，可以绑定
	if resp.Code == 30401 {
		inviteCode := user.InviteCode                                // 中台的code 主要用这个来同步绑定关系
		ret, err := baseLogic.BaseBindingInviter(userId, inviteCode) // 同步邀请关系
		if err != nil {
			return false, err
		}
		resp := &proto.S2C_HTTP{}
		err = json.Unmarshal(ret, resp)
		if err != nil {
			logger.Errorf("BaseBindingInviter Invitation to json", err.Error())
			return false, err
		}
		//绑定成功
		if resp.Code == 0 {
			//添加绑定关系
			id, _ := db_service.ActivityInvitationIns.AddUserInvit(userId, user.UserId)
			if id > 0 {
				//添加成功
				return true, nil
			} else {
				logger.Errorf("BindingInviter add err:", err.Error(), " userid:", userId, "invitId:", user.UserId)
				return false, errors.New("Binding failed")
			}
		} else {
			fmt.Println("同步绑定关系", err)
			return false, errors.New("Binding failed")
		}
	}
	return false, errors.New("Binding relationship already exists")
}

/**
手机注册保存用户基本信息
*/
func (this *UserInvitation) NewRegisterModelSave(registerModel proto.C2SRegisterMobile, msg []byte) (bool, error) {
	responseMessage := &proto.S2CCommon{}
	err := json.Unmarshal(msg, responseMessage)
	if err != nil {
		logger.Errorf("NewRegisterModelBindingInviter json.Unmarshal error, err=", err.Error())
		return false, err
	}
	// 判断是否注册成功
	if responseMessage.Code != errcode.MSG_SUCCESS {
		return false, errors.New("Registration failed")
	}
	// 获取用户ID
	var baseLogin BaseLogic
	msgdata, err := baseLogin.UserLogin(registerModel.Mobile, registerModel.Password, false, registerModel.CountryCode)
	if err != nil {
		return false, err
	}
	resp := &proto.S2C_HTTP{}
	err = json.Unmarshal(msgdata, resp)
	if err != nil {
		logger.Errorf("NewRegisterModelBindingInviter json.Unmarshal error, err=", err.Error())
		return false, err
	}
	if resp.Code != errcode.MSG_SUCCESS {
		return false, errors.New("BindingInviter failed")
	}
	userid := resp.Data["userId"].(string)
	// 保存用户信息
	userdata := model.User{
		SysType:       proto.SysType,
		UserId:        userid,
		UserType:      0,
		RoleId:        0,
		NickName:      registerModel.NickNname,
		Sex:           0,
		Level:         0,
		CountryCode:   registerModel.CountryCode,
		Mobile:        registerModel.Mobile,
		Email:         "",
		Token:         "",
		Status:        1,
		LocationId:    0,
		HouseNum:      0,
		ModifyNameNum: 1,
		InviterId:     "0",
		KycPassed:     0,
		KycStatus:     -1,
		TopLevel:      0,
		Point:         0,
	}
	_, err = db_service.UserIns.Add(&userdata)
	if err != nil {
		return false, err
	}
	// 添加邀请码
	return this.NewRegBingdingInviter(userid, registerModel.Inviter)
}

/**
邮箱注册保存用户基本信息
*/
func (this *UserInvitation) NewRegisterEmailSave(registerEmail proto.C2SRegisterEmail, msg []byte) (bool, error) {
	responseMessage := &proto.S2CCommon{}
	err := json.Unmarshal(msg, responseMessage)
	if err != nil {
		logger.Errorf("NewRegisterModelBindingInviter json.Unmarshal error, err=", err.Error())
		return false, err
	}
	// 判断是否注册成功
	if responseMessage.Code != errcode.MSG_SUCCESS {
		return false, errors.New("Registration failed")
	}
	// 获取用户ID
	var baseLogin BaseLogic
	msgdata, err := baseLogin.UserLogin(registerEmail.Email, registerEmail.Password, true, 0)
	if err != nil {
		return false, err
	}
	resp := &proto.S2C_HTTP{}
	err = json.Unmarshal(msgdata, resp)
	if err != nil {
		logger.Errorf("NewRegisterModelBindingInviter json.Unmarshal error, err=", err.Error())
		return false, err
	}
	if resp.Code != errcode.MSG_SUCCESS {
		return false, errors.New("BindingInviter failed")
	}
	userid := resp.Data["userId"].(string)
	// 保存用户信息
	userdata := model.User{
		SysType:       proto.SysType,
		UserId:        userid,
		UserType:      0,
		RoleId:        0,
		NickName:      registerEmail.NickNname,
		Sex:           0,
		Level:         0,
		Email:         registerEmail.Email,
		Token:         "",
		Status:        1,
		LocationId:    0,
		HouseNum:      0,
		ModifyNameNum: 1,
		InviterId:     "0",
		KycPassed:     0,
		KycStatus:     -1,
		TopLevel:      0,
		Point:         0,
	}
	_, err = db_service.UserIns.Add(&userdata)
	if err != nil {
		return false, err
	}
	// 同步邀请码
	return this.NewRegBingdingInviter(userid, registerEmail.Inviter)
}

/**
注册的可以直接绑定
*/
func (this *UserInvitation) NewRegBingdingInviter(userid string, invitaCode string) (bool, error) {
	// 获取用户信息
	invitUser, err := db_service.UserIns.GetDataByInvite(invitaCode) // 检查邀请码是否存在
	if err != nil {
		return false, errors.New("invitcode not found")
	}
	// 保存邀请关系
	id, err := db_service.ActivityInvitationIns.AddUserInvit(userid, invitUser.UserId)
	if id > 0 {
		//添加成功 需要创建角色才能增加积分
		//var rankList double_year.RankList
		//rankList.UpdateProp(invitUser.UserId, double_year.PropInviteFriend, 1) // 添加积分
		db.RedisMgr.SetAdd(USER_INVITER+userid, "1")                         // 标记用户以绑定
		db.RedisMgr.GetRedisClient().Del(USER_INVIT_LIST + invitUser.UserId) // 清除用户邀请列表
		return true, nil
	} else {
		logger.Errorf("NewRegBingdingInviter add err:", err, " userid:", userid, "invitId:", invitUser.UserId)
		return false, errors.New("Binding failed")
	}
}

/**
获取用户邀请列表
*/
func (this *UserInvitation) GetUserInvitList(userId string, isDay int) ([]map[string]interface{}, error) {
	if isDay == 1 {
		return this.GetUserDayInvitList(userId)
	}
	var invitList = make([]map[string]interface{}, 0)
	//key := USER_INVIT_LIST + userId
	//result := db.RedisMgr.Get(key)
	//if len(result) > 0 {
	//	err := json.Unmarshal([]byte(result), &invitList)
	//	if err != nil {
	//		logger.Errorf("GetUserInvitList json.Unmarshal error, err=", err.Error())
	//		return invitList, err
	//	}
	//	return invitList, nil
	//}
	//从数据库获取
	dataList, err := db_service.ActivityInvitationIns.GetUserInvitList(userId)
	if err != nil {
		logger.Errorf("GetUserInvitList db error, err=", err.Error())
		return invitList, err
	}
	for _, v := range dataList {
		data := map[string]interface{}{
			//"userId":     v.UserId,
			"nickName":   v.NickName,
			"createTime": v.UpdateTime.Format("2006/01/02"),
		}
		invitList = append(invitList, data)
	}
	//json, err := json.Marshal(invitList)
	//if err != nil {
	//	logger.Errorf("GetUserInvitList json.Marshal, err=", err.Error())
	//	return invitList, err
	//}
	//缓存一个小时
	//db.RedisMgr.GetRedisClient().Set(key, string(json), time.Second*3600)
	return invitList, nil
}

func (this *UserInvitation) GetUserDayInvitList(userId string) ([]map[string]interface{}, error) {
	datetime := time.Now()
	date := datetime.Format("2006-01-02")
	var invitList = make([]map[string]interface{}, 0)
	//从数据库获取
	dataList, err := db_service.ActivityInvitationIns.GetUserDayInvitList(userId, date+" 00:00:01", date+" 23:59:59")
	if err != nil {
		logger.Errorf("GetUserInvitList db error, err:[%v]", err.Error())
		return invitList, err
	}
	for _, v := range dataList {
		data := map[string]interface{}{
			//"userId":     v.UserId,
			"nickName":   v.NickName,
			"createTime": v.UpdateTime.Format("2006/01/02"),
		}
		invitList = append(invitList, data)
	}
	return invitList, nil
}

//检查邀请用户是否是第一次创建角色
func (this *UserInvitation) CheckUserInviteIsFirstCreateRole(userId string) (bool, string, error) {
	//判断用户是否是邀请用户
	inviteInfo, err := db_service.ActivityInvitationIns.GetInvitInfoByUserId(userId)
	if err != nil {
		logger.Errorf("CheckUserInviteIsCreateRole GetInvitInfoBy UserId:[%v] ,err:[%v] \n", userId, err.Error())
		return false, "", err
	}
	//判断是否有邀请记录
	if inviteInfo.Id == 0 {
		logger.Debugf("CheckUserInviteIsCreateRole invite is null inviteInfo:[%v]", inviteInfo)
		return false, "", nil
	}
	logger.Debugf("CheckUserInviteIsCreateRole userId:[%v] ,inviteInfo:[%v]", userId, inviteInfo)
	if inviteInfo.IsCreateRole == 1 {
		// 已经创建过角色了
		return false, "", nil
	}
	row, err := db_service.ActivityInvitationIns.UpdateIsCreateRole(userId, 1)
	if err != nil {
		logger.Errorf("CheckUserInviteIsCreateRole UpdateIsCreateRole falid userId:[%v] ,err:[%v] \n", userId, err.Error())
		return false, "", err
	}
	if row == 0 {
		return false, "", nil
	}
	return true, inviteInfo.InviterId, nil
}
