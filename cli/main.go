package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/eryk-vieira/next.go/cli/build"
	"github.com/eryk-vieira/next.go/cli/runner"
	"github.com/eryk-vieira/next.go/cli/types"
)

func main() {
	settings := parseSettings(os.Args[2:])
	command := getCommand(os.Args[1:])

	if command == "build" {
		done := make(chan bool)

		// Start a goroutine to update and print the timer
		go func() {
			start := time.Now()
			for {
				select {
				case <-done:
					return
				default:
					fmt.Printf("\rBuilding: %s", time.Since(start).Round(time.Millisecond))
				}
			}
		}()

		builder := build.NewBuilder(settings)
		builder.Build()

		done <- true

		fmt.Printf("\nBuild successfully!")

		return
	}

	if command == "run" {
		runner := runner.ServerRunner{
			Port:   settings.Server.Port,
			Router: "chi",
		}

		runner.Run()
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
