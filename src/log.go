package goSpider

import (
	"fmt"
	"time"
)

func (spider *Spider) logPrint() {
	spider.startLog()
	ticker := time.NewTicker(time.Millisecond * 10000)
	go func() {
		for t := range ticker.C {
			dataChanLen := len(spider.DataList)
			urlChanLen := len(spider.Urls)

			fmt.Println("Tick at", t)
			fmt.Println("当前待处理数据：", dataChanLen, "当前待抓取链接：", urlChanLen, "当前爬虫运行的任务（协程）：", spider.ops)
			if spider.ops == 0 {
				close(spider.DataList)
				close(spider.Urls)
				ticker.Stop()
				fmt.Println("--------------------------------------------------------------------------")
				fmt.Println("爬取完成……")
				fmt.Println("完成时间：", time.Now().Format("2006-01-02 15:04:05"))
			}
		}
	}()
}

func (spider *Spider) startLog() {
	fmt.Println("当前时间：", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("爬虫并发数量：%d\n", spider.config.get("spiderCount", 0).(int))
	fmt.Printf("数据并行处理数：%d\n", spider.config.get("dataHandleCount", 0).(int))
	fmt.Println("开启抓取数据……")
	fmt.Println("--------------------------------------------------------------------------")
}
