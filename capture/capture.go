package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

type Capture struct {
	Handle *pcap.Handle
	Events chan Event
}

func NewCaptureProvider(deviceName string) (*Capture, error) {
	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	bufferLen := 50000
	return &Capture{
		Handle: handle,
		Events: make(chan Event, bufferLen),
	}, nil
}

func (c *Capture) StartPacketCapture(packetIn chan<- gopacket.Packet) {
	source := gopacket.NewPacketSource(c.Handle, c.Handle.LinkType())

	var netLayer gopacket.Layer
	for packet := range source.Packets() {
		if packetIn != nil {
			select {
			case packetIn <- packet:
			default:
			}
		}

		netLayer = packet.NetworkLayer()
		if netLayer == nil {
			continue
		}

		if IsValidLayerType(netLayer.LayerType()) { // TODO update this to get ipv6 packets as well
			select {
			case c.Events <- NewEventFromPacket(packet):
			default:
				log.Println("Dropped packet (channel full)")
			}
		} else {
			log.Printf("Dropped packet (invalid network layer type %s)\n", netLayer.LayerType())
		}
	}
}

func (c *Capture) SetHandleBPFFilter(filter string) error {
	return c.Handle.SetBPFFilter(filter)
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
