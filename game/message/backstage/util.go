package backstage

import (
	"game_server/game/errcode"
	"game_server/game/proto"
	"net/http"
	"time"

	"game_server/game/message"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 返回错误结果
func JsonErrorCode(c *gin.Context, ecode int32) {
	var response proto.S2CCommon
	response.Code = ecode
	response.Message = errcode.ERROR_MSG[response.Code]
	c.JSON(http.StatusOK, response)
	return
}

// 返回成功结果
func JsonOK(c *gin.Context) {
	var response proto.S2CCommon
	response.Code = errcode.HTTP_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	c.JSON(http.StatusOK, response)
}

// 返回成功结果
func JsonData(c *gin.Context, data interface{}) {
	var response proto.S3C_HTTP
	response.Code = errcode.HTTP_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.Data = data
	c.JSON(http.StatusOK, response)
}

// 返回分页结果
func JsonPage(c *gin.Context, data interface{}, page, size, total int) {
	var response proto.S3C_PAGE
	response.Code = errcode.HTTP_SUCCESS
	response.Message = errcode.ERROR_MSG[response.Code]
	response.Data = data
	response.Page = page
	response.Size = size
	response.Total = total
	c.JSON(http.StatusOK, response)
}

var (
	ADMIN_UID = "uid"
)

// 创建token
func CreateToken(uid, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		ADMIN_UID: uid,
		"exp":     time.Now().Add(time.Duration(message.G_BaseCfg.Backstage.TokenExpirationTime) * time.Second).Unix(),
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析token
func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)[ADMIN_UID].(string), nil
}
