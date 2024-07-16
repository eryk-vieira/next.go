package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/eryk-vieira/next.go/cli/nextgo/build"
	"github.com/eryk-vieira/next.go/cli/nextgo/types"
)

type Hello string

func main() {
	settings := parseSettings(os.Args[2:])
	command := getCommand(os.Args[1:])

	if command == "build" {
		builder := build.NewBuilder(settings)
		builder.Build()

		fmt.Println("Build successfully!")

		return
	}

	if command == "run" {
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

	var configFile string = "nextgo.config.json"

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
