package elevator

/* Put in config!
var (
	elev         Elevator
	doorOpenTime = 3
)

const (
	btnPress Events = iota
	onFloorArrival
	timerTimedOut
)*/

func Fsm(
	ch_orderChan chan ButtonEvent,
	ch_doRequest chan bool,
	ch_floorArrival chan int,
	ch_newRequest chan bool,
	ch_Obstruction chan bool,
	channel_DoorTimer chan bool) {

	elev := InitElev()

	SetDoorOpenLamp(false)
	SetMotorDirection(MD_Down)

	for {
		floor := <-ch_floorArrival
		if floor != 0 {
			SetMotorDirection(MD_Down)
		} else {
			SetMotorDirection(MD_Stop)
			break
		}
	}

	for {
		select {
		case order := <-ch_orderChan:

			switch {
			case elev.Behaviour == ele:
				if elev.Floor == order.Floor {
					// Reset doortimer
				} else {
					// Set order at this point to "true"
				}
			case elev.Behaviour == Idle:
				// blablabla

			case elev.Behaviour == DoorOpen:

			}
		case floor := <-ch_floorArrival:
			elev.Floor = floor
			switch {
			case elev.Behaviour == Moving:

			}
		}

	}

}
