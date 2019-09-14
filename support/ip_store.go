package support

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"time"
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
	Db.Order("last_use_at asc , ms asc").First(&ip)
	Db.Model(&ip).Update("last_use_at", time.Now())
	return &ip
}

type ProxyIp struct {
	gorm.Model
	Ip       string `gorm:"unique_index:idx_ip_port"`
	Port     string `gorm:"unique_index:idx_ip_port"`
	LastUseAt time.Time
	Location string
	Ms       int64
}
