package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
	"networktrafficart/util"
	"time"
)

type Capture struct {
	Handle *pcap.Handle
	Events chan PacketEvent
}

func NewCaptureProvider(deviceName string) (*Capture, error) {
	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	bufferLen := 2500000
	return &Capture{
		Handle: handle,
		Events: make(chan PacketEvent, bufferLen),
	}, nil
}

func (c *Capture) StartPacketCapture(packetIn chan<- gopacket.Packet, WritePacketsToCSV bool) {
	source := gopacket.NewPacketSource(c.Handle, c.Handle.LinkType())

	for packet := range source.Packets() {
		if WritePacketsToCSV && packetIn != nil {
			select {
			case packetIn <- packet:
			default:
			}
		}

		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil { // TODO update this to get ipv6 packets as well
			ip, _ := ipLayer.(*layers.IPv4)

			event := NewPacketEvent(
				len(packet.Data()),
				ip.SrcIP,
				ip.DstIP,
			)

			select {
			case c.Events <- event:
			default:
				log.Println("Dropped packet (channel full)")
			}
		}
	}
}

func (c *Capture) MockPacketEventStream(delayMicros int, batchSize int) {
	micro := time.Duration(delayMicros) * time.Microsecond
	events := make([]PacketEvent, 0, batchSize)

	for {
		events = events[:0]
		for batch := 0; batch < batchSize; batch++ {
			event := NewPacketEvent(
				500,
				util.GenerateRandomIPv4(),
				util.GenerateRandomIPv4(),
			)
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

func GetInterfaceIPv4(deviceName string) (net.IP, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	for _, d := range devices {
		if d.Name != deviceName {
			continue
		}
		for _, address := range d.Addresses {
			if address.IP.To4() != nil {
				return address.IP.To4(), nil
			}
		}
		return nil, fmt.Errorf("device %s found but has no IPv4 address", deviceName)
	}

	return nil, fmt.Errorf("device %s not found", deviceName)
}
