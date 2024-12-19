package network

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"
	discovery "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"github.com/multiformats/go-multiaddr"
)

const (
	protocolID       = "/myapp/1.0.0"
	rendezvousString = "clipt-lobby"
)

var defaultBootstrapPeers = []string{
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
	"/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zhtx6BHJX9KiKNN6tpvbUcqanj75Nb",
}

// Global variables to track peers
var (
	connectedPeers = make(map[peer.ID]bool)
	peersMutex     sync.RWMutex
)

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	log.Printf("[🔍] Found peer via mDNS: %s", pi.ID.String()[:12])
	if pi.ID == n.h.ID() {
		return
	}

	log.Printf("[👥] Found local peer %s with %d addresses", pi.ID.String()[:12], len(pi.Addrs))
	for _, addr := range pi.Addrs {
		log.Printf("    - %s", addr)
	}

	if err := n.h.Connect(context.Background(), pi); err != nil {
		log.Printf("[❌] Failed to connect to local peer %s: %v", pi.ID.String()[:12], err)
		return
	}
	log.Printf("[🤝] Connected to local peer via mDNS: %s", pi.ID.String()[:12])
}

func RunNode() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	go func() {
		errChan <- runNodeInternal(ctx)
	}()

	// Start periodic peer list display
	go displayPeerList(ctx)

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
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",    // Regular TCP
			"/ip4/0.0.0.0/tcp/0/ws", // WebSocket
		),
		libp2p.NATPortMap(), // Add this line for NAT traversal
	)
	if err != nil {
		return fmt.Errorf("failed to create host: %v", err)
	}
	defer host.Close()

	log.Printf("\n[🟢] Node started with ID: %s", host.ID())
	log.Printf("[📡] Node addresses:")
	for _, addr := range host.Addrs() {
		log.Printf("    - %s/p2p/%s", addr, host.ID())
	}

	log.Printf("[🔍] Setting up mDNS discovery...")
	disc := mdns.NewMdnsService(host, rendezvousString, &discoveryNotifee{h: host})
	if err := disc.Start(); err != nil {
		return fmt.Errorf("failed to start mDNS discovery: %v", err)
	}

	// Create a new DHT client mode first
	dhtOpts := []dht.Option{
		dht.Mode(dht.ModeClient),
		dht.BootstrapPeers(
			func() []peer.AddrInfo {
				peers := []peer.AddrInfo{}
				for _, addr := range defaultBootstrapPeers {
					ma, err := multiaddr.NewMultiaddr(addr)
					if err != nil {
						log.Printf("[⚠️] Invalid bootstrap address: %s", addr)
						continue
					}
					pi, err := peer.AddrInfoFromP2pAddr(ma)
					if err != nil {
						log.Printf("[⚠️] Invalid peer info from address: %s", addr)
						continue
					}
					peers = append(peers, *pi)
				}
				return peers
			}()...,
		),
	}

	// Create a new DHT
	log.Printf("[🔄] Creating new DHT...")
	cliptDHT, err := dht.New(ctx, host, dhtOpts...)
	if err != nil {
		return fmt.Errorf("failed to create DHT: %v", err)
	}

	// Connect to bootstrap peers
	log.Printf("[🔄] Connecting to bootstrap peers...")
	for _, addr := range defaultBootstrapPeers {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			log.Printf("[⚠️] Invalid bootstrap address: %s", addr)
			continue
		}
		pi, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			log.Printf("[⚠️] Invalid peer info from address: %s", addr)
			continue
		}
		if err := host.Connect(ctx, *pi); err != nil {
			log.Printf("[⚠️] Failed to connect to bootstrap peer %s: %v", addr, err)
			continue
		}
		log.Printf("[✅] Connected to bootstrap peer: %s", addr)
	}

	// Bootstrap the DHT
	log.Printf("[🔄] Bootstrapping DHT...")
	if err = cliptDHT.Bootstrap(ctx); err != nil {
		return fmt.Errorf("failed to bootstrap DHT: %v", err)
	}
	log.Printf("[✅] DHT bootstrap complete")

	// Create a routing discovery instance
	routingDiscovery := routing.NewRoutingDiscovery(cliptDHT)

	// Advertise this node
	discovery.Advertise(ctx, routingDiscovery, rendezvousString)
	log.Printf("[📢] Advertising with rendezvous string: %s", rendezvousString)

	// Set up stream handler before finding peers
	host.SetStreamHandler(protocolID, handleStream)
	log.Printf("[🔌] Stream handler set up for protocol: %s", protocolID)

	// Start continuous peer discovery
	go continuouslyFindPeers(ctx, routingDiscovery, host)

	// Start heartbeat to maintain connections
	go maintainConnections(ctx, host)

	<-ctx.Done()
	return nil
}

func continuouslyFindPeers(ctx context.Context, routingDiscovery *routing.RoutingDiscovery, host host.Host) {
	log.Printf("[🔍] Starting peer discovery...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("[🔍] Looking for peers with rendezvous string: %s", rendezvousString)
			peerChan, err := routingDiscovery.FindPeers(ctx, rendezvousString)
			if err != nil {
				log.Printf("[❌] Failed to find peers: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for peer := range peerChan {
				if peer.ID == host.ID() {
					continue
				}

				if len(peer.Addrs) == 0 {
					log.Printf("[⚠️] Found peer %s but it has no addresses", peer.ID.String()[:12])
					continue
				}

				log.Printf("[👥] Found peer %s with %d addresses", peer.ID.String()[:12], len(peer.Addrs))
				for _, addr := range peer.Addrs {
					log.Printf("    - %s", addr)
				}

				if err := connectToPeer(ctx, host, peer); err != nil {
					log.Printf("[❌] Failed to connect to peer %s: %v", peer.ID.String()[:12], err)
				}
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func connectToPeer(ctx context.Context, host host.Host, peerInfo peer.AddrInfo) error {
	peersMutex.RLock()
	if connectedPeers[peerInfo.ID] {
		peersMutex.RUnlock()
		return nil
	}
	peersMutex.RUnlock()

	if err := host.Connect(ctx, peerInfo); err != nil {
		return err
	}

	peersMutex.Lock()
	connectedPeers[peerInfo.ID] = true
	peersMutex.Unlock()

	log.Printf("[🤝] Connected to peer: %s", peerInfo.ID.String()[:12])

	// Open a stream and send hello message
	stream, err := host.NewStream(ctx, peerInfo.ID, protocolID)
	if err != nil {
		log.Printf("[❌] Failed to open stream to peer %s: %v", peerInfo.ID.String()[:12], err)
		return err
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	_, err = rw.WriteString(fmt.Sprintf("Hello from %s\n", host.ID().String()[:12]))
	if err != nil {
		return err
	}
	err = rw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func handleStream(stream network.Stream) {
	remotePeer := stream.Conn().RemotePeer()
	log.Printf("[📨] New stream from peer: %s", remotePeer.String()[:12])

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	go readData(rw, remotePeer)
}

func readData(rw *bufio.ReadWriter, peer peer.ID) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			if !strings.Contains(err.Error(), "stream reset") {
				log.Printf("[❌] Error reading from peer %s: %v", peer.String()[:12], err)
			}
			return
		}
		if str = strings.TrimSpace(str); str != "" {
			log.Printf("[📩] %s says: %s", peer.String()[:12], str)
		}
	}
}

func maintainConnections(ctx context.Context, host host.Host) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peersMutex.Lock()
			for peerID := range connectedPeers {
				if host.Network().Connectedness(peerID) != network.Connected {
					delete(connectedPeers, peerID)
					log.Printf("[👋] Peer disconnected: %s", peerID.String()[:12])
				}
			}
			peersMutex.Unlock()
		}
	}
}

func displayPeerList(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			peersMutex.RLock()
			if len(connectedPeers) > 0 {
				log.Printf("\n[👥] Currently connected peers (%d):", len(connectedPeers))
				for peerID := range connectedPeers {
					log.Printf("    - %s", peerID.String()[:12])
				}
			} else {
				log.Printf("\n[👥] No peers currently connected")
			}
			peersMutex.RUnlock()
		}
	}
}
