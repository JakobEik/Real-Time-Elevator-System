package distributor

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
	"Project/watchdog"
)

func Distributor(
	ch_msgToDistributor <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage,
	// To local
	ch_executeOrder chan<- drv.ButtonEvent,
	ch_globalHallOrders chan<- [][]bool,
	// From Local
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_buttonPress <-chan drv.ButtonEvent,
	// Monitoring channels
	ch_failure <-chan bool,
	ch_peerTxEnable chan<- bool) {

	ch_peerTxEnable <- true // Enable peer transmitter by default
	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Distributor")

	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true
		// LOCAL CHANNELS
		case order := <-ch_buttonPress:
			//fmt.Println("Buttonpress:", order)
			if order.Button == drv.BT_Cab {
				// Local state will be updated by this and sent to master => not necessary send the order to master
				ch_executeOrder <- order
			} else {
				msg := utils.CreateMessage(c.MasterID, order, c.NEW_ORDER)
				ch_msgToPack <- msg
				//fmt.Println("order sent to master")
			}

		case state := <-ch_localStateUpdated:
			e.AcceptanceTests(state)
			msg := utils.CreateMessage(c.MasterID, state, c.LOCAL_STATE_CHANGED)
			ch_msgToPack <- msg

		// NETWORK MESSAGES
		case m := <-ch_msgToDistributor:
			//fmt.Println("DISTRIBUTOR RECEIVE:", msg.Type)
			msg := utils.DecodeMessage(m)
			switch msg.Type {
			case c.DO_ORDER:
				order := msg.Content.(drv.ButtonEvent)
				ch_executeOrder <- order
				//fmt.Println("EXECUTE ORDER:", order)
			case c.HALL_LIGHTS_UPDATE:
				orders := msg.Content.([][]bool)
				ch_globalHallOrders <- orders
			}

		case failure := <-ch_failure:
			println("FAILURE==========", failure)
			ch_peerTxEnable <- !failure
		}
	}
}
