package capture

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"math/rand"
	"net"
)

type Event struct {
	Size      int
	SrcIP     net.IP
	DstIP     net.IP
	IsInbound bool
}

func NewEvent(size int, srcIP, dstIP net.IP) Event {
	return Event{
		Size:      size,
		SrcIP:     srcIP,
		DstIP:     dstIP,
		IsInbound: rand.Intn(2) == 1,
	}
}

func NewEventFromPacket(packet gopacket.Packet, subnet *net.IPNet) Event {
	var srcIP, dstIP net.IP

	// TODO convert or not based on the actual given ip type
	switch layer := packet.NetworkLayer().(type) {
	case *layers.IPv4:
		srcIP = normalizeIP(layer.SrcIP)
		dstIP = normalizeIP(layer.DstIP)
	case *layers.IPv6:
		srcIP = normalizeIP(layer.SrcIP)
		dstIP = normalizeIP(layer.DstIP)
	default:
		fmt.Printf("Unknown layer type %s - check if layer type is valid before calling\n", packet.NetworkLayer().LayerType().String())
	}

	return Event{
		packet.Metadata().Length,
		srcIP,
		dstIP,
		subnet.Contains(dstIP),
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

// TODO find a better way to convert IPv4 and IPv6 to a common representation
func iPv6toIPv4Format(ip net.IP) net.IP {
	return ip[12:16].To4()
}

func isIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}

func normalizeIP(ip net.IP) net.IP {
	if isIPv6(ip) {
		return iPv6toIPv4Format(ip)
	} else {
		return ip
	}
}
