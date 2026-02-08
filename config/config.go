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
	EnableBPF                         bool
	BPFFilter                         string
	PacketEventWatcherAggressionCurve float64
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .config file")
	}
	config = &Config{
		Fullscreen:                        u.IsTrueStr(os.Getenv("FULLSCREEN")),
		EnableMockPacketEventStream:       u.IsTrueStr(os.Getenv("ENABLE_MOCK_PACKET_EVENT_STREAM")),
		MockPacketEventStreamDelayMicros:  u.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_STREAM_DELAY_MICROS")),
		MockPacketEventBatchSize:          u.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_BATCH_SIZE")),
		PacketEventWatcherMaxDelayMicros:  u.ParseToInt(os.Getenv("PACKET_EVENT_WATCHER_MAX_DELAY_MICROS")),
		WritePacketsToCSV:                 u.IsTrueStr(os.Getenv("WRITE_PACKETS_TO_CSV")),
		CsvName:                           os.Getenv("CSV_NAME"),
		EnableBPF:                         u.IsTrueStr(os.Getenv("ENABLE_BPF")),
		BPFFilter:                         strings.TrimSpace(os.Getenv("BPF_FILTER")),
		PacketEventWatcherAggressionCurve: u.ParseToFloat(os.Getenv("PACKET_EVENT_WATCHER_AGGRESSION_CURVE")),
	}
	fmt.Println("Config is initialized")
}

func GetConfig() *Config {
	return config
}
