package funcs

/**
 * 业务相关的通用函数库
 * 2021/05/08
 */

import (
	"net/http"
	"strings"
)

// IsMobile 判断是否手机访问
// 如果是百度爬虫 当作pc处理
func IsMobile(r *http.Request) bool {
	if !IsBaiduSpider(r) {
		userAgent := r.UserAgent()
		clients := []string{"Android", "Mobile", "baiduboxapp"}
		for _, v := range clients {
			if strings.Contains(userAgent, v) {
				return true
			}
		}
	}

	return false
}

// IsBaiduSpider 是否百度爬虫
func IsBaiduSpider(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.UserAgent()), "baiduspider")
}

// FormatLimit 返回实际 offset
func FormatLimit(pageNum, pageAmount int) int {
	var offset int
	if pageNum == 0 {
		pageNum = 1
	}
	offset = pageAmount*pageNum - pageAmount
	return offset
}
