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


	globalState := utils.InitGlobalState()
	println(globalState)
	localElevatorState := e.InitElev(0)
	fmt.Println(localElevatorState)
	// Ask network if they have a global state: true => globalState = this state, else
	if false {
		//TODO
	}
	for {
		select {
		case order := <-ch_buttonPress:
			if order.Button == drv.BT_Cab {
				ch_executeOrder <- order
			}else{
				ch_newLocalOrder <- order
				fmt.Println("order sent")
			}

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

