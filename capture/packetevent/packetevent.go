package packetevent

import "net"

type PacketEvent struct {
	Size  int
	SrcIP net.IP
	DstIP net.IP
}
