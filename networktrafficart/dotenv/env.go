package dotenv

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	u "networktrafficart/networktrafficart/util"
	"os"
	"strings"
	"sync"
)

var (
	env  *Dotenv
	once sync.Once
)

type Dotenv struct {
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

func GetDotenv() *Dotenv {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		env = &Dotenv{
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
		fmt.Println("Dotenv is initialized")
	})
	return env
}
