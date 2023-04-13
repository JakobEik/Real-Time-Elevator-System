package watchdog

import (
	"time"
)

func Watchdog(sec int, ch_wdstart chan bool, ch_wdstop chan bool, ch_bark chan bool) {
	wdTimer := time.NewTimer(time.Duration(sec) * time.Second)
	for {
		select {
		case <-ch_wdstart:
			wdTimer.Reset(time.Duration(sec) * time.Second)

		case <-ch_wdstop:
			wdTimer.Stop()

		case <-wdTimer.C:
			ch_bark <- true
			println("WATCHDOG BARK FROM WATCHDOG")
			wdTimer.Reset(time.Duration(sec) * time.Second)
		}
	}
}
