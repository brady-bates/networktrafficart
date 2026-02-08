package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	u "networktrafficart/util"
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
	BPF                               BerkeleyPacketFilter
	PacketEventWatcherAggressionCurve float64
}

type BerkeleyPacketFilter struct {
	Enable bool
	Filter string
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .config file")
	}

	bpf := BerkeleyPacketFilter{
		Enable: u.IsTrueStr(os.Getenv("ENABLE_BPF")),
		Filter: strings.TrimSpace(os.Getenv("BPF_FILTER")),
	}
	config = &Config{
		Fullscreen:                        u.IsTrueStr(os.Getenv("FULLSCREEN")),
		EnableMockPacketEventStream:       u.IsTrueStr(os.Getenv("ENABLE_MOCK_PACKET_EVENT_STREAM")),
		MockPacketEventStreamDelayMicros:  u.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_STREAM_DELAY_MICROS")),
		MockPacketEventBatchSize:          u.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_BATCH_SIZE")),
		PacketEventWatcherMaxDelayMicros:  u.ParseToInt(os.Getenv("PACKET_EVENT_WATCHER_MAX_DELAY_MICROS")),
		WritePacketsToCSV:                 u.IsTrueStr(os.Getenv("WRITE_PACKETS_TO_CSV")),
		CsvName:                           os.Getenv("CSV_NAME"),
		BPF:                               bpf,
		PacketEventWatcherAggressionCurve: u.ParseToFloat(os.Getenv("PACKET_EVENT_WATCHER_AGGRESSION_CURVE")),
	}
	fmt.Println("Config is initialized")
}

func GetConfig() *Config {
	return config
}
