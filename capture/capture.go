package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"networktrafficart/capture/packetevent"
	"networktrafficart/config"
	"networktrafficart/util"
	"time"
)

const (
	IPv4Range = 4294967295.0
)

type Capture struct {
	handle *pcap.Handle
	Events chan packetevent.PacketEvent
}

func NewCaptureProvider(deviceName string, bpfConfig config.BerkeleyPacketFilter) (*Capture, error) {
	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	ipv4, err := getInterfaceIPv4(deviceName)
	if err != nil {
		return nil, err
	}

	if bpfConfig.Enable {
		filter := fmt.Sprintf("%s %s", bpfConfig.Filter, ipv4.String())
		err = handle.SetBPFFilter(filter)
		if err != nil {
			return nil, err
		}
	}

	bufferLen := 25000
	return &Capture{
		handle: handle,
		Events: make(chan packetevent.PacketEvent, bufferLen),
	}, nil
}

func (c *Capture) StartPacketCapture(packetIn chan<- gopacket.Packet, WritePacketsToCSV bool) {
	source := gopacket.NewPacketSource(c.handle, c.handle.LinkType())

	for packet := range source.Packets() {
		if WritePacketsToCSV && packetIn != nil {
			select {
			case packetIn <- packet:
			default:
			}
		}

		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil { // TODO update this to get ipv6 packets as well
			ip, _ := ipLayer.(*layers.IPv4)

			event := packetevent.PacketEvent{
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

func (c *Capture) MockPacketEventStream(delayMicros int, batchSize int) {
	micro := time.Duration(delayMicros) * time.Microsecond
	events := make([]packetevent.PacketEvent, 0, batchSize)

	for {
		events = events[:0]
		for batch := 0; batch < batchSize; batch++ {
			event := packetevent.PacketEvent{
				Size:  500,
				SrcIP: util.GenerateRandomIPv4(),
				DstIP: util.GenerateRandomIPv4(),
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

func getInterfaceIPv4(deviceName string) (net.IP, error) {
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
