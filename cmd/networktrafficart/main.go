package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"net"
	"networktrafficart/capture"
	"networktrafficart/capture/mockeventstream"
	"networktrafficart/config"
	"networktrafficart/csv"
	"networktrafficart/display"
	"networktrafficart/simulation"
	"networktrafficart/util"
)

const (
	title             = "Network Traffic Art"
	captureDeviceName = "en0" // TODO add logic to get the best device to capture with
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}
	conf := config.GetConfig()

	capt, err := capture.NewCaptureProvider(captureDeviceName)
	if err != nil {
		log.Fatal(err)
	}

	if conf.EnablePacketCaptureFilter {
		var ipv4 net.IP
		if ipv4, err = capture.GetInterfaceIPv4(captureDeviceName); err != nil {
			log.Fatal(err)
		}

		filter := fmt.Sprintf("%s %s", conf.PacketCaptureFilter, ipv4.String())
		if err = capt.SetHandleBPFFilter(filter); err != nil {
			log.Println("Failed to set packet filter ", err)
		}
	}

	var csvWriterIn chan gopacket.Packet
	if conf.WritePacketsToCSV {
		csvWriterIn = make(chan gopacket.Packet)
		go csv.StreamToCSV(util.GetShutDownCtx(), csvWriterIn, conf.CsvName)
	}

	go capt.StartPacketCapture(csvWriterIn)

	if conf.EnableMockEventStream {
		go mockeventstream.Init(capt, conf.MockEventStreamDelayMicros, conf.MockEventBatchSize)
	}

	sim := simulation.NewSimulation(capt.Events)
	disp := display.NewDisplay(sim)

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle(title)
	ebiten.SetFullscreen(conf.Fullscreen)

	sim.Init(
		disp.ScreenWidth,
		disp.ScreenHeight,
		conf.ParticleBufferConsumerMaxDelayMicros,
		conf.ParticleBufferConsumerAggressionCurve,
	)
	if err = ebiten.RunGame(disp); err != nil {
		log.Fatal(err)
	}
}
