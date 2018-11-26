package main

import "github.com/ktpswjz/httpserver/http/proxy/host"

func main() {
	target := &Target{}
	server := host.NewHost("9090", target.GetUrl)
	server.Run()
}
