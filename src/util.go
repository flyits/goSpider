package goSpider

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
)

//  Description  动态调用结构体方法
//  @param object		interface{}		结构体实例
//  @param methodName	string    		方法名
//  @param args			...interface{}	其他参数
//  @return void
func (spider *Spider) CallUserFunc(object interface{}, methodName string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	atomic.AddUint64(&spider.Ops, 1)
	reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
	atomic.AddUint64(&spider.Ops, ^uint64(0))
	runtime.Gosched()
}

//  Description  去除字符串前后指定字符，默认换行、空字符
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

//  Description  	去除字符串前后指定字符，默认换行、空字符
//  @param strNum	string		源字符串
//  @return int
func Str2Int(strNum string) int {
	index, err := strconv.Atoi(strNum)
	if err != nil {
		fmt.Println(err)
	}
	return index
}
