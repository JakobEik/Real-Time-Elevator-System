//package state_machine

import (
	"fmt"
	"os"
	"time"
)

package main

import (
"fmt"
"math"
"os"

"./elevio"
"./fsm"
"./requests"
"./timer"
)

const (
	DoorOpenDuration = 3 * time.Second
)

var (
	elev Elevator
	outputDevice elevio.ElevatorIO
)
func fsm_onRequestButtonPress(btn_floor int, btn_type elevio.ButtonType{
func main() {
	var err error

	outputDevice, err = elevio.Initialize("localhost:15657", elevio.NumFloors)
	if err != nil {
		fmt.Println("failed to initialize elevator I/O:", err)
		os.Exit(1)
	}

	elev = elevatorUninitialized()
	con_load("elevator.con",
		con_val("doorOpenDuration_s", &elev.config.doorOpenDuration_s, "%lf"),
		con_enum("clearRequestVariant", &elev.config.clearRequestVariant,
			con_match(CVAll),
			con_match(CVInDirn),
		),
	)

	for {
		select {
		case btnPress := <-outputDevice.ButtonPressCh():
			btnFloor, btnType := btnPress.Floor, btnPress.ButtonType
			fmt.Printf("%s(%d, %s)\n", fsmOnRequestButtonPress, btnFloor, btnType)
			elevatorPrint(elev)

			switch elev.behaviour {
			case EBDoorOpen:
				if requests.ShouldClearImmediately(elev, btnFloor, btnType) {
					timer.Start(DoorOpenDuration)
				} else {
					elev.requests[btnFloor][btnType] = true
				}

			case EBMoving:
				elev.requests[btnFloor][btnType] = true

			case EBIdle:
				elev.requests[btnFloor][btnType] = true
				pair := requests.ChooseDirection(elev)
				elev.dirn, elev.behaviour = pair.dirn, pair.behaviour

				switch elev.behaviour {
				case EBDoorOpen:
					outputDevice.DoorLight(true)
					timer.Start(DoorOpenDuration)
					elev = requests.ClearAtCurrentFloor(elev)
				case EBMoving:
					outputDevice.MotorDirection(elev.dirn)
				case EBIdle:
				}
			}

			setAllLights(elev)

			fmt.Println("\nNew state:")
			elevatorPrint(elev)

		case floor := <-outputDevice.FloorSensorCh():
			fmt.Printf("%s(%d)\n", fsmOnFloorArrival, floor)
			elevatorPrint(elev)

			elev.floor = floor
			outputDevice.FloorIndicator(elev.floor)

			switch elev.behaviour {
			case EBMoving:
				if requests.ShouldStop(elev) {
					outputDevice.MotorDirection(elevio.MDStop)
					outputDevice.DoorLight(true)
					elev = requests.ClearAtCurrentFloor(elev)
					timer.Start(DoorOpenDuration)
					setAllLights(elev)
					elev.behaviour = EBDoorOpen
				}

			default:
			}

			fmt.Println("\nNew state:")
			elevatorPrint(elev)

		case <-timer.TimeoutCh():
			fmt.Printf("%s()\n", fsmOnDoorTimeout)
			elevatorPrint(elev)

			switch elev.behaviour {
			case EBDoorOpen:
				pair := requests.ChooseDirection(elev)
				elev.dirn, elev.behaviour = pair.dirn, pair.behaviour

				switch elev.behaviour {
				case EBDoorOpen:
					timer.Start(DoorOpenDuration)
					elev = requests.ClearAtCurrentFloor(elev)
					setAllLights(elev)
				case EBMoving, EBIdle:
					outputDevice.DoorLight(false)
					outputDevice.MotorDirection(elev.dirn)
				}

			default:
