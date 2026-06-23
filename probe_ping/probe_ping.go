package probe_ping

/*
Package used to provisioning and utilize probes for ICMP.
Consider outsource network services like DNS resolution and connection requests to
another package that can be called from multiple probes.
*/

import (
    "net"
    "time"
    "fmt"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
    "errors"
)

//Each probe will be a struct that will perform the operations
type Probe struct {
    Addr *net.IPAddr //Resolved IP of the hostname or an IP addresss
    Conn *icmp.PacketConn //Connect used to send the ICMP pings
    Alive bool //Not used right now. Will figure it out later.
    MessageBytes []byte //Message as bytes to send to icmp peer
    Message icmp.Message //Message to send icmp peer. Considering making as a separate struct.
}

//Returns a pointer to a new probe with some values set.
func NewProbe(target *net.IPAddr) (*Probe, error) {
    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil{
        return nil, errors.New("Unable to create a PacketConnection. Try using sudo.")
    }
    probe :=  Probe{
        Addr: target,
        Conn: conn,
    }

    err = probe.PrepMessage()
    if err != nil {
        return nil, err
    }
    return &probe, nil
}

//Function internal to the package used to create the
//ICMP message struct
func (p *Probe) createMessage() (error) {
    msg := icmp.Message{
        Type: ipv4.ICMPTypeEcho,
        Code: 0,
        Body: &icmp.Echo{
            ID: 1234,
            Seq: 1,
            Data: []byte("ping"),
        },
    }
    p.Message = msg
    return nil
}

//Function internal to the package used to marshal
//the ICMP message struct as bytes.
func (p *Probe) marshalMessage() (error) {
    msgBytes, err := p.Message.Marshal(nil)
    if err != nil{
        return errors.New("There was an error Marshalling the message to bytes")
    }
    
    p.MessageBytes = msgBytes
    return nil
}

//Public function used to organize prepartion of the
//ICMP message struct and marshalled bytes/
//Considering making this an internal function.
func (p *Probe) PrepMessage() (error){
    err := p.createMessage()
    if err != nil{
        return errors.New("Unable to create icmp Message!")
    }

    err = p.marshalMessage()
    if err != nil{
        return errors.New("Unable to Marshall ICMP message to bytes.")
    }
    return nil

}

//Function internal to the package used to send pings
func (p *Probe) sendPing() (error){
    _, err := p.Conn.WriteTo(p.MessageBytes, p.Addr)
    if err != nil{
        return errors.New("Unable to send pings")
    }
    
    return nil
}

//Internal function used to receive pings.
func (p *Probe) recvPing() (*icmp.Message, *net.Addr, error){
    reply := make([]byte, 1500)
    n, peer, err := p.Conn.ReadFrom(reply)
    if err != nil{
        return nil, nil ,errors.New("Unable to receive ping packets")
    }

    parsed, err := icmp.ParseMessage(1, reply[:n])
    if err != nil{
        return nil, nil, errors.New("Unable to parse ICMP message.")
    }
    return parsed, &peer, nil
}

//Public function used to aggregate send and recv
//ping functions.
//Right now the function only returns an error, but there is room for improvement.
//Perhaps we can return a struct and use the to validate aliveness.
func (p *Probe) Ping(repeat int, interval int64) (error) {

    for i := 0; i < repeat; i++ {

        err := p.sendPing()
        if err != nil{
            return err
        }

        rcvMessage, peer, err := p.recvPing()
        if err != nil{
            return err
        }
        
        //Switch case used to treat received ICMP messages differently.
        //There is room for improvement here. We can handle more message types
        switch rcvMessage.Type {
            case ipv4.ICMPTypeEchoReply:
                fmt.Printf("Got a reply from %v\n", *peer)
            default:
                fmt.Printf("Got an unexpected reply %v.\n", rcvMessage.Type)
        }

        time.Sleep(time.Duration(interval) * time.Second)

    }

    return nil

}
