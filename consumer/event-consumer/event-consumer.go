package event_consumer

import (
	"log"
	"sync"
	"tgBot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvent, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvent) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err = c.handleEvents(gotEvent); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c Consumer) handleEvents(events []events.Event) error {
	var wg sync.WaitGroup

	wg.Add(len(events))

	for _, event := range events {
		go c.processEvent(event, &wg)
	}

	return nil
}
func (c Consumer) processEvent(event events.Event, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("got new event: %s", event.Text)

	if err := c.processor.Process(event); err != nil {
		log.Printf("can't handle event: %s", err.Error())
	}
}
