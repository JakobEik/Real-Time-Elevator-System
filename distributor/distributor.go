package distributor

import (
	c "Project/config"
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
	"fmt"
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
	// watchdog
	ch_watchdogStuckBark <-chan bool,
) {

	for {
		select {
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
			msg := utils.CreateMessage(c.MasterID, state, c.LOCAL_STATE_CHANGED)
			ch_msgToPack <- msg

		// NETWORK MESSAGES
		case msg := <-ch_msgToDistributor:
			//fmt.Println("DISTRIBUTOR RECEIVE:", msg.Type)
			content := msg.Content
			switch msg.Type {
			case c.DO_ORDER:
				var order drv.ButtonEvent
				utils.DecodeContentToStruct(content, &order)
				ch_executeOrder <- order
				fmt.Println("EXECUTE ORDER:", order)
			case c.GLOBAL_HALL_ORDERS:
				var orders [][]bool
				utils.DecodeContentToStruct(content, &orders)
				ch_globalHallOrders <- orders
			}

		case <-ch_watchdogStuckBark:
			fmt.Println("WATCHDOG BARK FROM DISTRIBUTOR")
		}
	}
}
