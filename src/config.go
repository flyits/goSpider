package goSpider

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Data map[string]interface{}
}

func (config *Config) get(key string, defaultValue interface{}) interface{} {
	if _, exists := config.Data[key]; exists {
		if reflect.TypeOf(config.Data[key]).Kind() == reflect.Float64 {
			return int(config.Data[key].(float64))
		}
		return config.Data[key]
	}
	return defaultValue
}

// 加载配置
func (config *Config) init() {
	currentDirectory, err := os.Getwd() //获取当前目录
	if err != nil {
		panic(err)
	}

	configPath := currentDirectory + string(os.PathSeparator) + "Config" + string(os.PathSeparator) + "goSpider.json"
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		config.createDefaultConfigFile(currentDirectory)
	} else {
		err = json.Unmarshal(buf, &config.Data)
		if err != nil {
			log.Panicln("配置文件解析失败：", string(buf), err)
		}
	}
	config.flagParamInit()
}

// 获取命令行参数
func (config *Config) flagParamInit() {

	flag.Int("spiderCount", 30, "爬虫并发数量")
	flag.Int("sc", 30, "")
	flag.Int("dataHandleCount", 1000, "数据并行处理数")
	flag.Int("dhc", 1000, "")
	flag.Parse()

	for _, v := range os.Args {
		args := strings.Split(v, "=")
		argName := strings.Replace(args[0], "-", "", -1)
		switch argName {
		case "sc", "spiderCount":
			config.Data["spiderCount"] = Str2Int(args[1])
		case "dhc", "dataHandleCount":
			config.Data["dataHandleCount"] = Str2Int(args[1])
		}
	}
}

func (config *Config) createDefaultConfigFile(currentDirectory string) {
	configPath := currentDirectory + string(os.PathSeparator) + "Config"
	if !pathExists(configPath) {
		createDir(configPath)
	}
	configFile := configPath + string(os.PathSeparator) + "goSpider.json"
	content := "{\n  \"User-Agent\": \"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36\",\n  \"flagDevToolUrl\": \"http://127.0.0.1:9222/json/version\",\n  \"spiderCount\": 30,\n  \"dataHandleCount\": 1000\n}"

	var err error
	var file *os.File

	if !pathExists(configFile) {
		file, err = os.Create(configFile) //创建文件

	} else {
		file, err = os.OpenFile(configFile, os.O_APPEND, 0666) //打开文件
	}
	if err != nil {
		panic(err)
	}
	_, writeErr := io.WriteString(file, content) //写入文件(字符串)
	if writeErr != nil {
		panic(writeErr)
	}
	err = json.Unmarshal([]byte(content), &config.Data)
	if err != nil {
		log.Panicln("配置文件解析失败：", content, err)
	}
}
