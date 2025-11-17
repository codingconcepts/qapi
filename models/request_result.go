package models

import (
	"fmt"

	"github.com/rs/zerolog"
)

type RequestResult struct {
	StatusCode int
}

type RequestResults []RequestResult

func (rr RequestResults) Log(logger *zerolog.Logger) {
	event := logger.Info()

	var keys []int
	counts := map[int]int{}
	for _, r := range rr {
		if _, exists := counts[r.StatusCode]; !exists {
			keys = append(keys, r.StatusCode)
		}
		counts[r.StatusCode]++
	}

	for _, key := range keys {
		event = event.Int(fmt.Sprintf("%d", key), counts[key])
	}

	event.Msg("status_codes")
}
