package main

import (
	"go-tun/server"
	"go-tun/util"
	"log"
)

var rxTc = &server.TrafficPacketCallback{}
var txTc = &server.TrafficPacketCallback{}

func main() {
	startServer()
}

func startServer() {
	defer func() {
		log.Println(rxTc.T / 1024 / 1024)
		log.Println(txTc.T / 1024 / 1024)
	}()

	options := server.NewOptions()
	options.AddRxCallback(rxTc)
	options.AddTxCallback(txTc)

	s, err := server.NewServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()

	util.Serve()
}
