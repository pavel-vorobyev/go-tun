package main

import (
	"go-tun/server"
	"log"
)

func main() {
	options := server.CreateOptions()
	s, err := server.CreateServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	s.Start()
}
