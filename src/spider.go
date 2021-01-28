package goSpider

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
	downloader  *Downloader
	dataHandler *DataHandler
	Wg          sync.WaitGroup

	ops uint64
}

type SpiderHandlerFunc interface {
	getSpiderFunc() string
}

// 初始化
func (spider *Spider) Init(urls []UrlItem) {
	spider.Urls = make(chan UrlItem, 100000)
	spider.DataList = make(chan DataItem, 100000)
	spider.config.init()
	spider.ops = 0
	spider.saveData()
	spider.run(urls)
	spider.logPrint()

	spider.Wg.Wait()
}

func (spider *Spider) saveData() {
	for i := 0; i < spider.config.get("dataHandleCount", 1000).(int); i++ {
		go func() {
			spider.Wg.Add(1)
			defer spider.Wg.Done()
			for item := range spider.DataList {
				spider.CallUserFunc(item, item.GetHandler())
			}
		}()
	}
}

func (spider *Spider) run(urls []UrlItem) {
	for i := 0; i < spider.config.get("spiderCount", 1000).(int); i++ {
		go func() {
			spider.Wg.Add(1)
			defer spider.Wg.Done()
			for spiderItem := range spider.Urls {
				spider.downloader = &Downloader{
					spider:  spider,
					urlItem: &spiderItem,
				}
				atomic.AddUint64(&spider.downloader.ops, 1)
				spider.CallUserFunc(spider.downloader, spiderItem.getSpiderFunc())
				atomic.AddUint64(&spider.downloader.ops, ^uint64(0))
				runtime.Gosched()
			}
		}()
		time.Sleep(500000 * time.Nanosecond)
	}
	for _, url := range urls {
		spider.Urls <- url
	}
}

func (spider *Spider) getOps() uint64 {
	return spider.ops
}
