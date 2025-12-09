package main

import (
	"log"

	"github.com/M-oses340/decentralised_file_system/bin/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":3000") // FIXED
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
