package state_machine

import (
	"./elevator_io"
	"./requests"
	"fmt"
	"project-group-77/elevator"
	"time"
)

func main() {
	// initialize elevator and output device

	var e elevator.Elevator = elevator.UninitializedElevator()
	outputDevice := elevator_io.ElevioGetOutputDevice()

	// load config from file
	config, err := LoadConfig("elevator.con")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	e.Config = config

	// set initial motor direction and behavior
	outputDevice.MotorDirection(elevator_io.MdDown)
	e.Dirn = elevator_io.MdDown
	e.Behaviour = elevator.EbMoving

	// continuously listen for elevator events
	for {
		select {
		case buttonPress := <-elevator_io.ButtonPressCh:
			handleButtonPress(buttonPress, &e, outputDevice)
		case floor := <-elevator_io.FloorSensorCh:
			handleFloorArrival(floor, &e, outputDevice)
		case <-timer.TimeoutCh:
			handleDoorTimeout(&e, outputDevice)
		default:
			handleIdle(&e, outputDevice)
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// LoadConfig loads the elevator configuration from a file and returns it.
func LoadConfig(filename string) (elevator.Config, error) {
	config := elevator.Config{}
	err := LoadConfigFromFile(filename, &config.DoorOpenDuration_s, &config.ClearRequestVariant)
	return config, err
}

// LoadConfigFromFile loads the elevator configuration from a file and sets the given variables.
func LoadConfigFromFile(filename string, doorOpenDuration_s *float64, clearRequestVariant *requests.ClearRequestVariant) error {
	values, err := con_load.ReadFile(filename)
	if err != nil {
		return err
	}

	for key, val := range values {
		switch key {
		case "doorOpenDuration_s":
			*doorOpenDuration_s = val.(float64)
		case "clearRequestVariant":
			*clearRequestVariant = val.(requests.ClearRequestVariant)
		}
	}

	return nil
}

// handleButtonPress handles an elevator button press event.
func handleButtonPress(buttonPress elevator_io.ButtonPress, e *elevator.Elevator, outputDevice elevator_io.ElevOutputDevice) {
	btnFloor := buttonPress.Floor
	btnType := buttonPress.Button

	fmt.Printf("\n\n%s(%d, %s)\n", "handleButtonPress", btnFloor, btnType.String())
	elevator.Print(*e)

	switch e.Behaviour {
	case elevator.EbDoorOpen:
		if requests.ShouldClearImmediately(*e, btnFloor, btnType) {
			timer.Start(e.Config.DoorOpenDuration_s)
		} else {
			e.Requests[btnFloor][btnType] = true
		}
	case elevator.EbMoving:
		e.Requests[btnFloor][btnType] = true
	case elevator.EbIdle:
		e.Requests[btnFloor][btnType] = true
		pair := requests.ChooseDirection(*e)
		e.Dirn = pair.Dirn
		e.Behaviour = pair.Behaviour
		switch e.Behaviour {
		case elevator.EbDoorOpen:
			outputDevice.DoorLight(true)
			timer.Start(e.Config.DoorOpenDuration_s)
			*e = requests.ClearAtCurrentFloor(*e)
			outputDevice.ClearRequestLights(*e)
		case elevator.EbMoving:
			outputDevice.MotorDirection(e.Dirn)
		}
	}

	outputDevice.SetRequestButtonLights(*e)
}
