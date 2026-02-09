package capture

import "net"

type PacketEvent struct {
	Size  int
	SrcIP net.IP
	DstIP net.IP
}

func NewPacketEvent(size int, srcIP, dstIP net.IP) PacketEvent {
	return PacketEvent{
		Size:  size,
		SrcIP: srcIP,
		DstIP: dstIP,
	}
}
