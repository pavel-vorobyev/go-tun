package main

import (
	"go-tun/server"
	"go-tun/server/packet"
	"go-tun/util"
	"log"
)

func main() {
	options := server.CreateOptions()
	options.AddRxCallback(&packet.TrafficCallback{})

	s, err := server.CreateServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
	util.Serve()

	// simple.Run()
}
