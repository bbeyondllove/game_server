package backstage

import (
	"game_server/core/base"

	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"game_server/game/message"
)

// @title 管理后台http接口
// @version 1.0.0
// @description 管理后台http接口
// @license.name Apache 2.0
// @host http://47.106.234.171/:8092
// @BasePath /
func HttpServer() {
	AdminManage_Init()

	r := gin.Default()
	r.Use(message.Cors())

	AdminManage_Router(r)

	r.Run(":" + base.Setting.Admin.HttpPort)
}

func AdminManage_Init() {
	// 启动走马灯通知
	StartLanternsNotice()
}

// 运营管理后台
// 路由
func AdminManage_Router(router *gin.Engine) {
	group := router.Group("/v1")
	group.Use(ginprom.PromMiddleware(nil))
	group.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	group.POST("/login", AdminLogin)
	api := group.Group("/api")
	api.Use(AuthHandler())
	api.GET("/logout", AdminLimitHandler(""), AdminLogout)
	api.GET("/get/admin/users", AdminLimitHandler("4001"), GetAdminUsers)
	api.GET("/add/admin/user", AdminLimitHandler("4001"), AddAdminUser)
	api.GET("/update/admin/user", AdminLimitHandler("4001"), UpdateAdminUser)
	api.GET("/update/admin/user/status", AdminLimitHandler("4001"), UpdateAdminUserStatus)
	api.GET("/delete/admin/user", AdminLimitHandler("4001"), DeleteAdminUser)
	api.GET("/get/admin/roles", AdminLimitHandler("4001"), GetAdminRole)
	api.GET("/add/admin/role", AdminLimitHandler("4001"), AddAdminRole)
	api.GET("/update/admin/role", AdminLimitHandler("4001"), UpdateAdminRole)
	api.GET("/delete/admin/role", AdminLimitHandler("4001"), DeleteAdminRole)
	api.GET("/get/admin/limits", AdminLimitHandler("4001"), GetAdminLimit)

	api.GET("/get/users", AdminLimitHandler("5001"), GetUsers)
	api.GET("/frozen/user", AdminLimitHandler("5001"), FrozenUser)
	api.GET("/get/certification", AdminLimitHandler("5002"), GetCertifications)
	api.GET("/certification/user", AdminLimitHandler("5002"), CertificationUser)
	api.GET("/get/user/certification/record", AdminLimitHandler("5002"), GetCertificationRecord)

	api.GET("/get/tasks", AdminLimitHandler("5301"), GetTasks)
	api.GET("/get/items", AdminLimitHandler("5301"), GetItems)
	api.GET("/set/task/award", AdminLimitHandler("5301"), SetTaskAward)

	api.POST("/upload/image", AdminLimitHandler("5303"), UploadImage)
	api.GET("/get/notices", AdminLimitHandler("5303"), GetNotice)
	api.POST("/add/notice", AdminLimitHandler("5303"), AddNotice)
	api.POST("/update/notice", AdminLimitHandler("5303"), UpdateNotice)

	api.GET("/get/activity/info", AdminLimitHandler("5302"), GetActivityInfo)
	api.GET("/set/activity/info", AdminLimitHandler("5302"), SetActivityInfo)

	api.GET("/get/lastday/statistics", AdminLimitHandler("5101"), GetLastDayStatistics)
	api.GET("/get/day/statistics", AdminLimitHandler("5101"), GetDayStatistics)

	api.GET("/get/lastday/retained/statistics", AdminLimitHandler("5102"), GetLastStatisticsRetained)
	api.GET("/get/retained/statistics", AdminLimitHandler("5102"), GetStatisticsRetained)
	api.GET("/get/activecount/statistics", AdminLimitHandler("5102"), GetStatisticsActiveCount)

	api.GET("/get/lastday/treasurebox/statistics", AdminLimitHandler("5103"), GetLastStatisticsTreasureBox)
	api.GET("/get/treasurebox/statistics", AdminLimitHandler("5103"), GetStatisticsTreasureBox)
	api.GET("/get/realtreasurebox/statistics", AdminLimitHandler("5103"), GetStatisticsRealTreasureBox)

	api.GET("/get/doubleyearcdt/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearCdt)
	api.GET("/get/doubleyearcdt/user/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearUserCdt)

	api.GET("/get/doubleyear/fragment/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearFragment)
	api.GET("/get/doubleyear/user/fragment/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearUserFragment)

	api.GET("/get/doubleyear/dailyranking/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearDailyRanking)
	api.GET("/get/doubleyear/user/dailyranking/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearUserDailyRanking)

	api.GET("/get/doubleyear/totalranking/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearTotalRanking)
	api.GET("/get/doubleyear/user/totalranking/statistics", AdminLimitHandler("5103"), GetStatisticsDoubleYearUserTotalRanking)

	api.GET("/get/signin/statistics", AdminLimitHandler("5103"), GetStatisticsSignIn)
}
