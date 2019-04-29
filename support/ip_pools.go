package support

import (
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/resty.v1"
	"time"
)

func Add(ip string, port string, location string) {
	//test ip
	resty.RemoveProxy()
	server := "http://" + ip + ":" + port

	client := resty.SetProxy(server)
	client.SetHTTPMode()
	client.SetRetryMaxWaitTime(time.Duration(time.Second * 1))
	client.SetTimeout(time.Duration(time.Second * 2))
	resp, _ := client.R().Get("http://www.baidu.com")

	if resp.StatusCode() == 200 && resp.Time().Nanoseconds()/1e6 < 2000 {
		fmt.Printf("\nproxy ip: %v", server)
		fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
		fmt.Printf("\nResponse Time: %v", resp.Time().Nanoseconds()/1e6)

		Save(ip, port, location, resp.Time().Nanoseconds()/1e6)
	}

}

func Patch() {
	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		d := e.DOM
		ip := d.Find("td[data-title='IP']").Text()
		port := d.Find("td[data-title='PORT']").Text()
		location := d.Find("td[data-title='位置']").Text()
		Add(ip, port, location)

	})

	c.OnHTML("#listnav", func(e *colly.HTMLElement) {
		d := e.DOM
		page := d.Find("li a.active").Parent().Next().Find("li a").AttrOr("href", "")
		fmt.Println("\nVisiting", page)
		e.Request.Visit("https://www.kuaidaili.com" + page)
	})

	c.Visit("https://www.kuaidaili.com/free/inha/")

}
