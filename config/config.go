package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"networktrafficart/util"
	"os"
	"strings"
)

var (
	config *Config
)

type Config struct {
	Fullscreen                        bool
	EnableMockPacketEventStream       bool
	MockPacketEventStreamDelayMicros  int
	MockPacketEventBatchSize          int
	PacketEventWatcherMaxDelayMicros  int
	WritePacketsToCSV                 bool
	CsvName                           string
	PacketFilter                      PacketFilter
	PacketEventWatcherAggressionCurve float64
}

type PacketFilter struct {
	Enable bool
	Filter string
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .config file")
	}

	pf := PacketFilter{
		Enable: util.IsTrueStr(os.Getenv("ENABLE")),
		Filter: strings.TrimSpace(os.Getenv("FILTER")),
	}
	config = &Config{
		Fullscreen:                        util.IsTrueStr(os.Getenv("FULLSCREEN")),
		EnableMockPacketEventStream:       util.IsTrueStr(os.Getenv("ENABLE_MOCK_PACKET_EVENT_STREAM")),
		MockPacketEventStreamDelayMicros:  util.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_STREAM_DELAY_MICROS")),
		MockPacketEventBatchSize:          util.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_BATCH_SIZE")),
		PacketEventWatcherMaxDelayMicros:  util.ParseToInt(os.Getenv("PACKET_EVENT_WATCHER_MAX_DELAY_MICROS")),
		WritePacketsToCSV:                 util.IsTrueStr(os.Getenv("WRITE_PACKETS_TO_CSV")),
		CsvName:                           os.Getenv("CSV_NAME"),
		PacketFilter:                      pf,
		PacketEventWatcherAggressionCurve: util.ParseToFloat(os.Getenv("PACKET_EVENT_WATCHER_AGGRESSION_CURVE")),
	}
	fmt.Println("Config is initialized")
}

func GetConfig() *Config {
	return config
}
