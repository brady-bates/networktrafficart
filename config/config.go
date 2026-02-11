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
	EnableMockEventStream                 bool
	MockEventStreamDelayMicros            int
	MockEventBatchSize                    int
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
		EnableMockEventStream:                 util.IsTrueStr(os.Getenv("ENABLE_MOCK_EVENT_STREAM")),
		MockEventStreamDelayMicros:            util.ParseToInt(os.Getenv("MOCK_EVENT_STREAM_DELAY_MICROS")),
		MockEventBatchSize:                    util.ParseToInt(os.Getenv("MOCK_EVENT_BATCH_SIZE")),
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
