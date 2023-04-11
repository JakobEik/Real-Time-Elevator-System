package watchdog

import (
	"time"
)

func Watchdog(sec int, ch_alive chan bool, ch_dead chan bool) {
	wdTimer := time.NewTimer(time.Duration(sec) * time.Second)
	for {
		select {
		case <-ch_alive:
			wdTimer.Reset(time.Duration(sec) * time.Second)

		case <-wdTimer.C:
			ch_dead <- true
			wdTimer.Reset(time.Duration(sec) * time.Second)
		}
	}
}
