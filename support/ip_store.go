package support

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

var Db *gorm.DB

func init() {
	db, err := gorm.Open("sqlite3", "/app/sqlite/ipProxyPools.db")
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate(&ProxyIp{})
	Db = db
}
func Save(ip string, port string, location string, time int64) {
	Db.Create(&ProxyIp{
		Ip:       ip,
		Port:     port,
		Location: location,
		Ms:       time,
	})
}

func GetIPs() []*ProxyIp {
	var ips []*ProxyIp
	Db.Order("ms asc").Find(&ips)
	return ips
}

func GetFastIPs() *ProxyIp {
	ip := ProxyIp{}
	Db.Order("ms asc").First(&ip)
	return &ip
}

type ProxyIp struct {
	gorm.Model
	Ip       string `gorm:"unique_index:idx_ip_port"`
	Port     string `gorm:"unique_index:idx_ip_port"`
	Location string
	Ms       int64
}
