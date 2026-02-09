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

	capt, err := capture.NewCaptureProvider("en0", conf.PacketFilter)
	if err != nil {
		log.Fatal(err)
	}

	var csvWriterIn chan gopacket.Packet
	if conf.WritePacketsToCSV {
		csvWriterIn = make(chan gopacket.Packet)
		go csv.StreamToCSV(csvWriterIn, conf.CsvName)
	}

	go capt.StartPacketCapture(csvWriterIn, conf.WritePacketsToCSV)

	if conf.EnableMockPacketEventStream {
		go capt.MockPacketEventStream(conf.MockPacketEventStreamDelayMicros, conf.MockPacketEventBatchSize)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetFullscreen(false)

	d := display.NewDisplay(
		capt.Events,
		universe.NewUniverse(),
	)
	go d.WatchPacketEventChannel(conf.PacketEventWatcherAggressionCurve, conf.PacketEventWatcherMaxDelayMicros)

	ebiten.SetWindowTitle("NetworkTrafficArt")
	ebiten.SetFullscreen(conf.Fullscreen)
	err = ebiten.RunGame(d)
	if err != nil {
		log.Fatal(err)
	}
}
