package main

import (
	"fmt"
	"project-group-77/elevator"
	"project-group-77/elevio"
)

var elevatorState = elevator.UninitializedElevator()

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			onButtonPress(a)
		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

func onButtonPress(btn_event elevio.ButtonEvent) {
	fmt.Printf("%+v\n", btn_event)
	elevio.SetButtonLamp(btn_event.Button, btn_event.Floor, true)

	btnFloor := btn_event.Floor
	btnType := btn_event.Button

	fmt.Printf("\n\n%s(%d, %s)\n", "handleButtonPress", btnFloor, btnType.String())
	elevator.ElevPrint(elevatorState)
	if elevatorState.Floor > btnFloor {
		elevio.SetMotorDirection(elevio.MD_Down)
	} else if elevatorState.Floor < btnFloor {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Stop)
	}

}
