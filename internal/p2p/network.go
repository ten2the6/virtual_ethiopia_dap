package p2p

import (
    "encoding/json"
    "fmt"
    "net"
    "sync"
)

// Network represents the P2P network functionality
type Network struct {
    peers     map[string]net.Conn
    mu        sync.RWMutex
    listener  net.Listener
    isRunning bool
}

// NewNetwork creates a new P2P network instance
func NewNetwork() *Network {
    return &Network{
        peers:     make(map[string]net.Conn),
        isRunning: false,
    }
}

// Start initializes the P2P network
func (n *Network) Start(port string) error {
    listener, err := net.Listen("tcp", ":"+port)
    if err != nil {
        return fmt.Errorf("failed to start P2P network: %v", err)
    }

    n.listener = listener
    n.isRunning = true

    go n.listen()
    return nil
}

// Stop shuts down the P2P network
func (n *Network) Stop() error {
    n.mu.Lock()
    defer n.mu.Unlock()

    if !n.isRunning {
        return nil
    }

    if err := n.listener.Close(); err != nil {
        return fmt.Errorf("failed to close listener: %v", err)
    }

    for _, conn := range n.peers {
        conn.Close()
    }

    n.peers = make(map[string]net.Conn)
    n.isRunning = false
    return nil
}

// Connect establishes a connection with a peer
func (n *Network) Connect(address string) error {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return fmt.Errorf("failed to connect to peer: %v", err)
    }

    n.mu.Lock()
    n.peers[address] = conn
    n.mu.Unlock()

    go n.handlePeer(conn)
    return nil
}

// listen handles incoming connections
func (n *Network) listen() {
    for n.isRunning {
        conn, err := n.listener.Accept()
        if err != nil {
            if n.isRunning {
                fmt.Printf("Failed to accept connection: %v\n", err)
            }
            continue
        }

        go n.handlePeer(conn)
    }
}

// handlePeer processes messages from a peer
func (n *Network) handlePeer(conn net.Conn) {
    defer conn.Close()

    decoder := json.NewDecoder(conn)
    for {
        var message struct {
            Type string          `json:"type"`
            Data json.RawMessage `json:"data"`
        }

        if err := decoder.Decode(&message); err != nil {
            return
        }

        // Handle different message types here
        switch message.Type {
        case "block":
            n.handleBlockMessage(message.Data)
        case "transaction":
            n.handleTransactionMessage(message.Data)
        }
    }
}

func (n *Network) handleBlockMessage(data json.RawMessage) {
    // Implement block handling
    fmt.Println("Received block message")
}

func (n *Network) handleTransactionMessage(data json.RawMessage) {
    // Implement transaction handling
    fmt.Println("Received transaction message")
}

// Broadcast sends a message to all peers
func (n *Network) Broadcast(messageType string, data interface{}) error {
    message := struct {
        Type string      `json:"type"`
        Data interface{} `json:"data"`
    }{
        Type: messageType,
        Data: data,
    }

    messageJSON, err := json.Marshal(message)
    if err != nil {
        return err
    }

    n.mu.RLock()
    defer n.mu.RUnlock()

    for _, conn := range n.peers {
        if _, err := conn.Write(append(messageJSON, '\n')); err != nil {
            fmt.Printf("Failed to send to peer: %v\n", err)
        }
    }

    return nil
}