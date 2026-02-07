package main

import (
	"github.com/google/gopacket"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"networktrafficart/networktrafficart/capture"
	"networktrafficart/networktrafficart/csv"
	"networktrafficart/networktrafficart/display"
	"networktrafficart/networktrafficart/dotenv"
	"networktrafficart/networktrafficart/universe"
)

func main() {
	env := dotenv.GetDotenv()
	var packetChan chan gopacket.Packet

	provider, err := capture.NewCaptureProvider("en0")
	if err != nil {
		log.Fatal(err)
	}

	if env.WritePacketsToCSV {
		packetChan = make(chan gopacket.Packet)
		go csv.StreamToCSV(packetChan)
	}

	go provider.StartPacketCapture(packetChan)

	if env.EnableMockPacketEventStream {
		go provider.MockPacketEventStream()
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)

	d := display.NewDisplay(
		provider.Events,
		universe.NewUniverse(),
	)

	err = ebiten.RunGame(d)
	if err != nil {
		log.Fatal(err)
	}
}
