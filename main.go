package main

import "ip-proxy-pools/support"

func main() {

	go support.Patch()
	go support.Server("0.0.0.0:8080")
	support.ApiServerStart()
}
