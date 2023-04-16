package elevator

import (
	c "Project/config"
	drv "Project/driver"
	"Project/failroutine"
	"Project/watchdog"
	"time"
)

const motorLossTimeDuration = time.Second * 4
const doorOpenDuration = time.Second * 3

func Fsm(
	ch_executeOrder <-chan drv.ButtonEvent,
	ch_floorArrival chan int,
	ch_obstruction chan bool,
	ch_stop <-chan bool,
	ch_newLocalState chan<- c.ElevatorState,
	ch_hallLights <-chan [][]bool,
	ch_unavailable chan<- bool,
	ch_offNetwork <-chan bool) {

	offNetwork := false
	isObstructed := false

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Fsm")

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
		case hallLights := <-ch_hallLights:
			setHallLights(hallLights)
		//______________EXECUTE ORDER EVENT______________
		case order := <-ch_executeOrder:
			floor := order.Floor
			btn_type := order.Button
			switch elev.Behavior {
			case c.EB_DOOR_OPEN:
				if shouldClearImmediatly(elev, floor, btn_type) {
					elev.Orders[floor][btn_type] = false
					doorTimer.Reset(doorOpenDuration)
				} else {
					elev.Orders[floor][btn_type] = true
				}

			case c.EB_MOVING:
				elev.Orders[floor][btn_type] = true

			case c.EB_IDLE:
				if shouldClearImmediatly(elev, floor, btn_type) {
					drv.SetDoorOpenLamp(true)
					elev.Behavior = c.EB_DOOR_OPEN
					doorTimer.Reset(doorOpenDuration)
				}
				elev.Orders[floor][btn_type] = true
				elev.Direction, elev.Behavior = chooseElevDirection(elev)
				switch elev.Behavior {
				case c.EB_DOOR_OPEN:
					drv.SetDoorOpenLamp(true)
					doorTimer.Reset(doorOpenDuration)
					elev = clearAtCurrentFloor(elev)

				case c.EB_MOVING:
					drv.SetMotorDirection(elev.Direction)
					setMotorLossTimer(elev.Direction, motorLossTimer)

				}
			}
			ch_newLocalState <- elev

		//__________________FLOOR EVENT___________________
		case floor := <-ch_floorArrival:
			motorLossTimer.Reset(motorLossTimeDuration)
			elev.Floor = floor
			drv.SetFloorIndicator(floor)
			ch_unavailable <- false

			if shouldStop(elev) {
				drv.SetMotorDirection(drv.MD_Stop)
				setMotorLossTimer(drv.MD_Stop, motorLossTimer)
				elev.Behavior = c.EB_DOOR_OPEN
				elev = clearAtCurrentFloor(elev)
				drv.SetDoorOpenLamp(true)
				// If sleep is not here, sometimes the door lamp on the physical elevator wont turn on?????
				time.Sleep(time.Millisecond)
				doorTimer.Reset(doorOpenDuration)
				ch_obstruction <- isObstructed
			}
			ch_newLocalState <- elev

		//________________STOP EVENT________________
		case <-ch_stop:
			drv.SetMotorDirection(drv.MD_Stop)
			ch_newLocalState <- elev

		//________________OBSTRUCTION EVENT________________
		case obstruction := <-ch_obstruction:
			isObstructed = obstruction
			if obstruction && elev.Behavior == c.EB_DOOR_OPEN {
				doorTimer.Stop()
				elev = clearAllFloors(elev)
				time.Sleep(time.Millisecond * 50)
				ch_unavailable <- true
			} else {
				doorTimer.Reset(doorOpenDuration)
				ch_unavailable <- false
			}
			ch_newLocalState <- elev

		//_______________DOOR CLOSING EVENT______________
		case <-doorTimer.C:
			drv.SetDoorOpenLamp(false)
			elev.Behavior = c.EB_IDLE
			// Next Order
			elev.Direction, elev.Behavior = chooseElevDirection(elev)
			drv.SetMotorDirection(elev.Direction)
			setMotorLossTimer(elev.Direction, motorLossTimer)
			if elev.Behavior == c.EB_DOOR_OPEN {
				// If there is another hall order at his floor in a different direction but there are no other orders
				// for this elevator, this will open the door again and clear the order
				ch_floorArrival <- elev.Floor
			}
			ch_newLocalState <- elev

		//_________________MONITORING________________
		case value := <-ch_offNetwork:
			offNetwork = value
		case <-ch_bark:
			ch_pet <- true
		case <-motorLossTimer.C:
			ch_unavailable <- true
			println("======================= MOTOR LOSS ==============================")
			failroutine.FailRoutine()

		}

		//______THIS RUNS EVERY LOOP_______
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
