package goSpider

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Data map[string]interface{}
}

func (config *Config) get(key string, defaultValue interface{}) interface{} {
	if _, exists := config.Data[key]; exists {
		return config.Data[key]
	}
	return defaultValue
}

// 加载配置
func (config *Config) init() {
	pwd, err := os.Getwd() //获取当前目录
	if err != nil {
		panic(err)
	}

	configPath := pwd + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "goSpider.json"
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Panicln("配置文件加载失败：", err)
	}
	err = json.Unmarshal(buf, &config.Data)
	if err != nil {
		log.Panicln("配置文件解析失败：", string(buf), err)
	}

}

// 获取命令行参数
func (config *Config) flagParamInit() {
	flagTest := flag.Bool("test", false, "测试模式")
	flag.Bool("t", false, "")
	flagSpiderCount := flag.Int("spiderCount", 30, "爬虫并发数量")
	flag.Int("sc", 30, "")
	flagDataHandleCount := flag.Int("dataHandleCount", 1000, "数据并行处理数")
	flag.Int("dhc", 1000, "")
	flag.Parse()

	for _, v := range os.Args {
		args := strings.Split(v, "=")
		argName := strings.Replace(args[0], "-", "", -1)
		switch argName {
		case "t", "test":
			if args[1] == "true" || args[1] == "TRUE" {
				*flagTest = true
			} else {
				*flagTest = false
			}
		case "sc", "spiderCount":
			*flagSpiderCount = tool.Str2Int(args[1])
		case "dhc", "dataHandleCount":
			*flagDataHandleCount = tool.Str2Int(args[1])
		}
	}
}
