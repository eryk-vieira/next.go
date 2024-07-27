package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	build_tui "github.com/eryk-vieira/next.go/cli/nextgo/tui"
	"github.com/eryk-vieira/next.go/cli/nextgo/types"
)

func main() {
	settings := parseSettings(os.Args[2:])
	command := getCommand(os.Args[1:])

	if command == "build" {
		build_tui.Run(settings)

		return
	}

	if command == "run" {
		isOpen := raw_connect("localhost", settings.Server.Port)

		if isOpen {
			log.Fatalf("Port %s already in use", settings.Server.Port)

			return
		}

		fmt.Println(fmt.Sprintf("Runnning server on port: %s", settings.Server.Port))

		cmd := exec.Command("./server")
		cmd.Dir = ".dist/"
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		cmd.Run()
	}
}

func parseSettings(args []string) *types.Settings {
	var settings types.Settings

	var configFile string = "nextgo.json"

	if len(args) > 0 {
		configFile = args[0]
	}

	jsonFile, err := os.ReadFile(configFile)

	if err != nil {
		panic(fmt.Sprintf("%s file not found", configFile))
	}

	err = json.Unmarshal(jsonFile, &settings)

	if err != nil {
		panic(err)
	}

	return &settings
}

func getCommand(args []string) string {
	return args[0]
}

func raw_connect(host string, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)

	if err != nil {
		return false
	}

	if conn != nil {
		return true
	}

	return false
}
