package goSpider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	Request       = "Request"
	JsLoadRequest = "JsLoadRequest"
)

type Downloader struct {
	spider     *Spider
	urlItem    *UrlItem
	spiderFunc string

	ops uint64
}

type UrlItem struct {
	Url         string                                  // 爬取地址
	SpiderFunc  string                                  // 请求方法
	HandlerFunc func(response Response, spider *Spider) // 处理方法
	WaitExpr    string                                  // 需要 js 渲染的站点，需提前传递dom节点
	Expr        string                                  // 需要 js 渲染的站点，需提前传递dom节点
	Attach      interface{}                             // 额外参数
}

func (urlItem *UrlItem) getSpiderFunc() string {
	if urlItem.SpiderFunc == "" {
		return Request
	}
	return urlItem.SpiderFunc
}

var flagDevToolWsUrl *string

func (download *Downloader) Request() {
	var err error
	// 生成client客户端
	client := &http.Client{}
	// 生成Request对象
	req, err := http.NewRequest("GET", download.urlItem.Url, nil)
	if err != nil {
		fmt.Println(err)
	}
	userAgent := download.spider.Config.get("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36")
	// 添加Header
	req.Header.Add("User-Agent", userAgent.(string))
	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	// 请求失败重新放入爬虫通道 重试
	if resp == nil {
		//urlChan <- SpiderChanItem{Url: url, SpiderFunc: "Request", HandlerFunc: handle, Attach: attach}
		return
	}
	// 设定关闭响应体
	defer resp.Body.Close()
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	download.response(string(body))
}

func (download *Downloader) JsLoadRequest() {
	var htmlContent string

	if download.urlItem.Url == "" {
		fmt.Println("无效的地址：", download.urlItem.Url)
		return
	}

	if download.urlItem.Expr == "" {
		download.urlItem.Expr = "html"
	}
	if download.urlItem.WaitExpr == "" {
		download.urlItem.WaitExpr = "html"
	}

	if *getFlagDevToolWsUrl(download.spider) == "" {
		log.Fatal("must specify -devtools-ws-url")
	}

	// create allocator context for use with creating a browser context later
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), *getFlagDevToolWsUrl(download.spider))
	defer cancel()

	// create context
	ctxt, cancel2 := chromedp.NewContext(allocatorContext)
	defer cancel2()

	timeoutCtx, cancel3 := context.WithTimeout(ctxt, 20*time.Second)
	defer cancel3()

	// run task list
	if err := chromedp.Run(timeoutCtx,
		chromedp.Emulate(device.IPhone7Pluslandscape),
		chromedp.Navigate(download.urlItem.Url),
		chromedp.WaitVisible(download.urlItem.WaitExpr, chromedp.ByQuery),
		chromedp.OuterHTML(download.urlItem.Expr, &htmlContent, chromedp.ByQuery),
	); err != nil {
		download.spider.Urls <- *download.urlItem
		fmt.Println("请求失败："+download.urlItem.Url, err)
		return
	}
	download.response(htmlContent)
}

func (download *Downloader) response(htmlStr string) {
	doc, err := htmlquery.Parse(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println(err)
	}

	response := Response{
		HtmlStr:  htmlStr,
		HtmlNode: doc,
		Url:      download.urlItem.Url,
		Attach:   download.urlItem.Attach,
	}
	download.urlItem.HandlerFunc(response, download.spider)
}

// 自动获取远程无头浏览器地址
func getFlagDevToolWsUrl(spider *Spider) *string {
	if flagDevToolWsUrl == nil {
		url := spider.Config.get("flagDevToolUrl", "http://127.0.0.1:9222/json/version")
		req, err := http.NewRequest("GET", url.(string), nil)
		echoErr(err, true)
		res, err := http.DefaultClient.Do(req)
		echoErr(err, true)
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		echoErr(err, true)
		var responseData map[string]string
		if err := json.Unmarshal(body, &responseData); err != nil {
			log.Printf("chromedp/headless-shell 服务异常，请检查！")
			panic(err)
		}
		if _, exists := responseData["webSocketDebuggerUrl"]; !exists {
			log.Printf("chromedp/headless-shell 服务异常")
			panic("webSocketDebuggerUrl 获取失败")
		}
		webSocketDebuggerUrl := responseData["webSocketDebuggerUrl"]
		flagDevToolWsUrl = &webSocketDebuggerUrl
	}
	return flagDevToolWsUrl
}
