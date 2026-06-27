package main

import (
	"fmt"
	//"golang.org/x/net/ipv4"
	"github.com/subnetscripter/ping_island/probe_ping"
	"log"
)

func main() {
	fmt.Println("Initiating Program!")

	target := "google.com"
	pingProbe, err := probe_ping.NewProbe(target)
	if err != nil {
		log.Fatal(err)
	}

	err = pingProbe.ICMPHandler.Ping(10, 1)
	if err != nil {
		log.Fatal(err)
	}
}
