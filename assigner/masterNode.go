package assigner

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	p "Project/network/peers"
	"Project/utils"
	"fmt"
	//drv "Project/driver"
)

func MasterNode(
	ch_peerUpdate chan p.PeerUpdate,
	ch_peerTxEnable chan bool,
	ch_msgFromNetwork <-chan c.Packet,
	ch_msgToNetwork chan<- c.Packet) {

	masterID := 0
	globalState := utils.InitGlobalState()
	println("LENGTH GLOBALSTATE:", len(globalState))

	for packet := range ch_msgFromNetwork {
		if acceptPacket(packet, masterID) {
			content := packet.Msg.Content
			switch packet.Msg.MsgType {
			case c.NEW_ORDER:
				var order drv.ButtonEvent
				utils.CastToType(content, &order)
				lowestCostElevator := calculateCost(globalState, order)
				packet := utils.CreatePacket(lowestCostElevator, order, c.DO_ORDER)
				ch_msgToNetwork <- packet
				fmt.Println("ORDER to elevator:", lowestCostElevator)

			case c.LOCAL_STATE_CHANGED:
				var state e.ElevatorState
				utils.CastToType(content, &state)
				elevatorID := packet.Msg.SenderID
				globalState[elevatorID] = state
				packetState := utils.CreatePacket(c.ToEveryone, globalState, c.UPDATE_GLOBAL_STATE)
				ch_msgToNetwork <- packetState

				globalHallOrders := getGlobalHallOrders(globalState)
				packetHall := utils.CreatePacket(c.ToEveryone, globalHallOrders, c.GLOBAL_HALL_ORDERS)
				ch_msgToNetwork <- packetHall
			}
		}
	}

}

func acceptPacket(packet c.Packet, masterID int) bool {
	//TODO: CHECKSUM
	checksumCorrect := packet.Checksum == 0
	return checksumCorrect && masterID == c.ElevatorID

}

func setMaster(masterID *int) {
	//TODO: IMPLEMENT GO ROUTINE

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
