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
	var elevatorIDsOnNetwork []int

	ch_distributeLostOrders := make(chan []string, 10)

	for {
		select {
		// MASTER
		case msg := <-ch_msgToAssigner:
			//fmt.Println("MASTER RECEIVE:", msg.Type)
			content := msg.Content
			switch msg.Type {
			case c.NEW_ORDER:
				if c.ElevatorID != c.MasterID {
					continue
				}

				var order drv.ButtonEvent
				utils.DecodeContentToStruct(content, &order)
				lowestCostElevator := calculateCost(globalState, order, elevatorIDsOnNetwork)
				packet := utils.CreateMessage(lowestCostElevator, order, c.DO_ORDER)
				ch_msgToPack <- packet
				fmt.Println("ORDER to elevator:", lowestCostElevator)

			case c.LOCAL_STATE_CHANGED:
				var state e.ElevatorState
				utils.DecodeContentToStruct(content, &state)
				elevatorID := msg.SenderID
				globalState[elevatorID] = state
				e.PrintGlobalState(globalState)
				if c.ElevatorID != c.MasterID {
					continue
				}

				packet := utils.CreateMessage(c.ToEveryone, globalState, c.UPDATE_GLOBAL_STATE)
				ch_msgToPack <- packet

				globalHallOrders := getGlobalHallOrders(globalState)
				packet2 := utils.CreateMessage(c.ToEveryone, globalHallOrders, c.GLOBAL_HALL_ORDERS)
				ch_msgToPack <- packet2

			case c.UPDATE_GLOBAL_STATE:
				if c.ElevatorID == c.MasterID {
					continue
				}
				content := msg.Content.([]interface{})
				// Iterates through the array, converts each one to ElevatorState and updates the global state
				for i, value := range content {
					var state e.ElevatorState
					utils.DecodeContentToStruct(value, &state)
					globalState[i] = state
				}
				e.PrintGlobalState(globalState)

			}

		// SLAVE
		case update := <-ch_peerUpdate:
			fmt.Println("PEER UPDATE:", update)
			// Assign new master
			elevatorIDsOnNetwork = stringArrayToIntArray(update.Peers)
			c.MasterID = getMaster(elevatorIDsOnNetwork)
			fmt.Println("MASTER ID:", c.MasterID)

			// Distribute orders from lost peers if there are any
			fmt.Println("LENGTH:", len(update.Lost), "lost:", update.Lost)
			println(len(update.Lost) > 0)
			if len(update.Lost) > 0 {
				ch_distributeLostOrders <- update.Lost
			}

		case lostPeers := <-ch_distributeLostOrders:
			if c.ElevatorID == c.MasterID {
				println("MASTER WILL DISTRIBUTE")
				IDs := stringArrayToIntArray(lostPeers)
				for _, elevatorID := range IDs {
					println("Distribute for:", elevatorID)
					distributeOrders(globalState[elevatorID], ch_msgToPack)
				}
			}

		}

	}

}

// distributeOrders sends every order from the given elevator to the master of the system
func distributeOrders(elevator e.ElevatorState, ch_msgToPack chan<- c.NetworkMessage) {
	orders := elevator.Orders
	println("DISTRIBUTE THIS ELEVATOR")
	e.PrintState(elevator)
	for floor := range orders {
		for btn := 0; btn < c.N_BUTTONS-1; btn++ {
			if orders[floor][btn] == true {
				order := drv.ButtonEvent{Floor: floor, Button: drv.ButtonType(btn)}
				msg := utils.CreateMessage(c.MasterID, order, c.NEW_ORDER)
				fmt.Println("DISTRIBUTE ORDER:", order)
				ch_msgToPack <- msg
			}
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

func calculateCost(GlobalState []e.ElevatorState, order drv.ButtonEvent, elevatorIDs []int) int {
	var lowestCostID int
	cost := 9999
	for _, elevID := range elevatorIDs {

		ElevatorCost := Cost(GlobalState[elevID], order)
		//println("ID:", index, ", COST:", ElevatorCost)
		if ElevatorCost < cost {
			cost = ElevatorCost
			lowestCostID = elevID
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

func stringArrayToIntArray(strings []string) []int {
	ints := make([]int, len(strings))
	var err error
	for i, s := range strings {
		ints[i], err = strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error:", err)
			panic(err)
		}
	}
	return ints
}
