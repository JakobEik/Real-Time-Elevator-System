package utilities

import e "Project/elevator"

func initGlobalState() []e.ElevatorState {
	globalState := make([]e.ElevatorState, N_ELEVATORS)
	for i := 0; i < N_ELEVATORS-1; i++ {
		globalState = append(globalState, e.InitElev())
	}
	return globalState
}
