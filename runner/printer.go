package runner

import (
	"time"

	"github.com/codingconcepts/qapi/models"
	"github.com/rs/zerolog"
)

type Printer struct {
	events chan models.RequestResult
	logger *zerolog.Logger
}

func NewPrinter(events chan models.RequestResult, logger *zerolog.Logger) *Printer {
	return &Printer{
		events: events,
		logger: logger,
	}
}

func (p *Printer) Run() {
	printTicks := time.Tick(time.Second)

	var events models.RequestResults
	for {
		select {
		case rr := <-p.events:
			events = append(events, rr)

		case <-printTicks:
			events.Log(p.logger)
		}
	}
}
