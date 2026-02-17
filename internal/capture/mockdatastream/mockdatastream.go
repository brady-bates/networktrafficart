package mockdatastream

import (
	"networktrafficart/internal/capture"
	"networktrafficart/internal/util"
	"time"
)

func Init(c *capture.Capture, delayMicros int, batchSize int) {
	micro := time.Duration(delayMicros) * time.Microsecond
	for {
		for range batchSize {
			select {
			case c.Events <- capture.NewPacketData(500, util.GenerateRandomIPv4(), util.GenerateRandomIPv4()):
			default:
			}
		}
		time.Sleep(micro)
	}
}
