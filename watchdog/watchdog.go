package watchdog

import (
	"fmt"
	"time"
)

func Watchdog(sec int, ch_wdstart chan bool, ch_wdstop chan bool, ch_bark chan bool) {
	wdTimer := time.NewTimer(time.Duration(sec) * time.Second)
	for {
		select {
		case <-ch_wdstart:
			wdTimer.Reset(time.Duration(sec) * time.Second)
			fmt.Println("WATCHDOG RESET")

		case <-ch_wdstop:
			wdTimer.Stop()
			fmt.Println("WATCHDOG STOPPED")

		case <-wdTimer.C:
			ch_bark <- true
			wdTimer.Stop()
			println("WATCHDOG BARK FROM WATCHDOG")
		}
	}
}
