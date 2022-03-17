package main

import (
	"fmt"
	goSpider "github.com/Flyits/goSpider/src"
)

func main() {
	index := goSpider.UrlItem{
		SpiderFunc:  goSpider.JsLoadRequest,
		Url:         "https://car.m.yiche.com/aodia3-3999/peizhi/",
		WaitExpr:    "html",
		Expr:        "html",
		HandlerFunc: spiderCar,
	}
	spider := goSpider.Init()
	spider.AddJob(index)
	spider.Run()

}

func spiderCar(response goSpider.Response, spider *goSpider.Spider) {
	fmt.Println(response.HtmlStr)
	//movies := htmlquery.Find(response.HtmlNode, "//div[@class='screening-bd']//ul[@class='ui-slide-content']/li")
	//for k, _ := range movies {
	//	Car := &Car{
	//		Title:     htmlquery.SelectAttr(movies[k], "data-title"),
	//		Pic:       htmlquery.SelectAttr(htmlquery.FindOne(movies[k], "//img"), "src"),
	//		Star:      htmlquery.SelectAttr(movies[k], "data-rate"),
	//		DetailUrl: htmlquery.SelectAttr(htmlquery.FindOne(movies[k], "//a"), "a"),
	//	}
	//	// 此处需注意必须传递指针
	//	spider.DataList <- Car
	//}

}

// DataHandle 异步处理/存储回调方法  也可在 response 里直接同步处理/存储
func (Car *Car) DataHandle(spider *goSpider.Spider) {
	fmt.Println("数据：", Car)

	// todo 处理……
	//
	// todo 存储……
}

type Car struct {
	Title     string
	Pic       string
	Star      string
	DetailUrl string
}

func (Car) GetHandler() string {
	return "DataHandle"
}
