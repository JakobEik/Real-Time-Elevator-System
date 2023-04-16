package assigner

import (
	c "Project/config"
	drv "Project/driver"
	p "Project/network/peers"
	"Project/utils"
	"Project/watchdog"
	"fmt"
	"strconv"
)

func Assigner(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_packMessage chan<- c.NetworkMessage,
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
				bestElevatorID := getBestElevatorForOrder(globalState, order, peersOnline)
				doOrder := utils.CreateMessage(bestElevatorID, order, c.DO_ORDER)
				ch_packMessage <- doOrder
			// ============ MASTER ===========
			case c.LOCAL_STATE_CHANGED:
				if !isMaster {
					continue
				}
				state := msg.Content.(c.ElevatorState)
				elevatorID := msg.SenderID
				globalState[elevatorID] = state
				for _, ID := range peersOnline {
					globalStateUpdate := utils.CreateMessage(ID, globalState, c.UPDATE_GLOBAL_STATE)
					ch_packMessage <- globalStateUpdate

					globalHallOrders := getHallOrders(globalState, peersOnline)
					hallLightsUpdate := utils.CreateMessage(ID, globalHallOrders, c.HALL_LIGHTS_UPDATE)
					ch_packMessage <- hallLightsUpdate
				}
			// ============ MASTER ===========
			case c.NEW_MASTER:
				println("NEW MASTER")
				if !isMaster {
					continue
				}
				fmt.Println("UPDATE THIS TO MASTER")
				backupGlobalState := msg.Content.([]c.ElevatorState)
				globalState = backupGlobalState
				master := strconv.Itoa(c.ElevatorID)
				// When master goes online, this function is needed to update its cab calls
				// This is because when ch_peerUpdate receives the first time, the global state of the master is initialized.
				// The master needs to receive the backup state from the other slaves first and update
				newPeerUpdate(master, globalState, ch_packMessage, peersOnline)

			// ============ SLAVE ===========
			case c.UPDATE_GLOBAL_STATE:
				if isMaster {
					continue
				}
				states := msg.Content.([]c.ElevatorState)
				globalState = states
			}

		//________ PEER UPDATE ________
		// ============ MASTER AND SLAVE ===========
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
				ch_packMessage <- newMasterUpdate
			}
			// ============ MASTER ===========
			if isMaster {
				// HANDLE LOST PEERS
				lostPeersIDs := stringArrayToIntArray(update.Lost)
				lostOrders := getHallOrders(globalState, lostPeersIDs)
				ordersToDistribute := makeMessagesFromOrders(lostOrders)
				for _, message := range ordersToDistribute {
					ch_packMessage <- message
				}
				globalState = updateGlobalState(globalState, lostPeersIDs)

				// HANDLE NEW PEERS
				newPeerUpdate(update.New, globalState, ch_packMessage, peersOnline)

			}

		}
	}
}

// newPeerUpdate updates the state of the new peer.
// This includes its backup global state, hall lights, and cab orders if there are any
func newPeerUpdate(
	newPeer string,
	globalState []c.ElevatorState,
	ch_msgToPack chan<- c.NetworkMessage,
	peersOnline []int) {

	if len(newPeer) > 0 {
		elevatorID, _ := strconv.Atoi(newPeer)
		stateUpdate := utils.CreateMessage(elevatorID, globalState, c.UPDATE_GLOBAL_STATE)
		ch_msgToPack <- stateUpdate

		hallOrders := getHallOrders(globalState, peersOnline)
		hallLightsUpdate := utils.CreateMessage(elevatorID, hallOrders, c.HALL_LIGHTS_UPDATE)
		ch_msgToPack <- hallLightsUpdate

		cabCalls := getCabCalls(globalState[elevatorID])
		for _, order := range cabCalls {
			cabCallsUpdate := utils.CreateMessage(elevatorID, order, c.DO_ORDER)
			ch_msgToPack <- cabCallsUpdate
		}

	}

}
