package main

import (
	"Project/config"
	"Project/driver"
	"Project/elevator"
)

func main() {
	// channels for order assigner
	ch_requestLocalState := make(chan elevator.ElevatorState)
	ch_currentLocalState := make(chan elevator.ElevatorState)

	// channels for Network
	//ch_peerUpdate := make(chan bruh)
	//ch_stateToNetwork := make(chan globalState)
	//ch_stateFromNetwork := make(chan globalState)

	// channels for order distributor
	ch_localStateUpdated := make(chan config.Order)
	ch_doOrder := make(chan config.Order)
	ch_newOrder := make(chan config.Order)

	// channels for FSM
	ch_newCabCall := make(chan config.Order)
	ch_floorArrival := make(chan int)
	ch_obstruction := make(chan bool)
	ch_stop := make(chan bool)

	// channels for Elevio Driver
	ch_buttons := make(chan driver.ButtonEvent)

	//channel_DoorTimer := make(chan bool)

	go driver.PollButtons(ch_buttons)
	go driver.PollFloorSensor(ch_floorArrival)
	go driver.PollObstructionSwitch(ch_obstruction)
	go driver.PollStopButton(ch_stop)
	elevator.Fsm(ch_doOrder, ch_newCabCall, ch_floorArrival, ch_obstruction, ch_stop)
	//request_executor.fsm(ch_doRequest, ch_floorArrival, ch_newRequest, ch_Obstruction, channel_DoorTimer)
}

/*
	numFloors := 4

	elevio_driver.Init("localhost:15657", numFloors)

	var d elevio_driver.MotorDirection = elevio_driver.MD_Up
	//elevio_driver.SetMotorDirection(d)

	drv_buttons := make(chan elevio_driver.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	//timerChan := make(chan bool)
	//requestsChan := make(chan int, 20)
	//newRequestChan := make(chan int, 20)

	go elevio_driver.PollButtons(drv_buttons)
	go elevio_driver.PollFloorSensor(drv_floors)
	go elevio_driver.PollObstructionSwitch(drv_obstr)
	go elevio_driver.PollStopButton(drv_stop)
}
*/
//go onNewRequest(newRequestChan)

/*
	for {
		select {
		case a := <-drv_buttons:
			fmt.Println("BUTTON:", a)
			onButtonPress(a, requestsChan)
		case request := <-requestsChan:
			fmt.Println("NEW REQUEST", request)
			onNewRequest(request)
		case a := <-drv_floors:
			fmt.Printf("FLOOR %+v\n", a)
			onFloorEvent(a, requestsChan)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio_driver.SetMotorDirection(elevio_driver.MD_Stop)
			} else {
				elevio_driver.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio_driver.ButtonType(0); b < 3; b++ {
					elevio_driver.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func doorOpen() {
	time.Sleep(time.Second * 1)
}

func onNewRequest(floor int) {
	currentFloor := elevatorState.Floor
	fmt.Println("current:", currentFloor, ", req:", floor)

	if floor < currentFloor {
		elevio_driver.SetMotorDirection(elevio_driver.MD_Down)
	} else if floor > currentFloor {
		elevio_driver.SetMotorDirection(elevio_driver.MD_Up)
	}

}

func onFloorEvent(floor int, requestChan chan<- int) {
	elevatorState.Floor = floor
	if len(requestList) > 0 {
		requestedFloor := requestList[0]
		if floor == requestedFloor {
			elevio_driver.SetMotorDirection(elevio_driver.MD_Stop)
			elevio_driver.SetButtonLamp(elevio_driver.BT_HallUp, floor, false)
			elevio_driver.SetButtonLamp(elevio_driver.BT_HallDown, floor, false)
			elevio_driver.SetButtonLamp(elevio_driver.BT_Cab, floor, false)
			requestList = remove(requestList, 0)
			if len(requestList) > 0 {
				requestChan <- requestList[0]
			}
			doorOpen()
		}

	}

}

func onButtonPress(btn_event elevio_driver.ButtonEvent, requestChan chan<- int) {
	fmt.Printf("%+v\n", btn_event)
	elevio_driver.SetButtonLamp(btn_event.Button, btn_event.Floor, true)

	btnFloor := btn_event.Floor
	//btnType := btn_event.Button
	requestList = append(requestList, btnFloor)
	requestChan <- btnFloor
}
*/
