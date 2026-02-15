package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

type Capture struct {
	Handle      *pcap.Handle
	Events      chan Event
	localSubnet *net.IPNet
}

func NewCaptureProvider(deviceName string, subnet *net.IPNet) (*Capture, error) {
	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	bufferLen := 50000
	return &Capture{
		Handle:      handle,
		Events:      make(chan Event, bufferLen),
		localSubnet: subnet,
	}, nil
}

func (c *Capture) StartPacketCapture(packetIn chan<- gopacket.Packet) {
	source := gopacket.NewPacketSource(c.Handle, c.Handle.LinkType())

	for packet := range source.Packets() {
		if packetIn != nil {
			select {
			case packetIn <- packet:
			default:
			}
		}

		if packet.NetworkLayer() == nil {
			continue
		}

		if IsValidLayerType(packet.NetworkLayer().LayerType()) { // TODO update this to get ipv6 packets as well
			select {
			case c.Events <- NewEventFromPacket(packet, c.localSubnet):
			default:
				log.Println("Dropped packet (channel full)")
			}
		} else {
			log.Fatalf("Dropped packet (invalid network layer type %s)\n", packet.NetworkLayer().LayerType())
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

func GetIPv4SubnetRange() (*net.IPNet, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		for _, address := range device.Addresses {
			if address.IP == nil || address.IP.IsLoopback() || address.IP.To4() == nil {
				continue
			}

			ipNet := &net.IPNet{
				IP:   address.IP.Mask(address.Netmask),
				Mask: address.Netmask,
			}
			return ipNet, nil
		}
	}

	return nil, fmt.Errorf("could not find any IPv4 addresses")
}
