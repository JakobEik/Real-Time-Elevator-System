package network

import (
	a "Project/assigner"
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
)

func NetworkHandler(
	ch_msgFromNetwork <-chan c.Packet,
	ch_msgToNetwork chan<- c.Packet,
	// To Master
	ch_calculateNewOrder chan<- drv.ButtonEvent,
	ch_localStateToMaster chan<- a.StateMessage,
	// To Local
	ch_doOrder chan<- drv.ButtonEvent,
	ch_updateGlobalState chan<- []e.ElevatorState,
	// From Master
	ch_orderCalculated <-chan a.OrderMessage,
	ch_updateGlobalStateDemand <-chan []e.ElevatorState,
	//From Local
	ch_localStateFromLocal <-chan e.ElevatorState,
	ch_newLocalOrder <-chan drv.ButtonEvent) {

	masterID := 0

	for {
		select {
		case order := <-ch_newLocalOrder:
			packet := utils.CreatePacket(masterID, order, c.NewOrder)
			ch_msgToNetwork <- packet

		case order := <-ch_orderCalculated:
			packet := utils.CreatePacket(order.ReceiverID, order.Order, c.DoOrder)
			ch_msgToNetwork <- packet

		case state := <-ch_localStateFromLocal:
			packet := utils.CreatePacket(masterID, state, c.LocalStateChange)
			ch_msgToNetwork <- packet

		case states := <-ch_updateGlobalStateDemand:
			packet := utils.CreatePacket(c.ToEveryone, states, c.UpdateGlobalState)
			ch_msgToNetwork <- packet

		case packet := <-ch_msgFromNetwork:
			if true {
				//TODO: CHECKSUM
				switch packet.Msg.MsgType {
				case c.NewOrder:
					content := packet.Msg.Content.(map[string]interface{})
					var order drv.ButtonEvent
					utils.CastToType(content, &order)
					ch_calculateNewOrder <- order
				case c.DoOrder:
					content := packet.Msg.Content.(map[string]interface{})
					var order drv.ButtonEvent
					utils.CastToType(content, &order)
					if packet.Msg.ReceiverID == c.ElevatorID {
						ch_doOrder <- order
					}
				case c.LocalStateChange:
					content := packet.Msg.Content.(map[string]interface{})
					var state e.ElevatorState
					utils.CastToType(content, &state)
					ch_localStateToMaster <- a.StateMessage{State: state, SenderID: packet.Msg.SenderID}
				/*case c.UpdateGlobalState:
				content := content.([]interface{})
				fmt.Println(content[0])
				fmt.Println(reflect.TypeOf(content[0]))

				var states []e.ElevatorState
				panic("S")
				//utils.CastToType(content.(map[string]interface{}), &states)
				ch_updateGlobalState <- states*/
				case c.MsgReceived:
					//TODO
				}

			}
		}
	}

}
