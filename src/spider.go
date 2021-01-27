package goSpiderFarmwork

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Spider struct {
	Urls        chan UrlItem
	DataList    chan DataItem
	config      Config
	downloader  Downloader
	dataHandler DataHandler

	Ops uint64
}

type HandlerFunc interface {
	getHandlerFunc() string
}

type SpiderHandlerFunc interface {
	getSpiderFunc() string
}

var wg sync.WaitGroup

// 初始化
func (spider *Spider) init() {
	spider.config.init()
	spider.Ops = 0
	spider.saveData()
	spider.run()
	spider.logPrint()
	wg.Wait()
}

func (spider *Spider) saveData() {
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for item := range spider.DataList {
				atomic.AddUint64(&spider.dataHandler.Ops, 1)
				spider.CallUserFunc(spider.dataHandler, item.GetHandler())
				atomic.AddUint64(&spider.dataHandler.Ops, ^uint64(0))
				runtime.Gosched()
			}
		}()
	}
}

func (spider *Spider) run() {
	for i := 0; i < spider.config.get("spiderCount", 1000).(int); i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for spiderItem := range spider.Urls {
				spider.downloader = Downloader{
					spider:  spider,
					urlItem: &spiderItem,
				}
				atomic.AddUint64(&spider.downloader.Ops, 1)
				spider.CallUserFunc(spider.downloader, spiderItem.SpiderFunc)
				atomic.AddUint64(&spider.downloader.Ops, ^uint64(0))
				runtime.Gosched()
			}
		}()
		time.Sleep(500000 * time.Nanosecond)
	}

}
