package timer

import (
	"Project/utilities"
	"time"
)

func DoorTimer(channel_DoorTimer chan<- bool) {
	for obstruction {
		time.Sleep(time.Second * utilities.DoorOpenDuration)
	}
	channel_DoorTimer <- true
}

/*package timer

import (
	"syscall"
)

func getWallTime() float64 {
	var t syscall.Timeval
	syscall.Gettimeofday(&t)
	return float64(t.Sec) + float64(t.Usec)*0.000001
}

var (
	timerEndTime float64
	timerActive  bool
)

func timerStart(duration float64) {
	timerEndTime = getWallTime() + duration
	timerActive = true
}

func timerStop() {
	timerActive = false
}

func timerTimedOut() bool {
	return timerActive && getWallTime() > timerEndTime
}
*/
