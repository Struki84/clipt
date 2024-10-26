// In your p2p.go or similar file
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
)

const (
	protocolID       = "/myapp/1.0.0"
	rendezvousString = "clipt-lobby"
)

func RunNode() error {
	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to handle shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Create error channel
	errChan := make(chan error, 1)

	go func() {
		errChan <- runNodeInternal(ctx)
	}()

	// Wait for either interrupt or error
	select {
	case err := <-errChan:
		return err
	case <-interrupt:
		log.Println("Received interrupt signal, shutting down...")
		cancel()
		return <-errChan
	}
}

func runNodeInternal(ctx context.Context) error {
	// Create a new libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0/ws"),
		libp2p.Transport(websocket.New),
	)
	if err != nil {
		return fmt.Errorf("failed to create host: %v", err)
	}
	defer host.Close()

	// Print the node's addresses
	log.Printf("Node started with ID: %s", host.ID())
	log.Printf("Node addresses:")
	for _, addr := range host.Addrs() {
		log.Printf("  - %s/p2p/%s", addr, host.ID())
	}

	// Create a new DHT
	cliptDHT, err := dht.New(ctx, host)
	if err != nil {
		return fmt.Errorf("failed to create DHT: %v", err)
	}

	// Bootstrap the DHT with default bootstrap nodes
	if err = cliptDHT.Bootstrap(ctx); err != nil {
		return fmt.Errorf("failed to bootstrap DHT: %v", err)
	}

	// Create a routing discovery instance
	routingDiscovery := routing.NewRoutingDiscovery(cliptDHT)

	// Advertise this node
	discovery.Advertise(ctx, routingDiscovery, rendezvousString)
	log.Printf("Advertising with rendezvous string: %s", rendezvousString)

	// Find peers
	peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
	if err != nil {
		return fmt.Errorf("failed to find peers: %v", err)
	}

	// Handle peer discovery
	go func() {
		for peer := range peerChan {
			select {
			case <-ctx.Done():
				return
			default:
				if peer.ID == host.ID() {
					continue
				}
				log.Printf("Found peer: %s", peer.ID)
				if err := host.Connect(ctx, peer); err != nil {
					log.Printf("Failed to connect to peer %s: %v", peer.ID, err)
					continue
				}
				log.Printf("Connected to peer: %s", peer.ID)
			}
		}
	}()

	// Set up stream handler
	host.SetStreamHandler(protocolID, handleStream)
	log.Printf("Stream handler set up for protocol: %s", protocolID)

	// Keep running until context is canceled
	<-ctx.Done()
	return nil
}

func handleStream(stream network.Stream) {
	remotePeer := stream.Conn().RemotePeer()
	log.Printf("New stream from peer: %s", remotePeer)
	defer stream.Close()
}
