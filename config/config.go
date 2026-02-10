package config

import (
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
	Fullscreen                            bool
	EnableMockPacketEventStream           bool
	MockPacketEventStreamDelayMicros      int
	MockPacketEventBatchSize              int
	ParticleBufferConsumerMaxDelayMicros  int
	WritePacketsToCSV                     bool
	CsvName                               string
	PacketFilter                          PacketFilter
	ParticleBufferConsumerAggressionCurve float64
}

type PacketFilter struct {
	Enable bool
	Filter string
}

// TODO add handling for missing env values?
func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	pf := PacketFilter{
		Enable: util.IsTrueStr(os.Getenv("ENABLE")),
		Filter: strings.TrimSpace(os.Getenv("FILTER")),
	}
	config = &Config{
		Fullscreen:                            util.IsTrueStr(os.Getenv("FULLSCREEN")),
		EnableMockPacketEventStream:           util.IsTrueStr(os.Getenv("ENABLE_MOCK_PACKET_EVENT_STREAM")),
		MockPacketEventStreamDelayMicros:      util.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_STREAM_DELAY_MICROS")),
		MockPacketEventBatchSize:              util.ParseToInt(os.Getenv("MOCK_PACKET_EVENT_BATCH_SIZE")),
		ParticleBufferConsumerMaxDelayMicros:  util.ParseToInt(os.Getenv("PARTICLE_BUFFER_CONSUMER_MAX_DELAY_MICROS")),
		WritePacketsToCSV:                     util.IsTrueStr(os.Getenv("WRITE_PACKETS_TO_CSV")),
		CsvName:                               os.Getenv("CSV_NAME"),
		PacketFilter:                          pf,
		ParticleBufferConsumerAggressionCurve: util.ParseToFloat(os.Getenv("PARTICLE_BUFFER_CONSUMER_AGGRESSION_CURVE")),
	}
}

func GetConfig() *Config {
	return config
}
