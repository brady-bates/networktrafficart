package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"net"
)

type Event struct {
	Size  int
	SrcIP net.IP
	DstIP net.IP
}

func NewEvent(size int, srcIP, dstIP net.IP) Event {
	return Event{
		Size:  size,
		SrcIP: srcIP,
		DstIP: dstIP,
	}
}

func NewEventFromPacket(packet gopacket.Packet) Event {
	var srcIP, dstIP net.IP

	netLayer := packet.NetworkLayer()
	switch layer := netLayer.(type) {
	case *layers.IPv4:
		srcIP = layer.SrcIP
		dstIP = layer.DstIP
	case *layers.IPv6:
		// TODO find a better way to convert IPv4 and IPv6 to a common representation
		srcIP = layer.SrcIP[12:16]
		dstIP = layer.DstIP[12:16]
	default:
		log.Println("Unsupported layer type - setting to defaults", packet.NetworkLayer().LayerType())
		return Event{
			Size:  20,
			SrcIP: net.IPv4zero,
			DstIP: net.IPv4zero,
		}
	}

	return Event{
		packet.Metadata().Length,
		srcIP,
		dstIP,
	}
}

func IsValidLayerType(layer gopacket.LayerType) bool {
	switch layer {
	case layers.LayerTypeIPv4,
		layers.LayerTypeIPv6:
		return true
	default:
		return false
	}
}
