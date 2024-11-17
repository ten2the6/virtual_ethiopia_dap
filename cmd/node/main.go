package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "virtual_ethiopia_dap/internal/api"
    "virtual_ethiopia_dap/internal/blockchain"
    "virtual_ethiopia_dap/internal/p2p"
)

type Node struct {
    chain     *blockchain.Chain
    network   *p2p.Network
    api       *api.Server
    nodeID    string
    p2pPort   string
    apiPort   string
    isRunning bool
}

func NewNode() *Node {
    chain := blockchain.NewChain()
    return &Node{
        nodeID:  os.Getenv("NODE_ID"),
        p2pPort: os.Getenv("P2P_PORT"),
        apiPort: os.Getenv("API_PORT"),
        chain:   chain,
        network: p2p.NewNetwork(),
        api:     api.NewServer(chain),
    }
}

// Start initializes and starts all node services
func (n *Node) Start() error {
    // Validate configuration
    if n.nodeID == "" || n.p2pPort == "" || n.apiPort == "" {
        return fmt.Errorf("missing required environment variables: NODE_ID, P2P_PORT, or API_PORT")
    }

    // Start P2P network
    if err := n.network.Start(n.p2pPort); err != nil {
        return fmt.Errorf("failed to start P2P network: %v", err)
    }

    // Start API server
    go func() {
        if err := n.api.Start(n.apiPort); err != nil {
            log.Printf("API server error: %v", err)
        }
    }()

    n.isRunning = true
    log.Printf("Node %s started successfully. P2P Port: %s, API Port: %s\n",
        n.nodeID, n.p2pPort, n.apiPort)
    return nil
}

// Stop gracefully shuts down the node
func (n *Node) Stop() error {
    if !n.isRunning {
        return nil
    }

    // Stop P2P network
    if err := n.network.Stop(); err != nil {
        log.Printf("Error stopping P2P network: %v", err)
    }

    // Stop API server (TODO)
    n.isRunning = false
    log.Println("Node stopped successfully")
    return nil
}

func main() {
    // Create and start the node
    node := NewNode()
    if err := node.Start(); err != nil {
        log.Fatal(err)
    }

    // Set up graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down node...")

    // Stop the node gracefully
    if err := node.Stop(); err != nil {
        log.Printf("Error during shutdown: %v", err)
    }
}