package main

import (
	"fmt"
	"github.com/Flyits/goSpider/src"
	"github.com/antchfx/htmlquery"
)

/*
 * 创建时间：2021-2-18 14:55
 */

func main() {
	index := goSpider.UrlItem{
		SpiderFunc:  goSpider.Request,
		Url:         "https://movie.douban.com/",
		HandlerFunc: test,
	}
	spider := goSpider.Init()
	spider.AddJob(index)
	spider.Run()

}

func test(response goSpider.Response, spider *goSpider.Spider) {

	movies := htmlquery.Find(response.HtmlNode, "//div[@class='screening-bd']//ul[@class='ui-slide-content']/li")
	for k, _ := range movies {
		movie := &movie{
			Title:     htmlquery.SelectAttr(movies[k], "data-title"),
			Pic:       htmlquery.SelectAttr(htmlquery.FindOne(movies[k], "//img"), "src"),
			Star:      htmlquery.SelectAttr(movies[k], "data-rate"),
			DetailUrl: htmlquery.SelectAttr(htmlquery.FindOne(movies[k], "//a"), "a"),
		}
		// 此处需注意必须传递指针
		spider.DataList <- movie
	}

}

// DataHandle 异步处理/存储回调方法  也可在 response 里直接同步处理/存储
func (movie *movie) DataHandle(spider *goSpider.Spider) {
	fmt.Println("数据：", movie)

	// todo 处理……
	//
	// todo 存储……
}

type movie struct {
	Title     string
	Pic       string
	Star      string
	DetailUrl string
}

func (movie) GetHandler() string {
	return "DataHandle"
}
