package main

import "github.com/z-wyd/ip-proxy-pools/support"

func main() {
	go support.Patch()
	support.Server("0.0.0.0:8080")
}
