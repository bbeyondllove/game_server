package message

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/logic"
	"game_server/game/model"
	"game_server/game/proto"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	"game_server/core/logger"

	"github.com/gin-gonic/gin"
)

//允许跨域
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func IpLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Debugf("gin.HandlerFunc start")
		//从redis中获取ip限制数据
		whiteLists, err := db.RedisGame.Get("ip_limits").Result()
		if err != nil {
			logger.Errorf("Get('ip_limits') failed err=", err.Error())
		}
		vistorIp := c.ClientIP()
		if whiteLists != "*" || !strings.Contains(whiteLists, vistorIp) {
			c.Set("ipPass", 0)
			logger.Errorf("IpLimit() not pass vistorIp=, whiteLists=", vistorIp, whiteLists)
		}
		c.Set("ipPass", 1)
		logger.Debugf("gin.HandlerFunc end")
	}
}

// @title 游戏http接口
// @version 1.6.7
// @description 游戏http接口
// @license.name Apache 2.0
// @host http://47.106.234.171/:8092
// @BasePath /
func HttpServer() {
	r := gin.Default()
	r.Use(Cors())

	r.POST("/api/registerByMobile", IpLimit(), registerByMobile)
	r.POST("/api/registerByEmail", IpLimit(), registerByEmail)
	r.POST("/api/getVerificationCode", IpLimit(), getVerificationCode)
	r.POST("/account/transfer", IpLimit(), HandleTransfer)
	r.GET("/account/balance", IpLimit(), HandleBalance)
	r.GET("/api/changeRecordList", IpLimit(), ChangeRecordList)
	r.POST("/api/memberIsActive", IpLimit(), MemberIsActive)
	r.POST("/api/currencySupport", IpLimit(), CurrencySupport)
	r.GET("/api/getUserSum", IpLimit(), getUserSum)
	r.POST("/api/sendEmail", IpLimit(), SendEmail)
	r.POST("/api/sendNotice", IpLimit(), PushNotice)
	r.POST("/api/kycSync", IpLimit(), kycSync)
	r.POST("/api/updateTaskStatus", IpLimit(), updateTaskStatus)

	// AdminManage_Router(r)

	r.Run(":" + base.Setting.Server.HttpPort)

}

// @Summary 会员是否激活
// @Description 用于查询会员是否激活
// @Produce  json
// @Accept  json
// @Param data body model.MemberIsActiveReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /api/MemberIsActive [post]

func MemberIsActive(c *gin.Context) {
	logger.Debugf("MemberIsActive start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data = make(map[string]interface{}, 0)

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// ipPass, exist := c.Get("ipPass")
	// if !exist {
	// 	logger.Errorf("c.Get('ipPass'), err")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// 	return
	// }
	// pass := ipPass.(int)
	// if pass != 1 {
	// 	logger.Errorf("ipPass is 0")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// }

	var cdr model.MemberIsActiveReq
	err := c.BindJSON(&cdr)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	// genToken := genTonken(c)
	// if cdr.Signature != genToken {
	// 	responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
	// 	responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
	// 	c.JSON(http.StatusOK, responseMessage)
	// 	return
	// }

	userData, err := db_service.UserIns.GetDataByUid(cdr.UserId)

	if userData.UserId == "" || userData.Status == 0 {
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		responseMessage.Data["data"] = false
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data["data"] = true
	c.JSON(http.StatusOK, responseMessage)
	return
}

// @Summary 是否支持某币种的划转
// @Description 用于查询是否支持某币种的划转
// @Produce  json
// @Accept  json
// @Param data body model.CurrencySupportReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /api/CurrencySupport [post]
func CurrencySupport(c *gin.Context) {
	logger.Debugf("CurrencySupport start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// ipPass, exist := c.Get("ipPass")
	// if !exist {
	// 	logger.Errorf("c.Get('ipPass'), err")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// 	return
	// }
	// pass := ipPass.(int)
	// if pass != 1 {
	// 	logger.Errorf("ipPass is 0")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// }

	var cdr model.CurrencySupportReq
	err := c.BindJSON(&cdr)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	// genToken := genTonken(c)
	// if cdr.Signature != genToken {
	// 	responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
	// 	responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
	// 	c.JSON(http.StatusOK, responseMessage)
	// 	return
	// }

	if strings.ToUpper(cdr.Currency) == "SAN" {
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		responseMessage.Data["data"] = true
		c.JSON(http.StatusOK, responseMessage)
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data["data"] = false
	c.JSON(http.StatusOK, responseMessage)
	return
}

// @Summary 查询余额
// @Description 查询余额
// @Produce  json
// @Accept  json
// @Param data body model.BalanceReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":{}}
// @Router /account/balance [get]
func HandleBalance(c *gin.Context) {
	logger.Debugf("HandleBalance start")

	var responseMessage model.BalanceRsp
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data = make(map[string]decimal.Decimal, 0)
	responseMessage.Data["availableBalance"] = decimal.NewFromInt(0)
	responseMessage.Data["totalBalance"] = decimal.NewFromInt(0)

	// if c.Request.Body == nil {
	// 	c.JSON(http.StatusOK, responseMessage)
	// 	return
	// }

	// ipPass, exist := c.Get("ipPass")
	// if !exist {
	// 	logger.Errorf("c.Get('ipPass'), err")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// 	return
	// }
	// pass := ipPass.(int)
	// if pass != 1 {
	// 	logger.Errorf("ipPass is 0")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// }

	var br model.BalanceReq
	if err := c.ShouldBindQuery(&br); err != nil {
		logger.Debugf("ShouldBindJSON failed, err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	logger.Debugf("BalanceReq=", br)

	requestMap := make(map[string]interface{}, 0)
	requestMap["nonce"] = br.Nonce
	requestMap["userId"] = br.UserId
	requestMap["tokenCode"] = br.TokenCode

	genToken := utils.GenTonken(requestMap)
	logger.Debugf("genToken=", genToken)
	if br.Sign != genToken {
		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	wallet, err := db_service.GameWalletIns.GetAmount(br.UserId, br.TokenCode)
	if err != nil {
		logger.Errorf("HandleGetAmount db.GetAmount() failed(), err=", err.Error())
		responseMessage.Code = errcode.ERROR_MYSQL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		return
	}
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.Data["availableBalance"] = decimal.NewFromFloat(wallet.AmountAvailable)
	responseMessage.Data["totalBalance"] = decimal.NewFromFloat(wallet.Amount)

	logger.Infof("HandleBalance end %+v", responseMessage)
	c.JSON(http.StatusOK, responseMessage)
	return
}

// @Summary 划转资产通证
// @Description 该接口用于将用户的资产通证从业务系统划转到Base系统，或从Base系统划转到业务系统。接口需要保证幂等性，即同一个 txId 的多次请求，只是一次划转交易，而不是多次划转，且每次请求需要返回一样的结果，如该笔划转交易已成功，则多次请求的结果都返回成功。
//该接口也使用签名校验。
// @Produce  json
// @Accept  json
// @Param data body model.TransferReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":{}}
// @Router /account/transfer [post]
func HandleTransfer(c *gin.Context) {
	logger.Debugf("HandleTransfer start")

	var responseMessage model.HttpCommon
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// ipPass, exist := c.Get("ipPass")
	// if !exist {
	// 	logger.Errorf("c.Get('ipPass'), err")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// 	return
	// }
	// pass := ipPass.(int)
	// if pass != 1 {
	// 	logger.Errorf("ipPass is 0")
	// 	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	// 	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	// }

	var cdr model.TransferReq
	if err := c.ShouldBindJSON(&cdr); err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	logger.Debugf("TransferReq=", cdr)

	/*
		//token 检验
		user_info := db.RedisMgr.HGetAll(cdr.UserId)
		if len(user_info) == 0 {
			logger.Errorf("HandleTransfer not login:%+v", cdr.UserId)
			responseMessage.Code = errcode.ERROR_NOT_LOGIN
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			c.JSON(http.StatusOK, responseMessage)
			return
		}
	*/

	//参数校验
	bFound := false
	for _, value := range G_BaseCfg.TokenCode {
		if cdr.TokenCode == value {
			bFound = true
			break
		}
	}
	if !bFound {
		responseMessage.Code = errcode.ERROR_HTTP_FORBIDDEN
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	requestMap := make(map[string]interface{}, 0)
	requestMap["txNo"] = cdr.TxNo
	requestMap["userId"] = cdr.UserId
	requestMap["tokenCode"] = cdr.TokenCode
	requestMap["amount"] = cdr.Amount.String()
	requestMap["nonce"] = cdr.Nonce

	genToken := utils.GenTonken(requestMap)
	logger.Debugf("genToken=", genToken)
	if cdr.Sign != genToken {
		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	//从数据库中查询
	wallet, err := db_service.GameWalletIns.GetAmount(cdr.UserId, cdr.TokenCode)
	if err != nil {
		logger.Errorf("HandleGetAmount db.GetAmount() failed(), err=", err.Error())

		responseMessage.Code = errcode.ERROR_WALLET
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_WALLET]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	amountAvailable := wallet.AmountAvailable
	amount := wallet.Amount

	cdrAmount, ret := cdr.Amount.Float64()
	if !ret {
		logger.Errorf("cdr.Amount.Float64 failed cdr.Amount=", cdr.Amount)
	}
	if cdrAmount < 0 && amountAvailable < math.Abs(cdrAmount) { //扣款前需要看钱够不够
		logger.Errorf("amountAvailable is not enough=", err.Error())

		responseMessage.Code = errcode.ERROR_HTTP_BALANCE_NOT_ENOUGH
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	//是否是重复的txid
	record, err := db_service.ExchangeRecordIns.GetExchangeByTxid(cdr.TxNo)
	if err != nil || record.SysOrderSn != "" {
		responseMessage.Code = errcode.ERROR_HTTP_REPEAT_TXID
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	ret, session := db_service.GameWalletIns.UpdateMoney(cdr.UserId, cdr.TokenCode, cdrAmount)
	defer session.Close()
	if !ret {
		logger.Errorf("UpdateMoney failed uid=", cdr.UserId)
		responseMessage.Code = errcode.ERROR_UPDATE_MONEY
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_UPDATE_MONEY]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	//更新redis中值
	amountInfo := make(map[string]interface{})
	amountInfo["Amount"] = amount
	amountInfo["AmountAvailable"] = amountAvailable

	value, err := db.RedisGame.HMSet(cdr.UserId, amountInfo).Result()
	if err != nil {
		responseMessage.Code = errcode.ERROR_REDIS
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		session.Rollback()
		logger.Errorf("RedisGame.HMSet failed uid=,err=", cdr.UserId, err.Error())
		return
	}
	if value == "OK" {
		logger.Debugf("RedisGame.HMSet success, uid=", cdr.UserId)
	}

	orderSn := utils.CreateOrderSn("Game")
	status, ExchangeType := 1, "charge"
	if cdrAmount < 0 { //扣款
		status = 2
		ExchangeType = "deduce"
	}

	//增加交易记录
	er := &model.ExchangeRecord{
		SysOrderSn:      cdr.TxNo,
		OrderSn:         orderSn,
		UserId:          cdr.UserId,
		ExchangeType:    ExchangeType,
		CurrencyType:    cdr.TokenCode,
		Amount:          cdrAmount,
		Status:          status,
		UserAmount:      cdrAmount,
		AmountAvailable: amountAvailable,
		TargetAccount:   1,
		Desc:            "",
		AdminUser:       "base",
		AdminUserId:     1,
		CreateTime:      time.Now(),
		UpdateTime:      time.Now(),
	}
	_, err = db_service.ExchangeRecordIns.Add(er)
	if err != nil {
		logger.Errorf("ExchangeRecordIns.Add failed uid=,err=", cdr.UserId, err.Error())
		session.Rollback()
		responseMessage.Code = errcode.ERROR_ADD_EXCHANGE_RECORD
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_ADD_EXCHANGE_RECORD]
		c.JSON(http.StatusOK, responseMessage)

		return
	}
	err = session.Commit()

	if err != nil {
		logger.Errorf("ChargeDeduce() session.Commit() failed uid=, err=", cdr.UserId, err.Error())

		responseMessage.Code = errcode.ERROR_MYSQL_COMMIT
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_MYSQL_COMMIT]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	logger.Infof("HandleTransfer end %+v", responseMessage)
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[errcode.MSG_SUCCESS]
	c.JSON(http.StatusOK, responseMessage)
	return
}

// @Summary 交易对帐
// @Description 用于游戏子钱包交易对账，基于性能考虑，接口每次最多返回 500 条数据记录，需请求言循环调用。
// @Produce  json
// @Accept  json
// @Param data body model.ChangeRecordListReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /api/ChangeRecordList [get]
func ChangeRecordList(c *gin.Context) {
	logger.Debugf("ChangeRecordList start")
	var responseMessage proto.S2CHTTP

	//用中间件进行IP限制
	ipPass, exist := c.Get("ipPass")
	if !exist {
		logger.Errorf("c.Get('ipPass'), err")
		responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		return
	}
	pass := ipPass.(int)
	if pass != 1 {
		logger.Errorf("ipPass is 0")
		responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	}
	//todo，调用时间间隔限制

	if c.Request.Body == nil {
		responseMessage.Code = 1
		responseMessage.Message = "body is nil"
		c.JSON(http.StatusOK, responseMessage)

		return
	}

	var crr model.ChangeRecordListReq
	err := c.BindJSON(&crr)
	if err != nil {
		logger.Errorf("Decode failed uid=,err=", crr.UserId, err.Error())

		responseMessage.Code = 2
		responseMessage.Message = err.Error()
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// deadLineStr := strconv.FormatInt(crr.DeadLine, 10)
	// lastRecordIdStr := strconv.FormatInt(crr.LastRecordId, 10)

	var dataMap = make(map[string]interface{}, 0)
	msg := make([]byte, proto.MAX_BUF_SIZE)
	_, err = c.Request.Body.Read(msg)
	if err != nil {
		logger.Errorf("Read error=", err.Error())
		return
	}
	err = json.Unmarshal(msg, &dataMap)
	if err != nil {
		logger.Errorf("json.Unmarshal error")
		return
	}
	genToken := utils.GenTonken(dataMap)
	if crr.Signature != genToken {
		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	datas, err := db_service.ExchangeRecordIns.GetExchangeList(crr.UserId, crr.LastRecordId, crr.DeadLine)
	if err != nil {
		logger.Errorf("ExchangeRecordIns.GetExchangeList() failed uid=,err=", crr.UserId, err.Error())

		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[errcode.MSG_SUCCESS]
	responseMessage.Data = datas
	c.JSON(http.StatusOK, responseMessage)
	logger.Debugf("ChangeRecordList end")

}

// 手机注册
func registerByMobile(c *gin.Context) {
	logger.Debugf("registerByMobile start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	var msg proto.C2SRegisterMobile
	err := c.BindJSON(&msg)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	if msg.NickNname == "" {
		msg.NickNname = msg.Mobile
	}

	request_map := utils.StructToMap(msg)
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	buf, err := json.Marshal(request_map)

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.MobileRegisterUrl
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		logger.Errorf("HttpPost err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	logger.Debugf("registerByMobile send to client:", string(msgdata))
	if msg.Inviter != "" {
		//绑定邀请码
		var invitaLogic logic.UserInvitation
		invitaLogic.NewRegisterModelSave(msg, msgdata)
	}
	c.String(http.StatusOK, string(msgdata))

	logger.Debugf("registerByMobile end")
}

//邮箱注册
func registerByEmail(c *gin.Context) {
	logger.Debugf("registerByEmail start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	var msg proto.C2SRegisterEmail
	err := c.BindJSON(&msg)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	if msg.NickNname == "" {
		msg.NickNname = msg.Email
	}

	request_map := utils.StructToMap(msg)
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	buf, err := json.Marshal(request_map)

	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.EmailRegisterUrl
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	if err != nil {
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		logger.Errorf("HttpPost err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	logger.Debugf("registerByEmail send to client:", string(msgdata))
	if msg.Inviter != "" {
		//绑定邀请码
		var invitaLogic logic.UserInvitation
		invitaLogic.NewRegisterEmailSave(msg, msgdata)
	}
	c.String(http.StatusOK, string(msgdata))

	logger.Debugf("registerByEmail end")
}

//获取验证码
func getVerificationCode(c *gin.Context) {
	logger.Debugf("getVerificationCode start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	var msg proto.C2SGetVerificationCode
	err := c.BindJSON(&msg)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	if msg.Language == "" {
		msg.Language = "en"
	}
	Url := ""
	request_map := make(map[string]interface{}, 0)
	if msg.CodeType == proto.CODE_TYPE_EMAIL {
		if msg.Email == "" {
			logger.Errorf("Email == ''")
			c.JSON(http.StatusOK, responseMessage)
			return
		}

		request_map["email"] = msg.Email
		Url = base.Setting.Base.EmailCodeUrl
	} else {
		if msg.CountryCode == 0 || msg.Mobile == "" {
			logger.Errorf("CountryCode ==,Mobile ==  ", msg.CountryCode, msg.Mobile)
			c.JSON(http.StatusOK, responseMessage)
			return

		}
		request_map["mobile"] = msg.Mobile
		request_map["countryCode"] = msg.CountryCode
		Url = base.Setting.Base.MobileCodeUrl
	}

	request_map["useFor"] = msg.UseFor
	request_map["language"] = msg.Language
	request_map["sysType"] = msg.SysType
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["sign"] = utils.GenTonken(request_map)
	buf, _ := json.Marshal(request_map)
	responseMessage.Message = ""
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + Url
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)

	if err != nil {
		responseMessage.Code = errcode.ERROR_USER_SYSTEM
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		logger.Errorf("HttpPost err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	logger.Debugf("getVerificationCode send to client:", string(msgdata))
	c.String(http.StatusOK, string(msgdata))

	logger.Debugf("getVerificationCode end")
}

//同步中台
func kycSyncBase(userId string) {
	logger.Debugf("kycSyncBase in")
	data, err := db_service.CertificationIns.GetDataByUid(userId)
	if err != nil {
		logger.Errorf("GetDataByUid err=", err.Error())
		return
	}

	request_map := make(map[string]interface{}, 0)
	request_map["sysType"] = proto.SysType
	request_map["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	request_map["userId"] = userId
	request_map["nationality"] = data.Nationality
	request_map["firstName"] = data.FirstName
	request_map["lastName"] = data.LastName
	request_map["idType"] = data.IdType
	request_map["idNumber"] = data.IdNumber
	var imgUrl proto.CertificationPhoto
	err = json.Unmarshal([]byte(data.ObjectKey), &imgUrl)
	if err == nil {
		request_map["idPhotoFront"] = imgUrl.Front
		request_map["idPhotoBack"] = imgUrl.Back
		request_map["idPhotoHandheld"] = imgUrl.Other
	}
	request_map["sign"] = utils.GenTonken(request_map)
	buf, _ := json.Marshal(request_map)
	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.UserSaveKycUrl
	msgdata, err := utils.HttpPost(url, string(buf), proto.JSON)
	logger.Debugf("kycSyncBase end err[%+V],return[%+V]", err, string(msgdata))

}

//kyc同步
func kycSync(c *gin.Context) {
	logger.Debugf("kycSync start")

	var responseMessage proto.KycSyncRsp
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	var msg proto.KycSync
	err := c.BindJSON(&msg)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.KycSync = msg
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SEND_KYC_STATUS)
	rsp.WriteData(responseMessage)

	Sched.SendToUser(msg.UserId, rsp)
	logger.Debugf("kycSync send to client:", string(rsp.Bytes()))
	c.JSON(http.StatusOK, responseMessage)

	if msg.KycStatus == 2 {
		kycSyncBase(msg.UserId)
	}
	logger.Debugf("kycSync end")
}

//开启或关闭任务（后台推送给后端）
func updateTaskStatus(c *gin.Context) {
	logger.Debugf("updateTaskStatus start")

	var responseMessage proto.S2CCommon
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	var msg proto.ChangeTaskStatus
	err := c.BindJSON(&msg)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	logger.Errorf("updateTasks msg:%+V", msg)
	if !updateTasks(&msg) {
		logger.Errorf("updateTasks error:")
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_UPDATE_TASK_STATUS)
	rsp.WriteData(responseMessage)

	broadTaskInfo(&msg)
	logger.Debugf("updateTaskStatus send to client:", string(rsp.Bytes()))
	c.JSON(http.StatusOK, responseMessage)

	logger.Debugf("updateTaskStatus end")
}

//用户统计
func getUserSum(c *gin.Context) {
	logger.Debugf("getUserSum start")

	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	allCount, _ := db_service.UserIns.GetAllCount()
	timeCount, _ := db_service.UserIns.GetTimeCount(startTime, endTime)

	var data proto.S2CGetUserSum
	data.Code = errcode.MSG_SUCCESS
	data.Message = errcode.ERROR_MSG[data.Code]
	data.AllCount = allCount
	data.NewCount = timeCount
	logger.Debugf("getUserSum send to client:%+v", data)
	c.JSON(http.StatusOK, data)

}

// @Description 用于邮件发送
// @Produce  json
// @Accept  json
// @Param data body model.SendGameEmailReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /api/sendEmail [post]
func SendEmail(c *gin.Context) {
	logger.Debugf("send email")
	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	var cdr model.SendGameEmailReq
	err = json.Unmarshal(data, &cdr)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	var dataMap = make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		logger.Errorf("json.Unmarshal error")
		return
	}
	// 签名
	delete(dataMap, "signature") // 移除签名
	genToken := utils.GenTonken(dataMap)
	if cdr.Signature != genToken {
		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	// 检查是否已经实名过
	isReaName, err := CheckUserRealName(cdr.UserId, 2)
	if err != nil {
		logger.Debugf("checkUserRealName in err:=", err.Error())
		responseMessage.Code = errcode.ERROR_PUSH_EMAIL
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_PUSH_EMAIL]
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// 已经实名的不需要推送邮件
	if isReaName == true {
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = errcode.ERROR_MSG[errcode.MSG_SUCCESS]
		responseMessage.Data = make(map[string]interface{}, 0)
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	// 添加邮件
	email, err := db_service.EmailIns.AddEmail(cdr.UserId, cdr.EmailType, cdr.EmailTitle, cdr.EmailContent)
	if err != nil {
		responseMessage.Code = errcode.ERROR_PUSH_EMAIL
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_PUSH_EMAIL]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	// 推送邮件
	emailList := make([]proto.EmailItem, 1)
	emailList[0] = proto.EmailItem{
		Id:           email.Id,
		UserId:       email.UserId,
		EmailType:    email.EmailType,
		EmailTitle:   email.EmailTitle,
		EmailContent: email.EmailContent,
		IsRead:       email.IsRead,
		CreateTime:   email.CreateTime.Format("2006/01/02 15:04:05"),
		PrizeList:    make([]model.EmailPrize, 0),
	}
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_USER_MEIAL_RSP)
	pushMessage := &proto.S2CSendEmail{}
	pushMessage.IsPush = 1
	pushMessage.Code = errcode.MSG_SUCCESS
	pushMessage.Message = ""
	pushMessage.NewEmailList = emailList
	logger.Debugf("HandlePushUserEmail end")
	rsp.WriteData(pushMessage)
	Sched.SendToUser(cdr.UserId, rsp)

	// 推送邮件数量
	ret, err := db_service.EmailIns.GetReadAndUnreadNum(cdr.UserId)
	if err == nil {
		// 推送邮件的数量
		rsp := &utils.Packet{}
		rsp.Initialize(proto.MSG_COUNT_EMAIL_RSP)
		responseMessage := &proto.S2CSetEmailRead{}
		responseMessage.Code = errcode.MSG_SUCCESS
		responseMessage.Message = ""
		responseMessage.ReadNum = ret["read"]
		responseMessage.UnreadNum = ret["unread"]
		rsp.WriteData(responseMessage)
		// 推送邮件数量
		Sched.SendToUser(cdr.UserId, rsp)
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[errcode.MSG_SUCCESS]
	responseMessage.Data = make(map[string]interface{}, 0)
	c.JSON(http.StatusOK, responseMessage)
}

// 检测用户实名
func CheckUserRealName(userId string, sysType int) (bool, error) {
	userinfoRequest := make(map[string]interface{}, 0)
	userinfoRequest["userId"] = userId
	userinfoRequest["sysType"] = sysType
	userinfoRequest["nonce"] = utils.GetRandString(proto.RAND_STRING_LEN)
	userinfoRequest["sign"] = utils.GenTonken(userinfoRequest)

	url := base.Setting.Base.UserHost + ":" + base.Setting.Base.UserPort + base.Setting.Base.GetUserinfoUrl
	userdata, err := utils.HttpGet(url, userinfoRequest)
	userResp := &proto.S2C_HTTP{}
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(userdata, userResp)
	if err != nil {
		return false, err
	}
	return userResp.Data["kycPassed"].(bool), nil
}

// @Description 用于推送公告
// @Produce  json
// @Accept  json
// @Param data body model.PushGameNoticeReq true "请求参数"
// @Success 200 {string} json "{"code":0,"message":"OK", "data":[{}]}
// @Router /api/sendNotice [post]
func PushNotice(c *gin.Context) {
	logger.Debugf("push notice")
	var responseMessage proto.S2C_HTTP
	responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	if c.Request.Body == nil {
		c.JSON(http.StatusOK, responseMessage)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	var cdr model.PushGameNoticeReq
	err = json.Unmarshal(data, &cdr)
	if err != nil {
		logger.Errorf("BindJSON err=", err.Error())
		responseMessage.Code = errcode.ERROR_PARAM_ILEGAL
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	var dataMap = make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		logger.Errorf("json.Unmarshal error")
		return
	}
	// 签名
	delete(dataMap, "signature") // 移除签名
	genToken := utils.GenTonken(dataMap)
	if cdr.Signature != genToken {
		responseMessage.Code = errcode.ERROR_HTTP_SIGNATURE
		responseMessage.Message = errcode.ERROR_MSG[errcode.ERROR_HTTP_SIGNATURE]
		c.JSON(http.StatusOK, responseMessage)
		return
	}
	//推送公告信息
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_NOTICE_RSP)
	pushMessage := &proto.S2CNotice{}
	//推送给所有用户
	pushMessage.Code = errcode.MSG_SUCCESS
	pushMessage.Message = ""
	pushMessage.Notice = proto.NoticeInfo{
		NoticeType:    cdr.NoticeType,
		NoticeTitle:   cdr.NoticeTitle,
		NoticeContent: cdr.NoticeContent,
		Version:       cdr.Version,
	}
	logger.Debugf("HandlePushUserEmail end")
	rsp.WriteData(pushMessage)
	// 广播给在线的用户
	Sched.BroadCastMsg(int32(3302), "0", rsp)
	// 给客户发送公告
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[errcode.MSG_SUCCESS]
	responseMessage.Data = make(map[string]interface{}, 0)
	c.JSON(http.StatusOK, responseMessage)
	return
}
