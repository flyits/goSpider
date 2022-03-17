package goSpider

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

// CallUserFunc
//  @Description  动态调用结构体方法
//  @param object		interface{}		结构体实例
//  @param methodName	string    		方法名
//  @param args			...interface{}	其他参数
//  @return void
func (spider *Spider) CallUserFunc(object interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	instance := reflect.ValueOf(object)
	if instance.Kind() != reflect.Ptr {
		fmt.Println("ERROR: spider.DataList must be passed in the pointer type")
		os.Exit(-1)
	}
	atomic.AddUint64(&spider.ops, 1)
	if len(inputs) > 0 {
		instance.MethodByName(methodName).Call(inputs)
	} else {
		instance.MethodByName(methodName).Call(make([]reflect.Value, 0))
	}
	atomic.AddUint64(&spider.ops, ^uint64(0))
	runtime.Gosched()
}

// Trim
//  @Description  去除字符串前后指定字符，默认换行、空字符
//  @param str		string		源字符串
//  @param charList	[]string    需去除的前后字符
//  @return string
func Trim(str string, charList []string) string {
	if charList == nil {
		charList = []string{"\n", " "}
	}
	for _, char := range charList {
		str = strings.Replace(str, char, "", -1)
	}
	return str
}

// Str2Int
//  @Description  	去除字符串前后指定字符，默认换行、空字符
//  @param strNum	string		源字符串
//  @return int
func Str2Int(strNum string) int {
	index, err := strconv.Atoi(strNum)
	if err != nil {
		fmt.Println(err)
	}
	return index
}

func createDir(filePath string) {
	isExist := pathExists(filePath)
	if !isExist {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		}
	}
}

func pathExists(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func echoErr(err error, isPanic bool) {
	if err == nil {
		return
	}
	if isPanic {
		panic(err)
	} else {
		fmt.Println(err)
	}
}
