package goSpider

import "golang.org/x/net/html"

type Response struct {
	HtmlStr  string                      // 页面dom节点(html源码)
	HtmlNode *html.Node                  // 页面dom节点(xml树)
	Url      string                      // 爬取的页面链接
	Attach   map[interface{}]interface{} // 额外参数
}
