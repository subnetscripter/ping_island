package probe_ping

/*
Package used to provisioning and utilize probes for ICMP.
Consider outsource network services like DNS resolution and connection requests to
another package that can be called from multiple probes.
*/

import (
	"github.com/subnetscripter/ping_island/icmp_handler"
)

// Each probe will be a struct that will perform the operations
type Probe struct {
	Alive        bool                     //Not used right now. Will figure it out later.
    ICMPHandler *icmp_handler.Handler
}

// Returns a pointer to a new probe with some values set.
func NewProbe(target string) (*Probe, error) {

    icmpHandler, err := icmp_handler.NewHandler(target)
    if err != nil{
        return nil, err
    }

    return &Probe{
        ICMPHandler: icmpHandler,
    }, nil
}





