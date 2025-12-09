package main

import (
	"log"

	"github.com/M-oses340/decentralised_file_system/bin/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(p2p.TCPTransportOpts{
		ListenAddr: ":3000",
	})

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {} // keep running forever
}
