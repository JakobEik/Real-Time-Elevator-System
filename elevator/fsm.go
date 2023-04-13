package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"fmt"
	"time"
)

func Fsm(
	ch_executeOrder <-chan drv.ButtonEvent,
	ch_floorArrival chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- ElevatorState,
	ch_globalHallOrders <-chan [][]bool,
	ch_peerTxEnable chan<- bool) {

	doorTimer := time.NewTimer(1)
	watchdogTimer := time.NewTimer(1)

	elev := InitElev(c.N_FLOORS - 1)
	clearAllFloors(&elev)
	drv.SetMotorDirection(drv.MD_Down)

	for {
		select {
		case <-watchdogTimer.C:
			//panic("ELEVATOR OUT OF SERVICE")
		case order := <-ch_executeOrder:
			//fmt.Println("NEW ORDER:", order)
			onNewOrderEvent(order, &elev, doorTimer)
			ch_newLocalState <- elev

		case floor := <-ch_floorArrival:
			watchdogTimer.Stop()
			//println("floor:", floor)
			elev.Floor = floor
			drv.SetFloorIndicator(floor)
			if shouldStop(elev) {
				//fmt.Println("DOOR OPEN")
				drv.SetMotorDirection(drv.MD_Stop)
				elev.Behavior = c.DOOR_OPEN
				clearAtCurrentFloor(&elev)
				drv.SetDoorOpenLamp(true)
				doorTimer.Reset(c.DoorOpenDuration)
				//println("set door timer")
			} else {
				watchdogTimer.Reset(c.WatchdogTimerDuration)
			}
			ch_newLocalState <- elev

		case <-ch_stop:
			//println("STOP")
			clearAllFloors(&elev)
			firstFloor := drv.ButtonEvent{Floor: 0, Button: drv.BT_Cab}
			onNewOrderEvent(firstFloor, &elev, doorTimer)
			ch_newLocalState <- elev

		case obstruction := <-ch_obstruction:
			if obstruction && elev.Behavior == c.DOOR_OPEN {
				doorTimer.Stop()
				ch_peerTxEnable <- false
			} else {
				doorTimer.Reset(c.DoorOpenDuration)
				clearAllFloors(&elev)
				ch_peerTxEnable <- true
			}
			ch_newLocalState <- elev

		case <-doorTimer.C:
			println("DOOR CLOSE")
			drv.SetDoorOpenLamp(false)
			elev.Behavior = c.IDLE
			elev.Direction, elev.Behavior = chooseElevDirection(elev)
			drv.SetMotorDirection(elev.Direction)
			if elev.Behavior == c.DOOR_OPEN {
				// If there is another hall order at his floor in a different direction but there are no other orders
				// for this elevator, this will open the door again and clear the order
				ch_floorArrival <- elev.Floor
			}

		case hallOrders := <-ch_globalHallOrders:
			setHallLights(hallOrders)
		}
		//PrintState(elev)

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
		if shouldClearImmediatly(*e, floor, btn_type) {
			e.Orders[floor][btn_type] = false
			doorTimer.Reset(c.DoorOpenDuration)
		} else {
			e.Orders[floor][btn_type] = true
		}

	case c.MOVING:
		e.Orders[floor][btn_type] = true

	case c.IDLE:
		if shouldClearImmediatly(*e, floor, btn_type) {
			drv.SetDoorOpenLamp(true)
			e.Behavior = c.DOOR_OPEN
			doorTimer.Reset(c.DoorOpenDuration)
			return
		}
		e.Orders[floor][btn_type] = true
		e.Direction, e.Behavior = chooseElevDirection(*e)
		switch e.Behavior {
		case c.DOOR_OPEN:
			drv.SetDoorOpenLamp(true)
			doorTimer.Reset(c.DoorOpenDuration)
			clearAtCurrentFloor(e)

		case c.MOVING:
			drv.SetMotorDirection(e.Direction)
		}
	}

}

func PrintState(elev ElevatorState) {
	println("  UP   DOWN  CAB")
	fmt.Println(elev.Orders[3])
	fmt.Println(elev.Orders[2])
	fmt.Println(elev.Orders[1])
	fmt.Println(elev.Orders[0])
	println()
	//println("Direction: ", elev.direction, ", behavior: ", elev.behavior, " Floor: ", elev.floor)
	println()
}

func PrintGlobalState(elevs []ElevatorState) {
	for i, elev := range elevs {
		print("  UP   DOWN  CAB")
		println("ELEVATOR:", i)
		fmt.Println(elev.Orders[3])
		fmt.Println(elev.Orders[2])
		fmt.Println(elev.Orders[1])
		fmt.Println(elev.Orders[0])
	}
	println()

}
