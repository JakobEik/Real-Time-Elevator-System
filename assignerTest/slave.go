package assignerTest

import (
	c "Project/config"
	e "Project/elevator"
	"Project/utils"
)

type slaveNode struct {
	globalStateBackup []e.ElevatorState
}

func (s slaveNode) globalStateUpdate(msg c.NetworkMessage) {
	content := msg.Content.([]interface{})
	// Iterates through the array, converts each one to ElevatorState and updates the global state
	for i, value := range content {
		var state e.ElevatorState
		utils.DecodeContentToStruct(value, &state)
		s.globalStateBackup[i] = state
	}
}

// getUpdatedMaster returns the elevator with the lowest ID
func getUpdatedMaster(peersOnlineStr []string) int {
	peersOnline := stringArrayToIntArray(peersOnlineStr)
	return getMaster(peersOnline)

}

// getMaster returns the elevator with the lowest ID
func getMaster(peersOnline []int) int {
	if len(peersOnline) == 0 {
		panic("NO MORE ELEVATORS ON NETWORK")
	}
	masterID := peersOnline[0]
	for _, elev := range peersOnline[1:] {
		if elev < masterID {
			masterID = elev
		}
	}
	return masterID
}
