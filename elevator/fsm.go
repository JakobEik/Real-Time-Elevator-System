package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"fmt"
	"time"
)

func Fsm(
	ch_executeOrder <-chan drv.ButtonEvent,
	ch_floorArrival <-chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- ElevatorState) {

	doorTimer := time.NewTimer(1)
	//<-doorTimer.C
	obstruct := false

	elev := InitElev(c.N_FLOORS - 1)
	// Clears all orders and goes to first floor
	clearAllFloors(&elev)
	elev.Orders[0][2] = true
	nextOrder(&elev)
	for {
		select {
		case order := <-ch_executeOrder:
			fmt.Println("NEW ORDER:", order)
			if !obstruct {
				onNewOrderEvent(order, &elev, doorTimer)
			}

		case floor := <-ch_floorArrival:
			//println("floor:", floor)
			onFloorArrivalEvent(floor, &elev, doorTimer)
		case <-ch_stop:
			//println("STOP")
			onStopEvent(&elev, doorTimer)
		case obstruction := <-ch_obstruction:
			obstruct = obstruction
			onObstructionEvent(obstruction, elev, doorTimer)
		case <-doorTimer.C:
			doorCloseEvent(&elev)
		}
		setAllLights(elev)
		//PrintState(elev)
		ch_newLocalState <- elev

	}

}

func doorCloseEvent(e *ElevatorState) {
	//println("DOOR CLOSE")
	drv.SetDoorOpenLamp(false)
	e.Behavior = c.Idle
	//println("&&&&&&&&&&&&&&&&& BEHAVIOR:", e.Behavior)
	//println(ordersIsEmpty(e))
	if !ordersIsEmpty(e) {
		nextOrder(e)
	}
}

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, doorTimer *time.Timer) {
	floor := order.Floor
	btn_type := order.Button
	switch e.Behavior {
	case c.DoorOpen:
		if shouldClearImmediatly(e, floor, btn_type) {
			e.Orders[floor][btn_type] = false
			doorTimer.Reset(c.DoorOpenDuration)
		} else {
			e.Orders[floor][btn_type] = true
		}

	case c.Moving:
		e.Orders[floor][btn_type] = true

	case c.Idle:
		if shouldClearImmediatly(e, floor, btn_type) {
			drv.SetDoorOpenLamp(true)
			doorTimer.Reset(c.DoorOpenDuration)
			return
		}
		e.Orders[floor][btn_type] = true
		direction, behavior := chooseElevDirection(e)
		e.Direction = direction
		e.Behavior = behavior
		switch behavior {
		case c.DoorOpen:
			drv.SetDoorOpenLamp(true)
			doorTimer.Reset(c.DoorOpenDuration)
			clearAtCurrentFloor(e)

		case c.Moving:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func setAllLights(e ElevatorState) {
	for floor := 0; floor < c.N_FLOORS; floor++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			drv.SetButtonLamp(drv.ButtonType(btn), floor, e.Orders[floor][btn])
		}
	}
}

func nNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, doorTimer *time.Timer) {
	floor := order.Floor
	btn_type := order.Button
	e.Orders[floor][btn_type] = true
	if shouldClearImmediatly(e, floor, btn_type) {
		doorTimer.Reset(c.DoorOpenDuration)
		e.Orders[floor][btn_type] = false
		drv.SetButtonLamp(btn_type, floor, false)

		return
	}

	drv.SetButtonLamp(btn_type, floor, true)
	direction, behavior := chooseElevDirection(e)
	e.Direction = direction
	e.Behavior = behavior

	switch e.Behavior {
	case c.DoorOpen:
		drv.SetButtonLamp(btn_type, floor, false)
		drv.SetDoorOpenLamp(true)
		clearAtCurrentFloor(e)
		doorTimer.Reset(c.DoorOpenDuration)

	case c.Moving:
		drv.SetMotorDirection(e.Direction)
	}

}

func onFloorArrivalEvent(floor int, e *ElevatorState, doorTimer *time.Timer) {
	// Elevator has arrived at floor and should clear orders at this floor

	e.Floor = floor
	drv.SetFloorIndicator(floor)

	if shouldStop(e) {
		//fmt.Println("DOOR OPEN")
		drv.SetMotorDirection(drv.MD_Stop)
		e.Behavior = c.DoorOpen
		clearAtCurrentFloor(e)
		drv.SetDoorOpenLamp(true)
		doorTimer.Reset(c.DoorOpenDuration)

	}
}

func onStopEvent(e *ElevatorState, doorTimer *time.Timer) {
	clearAllFloors(e)
	firstFloor := drv.ButtonEvent{Floor: 0, Button: drv.BT_Cab}
	onNewOrderEvent(firstFloor, e, doorTimer)
}

func onObstructionEvent(obstruction bool, e ElevatorState, doorTimer *time.Timer) {
	if obstruction {
		//println("OBSTRUCT")
		doorTimer.Stop()
		e.Behavior = c.Unavailable
		//<-doorTimer.C
	} else {
		//println("OBSTR OFF")
		switch e.Behavior {

		case c.DoorOpen:
			//println("RESET")
			doorTimer.Reset(c.DoorOpenDuration)
		}
	}

}

func nextOrder(e *ElevatorState) {
	direction, behavior := chooseElevDirection(e)
	e.Direction = direction
	e.Behavior = behavior
	//println("NEHAVIOR NEXT ORDER:", behavior)
	if direction != drv.MD_Stop {
		drv.SetMotorDirection(e.Direction)
	}
}

func PrintState(elev ElevatorState) {
	println("   UP  DOWN  CAB")
	fmt.Println(elev.Orders[3])
	fmt.Println(elev.Orders[2])
	fmt.Println(elev.Orders[1])
	fmt.Println(elev.Orders[0])
	println()
	//println("Direction: ", elev.direction, ", behavior: ", elev.behavior, " Floor: ", elev.floor)
	println()
}
