package elevator

/*import (
	c "Project/config"
	drv "Project/driver"
	"Project/timer"
	"fmt"
	"time"
)

func Fsm(
	ch_doOrder <-chan drv.ButtonEvent,
	ch_floorArrival <-chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- ElevatorState) {
	
	ch_doorTimer := make(chan timer.TimerBehavior)
	ch_timerDone := make(chan bool)
	go timer.Timer(ch_doorTimer, ch_timerDone)

	Stop := false

	elev := InitElev()
	e_ptr := &elev

	drv.SetDoorOpenLamp(false)
	drv.SetMotorDirection(drv.MD_Down)

	for {
		floor := <-ch_floorArrival
		drv.SetFloorIndicator(floor)
		if floor != 0 {
			drv.SetMotorDirection(drv.MD_Down)
		} else {
			drv.SetMotorDirection(drv.MD_Stop)
			break
		}
	}

	for {
		select {
		case order := <-ch_doOrder:
			println("NEW BUTTONPRESS!")
			onNewOrderEvent(order, e_ptr, ch_doorTimer)
			printState(elev)
			ch_newLocalState <- elev
		case floor := <-ch_floorArrival:
			println("Floor arrival:", floor)
			if Stop {
				elev = InitElev()
				e_ptr = &elev

				if floor != 0 {
					drv.SetMotorDirection(drv.MD_Down)
				} else {
					drv.SetMotorDirection(drv.MD_Stop)
					Stop = false
					break
				}

			} else {
				onFloorArrivalEvent(Stop, floor, e_ptr, ch_doorTimer)
				nextOrder(elev)
			}
			printState(elev)
			ch_newLocalState <- elev

		case stop := <-ch_stop:
			Stop = true
			onStopEvent(stop, &elev, ch_floorArrival)
			ch_newLocalState <- elev
		case obstruction := <-ch_obstruction:
			onObstructionEvent(obstruction, &elev)
			ch_newLocalState <- elev
		case doorClosed := <-ch_timerDone:
			doorEvent(doorClosed, elev)
		}

	}

}

func doorEvent(doorClosed bool, e ElevatorState){
	if doorClosed {
		drv.SetDoorOpenLamp(false)
		nextOrder(e)
	}else{
		drv.SetDoorOpenLamp(true)
	}
}

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, ch_doorTimer chan timer.TimerBehavior) {
	floor := order.Floor
	btn_type := order.Button
	e.Orders[floor][btn_type] = true
	switch e.Behavior {
	case c.DoorOpen:
		if shouldClearImmediatly(e, floor, btn_type) {
			time.Sleep(time.Second * c.DoorOpenDuration)
			e.Orders[floor][btn_type] = false
		}
	case c.Idle:
		drv.SetButtonLamp(btn_type, floor, true)
		direction, behavior := chooseElevDirection(e)
		e.Direction = direction
		e.Behavior = behavior
		switch e.Behavior {
		case c.DoorOpen:
			ch_doorTimer <- timer.TimerBehavior{Start:true, Duration:time.Second*c.DoorOpenDuration}
			clearAtCurrentFloor(e)
		case c.Moving:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func onFloorArrivalEvent(stop bool, floor int, e *ElevatorState, ch_doorTimer chan timer.TimerBehavior) {
	// Elevetor has arrived at floor and should clear orders at this floor
	e.Floor = floor
	drv.SetFloorIndicator(floor)
	//println(e.direction)

	if shouldStop(e) {
		fmt.Println("STOPPING")
		drv.SetMotorDirection(drv.MD_Stop)
		e.Behavior = c.Idle
		e.Direction = drv.MD_Stop
		clearAtCurrentFloor(e)
		ch_doorTimer <- timer.TimerBehavior{Start:true, Duration:time.Second*c.DoorOpenDuration}

	}

}

func onStopEvent(stop bool, e *ElevatorState, a <-chan int) {
	if e.Floor != 0 {
		drv.SetDoorOpenLamp(false)
		drv.SetMotorDirection(drv.MD_Down)
		e.Direction = drv.MD_Down
	}
}

func onObstructionEvent(obstruction bool, e *ElevatorState) {
	//TODO: IMPLEMENT
}

func nextOrder(e ElevatorState) {
	direction		drv.SetDoorOpenLamp(true)
	// Reset doortimer
	time.Sleep(3 * time.Second)
	drv.SetDoorOpenLamp(false), behavior := chooseElevDirection(&e)
	e.Direction = direction
	e.Behavior = behavior
	if direction != drv.MD_Stop {
		drv.SetMotorDirection(e.Direction)
	}
}

func printState(elev ElevatorState) {
	println("   UP  DOWN  CAB")
	fmt.Println(elev.Orders[3])
	fmt.Println(elev.Orders[2])
	fmt.Println(elev.Orders[1])
	fmt.Println(elev.Orders[0])
	println()
	//println("Direction: ", elev.direction, ", behavior: ", elev.behavior, " Floor: ", elev.floor)
	println()
}*/
