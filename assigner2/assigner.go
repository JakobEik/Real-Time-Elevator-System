package assigner2

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	"Project/utils"
	"Project/watchdog"
	"fmt"
	"math/rand"
	"strconv"
)

func Assigner(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Assigner")
	//globalStateUpdated := false
	globalState := utils.InitGlobalState()
	println("LENGTH GLOBAL STATE:", len(globalState))
	var peersOnline []int
	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true

		case m := <-ch_msgToAssigner:
			//fmt.Println("MASTER RECEIVE:", msg.Type)
			msg := utils.DecodeMessage(m)
			switch msg.Type {
			// ============ MASTER ===========
			case c.NEW_ORDER:
				if isElevatorOffline(msg.SenderID, peersOnline) || c.ElevatorID != c.MasterID {
					continue
				}
				order := msg.Content.(drv.ButtonEvent)
				bestScoreElevator := getBestElevatorForOrder(globalState, order, peersOnline)
				doOrder := utils.CreateMessage(bestScoreElevator, order, c.DO_ORDER)
				ch_msgToPack <- doOrder
				//fmt.Println("ORDER to elevator:", bestScoreElevator)

			case c.LOCAL_STATE_CHANGED:
				if c.ElevatorID != c.MasterID {
					continue
				}
				state := msg.Content.(e.ElevatorState)
				elevatorID := msg.SenderID
				globalState[elevatorID] = state

				for _, ID := range peersOnline {
					globalStateUpdate := utils.CreateMessage(ID, globalState, c.UPDATE_GLOBAL_STATE)
					ch_msgToPack <- globalStateUpdate

					globalHallOrders := getGlobalHallOrders(globalState, peersOnline)
					hallLightsUpdate := utils.CreateMessage(ID, globalHallOrders, c.HALL_LIGHTS_UPDATE)
					ch_msgToPack <- hallLightsUpdate
				}

			case c.NEW_MASTER:
				println("NEW MASTER")
				if c.ElevatorID != c.MasterID {
					continue
				}
				println("UPDATE NEW MASTER")
				states := msg.Content.([]e.ElevatorState)
				globalState = states

			// ============ SLAVE ===========
			case c.UPDATE_GLOBAL_STATE:
				if c.ElevatorID == c.MasterID {
					continue
				}
				states := msg.Content.([]e.ElevatorState)
				globalState = states
			}

		case update := <-ch_peerUpdate:
			fmt.Println("PEER UPDATE:", update)
			peersOnline = stringArrayToIntArray(update.Peers)
			newMaster := getMaster(peersOnline)
			if c.MasterID != newMaster && isNewElevator(c.ElevatorID, update) {
				println("SEND STATE TO NEW MASTER")
				// If there is a new master, the slave elevators will send their backup globalstates to the new master
				newMasterUpdate := utils.CreateMessage(newMaster, globalState, c.NEW_MASTER)
				ch_msgToPack <- newMasterUpdate
			}
			// Assign new master
			c.MasterID = newMaster
			println("MASTER:", c.MasterID)

			if c.ElevatorID == c.MasterID {
				lostPeersUpdate(update, globalState, ch_msgToPack)
				newPeerUpdate(update, globalState, ch_msgToPack, peersOnline)
			}
		}
		//printElevFloors(globalState)

	}

}

func printElevFloors(globalState []e.ElevatorState) {
	println()
	println(rand.Int())
	for ID, state := range globalState {
		fmt.Println("ELEVATOR :", ID, ", FLOOR :", state.Floor)
	}
}

func lostPeersUpdate(update p.PeerUpdate, globalState []e.ElevatorState, ch_msgToPack chan<- c.NetworkMessage) {
	// Distribute orders from lost peers if there are any
	if len(update.Lost) > 0 {
		IDs := stringArrayToIntArray(update.Lost)
		for _, elevatorID := range IDs {
			distributeOrders(globalState[elevatorID], ch_msgToPack)
		}
	}

}

func newPeerUpdate(
	update p.PeerUpdate,
	globalState []e.ElevatorState,
	ch_msgToPack chan<- c.NetworkMessage,
	peersOnline []int) {

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
