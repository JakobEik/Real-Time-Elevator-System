package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	"Project/utils"
	"fmt"
	"strconv"
	//drv "Project/driver"
)

func MasterNode(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage) {

	globalState := utils.InitGlobalState()
	println("LENGTH GLOBAL STATE:", len(globalState))
	var elevatorIDs []int

	for {
		select {
		case msg := <-ch_msgToAssigner:
			if c.ElevatorID != c.MasterID {
				continue
			}
			content := msg.Content
			switch msg.Type {
			case c.NEW_ORDER:
				var order drv.ButtonEvent
				utils.DecodeContentToStruct(content, &order)
				lowestCostElevator := calculateCost(globalState, order)
				packet := utils.CreateMessage(lowestCostElevator, order, c.DO_ORDER)
				ch_msgToPack <- packet
				fmt.Println("ORDER to elevator:", lowestCostElevator)

			case c.LOCAL_STATE_CHANGED:
				var state e.ElevatorState
				utils.DecodeContentToStruct(content, &state)
				elevatorID := msg.SenderID
				globalState[elevatorID] = state
				packetState := utils.CreateMessage(c.ToEveryone, globalState, c.UPDATE_GLOBAL_STATE)
				ch_msgToPack <- packetState

				globalHallOrders := getGlobalHallOrders(globalState)
				packetHall := utils.CreateMessage(c.ToEveryone, globalHallOrders, c.GLOBAL_HALL_ORDERS)
				ch_msgToPack <- packetHall
			}
		case update := <-ch_peerUpdate:
			fmt.Println("PEER UPDATE:", update.Peers)
			elevatorIDs = stringArrToIntArr(update.Peers)
			c.MasterID = getMaster(elevatorIDs)
			fmt.Println("MASTER ID:", c.MasterID)
		}

	}

}

// getMaster returns the elevator with the lowest ID
func getMaster(elevatorIDs []int) int {
	if len(elevatorIDs) == 0 {
		panic("ERROR")
	}
	masterID := elevatorIDs[0]
	for _, elev := range elevatorIDs[1:] {
		if elev < masterID {
			masterID = elev
		}
	}
	return masterID

}

func calculateCost(GlobalState []e.ElevatorState, order drv.ButtonEvent) int {
	var lowestCostID int
	cost := 9999

	for index, localState := range GlobalState {

		ElevatorCost := Cost(localState, order)
		//println("ID:", index, ", COST:", ElevatorCost)
		if ElevatorCost < cost {

			cost = ElevatorCost
			lowestCostID = index
		}
	}

	return lowestCostID
}

func getGlobalHallOrders(globalState []e.ElevatorState) [][]bool {
	buttons := e.InitElev(0).Orders
	for _, state := range globalState {
		for floor := 0; floor < c.N_FLOORS; floor++ {
			for btn := 0; btn < c.N_BUTTONS-1; btn++ {
				if state.Orders[floor][btn] == true {
					buttons[floor][btn] = true
				}
			}
		}
	}
	return buttons
}

func stringArrToIntArr(strings []string) []int {
	ints := make([]int, len(strings))
	var err error
	for i, s := range strings {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error:", err)
			return nil
		}
	}
	return ints
}
