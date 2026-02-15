package mockeventstream

import (
	capture2 "networktrafficart/internal/capture"
	"networktrafficart/internal/util"
	"time"
)

func Init(c *capture2.Capture, delayMicros int, batchSize int) {
	micro := time.Duration(delayMicros) * time.Microsecond
	events := make([]capture2.Event, 0, batchSize)

	for {
		events = events[:0]
		for batch := 0; batch < batchSize; batch++ {
			event := capture2.NewEvent(
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
