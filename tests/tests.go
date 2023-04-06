package tests

import (
	"time"
)

func Wait(until func() bool, d time.Duration) {
	for {
		select {
		case <-time.After(d):
			return
		default:
			if until() {
				return
			}
		}
	}
}
