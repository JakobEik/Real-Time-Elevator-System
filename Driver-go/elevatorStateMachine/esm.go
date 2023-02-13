package elevatorStateMachine

import (
	"Driver-go/elevio"
	"ftm"
	"time"

	
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

