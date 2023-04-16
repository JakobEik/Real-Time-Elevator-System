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
)

func Assigner(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage,
	ch_offNetwork chan<- bool) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Assigner")
	globalState := utils.InitGlobalState()

	isMaster := true
	var peersOnline []int
	for {
		select {
		//________ WATCHDOG ________
		case <-ch_bark:
			ch_pet <- true
		//________ NEW MESSAGE ________
		case m := <-ch_msgToAssigner:
			//fmt.Println("MASTER RECEIVE:", m.Type)
			msg := utils.DecodeContent(m)
			switch msg.Type {
			// ============ MASTER ===========
			case c.NEW_ORDER:
				if !isMaster {
					continue
				}
				order := msg.Content.(drv.ButtonEvent)
				bestScoreElevator := getBestElevatorForOrder(globalState, order, peersOnline)
				doOrder := utils.CreateMessage(bestScoreElevator, order, c.DO_ORDER)
				ch_msgToPack <- doOrder
			// ============ MASTER ===========
			case c.LOCAL_STATE_CHANGED:
				if !isMaster {
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
			// ============ MASTER ===========
			case c.NEW_MASTER:
				if !isMaster {
					continue
				}
				backupGlobalState := msg.Content.([]e.ElevatorState)
				globalState = backupGlobalState
				master := strconv.Itoa(c.ElevatorID)
				// When master goes online, this function is needed to update its cab calls
				// This is because when ch_peerUpdate receives the first time, the global state of the master is initialized.
				// The master needs to receive the backup state from the other slaves first and update
				newPeerUpdate(master, globalState, ch_msgToPack, peersOnline)

			// ============ SLAVE ===========
			case c.UPDATE_GLOBAL_STATE:
				if isMaster {
					continue
				}
				states := msg.Content.([]e.ElevatorState)
				globalState = states
			}

		//________ PEER UPDATE ________
		case update := <-ch_peerUpdate:
			// If an elevator loses internet connection, this will notify the local FSM to update its own lights
			if len(update.Peers) == 0 {
				ch_offNetwork <- true
			} else {
				ch_offNetwork <- false
			}

			fmt.Println("PEER UPDATE:", update)
			peersOnline = stringArrayToIntArray(update.Peers)
			oldMaster := c.MasterID
			// Master is always assigned to the elevator with the lowest ID on the network.
			// This is to make sure everyone that is connected always agrees on whom the master is
			c.MasterID = getMaster(peersOnline)
			isMaster = c.ElevatorID == c.MasterID
			if c.MasterID != oldMaster && !isNewElevator(c.ElevatorID, update) {
				// If there is a new master, the slave elevators will send him their backup global states
				newMasterUpdate := utils.CreateMessage(c.MasterID, globalState, c.NEW_MASTER)
				ch_msgToPack <- newMasterUpdate
			}

			if isMaster {
				globalState = lostPeersUpdate(update.Lost, globalState, ch_msgToPack)
				newPeerUpdate(update.New, globalState, ch_msgToPack, peersOnline)
			}
		}
	}
}

// lost
func lostPeersUpdate(lostPeers []string, globalState []e.ElevatorState, ch_msgToPack chan<- c.NetworkMessage) []e.ElevatorState {
	newGlobalState := globalState
	// Distribute orders from lost peers if there are any
	if len(lostPeers) > 0 {
		IDs := stringArrayToIntArray(lostPeers)
		for _, elevatorID := range IDs {
			println("DISTRIBUTE ELEVATOR :", elevatorID)
			newGlobalState[elevatorID] = distributeOrders(globalState[elevatorID], ch_msgToPack)
		}
	}
	return newGlobalState

}

func newPeerUpdate(
	newPeer string,
	globalState []e.ElevatorState,
	ch_msgToPack chan<- c.NetworkMessage,
	peersOnline []int) {
	if len(newPeer) > 0 {
		elevatorID, _ := strconv.Atoi(newPeer)
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
