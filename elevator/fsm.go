package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"Project/watchdog"
	"fmt"
	"time"
)

const motorLossTimeDuration = time.Second * 5
const doorOpenDuration = time.Second * 3

func Fsm(
	ch_executeOrder <-chan drv.ButtonEvent,
	ch_floorArrival chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- ElevatorState,
	ch_globalHallOrders <-chan [][]bool,
	ch_failure chan<- bool) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "FSM")

	doorTimer := time.NewTimer(1)
	doorTimer.Stop()
	motorLossTimer := time.NewTimer(1)
	motorLossTimer.Stop()

	elev := InitElev(c.N_FLOORS - 1)
	clearAllFloors(&elev)
	drv.SetMotorDirection(drv.MD_Down)
	elev.Direction = drv.MD_Down
	motorLossTimer.Stop()

	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true
		case hallOrders := <-ch_globalHallOrders:
			setHallLights(hallOrders)

		case <-motorLossTimer.C:
			ch_failure <- true
			println("======================= MOTOR LOSS ==============================")
		case order := <-ch_executeOrder:
			//fmt.Println("NEW ORDER:", order)
			onNewOrderEvent(order, &elev, doorTimer, motorLossTimer)
			ch_newLocalState <- elev
		case floor := <-ch_floorArrival:
			//println("floor:", floor)
			motorLossTimer.Reset(motorLossTimeDuration)
			elev.Floor = floor
			drv.SetFloorIndicator(floor)
			ch_failure <- false
			setCabLights(elev.Orders)
			setCabLights(elev.Orders)

			if shouldStop(elev) {
				//fmt.Println("DOOR OPEN")
				drv.SetMotorDirection(drv.MD_Stop)
				elev.Direction = drv.MD_Stop
				setMotorLossTimer(drv.MD_Stop, motorLossTimer)
				elev.Behavior = c.DOOR_OPEN
				clearAtCurrentFloor(&elev)
				drv.SetDoorOpenLamp(true)
				doorTimer.Reset(doorOpenDuration)
				//println("set door timer")
			}
			ch_newLocalState <- elev

		case <-ch_stop:
			//println("STOP")
			clearAllFloors(&elev)
			firstFloor := drv.ButtonEvent{Floor: 0, Button: drv.BT_Cab}
			onNewOrderEvent(firstFloor, &elev, doorTimer, motorLossTimer)
			ch_newLocalState <- elev

		case obstruction := <-ch_obstruction:
			if obstruction && elev.Behavior == c.DOOR_OPEN {
				doorTimer.Stop()
				time.Sleep(time.Millisecond * 50)
				clearAllFloors(&elev)
				ch_failure <- true
			} else {
				doorTimer.Reset(doorOpenDuration)
				ch_failure <- false
			}
			ch_newLocalState <- elev

		case <-doorTimer.C:
			//println("DOOR TIMER")
			drv.SetDoorOpenLamp(false)
			elev.Behavior = c.IDLE
			// Next Order
			elev.Direction, elev.Behavior = chooseElevDirection(elev)
			drv.SetMotorDirection(elev.Direction)
			setMotorLossTimer(elev.Direction, motorLossTimer)
			if elev.Behavior == c.DOOR_OPEN {
				// If there is another hall order at his floor in a different direction but there are no other orders
				// for this elevator, this will open the door again and clear the order
				ch_floorArrival <- elev.Floor
			}
			ch_newLocalState <- elev
		}
		//PrintState(elev)
		setCabLights(elev.Orders)

	}

}

func setMotorLossTimer(dir drv.MotorDirection, timer *time.Timer) {
	if dir == drv.MD_Stop {
		timer.Stop()
	} else {
		timer.Reset(motorLossTimeDuration)
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

func onNewOrderEvent(order drv.ButtonEvent, e *ElevatorState, doorTimer *time.Timer, motorLossTimer *time.Timer) {
	floor := order.Floor
	btn_type := order.Button
	switch e.Behavior {
	case c.DOOR_OPEN:
		if shouldClearImmediatly(*e, floor, btn_type) {
			e.Orders[floor][btn_type] = false
			doorTimer.Reset(doorOpenDuration)
		} else {
			e.Orders[floor][btn_type] = true
		}

	case c.MOVING:
		e.Orders[floor][btn_type] = true

	case c.IDLE:
		if shouldClearImmediatly(*e, floor, btn_type) {
			drv.SetDoorOpenLamp(true)
			e.Behavior = c.DOOR_OPEN
			doorTimer.Reset(doorOpenDuration)
			return
		}
		e.Orders[floor][btn_type] = true
		e.Direction, e.Behavior = chooseElevDirection(*e)
		switch e.Behavior {
		case c.DOOR_OPEN:
			drv.SetDoorOpenLamp(true)
			doorTimer.Reset(doorOpenDuration)
			clearAtCurrentFloor(e)

		case c.MOVING:
			drv.SetMotorDirection(e.Direction)
			setMotorLossTimer(e.Direction, motorLossTimer)

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
