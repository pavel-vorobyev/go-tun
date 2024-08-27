package main

import (
	"go-tun/server"
	"go-tun/server/packet"
	"go-tun/server/storage/address"
	"go-tun/util"
	"log"
)

var rxTc = &packet.TrafficCallback{}
var txTc = &packet.TrafficCallback{}
var cAddrStore = &address.DefaultCAddrStore{}

func main() {
	defer func() {
		log.Println(rxTc.T / 1024 / 1024)
		log.Println(txTc.T / 1024 / 1024)
		cAddrStore.Summary()
	}()

	options := server.CreateOptions()
	options.AddRxCallback(rxTc)
	options.AddTxCallback(txTc)
	options.SetCustomSrcAddressStore(cAddrStore)

	s, err := server.CreateServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()

	util.Serve()
}
