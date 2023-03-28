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
	ch_newLocalState chan<- ElevatorState,
	ch_globalHallOrders <-chan [][]bool) {

	doorTimer := time.NewTimer(1)
	<-doorTimer.C
	obstruct := false

	elev := InitElev(c.N_FLOORS - 1)
	// Clears all orders and goes to first floor
	clearAllFloors(&elev)
	elev.Orders[0][2] = true
	nextOrder(&elev)
	for {
		select {
		case order := <-ch_executeOrder:
			//fmt.Println("NEW ORDER:", order)
			//println("OBSTRUCT:", obstruct)
			if !obstruct {
				//println("execute")
				onNewOrderEvent(order, &elev, doorTimer)
			}

		case floor := <-ch_floorArrival:
			//println("floor:", floor)
			elev.Floor = floor
			drv.SetFloorIndicator(floor)

			if shouldStop(&elev) {
				//fmt.Println("DOOR OPEN")
				drv.SetMotorDirection(drv.MD_Stop)
				elev.Behavior = c.DOOR_OPEN
				clearAtCurrentFloor(&elev)
				drv.SetDoorOpenLamp(true)
				doorTimer.Reset(c.DoorOpenDuration)
				//println("set door timer")

			}

		case <-ch_stop:
			//println("STOP")
			clearAllFloors(&elev)
			firstFloor := drv.ButtonEvent{Floor: 0, Button: drv.BT_Cab}
			onNewOrderEvent(firstFloor, &elev, doorTimer)

		case obstruction := <-ch_obstruction:
			obstruct = obstruction
			onObstructionEvent(obstruction, elev, doorTimer)

		case <-doorTimer.C:
			//println("DOOR CLOSE")
			drv.SetDoorOpenLamp(false)
			elev.Behavior = c.IDLE
			//println("NO MORE ORDERS?", ordersIsEmpty(e))
			if !ordersIsEmpty(&elev) {
				nextOrder(&elev)
			}
		case hallOrders := <-ch_globalHallOrders:
			setHallLights(hallOrders)
		}
		//PrintState(elev)
		ch_newLocalState <- elev
		setCabLights(elev.Orders)

	}

}

func setHallLights(buttons [][]bool) {
	for floor := 0; floor < c.N_FLOORS; floor++ {
		for btn := 0; btn < c.N_BUTTONS; btn++ {
			drv.SetButtonLamp(drv.ButtonType(btn), floor, buttons[floor][btn])
		}
	}
}

func setCabLights(orders [][]bool) {
	for floor := 0; floor < c.N_FLOORS; floor++ {
		CAB_btn := c.N_BUTTONS - 1
		drv.SetButtonLamp(drv.ButtonType(CAB_btn), floor, orders[floor][CAB_btn])
	}
}

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, doorTimer *time.Timer) {
	floor := order.Floor
	btn_type := order.Button
	switch e.Behavior {
	case c.DOOR_OPEN:
		if shouldClearImmediatly(e, floor, btn_type) {
			e.Orders[floor][btn_type] = false
			doorTimer.Reset(c.DoorOpenDuration)
		} else {
			e.Orders[floor][btn_type] = true
		}

	case c.MOVING:
		e.Orders[floor][btn_type] = true

	case c.IDLE:
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
		case c.DOOR_OPEN:
			drv.SetDoorOpenLamp(true)
			doorTimer.Reset(c.DoorOpenDuration)
			clearAtCurrentFloor(e)

		case c.MOVING:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func onObstructionEvent(obstruction bool, e ElevatorState, doorTimer *time.Timer) {
	if obstruction && e.Behavior == c.IDLE {
		//println("OBSTRUCT")
		doorTimer.Stop()
		e.Behavior = c.UNAVAILABLE
		<-doorTimer.C
	} else {
		//println("OBSTR OFF")
		switch e.Behavior {

		case c.DOOR_OPEN:
			//println("RESET")
			doorTimer.Reset(c.DoorOpenDuration)
		}
	}

}

func nextOrder(e *ElevatorState) {
	direction, behavior := chooseElevDirection(e)
	e.Direction = direction
	e.Behavior = behavior
	//println("BEHAVIOR NEXT ORDER:", behavior.String())
	//println("DIRECTION:", direction)
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
