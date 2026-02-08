package csv

import (
	"encoding/csv"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"networktrafficart/util"
	"networktrafficart/util/shutdown"
	"os"
	"reflect"
)

type PacketRecord struct {
	Timestamp int64
	Length    int32

	SrcMAC    string
	DstMAC    string
	EtherType string

	SrcIP      string
	DstIP      string
	IPVersion  int32
	TTL        int32
	IPProtocol string

	SrcPort  int32
	DstPort  int32
	Protocol string

	TCPSyn bool
	TCPAck bool
	TCPFin bool
	TCPRst bool
	TCPPsh bool

	DNSQuery        string
	DNSIsResponse   bool
	DNSResponseCode int32
}

func WriteCSVHeader(writer *csv.Writer) error {
	return writer.Write(ReflectPacketRecord())
}

func AppendPacketToCSV(writer *csv.Writer, packet gopacket.Packet) error {
	return writer.Write(NewPacketRecord(packet).ToStringArray())
}

func StreamToCSV(packetOut <-chan gopacket.Packet, filename string) {
	sd := shutdown.GetShutDownCtx()
	fileExists, _ := util.FileExists(filename)
	var file *os.File
	var err error
	var writer *csv.Writer
	var packet gopacket.Packet

	// TODO abort write if file header doesn't match
	if !fileExists {
		file, err = os.Create(filename)
		if writer = csv.NewWriter(file); err != nil {
			log.Fatal(err)
		}
		if err = WriteCSVHeader(writer); err != nil {
			log.Fatal(err)
		}
	}
	if file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644); err != nil {
		log.Fatal(err)
	}

	if writer = csv.NewWriter(file); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case packet = <-packetOut:
			if err = AppendPacketToCSV(writer, packet); err != nil {
				log.Fatal(err)
			}
		case <-sd.Ctx.Done():
			goto shutdown
		}
	}

shutdown:
	fmt.Println("Shutdown signal received")
	writer.Flush()
	if err = file.Sync(); err != nil {
		log.Println(err)
	}
	if err = file.Close(); err != nil {
		log.Println(err)
	}
}

func NewPacketRecord(packet gopacket.Packet) PacketRecord {
	r := PacketRecord{
		Timestamp: packet.Metadata().Timestamp.UnixMilli(),
		Length:    int32(packet.Metadata().Length),
	}

	if eth, ok := packet.LinkLayer().(*layers.Ethernet); ok {
		r.SrcMAC = eth.SrcMAC.String()
		r.DstMAC = eth.DstMAC.String()
		r.EtherType = eth.EthernetType.String()
	}

	if ip4 := packet.Layer(layers.LayerTypeIPv4); ip4 != nil {
		v4 := ip4.(*layers.IPv4)
		r.SrcIP = v4.SrcIP.String()
		r.DstIP = v4.DstIP.String()
		r.IPVersion = 4
		r.TTL = int32(v4.TTL)
		r.IPProtocol = v4.Protocol.String()
	}

	if ip6 := packet.Layer(layers.LayerTypeIPv6); ip6 != nil {
		v6 := ip6.(*layers.IPv6)
		r.SrcIP = v6.SrcIP.String()
		r.DstIP = v6.DstIP.String()
		r.IPVersion = 6
		r.TTL = int32(v6.HopLimit)
		r.IPProtocol = v6.NextHeader.String()
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		r.SrcPort = int32(tcp.SrcPort)
		r.DstPort = int32(tcp.DstPort)
		r.Protocol = "TCP"
		r.TCPSyn = tcp.SYN
		r.TCPAck = tcp.ACK
		r.TCPFin = tcp.FIN
		r.TCPRst = tcp.RST
		r.TCPPsh = tcp.PSH
	}

	if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		r.SrcPort = int32(udp.SrcPort)
		r.DstPort = int32(udp.DstPort)
		r.Protocol = "UDP"
	}

	if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		dns := dnsLayer.(*layers.DNS)
		r.DNSIsResponse = dns.QR
		r.DNSResponseCode = int32(dns.ResponseCode)
		if len(dns.Questions) > 0 {
			r.DNSQuery = string(dns.Questions[0].Name)
		}
	}

	return r
}

func (r PacketRecord) ToStringArray() []string {
	return []string{
		fmt.Sprintf("%d", r.Timestamp),
		fmt.Sprintf("%d", r.Length),
		r.SrcMAC,
		r.DstMAC,
		r.EtherType,
		r.SrcIP,
		r.DstIP,
		fmt.Sprintf("%d", r.IPVersion),
		fmt.Sprintf("%d", r.TTL),
		r.IPProtocol,
		fmt.Sprintf("%d", r.SrcPort),
		fmt.Sprintf("%d", r.DstPort),
		r.Protocol,
		fmt.Sprintf("%t", r.TCPSyn),
		fmt.Sprintf("%t", r.TCPAck),
		fmt.Sprintf("%t", r.TCPFin),
		fmt.Sprintf("%t", r.TCPRst),
		fmt.Sprintf("%t", r.TCPPsh),
		r.DNSQuery,
		fmt.Sprintf("%t", r.DNSIsResponse),
		fmt.Sprintf("%d", r.DNSResponseCode),
	}
}

func ReflectPacketRecord() []string {
	t := reflect.TypeOf(PacketRecord{})
	var headers []string
	for i := 0; i < t.NumField(); i++ {
		headers = append(headers, t.Field(i).Name)
	}
	return headers
}
