package support

import (
	"fmt"
	"github.com/gocolly/colly"
	"gopkg.in/resty.v1"
	"log"
	"time"
)

const TIME_OUT int64 = 2000

const TEST_TIMES int = 3

func Add(ip string, port string, location string) {
	//test ip

	ret, avgTime, code := testIP(ip, port)
	log.Printf("test ip :%v port:%v  location:%v ret:%v avgTime:%v  code:%v  \n", ip, port, location, ret, avgTime, code)
	if ret {
		Save(ip, port, location, avgTime)
	}
}

func testIP(ip string, port string) (bool, int64, int) {
	var sum int64 = 0
	for i := 0; i <= TEST_TIMES; i++ {
		resty.RemoveProxy()
		server := "http://" + ip + ":" + port

		client := resty.SetProxy(server)
		client.SetHTTPMode()
		client.SetRetryMaxWaitTime(time.Duration(time.Second * 1))
		client.SetTimeout(time.Duration(time.Second * 2))
		resp, _ := client.R().Get("http://www.baidu.com")

		sum += resp.Time().Nanoseconds() / 1e6
		if resp.StatusCode() != 200 || resp.Time().Nanoseconds()/1e6 >= TIME_OUT {

			return false, sum / int64(i+1), resp.StatusCode()
		}

	}
	return true, sum / int64(TEST_TIMES), 200
}

func Patch() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("table tbody tr", func(e *colly.HTMLElement) {
		d := e.DOM
		ip := d.Find("td[data-title='IP']").Text()
		port := d.Find("td[data-title='PORT']").Text()
		location := d.Find("td[data-title='位置']").Text()
		go Add(ip, port, location)
	})

	c.OnHTML("#listnav", func(e *colly.HTMLElement) {
		d := e.DOM
		page := d.Find("li a.active").Parent().Next().Find("li a").AttrOr("href", "")
		fmt.Println("\nVisiting", page)
		time.Sleep(time.Second * 2)
		e.Request.Visit("https://www.kuaidaili.com" + page)
	})

	c.Visit("https://www.kuaidaili.com/free/inha/")

}
