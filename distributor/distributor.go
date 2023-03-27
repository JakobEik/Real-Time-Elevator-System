package distributor

import (
	drv "Project/driver"
	e "Project/elevator"
	"Project/utils"
	"fmt"
	//"Project/network/peers"
)

func Distributor(
	ch_doOrder <-chan drv.ButtonEvent,
	ch_localStateUpdated <-chan e.ElevatorState,
	ch_buttonPress <-chan drv.ButtonEvent,
	ch_updateGlobalState <-chan []e.ElevatorState,
	ch_localStateFromLocal chan<- e.ElevatorState,
	ch_newLocalOrder chan<- drv.ButtonEvent,
	ch_executeOrder chan<- drv.ButtonEvent) {

<<<<<<< HEAD

	globalState := utils.InitGlobalState()
	println(globalState)
	localElevatorState := e.InitElev(0)
	fmt.Println(localElevatorState)
=======
	//globalState := c.InitGlobalState()
	localElevatorState := e.InitElev(0)
>>>>>>> e6828a179660c35014b6bcf5dab0d9db9fce18e5
	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	}
	for {
		select {
<<<<<<< HEAD
		case order := <-ch_buttonPress:
			if order.Button == drv.BT_Cab {
				ch_executeOrder <- order
			}else{
				ch_newLocalOrder <- order
				fmt.Println("order sent")
			}

=======
		case newLocalOrder := <-ch_newLocalOrder:
			newLocalOrderEvent(newLocalOrder, ch_messageToNetwork, masterID)
>>>>>>> e6828a179660c35014b6bcf5dab0d9db9fce18e5
		case newLocalState := <-ch_localStateUpdated:
			localElevatorState = newLocalState
			ch_localStateFromLocal <- newLocalState
		case state := <-ch_updateGlobalState:
			globalState = state
		case order := <-ch_doOrder:
			ch_executeOrder <- order	
		}

	}
}

<<<<<<< HEAD
=======
func newLocalOrderEvent(order drv.ButtonEvent, ch_messageToNetwork chan<- c.NetworkMessage, masterID int) {
	msg := utils.CreateMessage(masterID, masterID, order, c.NewOrder)
	fmt.Println("order sent:", msg)
	ch_messageToNetwork <- msg
}

func localStateUpdatedEvent(
	newState e.ElevatorState,
	oldState e.ElevatorState, //pointer???
	ch_messageToNetwork chan<- c.NetworkMessage,
	masterID int) {

	oldState = newState
	msg := utils.CreateMessage(masterID, masterID, newState, c.LocalStateChange)
	ch_messageToNetwork <- msg

}
func newMessageEvent(msg c.NetworkMessage, ch_doOrder chan<- drv.ButtonEvent) {

	if msg.ReceiverID == c.ElevatorID {
		switch msg.MsgType {
		case c.DoOrder:
			fmt.Println("Order received:", msg.Content)
			content := msg.Content.(map[string]interface{})
			var order drv.ButtonEvent
			utils.ConvertMapToStruct(content, &order)
			ch_doOrder <- order

		}

	}
}
>>>>>>> e6828a179660c35014b6bcf5dab0d9db9fce18e5
