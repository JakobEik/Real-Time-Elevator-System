package distributor

import (
	c "Project/config"
	"Project/driver"
	"Project/elevator"
)

func motorPower(e elevator.ElevatorState) {
	if e.Behavior == c.IDLE && driver.GetFloor() == -1 {
		panic("NO MOTOR POWER")
	}
}
