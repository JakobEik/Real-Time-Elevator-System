package watchdog

import (
	"Project/failroutine"
	"time"
)

const watchdogTimerDuration = time.Second * 1

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
				println("Watchdog timer limit reached for " + moduleName)
				failroutine.FailRoutine()
			}
			pet = false
			ch_bark <- true

		}
	}
}
