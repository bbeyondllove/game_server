package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"time"

	"github.com/go-xorm/xorm"

	"game_server/core/logger"
)

type User struct{}

const (
	user_nid  = 1
	UserTable = "t_user"
)

//获取自动添加序列号
func (this *User) NextSeqId(session *xorm.Session) (string, error) {
	var maxIdStr string
	has, err := db.Mysql.Table(&model.User{}).Select("max(id)").Get(&maxIdStr)
	if err != nil {
		return "", err
	}
	maxId, _ := strconv.Atoi(maxIdStr)
	if has && maxId != 0 {
		return strconv.Itoa(maxId + 1), nil
	} else {
		count, err := db.Mysql.Count(&model.User{})
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(count + user_nid)), nil
	}
}

//添加记录
func (this *User) Add(data_map *model.User) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	seqId, err := this.NextSeqId(session)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}

	data_map.Id = seqId
	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		logger.Errorf("[Error]: ", err)
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		logger.Errorf("[Error]: ", err)
		return false, err
	}
	return true, nil
}

//根据用户ID获取用户信息
//user_id 用户ID
func (this *User) GetDataByUid(userId string) (model.User, error) {
	data := model.User{}
	_, err := db.Mysql.Table(UserTable).Where("user_id= ? ", userId).Get(&data)
	return data, err
}

//根据邀请码获取用户信息
func (this *User) GetDataByInvite(inviteCode string) (model.User, error) {
	data := model.User{}
	_, err := db.Mysql.Table(UserTable).Where("invite_code= ? ", inviteCode).Get(&data)
	return data, err
}

//根据手机或邮箱获取用户信息
//nType 1:手机;2：邮箱S
//account用户手机号码或邮箱
func (this *User) GetDataByAccount(nType int, account, userId string) (model.User, error) {
	data := model.User{}
	var err error
	if nType == proto.CODE_TYPE_PHONE {
		_, err = db.Mysql.Table(UserTable).Where("(user_id=0 and mobile= ?) or (user_id>0 and user_id= ?)", account, userId).Get(&data)
	} else {
		_, err = db.Mysql.Table(UserTable).Where("(user_id=0 and email= ?) or (user_id>0 and user_id= ?)", account, userId).Get(&data)

	}
	return data, err
}

//根据用户昵称获取用户信息
//nick_name 用户昵称
func (this *User) GetDataByNickName(nick_name string) (model.User, error) {
	data := model.User{}
	_, err := db.Mysql.Table(UserTable).Where("nick_name= ?", nick_name).Get(&data)
	return data, err
}

//获取所有用户数据
func (this *User) GetAllData() ([]model.User, error) {
	var data []model.User
	err := db.Mysql.Table(UserTable).Where("1=1").Find(&data)
	return data, err
}

//获取所有用户数据
func (this *User) GetAllCount() (int, error) {
	count, err := db.Mysql.Table(UserTable).Where("1=1").Count()
	return int(count), err
}

//获取所有用户数据
func (this *User) GetAllCountByPlatform(platform int) (int, error) {
	count, err := db.Mysql.Table(UserTable).Where("platform = ?", platform).Count()
	return int(count), err
}

//获取指定时间段用户数据
//startTime 开始时间
//endTime 结束时间
func (this *User) GetTimeCount(startTime, endTime string) (int, error) {
	count, err := db.Mysql.Table(UserTable).Where("create_time>=? and create_time<=?", startTime, endTime).Count()
	return int(count), err
}

// 更新用户邀请码
func (this *User) UpdateUserGameCode(userId string, gameCode string) (int, error) {
	data := map[string]interface{}{
		"game_code": gameCode,
	}
	number, err := db.Mysql.Table(UserTable).Where("user_id= ? ", userId).Update(data)
	return int(number), err
}

// 通过用户名/用户ID/手机/邮箱获取用户
func (this *User) GetAllUserByKey(key string) ([]model.User, error) {
	var data []model.User
	err := db.Mysql.Table(UserTable).Where("user_id like '%"+key+"%'").
		Or("nick_name like '%"+key+"%'").
		Or("mobile like '%"+key+"%'").
		Or("email like '%"+key+"%'").
		Cols("user_id", "nick_name", "country_code", "mobile", "email", "create_time", "update_time", "kyc_passed", "cdt", "Status", "login_ip").
		Find(&data)
	return data, err
}
func (this *User) GetPageUserByKey(key string, page, size int) (int, int, []model.User, error) {
	if page <= 0 {
		page = 1
	}
	session := db.Mysql.Table(UserTable).Where("user_id like '%" + key + "%'").
		Or("nick_name like '%" + key + "%'").
		Or("mobile like '%" + key + "%'").
		Or("email like '%" + key + "%'")
	count, err := session.Count()
	if err != nil {
		return 0, 0, nil, err
	}
	var data []model.User
	session = db.Mysql.Table(UserTable).Where("user_id like '%" + key + "%'").
		Or("nick_name like '%" + key + "%'").
		Or("mobile like '%" + key + "%'").
		Or("email like '%" + key + "%'")
	err = session.
		Limit(size, (page-1)*size).
		Cols("user_id", "nick_name", "country_code", "mobile", "email", "create_time", "update_time", "kyc_passed", "cdt", "Status", "login_ip").
		Find(&data)
	return int(count), page, data, err
}

// 是否禁用用户
func (this *User) UpdateUserStatus(userId string, status int) (int, error) {
	data := map[string]interface{}{
		"Status": status,
	}
	number, err := db.Mysql.Table(UserTable).Where("user_id= ? ", userId).Update(data)
	return int(number), err
}

// 更新用户实名状态
func (this *User) UpdateUserKycStatus(userId string, kyc_status int) (int, error) {
	data := map[string]interface{}{
		"kyc_status": kyc_status,
	}
	if kyc_status == 2 {
		data["kyc_passed"] = 1
	} else {
		data["kyc_passed"] = 0
	}
	number, err := db.Mysql.Table(UserTable).Where("user_id= ? ", userId).Update(data)
	return int(number), err
}

// 总新增:通过微信或者邮箱注册了的用户并且进入游戏创建角色后，有角色数据的用户数总和
func (this *User) TotalNewlyAddedCount() (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("role_id != 0").Count()
	return int(number), err
}

func (this *User) TotalNewlyAddedCountByPlatform(platform int) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("role_id != 0").And("platform = ?", platform).Count()
	return int(number), err
}

// day 格式为:'2017-06-16'
func (this *User) CountOfDay(day string) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("DATE_FORMAT(create_time,'%Y-%m-%d') = ?", day).Count()
	return int(number), err
}

// day 格式为:'2017-06-16'
func (this *User) CountOfDayAndPlatform(day string, platform int) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("DATE_FORMAT(create_time,'%Y-%m-%d') = ?", day).And("platform = ?", platform).Count()
	return int(number), err
}

// 新增
func (this *User) NewlyAddedCount(day string) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("role_id != 0").And("DATE_FORMAT(create_time,'%Y-%m-%d') = ?", day).Count()
	return int(number), err
}

func (this *User) NewlyAddedCountByPlatform(day string, platform int) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("role_id != 0").And("DATE_FORMAT(create_time,'%Y-%m-%d') = ?", day).And("platform = ?", platform).Count()
	return int(number), err
}

// 注册数
func (this *User) TotalRegisteCount(type_ string) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where(type_ + " != ''").Count()
	return int(number), err
}

func (this *User) TotalRegisteCountByPlatform(type_ string, platform int) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where(type_+" != ''").And("platform = ?", platform).Count()
	return int(number), err
}

// 实名用户数
func (this *User) RealNameCount() (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("kyc_passed = 1").Count()
	return int(number), err
}

func (this *User) RealNameCountByPlatform(platform int) (int, error) {
	number, err := db.Mysql.Table(UserTable).Where("kyc_passed = 1").And("platform = ?", platform).Count()
	return int(number), err
}

// 留存数,day 格式为:'2017-06-16'
func (this *User) GetRetainedCount(day string, diff int) (int, error) {
	sql := "SELECT count(*) FROM (SELECT *,timestampdiff(DAY,DATE_FORMAT(role_create_time,'%Y-%m-%d'),DATE_FORMAT(update_time,'%Y-%m-%d')) as DIFF FROM " + UserTable + " WHERE DATE_FORMAT(update_time,'%Y-%m-%d') = '" + day + "' ) a WHERE a.DIFF = " + strconv.Itoa(diff)
	number, err := db.Mysql.Sql(sql).Count()
	return int(number), err
}
