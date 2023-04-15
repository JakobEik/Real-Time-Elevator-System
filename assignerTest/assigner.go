package assignerTest

import (
	c "Project/config"
	p "Project/network/peers"
	"Project/utils"
	"Project/watchdog"
	"fmt"
	//drv "Project/driver"
)

//var globalState = utils.InitGlobalState()

func Assigner(
	ch_peerUpdate chan p.PeerUpdate,
	ch_msgToAssigner <-chan c.NetworkMessage,
	ch_msgToPack chan<- c.NetworkMessage) {

	ch_bark := make(chan bool)
	ch_pet := make(chan bool)
	go watchdog.Watchdog(ch_bark, ch_pet, "Assigner")

	master := masterNode{
		ch_msgToPack: ch_msgToPack,
		peersOnline:  make([]int, c.N_ELEVATORS),
		globalState:  utils.InitGlobalState(),
	}
	slave := slaveNode{globalStateBackup: utils.InitGlobalState()}

	for {
		select {
		//Watchdog
		case <-ch_bark:
			ch_pet <- true

		case update := <-ch_peerUpdate:
			fmt.Println("PEER UPDATE:", update)
			if !isElevatorOnlineStr(c.MasterID, update.Peers) {
				// Master went offline, assign new master
				c.MasterID = getUpdatedMaster(update.Peers)
				println("MASTER OFFLINE, NEW MASTER:", c.MasterID)
			}
			if c.ElevatorID == c.MasterID {
				master.peerUpdate(update)
			} else {
				master.globalState = slave.globalStateBackup
			}

		case msg := <-ch_msgToAssigner:
			switch msg.Type {
			// MASTER
			case c.NEW_ORDER:
				if c.ElevatorID == c.MasterID {
					master.newOrderEvent(msg)
				}

			case c.LOCAL_STATE_CHANGED:
				if c.ElevatorID == c.MasterID {
					master.newLocalStateEvent(msg)
				}

			// SLAVE
			case c.UPDATE_GLOBAL_STATE:
				if c.ElevatorID != c.MasterID {
					slave.globalStateUpdate(msg)

				}

			}

		}

	}
}
