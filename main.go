package main

import (
	"fmt"
	"project-group-77/elevator"
	"project-group-77/elevio"
	"time"
)

var elevatorState = elevator.UninitializedElevator()

var requestList []int

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	requestsChan := make(chan int, 20)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	for {
		select {
		case a := <-drv_buttons:
			fmt.Println("BUTTON:", a)
			onButtonPress(a, requestsChan)
		case request := <- requestsChan:
			fmt.Println("NEW REQUEST", request)
			onNewRequest(request)
		case a := <-drv_floors:
			fmt.Printf("FLOOR %+v\n", a)
			onFloorEvent(a, requestsChan)


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

func remove(slice []int, s int) []int {
    return append(slice[:s], slice[s+1:]...)
}

func doorOpen(){
	time.Sleep(time.Second * 1)
}

func onNewRequest(floor int){
	currentFloor := elevatorState.Floor
	fmt.Println("current:", currentFloor, ", req:", floor)
	if floor < currentFloor {
		elevio.SetMotorDirection(elevio.MD_Down)
	}else if floor > currentFloor{
		elevio.SetMotorDirection(elevio.MD_Up)
	}
	
}

func onFloorEvent(floor int, requestChan chan<- int){
	elevatorState.Floor = floor
	if len(requestList) > 0 {
		requestedFloor := requestList[0]
		if floor == requestedFloor {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
			requestList = remove(requestList, 0)
			if len(requestList) > 0{
				requestChan <- requestList[0]
			}
			doorOpen()
		}



	}



}

func onButtonPress(btn_event elevio.ButtonEvent, requestChan chan<- int) {
	fmt.Printf("%+v\n", btn_event)
	elevio.SetButtonLamp(btn_event.Button, btn_event.Floor, true)

	btnFloor := btn_event.Floor
	//btnType := btn_event.Button
	requestList = append(requestList, btnFloor)
	requestChan <- btnFloor
}