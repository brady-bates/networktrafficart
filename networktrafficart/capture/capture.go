package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"networktrafficart/networktrafficart/dotenv"
	"networktrafficart/networktrafficart/util"
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

func NewCaptureProvider(deviceName string) (*CaptureProvider, error) {
	env := dotenv.GetDotenv()

	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	ipv4, err := getNetInterfaceIPv4(deviceName)
	if err != nil {
		return nil, err
	}

	if env.EnableBPF {
		bpfFilter := fmt.Sprintf("%s %s", env.BPFFilter, ipv4.String())
		err = handle.SetBPFFilter(bpfFilter)
	}
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

func getNetInterfaceIPv4(deviceName string) (net.IP, error) {
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
