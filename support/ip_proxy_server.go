package support

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

func Server(listen string) {
	http.Handle("/", &Pxy{})
	http.ListenAndServe(listen, nil)
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
