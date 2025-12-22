package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/verda-cloud/model-templates/pkg/engine"
	"github.com/verda-cloud/model-templates/pkg/template"
)

var (
	silent      bool
	checkModels bool
)

var rootCmd = &cobra.Command{
	Use:   "validate [files...]",
	Short: "Validate model template configuration files",
	Long:  `Validate one or more JSON template configuration files. If no files are specified, all templates/*.json files will be validated.`,
	RunE:  runValidate,
}

func init() {
	rootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent mode - only return exit code")
	rootCmd.Flags().BoolVar(&checkModels, "check-models", false, "Verify that vLLM and SGLang models exist on HuggingFace")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	// Get files to verify
	var files []string

	if len(args) == 0 {
		// No arguments - test all templates/*.json
		matches, err := filepath.Glob("templates/*.json")
		if err != nil {
			return fmt.Errorf("error finding example files: %v", err)
		}
		files = matches
	} else {
		// Use specified file(s)
		files = args
	}

	if len(files) == 0 {
		return fmt.Errorf("no files to verify")
	}

	hasErrors := false
	for _, file := range files {
		if err := verifyFile(file, silent, checkModels); err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying %s: %v\n", file, err)
			hasErrors = true
		}
	}

	if hasErrors {
		return fmt.Errorf("validation failed")
	}

	return nil
}

func verifyFile(filePath string, silent bool, checkModels bool) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file %s: %v", filePath, err)
	}

	var cfg template.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to parse file %s: %v", filePath, err)
	}

	_, err = engine.GenerateCommandFromTemplate(cfg)
	if err != nil {
		return fmt.Errorf("invalid file format in file %s: %v", filePath, err)
	}

	// Check if model exists on HuggingFace (only for vLLM and SGLang engines)
	if checkModels && (cfg.Engine == "vllm" || cfg.Engine == "sglang") && cfg.Model != "" {
		if !silent {
			fmt.Printf("  Checking model %s on HuggingFace... ", cfg.Model)
		}

		exists, err := checkModelExists(cfg.Model)
		if err != nil {
			if !silent {
				fmt.Printf("⚠ (error: %v)\n", err)
			}
			return fmt.Errorf("failed to check model %s: %v", cfg.Model, err)
		}

		if !exists {
			if !silent {
				fmt.Printf("✗ NOT FOUND\n")
			}
			return fmt.Errorf("model %s not found on HuggingFace", cfg.Model)
		}

		if !silent {
			fmt.Printf("✓\n")
		}
	}

	if !silent {
		fileName := filepath.Base(filePath)
		if cfg.ShortExplanation != "" {
			fmt.Printf("✓ %-40s %s\n", fileName, cfg.ShortExplanation)
		} else {
			fmt.Printf("✓ %s\n", fileName)
		}
	}

	return nil
}

// checkModelExists verifies if a model exists on HuggingFace
func checkModelExists(modelPath string) (bool, error) {
	url := fmt.Sprintf("https://huggingface.co/%s", modelPath)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 200 = model exists (both public and gated models return 200)
	// 401/404 = model does not exist or is private
	// HuggingFace returns 401 for non-existent models as a security feature
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}
