package main

import (
	"go-tun/server"
	"go-tun/server/packet"
	"go-tun/util"
	"log"
)

func main() {
	defer func() {
		log.Println()
	}()

	rxTc := &packet.TrafficCallback{}
	txTc := &packet.TrafficCallback{}

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
