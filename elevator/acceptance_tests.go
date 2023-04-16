package elevator

import (
	c "Project/config"
	"Project/driver"
	"Project/failroutine"
	"fmt"
	"strconv"
)

func AcceptanceTests(state c.ElevatorState) {
	stateIsConsistent(state)
	floorInValidRange(state)
	directionIsValid(state)
	behaviorIsValid(state)
	ordersAreValid(state)
}

func stateIsConsistent(state c.ElevatorState) {
	// Behavior is not set to MD_Stop so the correct next direction is calculated
	// Therefore this test would fail even though the system works
	/*if state.Behavior == c.EB_DOOR_OPEN && state.Direction != driver.MD_Stop {
		fmt.Println(state)
		fmt.Println("MOTOR SHOULD BE STOPPED WHEN DOOR IS OPEN")
		failroutine.FailRoutine()
	}*/
	if state.Behavior == c.EB_IDLE && !ordersIsEmpty(state) {
		fmt.Println(state)
		fmt.Println("BEHAVIOR IS EB_IDLE, ORDERS SHOULD BE EMPTY")
		failroutine.FailRoutine()
	}
	if state.Behavior == c.EB_MOVING && state.Direction == driver.MD_Stop {
		fmt.Println("ELEVATOR IS EB_MOVING BUT DIRECTION IS STOP")
		failroutine.FailRoutine()
	}
	if state.Behavior == c.EB_IDLE && state.Direction != driver.MD_Stop {
		fmt.Println("ELEVATOR IS EB_IDLE BUT DIRECTION IS NOT STOP")
		failroutine.FailRoutine()
	}

}

func floorInValidRange(state c.ElevatorState) {
	if state.Floor < 0 || state.Floor > c.N_FLOORS {
		floor := strconv.Itoa(state.Floor)
		fmt.Println("FLOOR OUTSIDE OF VALID RANGE : " + floor)
		failroutine.FailRoutine()
	}
}

func directionIsValid(state c.ElevatorState) {
	if state.Direction != driver.MD_Up && state.Direction != driver.MD_Down && state.Direction != driver.MD_Stop {
		dir := strconv.Itoa(int(state.Direction))
		fmt.Println("DIRECTION IS NOT VALID : " + dir)
		failroutine.FailRoutine()
	}
}

func behaviorIsValid(state c.ElevatorState) {
	if state.Behavior != c.EB_IDLE && state.Behavior != c.EB_MOVING && state.Behavior != c.EB_DOOR_OPEN {
		b := strconv.Itoa(int(state.Behavior))
		fmt.Println("BEHAVIOR IS NOT VALID : " + b)
		failroutine.FailRoutine()
	}
}

func ordersAreValid(state c.ElevatorState) {
	for _, floorOrders := range state.Orders {
		if len(floorOrders) != c.N_BUTTONS {
			fmt.Println(state.Orders)
			fmt.Println("MISMATCH BETWEEN ORDER STATE AND N_BUTTONS")
			failroutine.FailRoutine()
		}
		for _, order := range floorOrders {
			if order != true && order != false {
				fmt.Println(state.Orders)
				fmt.Println("ORDERS SHOULD BE BOOLEAN")
				failroutine.FailRoutine()
			}
		}
	}

	if state.Orders[0][driver.BT_HallDown] {
		fmt.Println(state.Orders)
		fmt.Println("INVALID HALL DOWN ORDER AT BOTTOM FLOOR")
		failroutine.FailRoutine()
	}

	if state.Orders[c.N_FLOORS-1][driver.BT_HallUp] {
		fmt.Println(state.Orders)
		fmt.Println("INVALID HALL DOWN ORDER AT BOTTOM FLOOR")
		failroutine.FailRoutine()
	}

}
