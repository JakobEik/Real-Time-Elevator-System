package failroutine

import (
	"Project/config"
	drv "Project/driver"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func FailRoutine() {

	port := os.Args[1]
	config.ElevatorID, _ = strconv.Atoi(os.Args[2])
	id := os.Args[2]
	oSys := runtime.GOOS

	// wait to restart so that the elevator have time to send false on ch_peerTxEnable before restart
	time.Sleep(time.Second)

	switch oSys {
	case "windows":
		fmt.Println("Windows")
		drv.SetMotorDirection(drv.MD_Stop)
		err := exec.Command("cmd", "/C", "start", "powershell", "go", "run", fmt.Sprintf("main.go %s %s", port, id)).Run()
		if err != nil {
			fmt.Println("Unable to reboot process, crashing...")
		}
		fmt.Println("HALLO FRA FAILROUTINE!")
		// fmt.Println("Program killed !")
		os.Exit(0)

	case "darwin":
		fmt.Println("MAC operating system")
		os.Exit(0)

	case "linux":
		fmt.Println("Linux")
		drv.SetMotorDirection(drv.MD_Stop)
		err := exec.Command("gnome-terminal", "-x", "sh", "-c", fmt.Sprintf("go run main.go %s %s; exec bash", port, id)).Run()
		if err != nil {
			fmt.Println("Unable to reboot process, crashing...")
		}
		fmt.Println("HALLO FRA FAILROUTINE!")
		os.Exit(0)
	default:
		fmt.Println("FUBAR")
		os.Exit(0)
	}
}
