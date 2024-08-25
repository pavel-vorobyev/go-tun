package main

import (
	"go-tun/server"
	"go-tun/server/packet"
	"go-tun/util"
	"log"
)

var rxTc = &packet.TrafficCallback{}
var txTc = &packet.TrafficCallback{}

func main() {
	defer func() {
		log.Println(rxTc.T)
		log.Println(txTc.T)
	}()

	options := server.CreateOptions()
	options.AddRxCallback(rxTc)
	options.AddTxCallback(txTc)

	s, err := server.CreateServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
	util.Serve()

	// simple.Run()
}
