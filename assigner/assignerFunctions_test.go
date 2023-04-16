package assigner

import (
	"Project/driver"
	e "Project/elevator"
	"Project/network/peers"
	"Project/utils"
	"reflect"
	"testing"
)

func Test_updateGlobalState(t *testing.T) {

	globalState := utils.InitGlobalState()
	e1 := globalState[0]
	e1.Orders[0][0] = true
	e1.Orders[1][1] = true
	e1.Orders[2][2] = true

	e2 := globalState[1]
	e2.Orders[0][0] = true
	e2.Orders[1][1] = true
	e2.Orders[2][2] = true

	lostPeers := []int{0, 1}

	want := utils.InitGlobalState()
	want[0].Orders[2][2] = true
	want[1] = e2

	got := updateGlobalState(globalState, lostPeers)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("updateGlobalState(%v, %v) = %v; want %v", globalState, lostPeers, got, want)
	}
}

func Test_getCabCalls(t *testing.T) {
	elev := e.InitElev(0)
	elev.Orders[0][0] = true
	elev.Orders[1][1] = true
	elev.Orders[2][2] = true

	var want []driver.ButtonEvent
	want = append(want, driver.ButtonEvent{Floor: 2, Button: 2})

	got := getCabCalls(elev)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getCabCalls(%v) = %v; want %v", elev, got, want)
	}
}

func Test_getGlobalHallOrders(t *testing.T) {
	globalState := utils.InitGlobalState()
	e1 := globalState[0]
	e1.Orders[2][1] = true
	e1.Orders[3][1] = true
	e1.Orders[2][2] = true

	e2 := globalState[1]
	e2.Orders[0][0] = true
	e2.Orders[1][1] = true
	e2.Orders[2][2] = true

	lostPeers := []int{0}

	want := e.InitElev(0).Orders
	want[2][1] = true
	want[3][1] = true

	got := getHallOrders(globalState, lostPeers)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("getGlobalHallOrders() = %v; want %v", got, want)
	}

}

func Test_removeHallOrders(t *testing.T) {
	orders := [][]bool{{true, true, true}, {true, true, true}}
	want := [][]bool{{false, false, true}, {false, false, true}}
	got := removeHallOrders(orders)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("removeHallOrders() = %v; want %v", got, want)
	}
}

func Test_isNewElevator(t *testing.T) {
	elevID := 1
	update := peers.PeerUpdate{
		Peers: []string{"0", "2"},
		New:   "1",
		Lost:  nil,
	}

	got := isNewElevator(elevID, update)
	if !reflect.DeepEqual(got, true) {
		t.Errorf("isNewElevator() = %v; want %v", got, true)
	}

	elevID = 1
	update = peers.PeerUpdate{
		Peers: nil,
		New:   "2",
		Lost:  nil,
	}

	got = isNewElevator(elevID, update)
	if !reflect.DeepEqual(got, false) {
		t.Errorf("isNewElevator() = %v; want %v", got, false)
	}
}

func Test_isElevatorOffline(t *testing.T) {
	elevID := 1
	peersOnline := []int{0, 2}

	got := isElevatorOffline(elevID, peersOnline)
	if !reflect.DeepEqual(got, true) {
		t.Errorf("isNewElevator() = %v; want %v", got, true)
	}
}
