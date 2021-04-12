package base

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var Cfg Config

type Config struct {
	LogLevel    string        `yaml:"LogLevel"`
	WSAddr      string        `yaml:"WSAddr"`
	TCPAddr     string        `yaml:"TCPAddr"`
	MaxConnNum  int           `yaml:"MaxConnNum"`
	HttpTimeout time.Duration `yaml:"HTTPTimeout"`

	MaxMsgLen  uint32 `yaml:"MaxMsgLen"`
	CertFile   string `yaml:"CertFile"`
	KeyFile    string `yaml:"KeyFile"`
	HTTPListen string `yaml:"HTTPListen"`
}

func init() {

	//ReadFile函数会读取文件的全部内容，并将结果以[]
	filename := GetExecpath() + "/server.yml"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("读取配置文件错误")
	}

	//读取的数据为json格式，需要进行解码
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		log.Fatal("解析配置文件错误")
	}
}

// 获取当前程序运行目录
func GetExecpath() string {
	execpath, _ := os.Executable() // 获得程序路径
	path := filepath.Dir(execpath)
	return strings.Replace(path, "\\", "/", -1)
}
