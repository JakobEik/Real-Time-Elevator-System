package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"Project/failroutine"
	"Project/watchdog"
	"fmt"
	"time"
)

const motorLossTimeDuration = time.Second * 4
const doorOpenDuration = time.Second * 3

func Fsm(
	ch_executeOrder <-chan drv.ButtonEvent,
	ch_floorArrival chan int,
	ch_obstruction <-chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- c.ElevatorState,
	ch_hallLights <-chan [][]bool,
	ch_unavailable chan<- bool,
	ch_offNetwork <-chan bool) {

	offNetwork := false

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "FSM")

	doorTimer := time.NewTimer(1)
	doorTimer.Stop()
	motorLossTimer := time.NewTimer(1)
	motorLossTimer.Stop()

	elev := InitElev(c.N_FLOORS)
	elev = clearAllFloors(elev)
	drv.SetMotorDirection(drv.MD_Down)
	elev.Direction = drv.MD_Down
	motorLossTimer.Stop()

	for {
		select {
		case value := <-ch_offNetwork:
			offNetwork = value
			println("OFF NETTWORK : ", value)
		case <-ch_bark:
			ch_pet <- true
		case hallOrders := <-ch_hallLights:
			setHallLights(hallOrders)
		case <-motorLossTimer.C:
			ch_unavailable <- true
			println("======================= MOTOR LOSS ==============================")
			failroutine.FailRoutine()

		case order := <-ch_executeOrder:
			floor := order.Floor
			btn_type := order.Button
			switch elev.Behavior {
			case c.DOOR_OPEN:
				if shouldClearImmediatly(elev, floor, btn_type) {
					elev.Orders[floor][btn_type] = false
					doorTimer.Reset(doorOpenDuration)
				} else {
					elev.Orders[floor][btn_type] = true
				}

			case c.MOVING:
				elev.Orders[floor][btn_type] = true

			case c.IDLE:
				println("3")
				if shouldClearImmediatly(elev, floor, btn_type) {
					drv.SetDoorOpenLamp(true)
					elev.Behavior = c.DOOR_OPEN
					doorTimer.Reset(doorOpenDuration)

				}
				elev.Orders[floor][btn_type] = true
				elev.Direction, elev.Behavior = chooseElevDirection(elev)
				switch elev.Behavior {
				case c.DOOR_OPEN:
					drv.SetDoorOpenLamp(true)
					doorTimer.Reset(doorOpenDuration)
					elev = clearAtCurrentFloor(elev)

				case c.MOVING:
					drv.SetMotorDirection(elev.Direction)
					setMotorLossTimer(elev.Direction, motorLossTimer)

				}
			}
			ch_newLocalState <- elev
		case floor := <-ch_floorArrival:
			//println("floor:", floor)
			motorLossTimer.Reset(motorLossTimeDuration)
			elev.Floor = floor
			drv.SetFloorIndicator(floor)
			ch_unavailable <- false
			//setCabLights(elev.Orders)

			if shouldStop(elev) {
				//fmt.Println("DOOR OPEN")
				drv.SetMotorDirection(drv.MD_Stop)
				setMotorLossTimer(drv.MD_Stop, motorLossTimer)
				elev.Behavior = c.DOOR_OPEN
				elev = clearAtCurrentFloor(elev)
				drv.SetDoorOpenLamp(true)
				doorTimer.Reset(doorOpenDuration)
				//println("set door timer")
			}
			ch_newLocalState <- elev

		case <-ch_stop:
			drv.SetMotorDirection(drv.MD_Stop)
			ch_newLocalState <- elev

		case obstruction := <-ch_obstruction:
			if obstruction && elev.Behavior == c.DOOR_OPEN {
				doorTimer.Stop()
				time.Sleep(time.Millisecond * 50)
				elev = clearAllFloors(elev)
				ch_unavailable <- true
			} else {
				doorTimer.Reset(doorOpenDuration)
				ch_unavailable <- false
			}
			ch_newLocalState <- elev

		case <-doorTimer.C:
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
		setCabLights(elev.Orders)
		// Hall lights are usually set by master.
		// If this elevator goes offline, it will set its own hall lights
		if offNetwork {
			setHallLights(elev.Orders)
		}

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

func PrintState(elev c.ElevatorState) {
	println("  UP   DOWN  CAB")
	fmt.Println(elev.Orders[3])
	fmt.Println(elev.Orders[2])
	fmt.Println(elev.Orders[1])
	fmt.Println(elev.Orders[0])
	println()
	//println("Direction: ", elev.direction, ", behavior: ", elev.behavior, " Floor: ", elev.floor)
	println()
}

func PrintGlobalState(elevs []c.ElevatorState) {
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
