package main

import "ip-proxy-pools/support"

const port int = 8080

func main() {

	go support.Patch()

	//代理服务器
	go support.Server(port)

	support.SetSysProxy("127.0.0.1", port)

	// api 服务器
	support.ApiServerStart()
}
