package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/verda-cloud/model-templates/pkg/mapping"
	"github.com/verda-cloud/model-templates/pkg/template"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run examples/mapping_usage.go <config-file.json>")
		fmt.Println("\nExample:")
		fmt.Println("  go run examples/mapping_usage.go templates/qwen3-sglang.json")
		os.Exit(1)
	}

	configPath := os.Args[1]

	// Read the configuration file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// Parse the JSON configuration into the template.Config struct
	var config template.Config
	if err := json.Unmarshal(configData, &config); err != nil {
		log.Fatalf("Failed to parse config JSON: %v", err)
	}

	// Get the engine type
	engine := string(config.Engine)
	if engine == "" {
		log.Fatal("Config must specify an 'engine' field")
	}

	fmt.Printf("Loading mappings for engine: %s\n", engine)

	// Create a command builder for the specified engine
	builder, err := mapping.NewCommandBuilder(engine)
	if err != nil {
		log.Fatalf("Failed to create command builder: %v", err)
	}

	// Build the command
	command, err := builder.Build(config)
	if err != nil {
		log.Fatalf("Failed to build command: %v", err)
	}

	// Display the generated command
	fmt.Println("\nGenerated Command:")
	fmt.Println(strings.Join(command, " "))

	// Display as a formatted multi-line command
	fmt.Println("\nFormatted Command:")
	if len(command) > 0 {
		fmt.Print(command[0])
		for _, arg := range command[1:] {
			if strings.HasPrefix(arg, "-") {
				fmt.Printf(" \\\n  %s", arg)
			} else {
				fmt.Printf(" %s", arg)
			}
		}
		fmt.Println()
	}
}
