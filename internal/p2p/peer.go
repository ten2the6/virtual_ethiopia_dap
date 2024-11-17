package p2p

import (
    "net"
)

// Peer represents a node in the P2P network
type Peer struct {
    Address string
    Port    string
    Conn    net.Conn
}

// NewPeer creates a new peer instance
func NewPeer(address, port string) *Peer {
    return &Peer{
        Address: address,
        Port:    port,
    }
}

// Connect establishes a connection with the peer
func (p *Peer) Connect() error {
    conn, err := net.Dial("tcp", p.Address+":"+p.Port)
    if err != nil {
        return err
    }
    p.Conn = conn
    return nil
}

// Disconnect closes the connection with the peer
func (p *Peer) Disconnect() error {
    if p.Conn != nil {
        return p.Conn.Close()
    }
    return nil
}