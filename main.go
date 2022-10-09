package main

import (
	"picture-oss-proxy/conf"
	"picture-oss-proxy/routes"
)

func main() {
	// Ek1+Ep1==Ek2+Ep2
	conf.Init()
	r := routes.NewRouter()
	_ = r.Run(conf.HttpPort)
}
