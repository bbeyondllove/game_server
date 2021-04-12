package base

import (
	"game_server/core/logger"
)

func LogInit(debug bool, svcName string) {
	// 设置输出调用信息
	logger.SetCallInfo(true)
	if debug {
		logger.SetLevel(logger.DEBUG)
		logger.SetConsole(true)
	} else {
		logger.SetLevel(logger.INFO)
		logger.SetConsole(false)
	}
}
