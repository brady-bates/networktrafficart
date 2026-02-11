package mockeventstream

import (
	"networktrafficart/capture"
	"networktrafficart/util"
	"time"
)

func Init(c *capture.Capture, delayMicros int, batchSize int) {
	micro := time.Duration(delayMicros) * time.Microsecond
	events := make([]capture.Event, 0, batchSize)

	for {
		events = events[:0]
		for batch := 0; batch < batchSize; batch++ {
			event := capture.NewEvent(
				500,
				util.GenerateRandomIPv4(),
				util.GenerateRandomIPv4(),
			)
			events = append(events, event)
		}

		for _, event := range events {
			select {
			case c.Events <- event:
			default:
			}
		}

		time.Sleep(micro)
	}
}
