package capture

import "net"

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
