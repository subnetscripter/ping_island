package icmp_handler

import(
    "errors"
    "fmt"
    "golang.org/x/net/icmp"
    "golang.org/x/net/ipv4"
    "github.com/subnetscripter/ping_island/network_handler"
    "time"
)


type Handler struct{
    // NetHandler *network_handler.Handler  //We may not necessarily need this. We can just use the Addr type
    Addr network_handler.IPAddress  //Aliased from net.Addr package
    Conn *icmp.PacketConn //Connect used to send the icmp packets.
    MessageBytes []byte //Message as bytes to send to icmp peer
    Message icmp.Message //Message to send icmp peer.
}

//Return a new ICMP handler
func NewHandler(target string) (*Handler, error){
    
    //Leverages the network_handler package to get temporarily
    //for address resolution. 
    //Was unsure if to keep the struct at a value of Handler
    //Or continue with the temp method. This can be changed.
    netHandler, err := network_handler.NewHandler(target)
    if err != nil {
        return nil, err
    }

    conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        return nil, errors.New("Unable to create a PacketConnection. Trying running as root or sudo.")
    }
    
    icmpMessage := icmp.Message{
        Type: ipv4.ICMPTypeEcho,
        Code: 0,
        Body: &icmp.Echo{
            ID: 1234,
            Seq: 1,
            Data: []byte("ping"),
        },
    }

    icmpMsgBytes, err := icmpMessage.Marshal(nil)
    if err != nil {
        return nil, errors.New("There was an error Marshalling the icmp message to bytes.")
    }

    return &Handler{
        Addr: netHandler.ReturnIPAddr(),
        Conn: conn,
        MessageBytes: icmpMsgBytes,
        Message: icmpMessage,
    }, nil
}


func (h *Handler) sendPing() error {
    _, err := h.Conn.WriteTo(h.MessageBytes, h.Addr)
    if err != nil {
        return errors.New("Unable to send pings")
    }

    return nil
}

//Internal function used to receive pings
func (h *Handler) recvPing() (*icmp.Message, *network_handler.IPAddress, error){
    reply := make([]byte, 1500)
    n, peer, err := h.Conn.ReadFrom(reply)
    if err != nil {
        return nil, nil, errors.New("Unable to receive icmp packets")
    }

    parsed, err := icmp.ParseMessage(1, reply[:n])
    if err != nil {
        return nil, nil, errors.New("Unable to parse ICMP message.")
    }

    return parsed, &peer, nil
}

/*
Public function used to aggregate send and recv ping functions
Right now the functions only returns an error, but there is room for improvement.
Perhaps we can return a struct in the future.
We also need to add checks to verify that the pings received come from the same peer sent to.
Conside that for pings you can return a slice of structs
*/

func (h *Handler) Ping(repeat int, interval int) (error) {
    for i := 0; i < repeat; i++ {
        
        err := h.sendPing()
        if err != nil{
            return err
        }

        rcvMessage, peer, err := h.recvPing()
        if err != nil{
            return err
        }

        //Switch case used to treat received ICMP messages differently.
        //There is room for improvement here. We can handle more message types.
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
