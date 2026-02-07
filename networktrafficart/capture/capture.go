package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"math/rand"
	"net"
	"networktrafficart/networktrafficart/dotenv"
	"time"
)

type PacketEvent struct {
	Size  int
	SrcIP net.IP
	DstIP net.IP
}

type CaptureProvider struct {
	handle *pcap.Handle
	Events chan PacketEvent
}

func NewCaptureProvider(device string) (*CaptureProvider, error) {
	handle, err := pcap.OpenLive(device, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	bufferLen := 25000
	return &CaptureProvider{
		handle: handle,
		Events: make(chan PacketEvent, bufferLen),
	}, nil
}

func (c *CaptureProvider) StartPacketCapture(packetChan chan<- gopacket.Packet) {
	env := dotenv.GetDotenv()
	source := gopacket.NewPacketSource(c.handle, c.handle.LinkType())

	for packet := range source.Packets() {
		if env.WritePacketsToCSV && packetChan != nil {
			select {
			case packetChan <- packet:
			default:
			}
		}

		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil { // TODO update this to get ipv6 packets as well
			ip, _ := ipLayer.(*layers.IPv4)

			event := PacketEvent{
				Size:  len(packet.Data()),
				SrcIP: ip.SrcIP,
				DstIP: ip.DstIP,
			}

			select {
			case c.Events <- event:
			default:
			}
		}
	}
}

func (c *CaptureProvider) MockPacketEventStream() {
	env := dotenv.GetDotenv()
	micro := time.Duration(env.MockPacketEventStreamDelayMicros) * time.Microsecond
	events := make([]PacketEvent, 0, env.MockPacketEventBatchSize)

	for {
		events = events[:0]
		for batch := 0; batch < env.MockPacketEventBatchSize; batch++ {
			event := PacketEvent{
				Size:  500,
				SrcIP: generateRandomIPv4(),
				DstIP: generateRandomIPv4(),
			}
			events = append(events, event)
		}

		for _, event := range events {
			select {
			case c.Events <- event:
			default:
			}
		}

		time.Sleep(micro)
	}
}

func generateRandomIPv4() net.IP {
	o1 := byte(rand.Intn(256))
	o2 := byte(rand.Intn(256))
	o3 := byte(rand.Intn(256))
	o4 := byte(rand.Intn(256))

	return net.IPv4(o1, o2, o3, o4).To4()
}
