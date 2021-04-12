package backstage

import (
	"game_server/db"
	ecode "game_server/game/errcode"
	"game_server/game/message"
	"net/http"
	"strings"

	"game_server/core/logger"

	"github.com/gin-gonic/gin"
)

// 登陆授权检查
func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		if message.G_BaseCfg.Backstage.TokenSwitch == 1 {
			token := c.Request.Header.Get("token")
			uid, err := ParseToken(token, message.G_BaseCfg.Backstage.TokenSecret)
			if err != nil {
				logger.Error(err)
				JsonErrorCode(c, ecode.ERROR_NOT_LOGIN)
				c.Abort()
				return
			}
			if len(uid) == 0 {
				logger.Error("uid is empty")
				JsonErrorCode(c, ecode.ERROR_NOT_LOGIN)
				c.Abort()
				return
			}
			c.Set(ADMIN_UID, uid)
		}
		c.Next()
	}
}

// 用户权限匹配
func AdminLimitHandler(limit string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(limit) > 0 {
			uid := c.GetString(ADMIN_UID)
			if len(uid) == 0 {
				logger.Errorf("uid is empty")
				JsonErrorCode(c, ecode.ERROR_NOT_LOGIN)
				c.Abort()
				return
			}
			limits, err := db.RedisGame.HGet(ADMIN_USER+uid, "role_limits").Result()
			if err != nil && err.Error() != "redis: nil" {
				logger.Errorf(ADMIN_USER+uid+" role_limits failed:%+v", err.Error())
				JsonErrorCode(c, ecode.ERROR_NOT_LOGIN)
				c.Abort()
				return
			}
			in_array := func(arr []string, item string) bool {
				for _, it := range arr {
					if it == item {
						return true
					}
				}
				return false
			}
			arr_limits := strings.Split(limits, ",")
			if !in_array(arr_limits, limit) {
				logger.Errorf("limit is not exists", arr_limits, limit)
				JsonErrorCode(c, ecode.ERROR_NOT_LOGIN)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
