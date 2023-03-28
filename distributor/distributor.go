package distributor

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
	"fmt"
	//"Project/network/peers"
)

func Distributor(
	ch_msgFromNetwork <-chan c.Packet,
	ch_msgToNetwork chan<- c.Packet,
// To local
	ch_executeOrder chan<- drv.ButtonEvent,
	ch_globalHallOrders chan<- [][]bool,
//From Local
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_buttonPress <-chan drv.ButtonEvent) {

	globalState := utils.InitGlobalState()
	println(globalState)
	masterID := 0

	for {
		select {
		case order := <-ch_buttonPress:
			//fmt.Println("Buttonpress:", order)
			if order.Button == drv.BT_Cab {
				ch_executeOrder <- order
			} else {
				packet := utils.CreatePacket(masterID, order, c.NEW_ORDER)
				ch_msgToNetwork <- packet
				//fmt.Println("order sent to master")

			}

		case state := <-ch_localStateUpdated:
			packet := utils.CreatePacket(masterID, state, c.LOCAL_STATE_CHANGED)
			ch_msgToNetwork <- packet

		case packet := <-ch_msgFromNetwork:
			if acceptPacket(packet) {
				content := packet.Msg.Content
				switch packet.Msg.MsgType {
				case c.DO_ORDER:
					var order drv.ButtonEvent
					utils.CastToType(content, &order)
					if packet.Msg.ReceiverID == c.ElevatorID {
						ch_executeOrder <- order
						fmt.Println("EXECUTE ORDER:", order)
					}

				case c.UPDATE_GLOBAL_STATE:
					content := packet.Msg.Content.([]interface{})
					// Iterates through the array, converts each one to ElevatorState and updates the global state
					for i, value := range content {
						var state e.ElevatorState
						utils.CastToType(value, &state)
						globalState[i] = state
					}

				case c.GLOBAL_HALL_ORDERS:
					var buttons [][]bool
					utils.CastToType(content, &buttons)
					ch_globalHallOrders <- buttons

				}
			}
		}

	}
}

func acceptPacket(packet c.Packet) bool {
	//TODO: CHECKSUM

	checksumCorrect := true
	receiverID := packet.Msg.ReceiverID
	return checksumCorrect && (receiverID == c.ElevatorID || receiverID == c.ToEveryone)

}
