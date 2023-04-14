package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	"Project/utils"
	"Project/watchdog"
	"fmt"
	"strconv"
	//drv "Project/driver"
)

func MasterNode(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Assigner")

	globalState := utils.InitGlobalState()
	println("LENGTH GLOBAL STATE:", len(globalState))
	var peersOnline []int
	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true
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
				lowestCostElevator := calculateCost(globalState, order, peersOnline)
				packet := utils.CreateMessage(lowestCostElevator, order, c.DO_ORDER)
				ch_msgToPack <- packet
				fmt.Println("ORDER to elevator:", lowestCostElevator)

			case c.LOCAL_STATE_CHANGED:
				var state e.ElevatorState
				utils.DecodeContentToStruct(content, &state)
				elevatorID := msg.SenderID
				globalState[elevatorID] = state
				if c.ElevatorID == c.MasterID {
					for _, ID := range peersOnline {
						globalStateUpdate := utils.CreateMessage(ID, globalState, c.UPDATE_GLOBAL_STATE)
						ch_msgToPack <- globalStateUpdate

						globalHallOrders := getGlobalHallOrders(globalState, peersOnline)
						hallLightsUpdate := utils.CreateMessage(ID, globalHallOrders, c.HALL_LIGHTS_UPDATE)
						ch_msgToPack <- hallLightsUpdate
					}

				}

			// SLAVE
			case c.UPDATE_GLOBAL_STATE:
				if c.ElevatorID != c.MasterID {
					content := msg.Content.([]interface{})
					// Iterates through the array, converts each one to ElevatorState and updates the global state
					for i, value := range content {
						var state e.ElevatorState
						utils.DecodeContentToStruct(value, &state)
						globalState[i] = state
					}
				}
			}

		case update := <-ch_peerUpdate:
			fmt.Println("PEER UPDATE:", update)
			// Assign new master
			peersOnline = stringArrayToIntArray(update.Peers)
			c.MasterID = getMaster(peersOnline)
			if c.ElevatorID != c.MasterID {
				continue
			}
			// Distribute orders from lost peers if there are any
			if len(update.Lost) > 0 {
				IDs := stringArrayToIntArray(update.Lost)
				for _, elevatorID := range IDs {
					distributeOrders(globalState[elevatorID], ch_msgToPack)
				}
			}

			if len(update.New) > 0 {
				elevatorID, _ := strconv.Atoi(update.New)
				stateUpdate := utils.CreateMessage(elevatorID, globalState, c.UPDATE_GLOBAL_STATE)
				ch_msgToPack <- stateUpdate

				hallOrders := getGlobalHallOrders(globalState, peersOnline)
				hallOrdersUpdate := utils.CreateMessage(elevatorID, hallOrders, c.HALL_LIGHTS_UPDATE)
				ch_msgToPack <- hallOrdersUpdate

				cabCalls := getCabCalls(globalState[elevatorID])
				for _, order := range cabCalls {
					cabCallsUpdate := utils.CreateMessage(elevatorID, order, c.DO_ORDER)
					ch_msgToPack <- cabCallsUpdate
				}

			}
		}

	}

}

func getCabCalls(e e.ElevatorState) []drv.ButtonEvent {
	orders := e.Orders
	var cabColumn []bool
	var cabCalls []drv.ButtonEvent
	for _, row := range orders {
		cabColumn = append(cabColumn, row[len(row)-1])
	}
	for floor, cabCall := range cabColumn {
		if cabCall {
			order := drv.ButtonEvent{Floor: floor, Button: drv.BT_Cab}
			cabCalls = append(cabCalls, order)
		}
	}
	return cabCalls
}

// distributeOrders sends every order from the given elevator to the master of the system
func distributeOrders(elevator e.ElevatorState, ch_msgToPack chan<- c.NetworkMessage) {
	orders := elevator.Orders
	//println("DISTRIBUTE THIS ELEVATOR")
	//e.PrintState(elevator)
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
	for index, elevID := range elevatorIDs {
		ElevatorCost := Cost(GlobalState[elevID], order)
		println("ID:", index, ", COST:", ElevatorCost)
		if ElevatorCost < cost {
			cost = ElevatorCost
			lowestCostID = elevID
		}
	}

	return lowestCostID
}

func getGlobalHallOrders(globalState []e.ElevatorState, onlineElevs []int) [][]bool {
	buttons := e.InitElev(0).Orders
	for _, ID := range onlineElevs {
		for floor := 0; floor < c.N_FLOORS; floor++ {
			for btn := 0; btn < c.N_BUTTONS-1; btn++ {
				if globalState[ID].Orders[floor][btn] == true {
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
