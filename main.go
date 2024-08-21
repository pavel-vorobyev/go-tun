package main

import (
	"go-tun/server"
	"go-tun/util"
	"log"
)

func main() {
	options := server.CreateOptions()
	s, err := server.CreateServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
	util.Serve()
}
