package message

import (
	"game_server/core/base"
	"io/ioutil"

	"game_server/core/logger"

	"gopkg.in/yaml.v2"
)

func ReloadConfig() {
	logger.Debugf("ReloadConfig in")
	ReadCfg(G_BaseCfg)
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err == nil {
		yaml.Unmarshal(yamlFile, base.Setting)
	}

	base.LogInit(base.Setting.Server.Debug, "game_server")
	logger.Infof("base.json value:", G_BaseCfg)
	logger.Infof("config.yml value:", base.Setting)
	ReloadDb()
	logger.Debugf("ReloadConfig end")
}
