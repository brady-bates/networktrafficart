package dotenv

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	env  *Dotenv
	once sync.Once
)

type Dotenv struct {
	EnableMockPacketEventStream      bool
	MockPacketEventStreamDelayMicros int
	MockPacketEventBatchSize         int
	PacketEventWatcherMaxDelayMicros int
	WritePacketsToCSV                bool
	CsvName                          string
	EnableBPF                        bool
	BPFFilter                        string
}

func GetDotenv() *Dotenv {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		env = &Dotenv{
			EnableMockPacketEventStream:      isTrueStr(os.Getenv("ENABLE_MOCK_PACKET_EVENT_STREAM")),
			MockPacketEventStreamDelayMicros: parseToInt(os.Getenv("MOCK_PACKET_EVENT_STREAM_DELAY_MICROS")),
			MockPacketEventBatchSize:         parseToInt(os.Getenv("MOCK_PACKET_EVENT_BATCH_SIZE")),
			PacketEventWatcherMaxDelayMicros: parseToInt(os.Getenv("PACKET_EVENT_WATCHER_MAX_DELAY_MICROS")),
			WritePacketsToCSV:                isTrueStr(os.Getenv("WRITE_PACKETS_TO_CSV")),
			CsvName:                          os.Getenv("CSV_NAME"),
			EnableBPF:                        isTrueStr(os.Getenv("ENABLE_BPF")),
			BPFFilter:                        strings.TrimSpace(os.Getenv("BPF_FILTER")),
		}
		fmt.Println("Dotenv is initialized")
	})
	return env
}

func isTrueStr(s string) bool {
	return s == "true"
}

func parseToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}
