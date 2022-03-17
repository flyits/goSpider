package goSpider

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type Spider struct {
	Urls       chan UrlItem
	DataList   chan DataItem
	Config     Config
	downloader *Downloader
	wg         sync.WaitGroup

	ops uint64
}

type SpiderHandlerFunc interface {
	getSpiderFunc() string
}

// Init 初始化
func Init() *Spider {
	spider := &Spider{}
	spider.Urls = make(chan UrlItem, 100000)
	spider.DataList = make(chan DataItem, 100000)
	spider.ops = 0
	return spider
}

// Run 执行此方法后，主线程将被阻塞至爬虫队列运行完毕
func (spider *Spider) Run() {
	spider.Config.init()
	spider.response()
	spider.startUp()
	spider.logPrint()
	spider.wg.Wait()
	return
}

func (spider *Spider) response() {
	for i := 0; i < spider.Config.get("dataHandleCount", 1000).(int); i++ {
		go func() {
			spider.wg.Add(1)
			defer spider.wg.Done()
			for item := range spider.DataList {
				spider.CallUserFunc(item, item.GetHandler(), spider)
			}
		}()
	}
}

func (spider *Spider) startUp() {
	for i := 0; i < spider.Config.get("spiderCount", 1000).(int); i++ {
		go func() {
			spider.wg.Add(1)
			defer spider.wg.Done()
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

}

func (spider *Spider) getOps() uint64 {
	return spider.ops
}

func (spider *Spider) AddJob(url UrlItem) *Spider {
	spider.Urls <- url
	return spider
}

func (spider *Spider) AddJobs(urls []UrlItem) *Spider {
	for _, url := range urls {
		spider.Urls <- url
	}
	return spider
}
