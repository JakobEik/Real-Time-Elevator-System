package elevator

import (
	c "Project/config"
	"Project/driver"
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
		panic("MOTOR SHOULD BE STOPPED WHEN DOOR IS OPEN")
	}
	if state.Behavior == c.IDLE && !ordersIsEmpty(state) {
		fmt.Println(state)
		panic("BEHAVIOR IS IDLE, ORDERS SHOULD BE EMPTY")
	}
	if state.Behavior == c.MOVING && state.Direction == driver.MD_Stop {
		panic("ELEVATOR IS MOVING BUT DIRECTION IS STOP")
	}
	if state.Behavior == c.IDLE && state.Direction != driver.MD_Stop {
		panic("ELEVATOR IS IDLE BUT DIRECTION IS NOT STOP")
	}

	if state.Behavior == c.DOOR_OPEN && state.Direction != driver.MD_Stop {
		panic("DOOR IS OPEN BUT DIRECTION IS NOT STOP")
	}
}

func floorInValidRange(state ElevatorState) {
	if state.Floor < 0 || state.Floor > c.N_FLOORS {
		floor := strconv.Itoa(state.Floor)
		panic("FLOOR OUTSIDE OF VALID RANGE : " + floor)
	}
}

func directionIsValid(state ElevatorState) {
	if state.Direction != driver.MD_Up && state.Direction != driver.MD_Down && state.Direction != driver.MD_Stop {
		dir := strconv.Itoa(int(state.Direction))
		panic("DIRECTION IS NOT VALID : " + dir)
	}
}

func behaviorIsValid(state ElevatorState) {
	if state.Behavior != c.IDLE && state.Behavior != c.MOVING && state.Behavior != c.DOOR_OPEN {
		b := strconv.Itoa(int(state.Behavior))
		panic("BEHAVIOR IS NOT VALID : " + b)
	}
}

func ordersAreValid(state ElevatorState) {
	for _, floorOrders := range state.Orders {
		if len(floorOrders) != c.N_BUTTONS {
			fmt.Println(state.Orders)
			panic("MISMATCH BETWEEN ORDER STATE AND N_BUTTONS")
		}
		for _, order := range floorOrders {
			if order != true && order != false {
				fmt.Println(state.Orders)
				panic("ORDERS SHOULD BE BOOLEAN")
			}
		}
	}

	if state.Orders[0][driver.BT_HallDown] {
		fmt.Println(state.Orders)
		panic("INVALID HALL DOWN ORDER AT BOTTOM FLOOR")
	}

	if state.Orders[c.N_FLOORS-1][driver.BT_HallUp] {
		fmt.Println(state.Orders)
		panic("INVALID HALL DOWN ORDER AT BOTTOM FLOOR")
	}

}
