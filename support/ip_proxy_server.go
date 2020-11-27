package support

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"log"
	"os/exec"
	"strings"
	"syscall"

	"net"
	"net/http"
	"net/url"

	"strconv"
	"time"
)

func Server(port int) {
	http.Handle("/", &Pxy{})
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
}

func SetSysProxy(ip string, port int) {

	k, err := registry.OpenKey(registry.CURRENT_USER, "Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	err = k.SetDWordValue("ProxyEnable", uint32(1))

	if err != nil {
		log.Fatal(err)
	}

	err = k.SetExpandStringValue("ProxyServer", ip+":"+strconv.Itoa(port))

	if err != nil {
		log.Fatal(err)
	}

	err = k.SetExpandStringValue("ProxyOverride", "*.ffcs.cn;192.168.*;125.0.0.1;<local>")

	if err != nil {
		log.Fatal(err)
	}

}

func out(byte []byte) {

	decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)

	log.Println(string(decodeBytes))
}

func ExecCmd(cmd string) {
	args := strings.Split(cmd, " ")
	var c *exec.Cmd

	if len(args) > 1 {
		cmdArgs := args[1:]
		c = exec.Command(args[0], cmdArgs...)
	} else {
		c = exec.Command(args[0])
	}
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	var out bytes.Buffer

	c.Stdout = &out
	c.Run()
	//out,_ := c.Out()
	log.Println(out.String())
}

type Pxy struct {
}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ip := GetFastIPs()
	proxyServer := "http://" + ip.Ip + ":" + ip.Port
	fmt.Printf("Received request %v  %s %s  proxy: %s \n", ip, req.Method, req.Host, proxyServer)
	// step 1
	outReq := new(http.Request)
	*outReq = *req // this only does shallow copies of maps

	proxy, err := url.Parse(proxyServer)

	zTransport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	//if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
	//	//if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
	//	//	clientIP = strings.Join(prior, ", ") + ", " + clientIP
	//	//}
	//	outReq.Header.Set("X-Forwarded-For", clientIP)
	//}

	// step 2
	res, err := zTransport.RoundTrip(outReq)

	fmt.Println(err)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)

		return
	}

	// step 3
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
	res.Body.Close()
}
