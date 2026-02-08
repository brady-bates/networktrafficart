package main

import (
	"github.com/google/gopacket"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"networktrafficart/capture"
	"networktrafficart/config"
	"networktrafficart/csv"
	"networktrafficart/display"
	"networktrafficart/universe"
)

func main() {
	config.LoadConfig()
	conf := config.GetConfig()

	provider, err := capture.NewCaptureProvider("en0")
	if err != nil {
		log.Fatal(err)
	}

	var csvWriterIn chan gopacket.Packet
	if conf.WritePacketsToCSV {
		csvWriterIn = make(chan gopacket.Packet)
		go csv.StreamToCSV(csvWriterIn)
	}

	go provider.StartPacketCapture(csvWriterIn)

	if conf.EnableMockPacketEventStream {
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
