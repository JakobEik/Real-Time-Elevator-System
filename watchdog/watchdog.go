package watchdog

import (
	c "Project/config"
	"fmt"
	"time"
)

func Watchdog(ch_wdstart chan bool, ch_wdstop chan bool, ch_bark chan bool) {
	wdTimer := time.NewTimer(c.WatchdogTimerDuration)
	for {
		select {
		case <-ch_wdstart:
			wdTimer.Reset(c.WatchdogTimerDuration)
			fmt.Println("WATCHDOG RESET")

		case <-ch_wdstop:
			wdTimer.Stop()
			fmt.Println("WATCHDOG STOPPED")

		case <-wdTimer.C:
			println("WATCHDOG BARK FROM WATCHDOG")
			ch_bark <- true
			wdTimer.Stop()
			panic("PANIC: HA DET")
		}
	}
}
