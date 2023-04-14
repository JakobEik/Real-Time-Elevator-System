package watchdog

import (
	"time"
)

const watchdogTimerDuration = time.Millisecond * 500

func Watchdog(ch_bark chan<- bool, ch_pet <-chan bool, moduleName string) {
	wdTimer := time.NewTicker(watchdogTimerDuration)
	ch_bark <- true
	pet := false
	for {
		select {
		case value := <-ch_pet:
			pet = value
		case <-wdTimer.C:
			if pet == false {
				panic("Watchdog timer limit reached for " + moduleName)
			}
			pet = false
			ch_bark <- true

		}
	}
}
