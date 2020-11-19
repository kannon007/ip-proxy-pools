package support

import (
	"github.com/asdine/storm/v3"
	"log"
	"time"
)

var Db *storm.DB

func init() {
	db, err := storm.Open("my.db")
	if err != nil {
		log.Panic(err)
	}

	Db = db
}
func Save(ip string, port string, location string, time int64) {
	log.Printf("save ip :%v  port:%v  location:%v  avgMs:%v \n", ip, port, location, time)
	err := Db.Save(&ProxyIp{
		Ip:       ip,
		Port:     port,
		Location: location,
		Ms:       time,
	})
	if err != nil {
		log.Println(err)
	}
}

func GetIPs() []*ProxyIp {
	var ips []*ProxyIp
	query := Db.Select()
	query.OrderBy("Ms")
	query.Find(&ips)
	return ips
}

func GetFastIPs() *ProxyIp {

	ip := ProxyIp{}
	query := Db.Select()
	query.OrderBy("LastUseAt Ms")
	query.First(&ip)

	Db.UpdateField(&ip, "LastUseAt", time.Now())
	return &ip
}

//Ip       string `gorm:"unique_index:idx_ip_port"`
//Port     string `gorm:"unique_index:idx_ip_port"`
type ProxyIp struct {
	//gorm.Model
	ID        int    `storm:"id,increment"`
	Ip        string `storm:"unique"`
	Port      string
	LastUseAt time.Time `storm:"index"`
	Location  string
	Ms        int64 `storm:"index"`
}
