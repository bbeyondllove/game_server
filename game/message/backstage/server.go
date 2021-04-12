package backstage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/db"
	dao "game_server/game/db_service"
	"game_server/game/errcode"
	ecode "game_server/game/errcode"
	"game_server/game/message"
	"game_server/game/model"
	"game_server/game/proto"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"game_server/core/logger"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var (
	ADMIN_USER          = "admin:user:"
	DEFAULT_ADMIN_LIMIT = "4001,4002,4003"
)

// @Summary 管理员用户登陆
// @Description 管理员用户登陆
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /v1/login [get]
func AdminLogin(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	var param map[string]string
	err = json.Unmarshal(data, &param)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	name := param["username"]
	password := param["password"]
	// name := c.PostForm("username")
	// password := c.PostForm("password")

	user, err := dao.AdminUserIns.GetDataByUname(name)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	// 账号不存在
	if user.Account != name {
		JsonErrorCode(c, ecode.ERROR_HTTP_USER_NOT_EXIST)
		return
	}
	// 账号被禁用
	if user.Status == 0 {
		JsonErrorCode(c, ecode.ERROR_HTTP_USER_NOT_ALLOW)
		return
	}
	// 账号密码错误
	// 加密处理
	// h := md5.New()
	// h.Write([]byte(password))
	// value := string(h.Sum(nil))
	if user.Password != password {
		JsonErrorCode(c, ecode.ERROR_HTTP_ACCOUNT_ERROR)
		return
	}
	token, err := CreateToken(strconv.FormatInt(user.Id, 10), message.G_BaseCfg.Backstage.TokenSecret)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_SYSTEM)
		return
	}
	// 获取用户权限
	role, err := dao.AdminRoleIns.GetDataById(user.RoleId)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	// 获取创建者
	creater_name := ""
	if user.Creater == 0 {
		creater_name = "system"
	} else {
		creater, err := dao.AdminUserIns.GetDataById(user.Creater)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_MYSQL)
			return
		}
		creater_name = creater.Username
	}
	user_data := map[string]interface{}{
		"username":    user.Username,
		"account":     user.Account,
		"department":  user.Department,
		"creater":     creater_name,
		"token":       token,
		"role_limits": role.Limits,
	}
	// 设置用户的信息至redis
	_, err = db.RedisGame.HMSet(ADMIN_USER+strconv.FormatInt(user.Id, 10), user_data).Result()
	if err != nil {
		logger.Error(err)
	}
	limits := strings.Split(role.Limits, ",")
	permissionList, permissions, _ := getAdminLimitList(limits)
	user_data["role_limits"] = limits
	user_data["permissionList"] = permissionList
	user_data["permissions"] = permissions
	JsonData(c, user_data)
}

// @Summary 管理员用户登出
// @Description 管理员用户登出
// @Produce  json
// @Accept  json
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/login [get]
func AdminLogout(c *gin.Context) {
	uid := c.GetString(ADMIN_UID)
	_, err := db.RedisGame.Del(ADMIN_USER + uid).Result()
	if err != nil {
		logger.Error(err)
	}
	JsonOK(c)
}

// 获取管理者名称
func getAdminName(id int64) string {
	if id <= 0 {
		return ""
	}
	user, err := dao.AdminUserIns.GetDataById(id)
	if err != nil {
		return ""
	}
	return user.Username
}

func GetAdminUsers(c *gin.Context) {
	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.AdminUserIns.GetPageData(page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for i := 0; i < len(data); i++ {
		role, _ := dao.AdminRoleIns.GetDataById(data[i].RoleId)
		data[i].RoleName = role.Name
		if data[i].Creater != 0 {
			creater, _ := dao.AdminUserIns.GetDataById(data[i].Creater)
			data[i].CreaterName = creater.Username
		} else {
			data[i].CreaterName = "system"
		}
		item := map[string]interface{}{
			"id":           data[i].Id,
			"role_id":      data[i].RoleId,
			"role_name":    data[i].RoleName,
			"department":   data[i].Department,
			"username":     data[i].Username,
			"password":     data[i].Password,
			"account":      data[i].Account,
			"creater_name": data[i].CreaterName,
			"creater":      data[i].Creater,
			"status":       data[i].Status,
			"create_time":  utils.Time2Str(data[i].CreateTime),
			"update_time":  utils.Time2Str(data[i].UpdateTime),
		}
		result = append(result, item)
	}
	JsonPage(c, result, page, size, total)
}

func AddAdminUser(c *gin.Context) {
	account := c.Query("account")
	username := c.Query("username")
	password := c.Query("password")
	department := c.Query("department")
	role_name := c.Query("role_name")

	if len(account) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	if len(password) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}

	// 角色
	var RoleId int64 = 0
	if len(role_name) > 0 {
		role, err := dao.AdminRoleIns.GetDataByName(role_name)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if role.Name != role_name {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		RoleId = role.Id
	}

	// 创建者
	operation_userid := c.GetString(ADMIN_UID)
	operation, err := dao.AdminUserIns.GetDataByUid(operation_userid)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	_, err = dao.AdminUserIns.Add(&model.AdminUser{
		RoleId:     RoleId,
		Account:    account,
		Password:   password,
		Department: department,
		Username:   username,
		Creater:    operation.Id,
		Status:     1,
	})
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonOK(c)
}

func UpdateAdminUser(c *gin.Context) {
	account := c.Query("account")
	var username_ptr *string = nil
	var password_ptr *string = nil
	var department_ptr *string = nil
	var role_name_ptr *string = nil

	username, exists := c.GetQuery("username")
	if exists {
		username_ptr = &username
	}
	password, exists := c.GetQuery("password")
	if exists {
		if len(password) == 0 {
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		password_ptr = &password
	}
	department, exists := c.GetQuery("department")
	if exists {
		department_ptr = &department
	}
	role_name, exists := c.GetQuery("role_name")
	if exists {
		role_name_ptr = &role_name
	}

	if len(account) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}

	var RoleId *int64 = nil
	if role_name_ptr != nil {

		if len(role_name) > 0 {
			role, err := dao.AdminRoleIns.GetDataByName(role_name)
			fmt.Println(role)
			if err != nil {
				logger.Error(err)
				JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
				return
			}
			if role.Name != role_name {
				logger.Error(err)
				JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
				return
			}
			RoleId = &role.Id
		} else {
			var Rid int64 = 0
			RoleId = &Rid
		}
	}
	_, err := dao.AdminUserIns.Update(account, username_ptr, password_ptr, department_ptr, RoleId)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonOK(c)
}

func UpdateAdminUserStatus(c *gin.Context) {
	account := c.Query("account")
	status_ := c.Query("status")

	status, err := strconv.ParseInt(status_, 10, 64)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	_, err = dao.AdminUserIns.UpdateStatus(account, status)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonOK(c)
}

func DeleteAdminUser(c *gin.Context) {
	account := c.Query("account")

	// 不允许删除管理员
	if account == "admin" {
		JsonErrorCode(c, ecode.ERROR_REQUEST_NOT_ALLOW)
		return
	}
	num, err := dao.AdminUserIns.Delete(account)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonData(c, num)
}

func AddAdminRole(c *gin.Context) {
	role_name := c.Query("role_name")
	remark := c.Query("remark")
	limit := c.Query("limit") // 权限需要使用,拼接所有权限ID

	if len(role_name) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}

	// operation_userid := c.GetString(ADMIN_UID)
	// operation, err := dao.AdminUserIns.GetDataByUid(operation_userid)
	// if err != nil {
	// 	logger.Error(err)
	// 	JsonErrorCode(c, ecode.ERROR_MYSQL)
	// 	return
	// }

	_, err := dao.AdminRoleIns.Add(&model.AdminRole{
		Name:   role_name,
		Remark: remark,
		Limits: limit,
	})
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonOK(c)
}

func UpdateAdminRole(c *gin.Context) {
	role_name := c.Query("role_name")
	remark := c.Query("remark")
	limit := c.Query("limit") // 权限需要使用,拼接所有权限ID

	if role_name == "超级管理员" {
		if len(limit) > 0 {
			limit = DEFAULT_ADMIN_LIMIT + "," + limit
		} else {
			limit = DEFAULT_ADMIN_LIMIT
		}
	}
	_, err := dao.AdminRoleIns.Update(role_name, model.AdminRole{
		Name:   role_name,
		Remark: remark,
		Limits: limit,
	})
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonOK(c)
}

func GetAdminRole(c *gin.Context) {
	data, err := dao.AdminRoleIns.GetAllData()
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]interface{}, 0)
	for _, item := range data {
		num_people, _ := dao.AdminUserIns.Count(item.Id)
		result = append(result, map[string]interface{}{
			"role_id":   item.Id,
			"role_name": item.Name,
			"remark":    item.Remark,
			"limits":    strings.Split(item.Limits, ","),
			"num":       num_people,
		})
	}

	JsonData(c, result)
}

func DeleteAdminRole(c *gin.Context) {
	role_name := c.Query("role_name")

	// 不允许删除管理员角色
	if role_name == "超级管理员" {
		JsonErrorCode(c, ecode.ERROR_REQUEST_NOT_ALLOW)
		return
	}
	num, err := dao.AdminRoleIns.Delete(role_name)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonData(c, num)
}

func getAdminLimitList(limits []string) (interface{}, interface{}, error) {
	in_array := func(arr []string, item string) bool {
		for _, it := range arr {
			if it == item {
				return true
			}
		}
		return false
	}
	data, err := dao.AdminLimitIns.GetAllData()
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}
	table_codes := make([]string, 0)
	table_names := make([]string, 0)
	table := make(map[string][]model.AdminLimit, 0)
	for _, item := range data {
		if !in_array(limits, strconv.FormatInt(item.Id, 10)) {
			continue
		}
		value, ok := table[item.Group]
		if ok {
			value = append(value, item)
			table[item.Group] = value
			table_codes = append(table_codes, item.Code)
		} else {
			table_names = append(table_names, item.Group)
			value = make([]model.AdminLimit, 0)
			value = append(value, item)
			table[item.Group] = value
			table_codes = append(table_codes, strings.Split(item.Code, ":")[0])
			table_codes = append(table_codes, item.Code)
		}
	}
	index := 1
	result := make([]interface{}, 0)
	for _, name := range table_names {
		value := table[name]
		result = append(result, map[string]interface{}{
			"id":       index,
			"label":    name,
			"children": value,
		})
		index++
	}

	return result, table_codes, nil
}

func GetAdminLimit(c *gin.Context) {
	data, err := dao.AdminLimitIns.GetData(1)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	table_names := make([]string, 0)
	table := make(map[string][]model.AdminLimit, 0)
	for _, item := range data {
		value, ok := table[item.Group]
		if ok {
			value = append(value, item)
			table[item.Group] = value
		} else {
			table_names = append(table_names, item.Group)
			value = make([]model.AdminLimit, 0)
			value = append(value, item)
			table[item.Group] = value
		}
	}
	index := 1
	result := make([]interface{}, 0)
	for _, name := range table_names {

		value := table[name]
		result = append(result, map[string]interface{}{
			"id":       index,
			"label":    name,
			"children": value,
		})
		index++
	}
	JsonData(c, result)
}

// @Summary 查询用户
// @Description 用于查询用户
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /v1/api/get/users [get]
func GetUsers(c *gin.Context) {
	logger.Debugf("SearchUser start")

	key := c.Query("key")
	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.UserIns.GetPageUserByKey(key, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]model.UserRsp, 0)
	for _, item := range data {
		result = append(result, model.UserRsp{
			UserId:        item.UserId,
			NickName:      item.NickName,
			CountryCode:   item.CountryCode,
			Mobile:        item.Mobile,
			Email:         item.Email,
			Status:        item.Status,
			KycPassed:     item.KycPassed,
			RegisteIp:     item.LoginIp,
			Cdt:           decimal.NewFromFloat32(item.Cdt),
			CreateTime:    utils.Time2Str(item.CreateTime),
			LastLoginTime: utils.Time2Str(item.UpdateTime),
		})
	}
	JsonPage(c, result, page, size, total)
	return
}

// @Summary 冻结/解冻用户
// @Description 用于冻结/解冻用户
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/frozen/user [get]
func FrozenUser(c *gin.Context) {
	user_id := c.Query("user_id")
	status_ := c.Query("status") // 是否禁用,0 禁用 1 未禁用

	if len(user_id) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	status, err := strconv.Atoi(status_)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	_, err = dao.UserIns.UpdateUserStatus(user_id, status)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	item, err := dao.UserIns.GetDataByUid(user_id)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_USER_NOT_EXIT)
		return
	}
	if item.UserId != user_id {
		JsonErrorCode(c, ecode.ERROR_USER_NOT_EXIT)
		return
	}
	JsonData(c, model.UserRsp{
		UserId:        item.UserId,
		NickName:      item.NickName,
		CountryCode:   item.CountryCode,
		Mobile:        item.Mobile,
		Email:         item.Email,
		Status:        item.Status,
		KycPassed:     item.KycPassed,
		RegisteIp:     item.LoginIp,
		Cdt:           decimal.NewFromFloat32(item.Cdt),
		CreateTime:    utils.Time2Str(item.CreateTime),
		LastLoginTime: utils.Time2Str(item.UpdateTime),
	})
}

// @Summary 获取实名认证用户
// @Description 获取实名认证用户
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/get/certifications [get]
func GetCertifications(c *gin.Context) {
	// logger.Debugf("SearchCertification start")

	key := c.Query("key")
	status := c.Query("status") //status : 0：首次审核，1：审核不通过，2：审核通过 3：再次提交审核

	// if len(key) == 0 {
	// 	JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
	// 	return
	// }

	data, err := dao.UserIns.GetAllUserByKey(key)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	user_info := make(map[string]model.User, 0)
	user_ids := make([]string, 0)
	for _, item := range data {
		user_ids = append(user_ids, item.UserId)
		user_info[item.UserId] = item
	}
	var certifications []model.Certification
	if len(status) == 0 {
		certifications, err = dao.CertificationIns.GetDataByUids(user_ids)
	} else {
		certifications, err = dao.CertificationIns.GetDataByUidsAndStatus(user_ids, status)
	}

	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]model.CertificationRsp, 0)
	for _, item := range certifications {
		ExamineTime := ""
		record, err := dao.CertificationRecordIns.GetLastDataByUid(item.UserId)
		if err == nil {
			if record.UserId == item.UserId {
				ExamineTime = record.CreateTime.Local().String()
			}
		}
		user := user_info[item.UserId]
		result = append(result, model.CertificationRsp{
			UserId:      item.UserId,
			NickName:    user.NickName,
			CountryCode: user.CountryCode,
			Mobile:      user.Mobile,
			Email:       user.Email,

			Nationality: item.Nationality,
			FirstName:   item.FirstName,
			LastName:    item.LastName,
			IdType:      item.IdType,
			IdNumber:    item.IdNumber,
			ObjectKey:   item.ObjectKey,
			Suggestion:  item.Reson,
			Status:      item.Status,
			ApplyTime:   utils.Time2Str(item.UpdateTime),
			ExamineTime: ExamineTime,
		})
	}
	JsonData(c, result)
}

// @Summary 认证用户
// @Description 认证用户
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/certification/user [get]
func CertificationUser(c *gin.Context) {
	user_id := c.Query("user_id")
	status_ := c.Query("status")        //status : 0：首次审核，1：审核不通过，2：审核通过 3：再次提交审核
	suggestion := c.Query("suggestion") //审核意见

	if len(user_id) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	status, err := strconv.Atoi(status_)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	if status == 1 {
		if len(suggestion) == 0 {
			logger.Error("审核不通过,未提交审核已经")
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
	}
	certification, err := dao.CertificationIns.GetDataByUid(user_id)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if certification.UserId != user_id {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_REQUEST_NOT_ALLOW)
		return
	}
	_, err = dao.CertificationIns.UpdateStatus(user_id, status, suggestion)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if status == 1 || status == 2 {
		_, err = dao.UserIns.UpdateUserStatus(user_id, status)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_MYSQL)
			return
		}
		key_passed := 0
		if status == 2 {
			key_passed = 1
		}
		// 更新用户状态
		db.RedisGame.HSet(user_id, "kyc_status", status)
		db.RedisGame.HSet(user_id, "kyc_passed", key_passed)
	}

	// 登陆用胡
	ExamineUserId := c.GetString(ADMIN_UID)

	user, _ := dao.AdminUserIns.GetDataByUid(ExamineUserId)

	dao.CertificationRecordIns.Add(&model.CertificationRecord{
		ExamineUserId:   ExamineUserId,
		ExamineUserName: user.Username,
		UserId:          user_id,
		Status:          certification.Status,
		Suggestion:      suggestion,
		ExamineStatus:   status,
	})

	JsonOK(c)
}

// @Summary 获取认证用户历史记录
// @Description 获取认证用户历史记录
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/get/user/certification/record [get]
func GetCertificationRecord(c *gin.Context) {
	user_id := c.Query("user_id")
	data, err := dao.CertificationRecordIns.GetAllData(user_id)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonData(c, data)
}

// @Summary 获取任务列表
// @Description 获取任务列表
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK"}
// @Router /v1/api/get/tasks [get]
func GetTasks(c *gin.Context) {
	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	task_type_ := c.Query("task_type")
	task_type, err := strconv.Atoi(task_type_)
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	total, page, tasks, err := dao.TasksIns.GetPageDataByTaskType(task_type, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	TaskAward, twerr := dao.TaskAwardIns.GetAllData()
	if twerr != nil {
		return
	}
	getAward := func(taskId int) (awardId, awardNum string) {
		for _, v := range TaskAward {
			if v.TaskId == taskId {
				return v.AwardId, v.AwardNum
			}
		}
		return "", ""
	}

	task_item := make([]proto.TaskItem, 0)
	for i := 0; i < len(tasks); i++ {
		item := proto.TaskItem{}
		err := utils.CopyFields(&item, *tasks[i])
		if err != nil {
			logger.Error(err)
		}
		awardId, awardNum := getAward(tasks[i].Id)
		if awardId == "" {
			continue
		}
		awardInfo := message.GetAwardInfo(awardId, awardNum)

		item.Awards = awardInfo
		task_item = append(task_item, item)
	}
	JsonPage(c, task_item, page, size, total)
}

func GetItems(c *gin.Context) {
	items, err := dao.ItemIns.GetAllData()
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonData(c, items)
}

func SetTaskAward(c *gin.Context) {
	task_id := c.Query("task_id")
	award_id := c.Query("award_id")
	award_num_ := c.Query("award_num")

	if len(task_id) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_EMPTY)
		return
	}
	if len(award_id) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_EMPTY)
		return
	}
	if len(award_num_) == 0 {
		JsonErrorCode(c, ecode.ERROR_PARAM_EMPTY)
		return
	}
	award_num, err := strconv.Atoi(award_num_)
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}

	num, err := dao.TaskAwardIns.UpdateDataById(task_id, award_id, award_num)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	JsonData(c, num)
}

func ossUpload(fileName string, fileByte []byte) (url string, err error) {
	/*
	   oss 的相关配置信息
	*/
	bucketName := base.Setting.Admin.Oss["Bucket"]
	endpoint := base.Setting.Admin.Oss["Endpoint"]
	accessKeyId := base.Setting.Admin.Oss["AccessKeyId"]
	accessKeySecret := base.Setting.Admin.Oss["AccessKeySecret"]
	domain := base.Setting.Admin.Oss["Domain"]
	uploadPath := base.Setting.Admin.Oss["Path"]

	//创建OSSClient实例
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.EnableCRC(false))
	if err != nil {
		return url, err
	}

	// 获取存储空间
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return url, err
	}

	//上传阿里云路径
	folderName := time.Now().Format("20060102")
	yunFileTmpPath := filepath.Join(uploadPath, folderName) + "/" + fileName

	// 上传Byte数组
	err = bucket.PutObject(yunFileTmpPath, bytes.NewReader([]byte(fileByte)))
	// err = bucket.PutObject(yunFileTmpPath, strings.NewReader(string(fileByte)))
	if err != nil {
		return url, err
	}

	return domain + "/" + yunFileTmpPath, nil
}

func UploadImage(c *gin.Context) {
	// https://zifu-admin-client.oss-cn-shenzhen.aliyuncs.com/test/common/1611731418682.png
	//获取表单数据 参数为name值
	var (
		err error
	)
	file, err := c.FormFile("file")
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}

	if message.G_BaseCfg.Backstage.MaxUploadImageSize != 0 && file.Size > message.G_BaseCfg.Backstage.MaxUploadImageSize {
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	fileHandle, err := file.Open() //打开上传文件
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	defer fileHandle.Close()
	fileByte, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_SYSTEM)
		return
	}

	url, err := ossUpload(strconv.FormatInt(time.Now().Unix(), 10)+".png", fileByte)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_SYSTEM)
		return
	}
	JsonData(c, gin.H{
		"url": url,
	})
}

func GetNotice(c *gin.Context) {
	notice_type_ := c.Query("notice_type")

	notice_type, err := strconv.Atoi(notice_type_)
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.NoticeIns.GetPageData(notice_type, start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for i := 0; i < len(data); i++ {
		item := map[string]interface{}{
			"id":             data[i].Id,
			"notice_type":    data[i].NoticeType,
			"notice_title":   data[i].NoticeTitle,
			"notice_content": data[i].NoticeContent,
			"notice_url":     data[i].NoticeUrl,
			"version":        data[i].Version,
			"remark":         data[i].Remark,
			"is_noticed":     data[i].IsNoticed,
			"notice_time":    data[i].NoticeTime,
			"create_time":    utils.Time2Str(data[i].CreateTime),
			"update_time":    utils.Time2Str(data[i].UpdateTime),
		}

		item["operator_name"] = getAdminName(data[i].Operator)
		data, err := dao.TimeNoticeIns.GetDataByNoticeId(data[i].Id)
		if err != nil {
			logger.Error(err)
		}

		item["walking_lanterns"] = &data
		result = append(result, item)
	}
	JsonPage(c, result, page, size, total)
}

// 激活走马灯通告
func StartLanternsNotice() {
	records, err := dao.NoticeIns.GetUnnoticed()
	if err != nil {
		return
	}
	for _, notice := range records {
		if notice.NoticeType != proto.NOTICE_BEFOR_UPDATE {
			continue
		}
		notices, err := dao.TimeNoticeIns.GetDataByNoticeId(notice.Id)
		if err != nil {
			continue
		}
		notice_time, err := utils.ParseTime(notice.NoticeTime)
		if err != nil {
			continue
		}
		for _, item := range notices {
			item.StartTime = notice_time.Add(-1 * time.Second * time.Duration(item.IntervalTime))
			message.G_WalkingLantenrns.AddTimeNotice(item.NoticeId, item)
		}
		message.G_WalkingLantenrns.AddNotice(notice_time, notice)
	}
	message.G_WalkingLantenrns.Start()
}

func BroadCastNotice(notice_info model.Notice) {
	logger.Debugf("BroadCastNotice:", notice_info)
	// 广播公告
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_NOTICE_RSP)
	pushMessage := &proto.S2CNotices{}
	//推送给所有用户
	pushMessage.Code = errcode.MSG_SUCCESS
	pushMessage.Message = ""
	pushMessage.Notice = []proto.NoticeInfo{
		proto.NoticeInfo{
			NoticeType:    notice_info.NoticeType,
			NoticeTitle:   notice_info.NoticeTitle,
			NoticeContent: notice_info.NoticeContent,
			NoticeUrl:     notice_info.NoticeUrl,
			Version:       notice_info.Version,
		},
	}
	rsp.WriteData(pushMessage)
	// Sched.BroadCastMsg(int32(3302), "0", rsp)
	// 广播给在线的用户
	if notice_info.NoticeType == proto.NOTICE_NORNAL {
		message.Sched.BroadCastMsgByVersion(int32(3302), "0", rsp, "")
	} else {
		message.Sched.BroadCastMsgByVersion(int32(3302), "0", rsp, notice_info.Version)
	}
}

func AddNotice(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	var notice proto.NoticeReq
	err = json.Unmarshal(data, &notice)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	notice_info := model.Notice{}
	err = utils.CopyFields(&notice_info, notice)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}

	// 检查公告时间
	var notice_time_ptr *time.Time = nil
	if notice_info.NoticeType == proto.NOTICE_BEFOR_UPDATE {
		if len(notice_info.NoticeTime) > 0 {
			time_value, err := utils.ParseTime(notice_info.NoticeTime)
			if err != nil {
				logger.Error(err)
				JsonErrorCode(c, ecode.ERROR_ILLEGAL)
				return
			}
			notice_time_ptr = &time_value
		} else {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
	}

	// 公告入库
	notice_info.Id = 0
	operation_userid := c.GetString(ADMIN_UID)
	notice_info.Operator, _ = strconv.ParseInt(operation_userid, 10, 64)

	_, err = dao.NoticeIns.AddLanterns(&notice_info, &notice.WalkingLanterns)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	fmt.Println(notice_time_ptr, notice_info)
	// 公告是否直接推送
	if notice_time_ptr == nil || notice_info.NoticeType != proto.NOTICE_BEFOR_UPDATE {
		// 清理旧的公告及走马灯
		message.G_WalkingLantenrns.Clean(notice_info.Id)
		// 广播
		BroadCastNotice(notice_info)
		// 更新公告为完成
		_, err := dao.NoticeIns.UpdateStatus(notice_info.Id, 1)
		if err != nil {
			logger.Error(err)
		}

	} else {
		// 公告延迟推送
		logger.Debugf("AddNotice:", notice_info.Id, notice.WalkingLanterns)
		// 清理旧的公告及走马灯
		message.G_WalkingLantenrns.Clean(notice_info.Id)
		notice_count := 0
		// 启动定时器,定时发送走马灯消息
		if len(notice.WalkingLanterns) > 0 {
			notice_time := *notice_time_ptr

			for _, item := range notice.WalkingLanterns {
				item.StartTime = notice_time.Add(-1 * time.Second * time.Duration(item.IntervalTime))
				ret := message.G_WalkingLantenrns.AddTimeNotice(item.NoticeId, item)
				if ret != -1 {
					notice_count++
				}
			}
		}
		// 没用可用跑马灯
		if notice_count <= 0 {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_ITMES_EXPIRED)
			return
		}
		// 启动公告定时器
		message.G_WalkingLantenrns.AddNotice(*notice_time_ptr, notice_info)
	}

	JsonData(c, notice_info.Id)
}

func UpdateNotice(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	var notice proto.NoticeReq
	err = json.Unmarshal(data, &notice)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	notice_info := model.Notice{}
	err = utils.CopyFields(&notice_info, notice)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	var notice_time_ptr *time.Time = nil
	if notice_info.NoticeType == proto.NOTICE_BEFOR_UPDATE {
		if len(notice_info.NoticeTime) > 0 {
			time_value, err := utils.ParseTime(notice_info.NoticeTime)
			if err != nil {
				logger.Error(err)
				JsonErrorCode(c, ecode.ERROR_ILLEGAL)
				return
			}
			notice_time_ptr = &time_value
		} else {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
	}

	_, err = dao.NoticeIns.UpdateLanterns(&notice_info, &notice.WalkingLanterns)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	// 检查是否直接通知
	if notice_time_ptr == nil || notice_info.NoticeType != proto.NOTICE_BEFOR_UPDATE {
		// 清理旧的公告及走马灯
		message.G_WalkingLantenrns.Clean(notice_info.Id)

		// 广播公告
		// 广播
		BroadCastNotice(notice_info)

		// 更新公告为完成
		_, err = dao.NoticeIns.UpdateStatus(notice_info.Id, 1)
		if err != nil {
			logger.Error(err)
		}
	} else {
		// 清理旧的公告及走马灯
		message.G_WalkingLantenrns.Clean(notice_info.Id)
		notice_count := 0
		// 启动定时器,定时发送走马灯消息
		if len(notice.WalkingLanterns) > 0 {
			notice_time := *notice_time_ptr
			for _, item := range notice.WalkingLanterns {
				item.StartTime = notice_time.Add(-1 * time.Second * time.Duration(item.IntervalTime))
				ret := message.G_WalkingLantenrns.AddTimeNotice(item.NoticeId, item)
				if ret != -1 {
					notice_count++
				}
			}
		}
		// 没用可用跑马灯
		if notice_count <= 0 {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_ITMES_EXPIRED)
			return
		}
		// 启动公告定时器
		message.G_WalkingLantenrns.AddNotice(*notice_time_ptr, notice_info)
	}

	JsonData(c, notice_info.Id)
}

func GetActivityInfo(c *gin.Context) {
	data, err := dao.ActivityConfigIns.GetAllData()
	if err != nil {
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]proto.ActivityStatusRsp, 0)
	for _, item := range data {
		result = append(result, proto.ActivityStatusRsp{
			Name:      item.Name,
			Type:      item.ActivityType,
			StartTime: utils.Time2Str(item.StartTime),
			EndTime:   utils.Time2Str(item.FinishTime),
			Status:    message.G_ActivityManage.GetActivityStatus(item.ActivityType),
		})
	}
	JsonData(c, result)
}

func SetActivityInfo(c *gin.Context) {
	activity_type_ := c.Query("activity_type")

	activity_type, err := strconv.Atoi(activity_type_)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_ILLEGAL)
		return
	}
	start_time_ := c.Query("start_time")
	var start_time time.Time = time.Now()
	if len(start_time_) > 0 {
		start_time, err = utils.ParseTime(start_time_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
	}

	end_time_ := c.Query("end_time")
	var end_time time.Time
	if len(end_time_) > 0 {
		end_time, err = utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
	} else {
		end_time = time.Now().Add(time.Duration(message.G_BaseCfg.Backstage.MaxActivityYear) * time.Hour * 24 * 365)
	}

	config := model.ActivityConfig{
		ActivityType: activity_type,
		StartTime:    start_time,
		FinishTime:   end_time,
	}

	activity, err := dao.ActivityConfigIns.GetData(activity_type)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if activity.ActivityType != activity_type {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_OBJ_NOT_EXISTS)
		return
	}
	logger.Debug(config)
	_, err = dao.ActivityConfigIns.Update(config)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	err = message.G_ActivityManage.ChangeTime(config, true)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
		return
	}
	data, err := dao.ActivityConfigIns.GetAllData()
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]proto.ActivityStatusRsp, 0)
	for _, item := range data {
		result = append(result, proto.ActivityStatusRsp{
			Name:      item.Name,
			Type:      item.ActivityType,
			StartTime: utils.Time2Str(item.StartTime),
			EndTime:   utils.Time2Str(item.FinishTime),
			Status:    message.G_ActivityManage.GetActivityStatus(item.ActivityType),
		})
	}
	JsonData(c, result)
}

func GetLastDayStatistics(c *gin.Context) {
	result, err := dao.StatisticsDayIns.GetLastData(0)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if result.Id == 0 {
		JsonData(c, nil)
		return
	}
	JsonData(c, gin.H{
		"id":                  result.Id,
		"date":                utils.GetTimeDay(result.Date),
		"platform":            result.Platform,
		"total_reg_count":     result.TotalRegCount,
		"total_newly_added":   result.TotalNewlyAdded,
		"day_registe_count":   result.DayRegisteCount,
		"day_newly_added":     result.DayNewlyAdded,
		"day_login_count":     result.DayLoginCount,
		"active_count":        result.ActiveCount,
		"online_count":        result.OnlineTime,
		"avg_online_count":    result.AvgOnlineTime,
		"phone_registe_count": result.PhoneRegisteCount,
		"email_registe_count": result.EmailRegisteCount,
		"real_name_count":     result.RealNameCount,
		"day_cdt":             result.DayCdt,
		"total_cdt":           result.TotalCdt,
	})
}

func GetDayStatistics(c *gin.Context) {
	platform_ := c.Query("platform")
	var platform_ptr *int = nil
	if len(platform_) > 0 {
		platform, err := strconv.Atoi(platform_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		platform_ptr = &platform
	}

	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDayIns.GetData(start_time_ptr, end_time_ptr, platform_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		result = append(result, gin.H{
			"id":                  item.Id,
			"date":                utils.GetTimeDay(item.Date),
			"platform":            item.Platform,
			"total_reg_count":     item.TotalRegCount,
			"total_newly_added":   item.TotalNewlyAdded,
			"day_registe_count":   item.DayRegisteCount,
			"day_newly_added":     item.DayNewlyAdded,
			"day_login_count":     item.DayLoginCount,
			"active_count":        item.ActiveCount,
			"online_count":        item.OnlineTime,
			"avg_online_count":    item.AvgOnlineTime,
			"phone_registe_count": item.PhoneRegisteCount,
			"email_registe_count": item.EmailRegisteCount,
			"real_name_count":     item.RealNameCount,
			"day_cdt":             item.DayCdt,
			"total_cdt":           item.TotalCdt,
		})
	}
	JsonPage(c, result, page, size, total)
}

func GetLastStatisticsRetained(c *gin.Context) {
	data, err := dao.StatisticsRetainedIns.GetLastData()
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if data.Id == 0 {
		JsonData(c, nil)
		return
	}
	result := utils.StructToJsonMap(data)
	result["date"] = utils.GetTimeDay(result["date"].(time.Time))
	result["date1"] = utils.GetTimeDay(result["date1"].(time.Time))
	result["date3"] = utils.GetTimeDay(result["date3"].(time.Time))
	result["date7"] = utils.GetTimeDay(result["date7"].(time.Time))
	result["date15"] = utils.GetTimeDay(result["date15"].(time.Time))
	result["date30"] = utils.GetTimeDay(result["date30"].(time.Time))
	result["retained1_ratio"] = strconv.Itoa(result["retained1_ratio"].(int)) + "%"
	result["retained3_ratio"] = strconv.Itoa(result["retained3_ratio"].(int)) + "%"
	result["retained7_ratio"] = strconv.Itoa(result["retained7_ratio"].(int)) + "%"
	result["retained15_ratio"] = strconv.Itoa(result["retained15_ratio"].(int)) + "%"
	result["retained30_ratio"] = strconv.Itoa(result["retained30_ratio"].(int)) + "%"

	JsonData(c, result)
}
func GetStatisticsRetained(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsRetainedIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		value["date1"] = utils.GetTimeDay(value["date1"].(time.Time))
		value["date3"] = utils.GetTimeDay(value["date3"].(time.Time))
		value["date7"] = utils.GetTimeDay(value["date7"].(time.Time))
		value["date15"] = utils.GetTimeDay(value["date15"].(time.Time))
		value["date30"] = utils.GetTimeDay(value["date30"].(time.Time))
		value["retained1_ratio"] = strconv.Itoa(value["retained1_ratio"].(int)) + "%"
		value["retained3_ratio"] = strconv.Itoa(value["retained3_ratio"].(int)) + "%"
		value["retained7_ratio"] = strconv.Itoa(value["retained7_ratio"].(int)) + "%"
		value["retained15_ratio"] = strconv.Itoa(value["retained15_ratio"].(int)) + "%"
		value["retained30_ratio"] = strconv.Itoa(value["retained30_ratio"].(int)) + "%"
		result = append(result, value)
	}
	JsonPage(c, result, page, size, total)
}
func GetStatisticsActiveCount(c *gin.Context) {
	data, err := dao.StatisticsActiveCountIns.GetData(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}
	JsonData(c, result)
}

func GetLastStatisticsTreasureBox(c *gin.Context) {
	data, err := dao.StatisticsTreasureBoxIns.GetLastData()
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	if data.Id == 0 {
		JsonData(c, nil)
		return
	}
	result := utils.StructToJsonMap(data)
	result["date"] = utils.GetTimeDay(result["date"].(time.Time))
	JsonData(c, result)
}
func GetStatisticsTreasureBox(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsTreasureBoxIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}
	JsonPage(c, result, page, size, total)
}
func GetStatisticsRealTreasureBox(c *gin.Context) {
	data, err := dao.StatisticsRealTreasureBoxIns.GetData(time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}
	JsonData(c, result)
}

func GetStatisticsSignIn(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsSignInIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}
	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearCdt(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}
	pv, err := dao.StatisticsDoubleYearCdtIns.GetTotalPV(start_time_ptr, end_time_ptr)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	uv, err := dao.StatisticsDoubleYearCdtIns.GetTotalUV(start_time_ptr, end_time_ptr)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	cdt, err := dao.StatisticsDoubleYearCdtIns.GetTotalCDT(start_time_ptr, end_time_ptr)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}
	total, page, data, err := dao.StatisticsDoubleYearCdtIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	var response proto.StatisticsDoubleYearCdtRsp
	response.Code = errcode.HTTP_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.Data = result
	response.Page = page
	response.Size = size
	response.Total = total
	response.Pv = pv
	response.Uv = uv
	response.Cdt = cdt
	c.JSON(http.StatusOK, response)
	// JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearUserCdt(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	user_info := c.Query("key")

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearUserDayCdtIns.GetData(user_info, start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	// for i := 0; i < len(result); i++ {
	// 	result[i].TotalCdt, _ = dao.StatisticsDoubleYearUserDayCdtIns.GetTotalCdt(result[i].UserId, start_time_ptr, end_time_ptr)
	// }
	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearFragment(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearFragmentIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearUserFragment(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	user_info := c.Query("key")

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearUserFragmentIns.GetData(user_info, start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearDailyRanking(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearDailyRankingIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearUserDailyRanking(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	user_info := c.Query("key")

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearUserDailyRankingIns.GetData(user_info, start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearTotalRanking(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearTotalRankingIns.GetData(start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}

func GetStatisticsDoubleYearUserTotalRanking(c *gin.Context) {
	start_time_ := c.Query("start_time")
	var start_time_ptr *time.Time = nil
	if len(start_time_) > 0 {
		start_time, err := utils.ParseTime(start_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		start_time_ptr = &start_time
	}
	end_time_ := c.Query("end_time")
	var end_time_ptr *time.Time = nil
	if len(end_time_) > 0 {
		end_time, err := utils.ParseTime(end_time_)
		if err != nil {
			JsonErrorCode(c, ecode.ERROR_ILLEGAL)
			return
		}
		end_time_ptr = &end_time
	}

	user_info := c.Query("key")

	page_ := c.Query("page")
	size_ := c.Query("size")

	page := 1
	if len(page_) > 0 {
		v, err := strconv.Atoi(page_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			page = v
		}
	}
	size := 10
	if len(size_) > 0 {
		v, err := strconv.Atoi(size_)
		if err != nil {
			logger.Error(err)
			JsonErrorCode(c, ecode.ERROR_PARAM_ILEGAL)
			return
		}
		if v > 0 {
			size = v
		}
	}

	total, page, data, err := dao.StatisticsDoubleYearUserTotalRankingIns.GetData(user_info, start_time_ptr, end_time_ptr, page, size)
	if err != nil {
		logger.Error(err)
		JsonErrorCode(c, ecode.ERROR_MYSQL)
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range data {
		value := utils.StructToJsonMap(item)
		value["date"] = utils.GetTimeDay(value["date"].(time.Time))
		result = append(result, value)
	}

	JsonPage(c, result, page, size, total)
}
