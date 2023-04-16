package distributor

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
	"Project/watchdog"
)

// Distributor distributes all events between FSM and PacketDistributor
func Distributor(
	ch_msgToDistributor <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage,
	// To local
	ch_executeOrder chan<- drv.ButtonEvent,
	ch_hallLights chan<- [][]bool,
	// From Local
	ch_localStateUpdated <-chan c.ElevatorState,
	ch_buttonPress <-chan drv.ButtonEvent,
	// Monitoring channels
	ch_unavailable <-chan bool,
	ch_peerTxEnable chan<- bool) {

	ch_peerTxEnable <- true
	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Distributor")

	for {
		select {
		//__________LOCAL CHANNELS_____________
		case order := <-ch_buttonPress:
			if order.Button == drv.BT_Cab {
				ch_executeOrder <- order
			} else {
				msg := utils.CreateMessage(c.MasterID, order, c.NEW_ORDER)
				ch_msgToPack <- msg
			}

		case state := <-ch_localStateUpdated:
			e.AcceptanceTests(state)
			msg := utils.CreateMessage(c.MasterID, state, c.LOCAL_STATE_CHANGED)
			ch_msgToPack <- msg

		// ____________NETWORK MESSAGES_____________
		case m := <-ch_msgToDistributor:
			msg := utils.DecodeContent(m)
			switch msg.Type {
			case c.DO_ORDER:
				order := msg.Content.(drv.ButtonEvent)
				ch_executeOrder <- order
			case c.HALL_LIGHTS_UPDATE:
				orders := msg.Content.([][]bool)
				ch_hallLights <- orders
			}
		// ____________MONITORING___________
		case <-ch_bark:
			ch_pet <- true
		case failure := <-ch_unavailable:
			ch_peerTxEnable <- !failure
		}
	}
}
