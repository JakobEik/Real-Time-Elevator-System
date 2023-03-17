package timer


import "fmt"
import "time"

type TimerBehavior struct{
	Start bool
	Duration time.Duration
}

func Timer(ch_timer <-chan TimerBehavior, ch_timerDone chan<- bool) {
	timer := time.NewTimer(time.Second *2)
	timer.Stop()
	
	for{
		select{
		case behavior := <- ch_timer:
			if behavior.Start {
				timer.Stop()
				timer.Reset(behavior.Duration)
				fmt.Println("New Timer")
				ch_timerDone <- false
			}else{
				timer.Stop()
				fmt.Println("STOP")
				ch_timerDone <- false
			}

		case <- timer.C:
			fmt.Println("finished timer")
			ch_timerDone <- true
		}
		
	}
}
