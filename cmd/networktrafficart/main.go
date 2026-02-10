package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net"
	"networktrafficart/capture"
	"networktrafficart/config"
	"networktrafficart/csv"
	"networktrafficart/display"
	"networktrafficart/simulation"
)

func main() {
	config.LoadConfig()
	conf := config.GetConfig()

	captureDeviceName := "en0" // TODO add logic to get the best device to capture with

	capt, err := capture.NewCaptureProvider(captureDeviceName)
	if err != nil {
		log.Fatal(err)
	}

	if conf.PacketFilter.Enable {
		var ipv4 net.IP
		if ipv4, err = capture.GetInterfaceIPv4(captureDeviceName); err != nil {
			log.Fatal(err)
		}

		filter := fmt.Sprintf("%s %s", conf.PacketFilter.Filter, ipv4.String())
		if err = capt.Handle.SetBPFFilter(filter); err != nil {
			log.Println("Failed to set packet filter ", err)
		}
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

	sim := simulation.NewSimulation(capt.Events)
	disp := display.NewDisplay(sim)

	ebiten.SetWindowTitle("NetworkTrafficArt")
	ebiten.SetFullscreen(conf.Fullscreen)

	sim.Init(
		disp.ScreenWidth,
		disp.ScreenHeight,
		conf.ParticleBufferConsumerAggressionCurve,
		conf.ParticleBufferConsumerMaxDelayMicros,
	)
	if err = ebiten.RunGame(disp); err != nil {
		log.Fatal(err)
	}
}
