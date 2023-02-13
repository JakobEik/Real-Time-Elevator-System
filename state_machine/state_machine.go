package state_machine

import (
	"project-group-77/elevio"


)

//type behavior int

//const


//type Elevator struct {
//	floor		int
//	directon	elevio.MotorDirection
//	Requests	[][]int
//	state
//}


func RunElevator(ch StateMachineChannels){

	elevator := Elevator{
		floor : elevio.GetFloor(),
		state : Idle,
		dir   : elevio.MD_Stop,
		elev_queue:

	}



	for {
		select{
		case a := <- new_order


			switch elevator.state
		case
		}

case a := <- floor_arrival

case a := <- timed_out

}
}


