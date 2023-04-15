package elevator

import (
	c "Project/config"
	"Project/driver"
	"Project/failroutine"
	"fmt"
	"strconv"
)

func AcceptanceTests(state ElevatorState) {
	stateIsConsistent(state)
	floorInValidRange(state)
	directionIsValid(state)
	behaviorIsValid(state)
	ordersAreValid(state)
}

func stateIsConsistent(state ElevatorState) {
	if state.Behavior == c.DOOR_OPEN && state.Direction != driver.MD_Stop {
		fmt.Println(state)
		fmt.Println("MOTOR SHOULD BE STOPPED WHEN DOOR IS OPEN")
		failroutine.FailRoutine()
	}
	if state.Behavior == c.IDLE && !ordersIsEmpty(state) {
		fmt.Println(state)
		fmt.Println("BEHAVIOR IS IDLE, ORDERS SHOULD BE EMPTY")
		failroutine.FailRoutine()
	}
	if state.Behavior == c.MOVING && state.Direction == driver.MD_Stop {
		fmt.Println("ELEVATOR IS MOVING BUT DIRECTION IS STOP")
		failroutine.FailRoutine()
	}
	if state.Behavior == c.IDLE && state.Direction != driver.MD_Stop {
		fmt.Println("ELEVATOR IS IDLE BUT DIRECTION IS NOT STOP")
		failroutine.FailRoutine()
	}

	if state.Behavior == c.DOOR_OPEN && state.Direction != driver.MD_Stop {
		fmt.Println("DOOR IS OPEN BUT DIRECTION IS NOT STOP")
		failroutine.FailRoutine()
	}

}

func floorInValidRange(state ElevatorState) {
	if state.Floor < 0 || state.Floor > c.N_FLOORS {
		floor := strconv.Itoa(state.Floor)
		fmt.Println("FLOOR OUTSIDE OF VALID RANGE : " + floor)
		failroutine.FailRoutine()
	}
}

func directionIsValid(state ElevatorState) {
	if state.Direction != driver.MD_Up && state.Direction != driver.MD_Down && state.Direction != driver.MD_Stop {
		dir := strconv.Itoa(int(state.Direction))
		fmt.Println("DIRECTION IS NOT VALID : " + dir)
		failroutine.FailRoutine()
	}
}

func behaviorIsValid(state ElevatorState) {
	if state.Behavior != c.IDLE && state.Behavior != c.MOVING && state.Behavior != c.DOOR_OPEN {
		b := strconv.Itoa(int(state.Behavior))
		fmt.Println("BEHAVIOR IS NOT VALID : " + b)
		failroutine.FailRoutine()
	}
}

func ordersAreValid(state ElevatorState) {
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
