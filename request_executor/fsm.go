package finiteStateMachine

import (
	elevator_state "project-group-77/elevator"
	"project-group-77/elevio_driver"
)

var (
	elev         Elevator
	doorOpenTime = 3
)

const (
	btnPress Events = iota
	onFloorArrival
	timerTimedOut
)

func Fsm(
	ch_orderChan chan elevio_driver.ButtonEvent,
	ch_doRequest chan bool,
	ch_floorArrival chan int,
	ch_newRequest chan bool,
	ch_Obstruction chan bool,
	channel_DoorTimer chan bool) {

	elev := elevator_state.InitElev()

	elevio_driver.SetDoorOpenLamp(false)
	elevio_driver.SetMotorDirection(elevio_driver.MD_Down)

	for {
		floor := <-ch_floorArrival
		if floor != 0 {
			elevio_driver.SetMotorDirection(elevio_driver.MD_Down)
		} else {
			elevio_driver.SetMotorDirection(elevio_driver.MD_Stop)
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
			case elev.Behaviour == elevator_state.Idle:
				// blablabla

			case elev.Behaviour == elevator_state.DoorOpen:

			}
		case floor := <-ch_floorArrival:
			elev.Floor = floor
			switch {
			case elev.Behaviour == elevator_state.Moving:

			}
		}

	}

}

/*
func fsm() {
	// initialize elevator and output device

	var e elevator.Elevator = elevator.UninitializedElevator()
	//outputDevice := elevator_io.ElevioGetOutputDevice()

	// load config from file
	//config, err := LoadConfig("elevator.con")
	//if err != nil {
	//	fmt.Println("Error loading config:", err)
	//	return
	//}

	e.Config = config

	// set initial motor direction and behavior
	outputDevice.MotorDirection(elevator_io.MdDown)
	e.Dirn = elevator_io.MdDown
	e.Behaviour = elevator.EbMoving

	// continuously listen for elevator events
	for {
		select {
		case buttonPress := <-elevator_io.ButtonPressCh:
			handleButtonPress(buttonPress, &e, outputDevice)
		case floor := <-elevator_io.FloorSensorCh:
			handleFloorArrival(floor, &e, outputDevice)
		case <-timer.TimeoutCh:
			handleDoorTimeout(&e, outputDevice)
		default:
			handleIdle(&e, outputDevice)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func requests_shouldStop(e Elevator) bool {
	switch e.dirn {
	case D_Down:
		return (
			e.requests[e.floor][B_HallDown] ||
				e.requests[e.floor][B_Cab] ||
				!requests_below(e))
	case D_Up:
		return (
			e.requests[e.floor][B_HallUp] ||
				e.requests[e.floor]const (
					EB_Idle
				)
		fallthrough
	default:
		return true
	}
}










//func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType{

*/
