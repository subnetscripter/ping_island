package main

import (
    "fmt"
    "net"
    //"golang.org/x/net/ipv4"
    "github.com/subnetscripter/ping_island/probe_ping"
    "log"
)

func main(){
    fmt.Println("Initiating Program!")

    target:= "8.8.8.8"
    addr, err := net.ResolveIPAddr("ip4", target) 
    if err != nil{
        log.Fatal(err)
    }

    pingProbe, err := probe_ping.NewProbe(addr)
    if err != nil{
        log.Fatal(err)
    }
    

    err = pingProbe.Ping(10, 1)
    if err != nil{
        log.Fatal(err)
    }
}
