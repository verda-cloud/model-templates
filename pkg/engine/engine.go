package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/verda-cloud/model-templates/pkg/template"
)

func GenerateCommandFromTemplate(cfg template.Config) ([]string, error) {
	switch template.EngineOption(strings.ToLower(string(cfg.Engine))) {
	case template.EngineVLLM:
		return generateVLLMCommand(cfg), nil
	case template.EngineSGLang:
		return generateSGLangCommand(cfg), nil
	case template.EngineCustom:
		return generateCustomCommand(cfg)
	default:
		return nil, fmt.Errorf("unknown engine option: %s", string(cfg.Engine))
	}
}

// appendCommonOptions adds common configuration options shared by both vLLM and SGLang.
// As some flags differ, the engine-specific ones are provided in the function call
func appendCommonOptions(parts []string, cfg template.Config, tokenizerFlag, seedFlag string) []string {
	if cfg.Host != nil {
		parts = append(parts, "--host", fmt.Sprint(*cfg.Host))
	}
	if cfg.Port != nil {
		parts = append(parts, "--port", fmt.Sprint(*cfg.Port))
	}

	if cfg.Quantization != nil {
		parts = append(parts, "--quantization", *cfg.Quantization)
	}

	if cfg.LoadFormat != nil {
		parts = append(parts, "--load-format", *cfg.LoadFormat)
	}
	if len(cfg.ModelLoaderExtraConfig) > 0 {
		configJSON, err := json.Marshal(cfg.ModelLoaderExtraConfig)
		if err == nil {
			parts = append(parts, "--model-loader-extra-config", string(configJSON))
		}
	}

	if cfg.TrustRemoteCode != nil && *cfg.TrustRemoteCode {
		parts = append(parts, "--trust-remote-code")
	}

	if cfg.Seed != nil {
		parts = append(parts, seedFlag, fmt.Sprint(*cfg.Seed))
	}
	if cfg.Tokenizer != nil {
		parts = append(parts, tokenizerFlag, *cfg.Tokenizer)
	}

	return parts
}

func generateCustomCommand(cfg template.Config) ([]string, error) {
	if cfg.Custom == nil || cfg.Custom.BaseCommand == "" {
		return nil, errors.New("custom engine requires 'custom' config with 'base_command'")
	}

	// Validate image/tag requirements for custom config
	if cfg.Custom.Image != nil && *cfg.Custom.Image != "" {
		hasTag := strings.Contains(*cfg.Custom.Image, ":")
		hasImageTag := cfg.ImageTag != nil && *cfg.ImageTag != ""

		if !hasTag && !hasImageTag {
			return nil, errors.New("custom config with 'image' must either include a tag (e.g., 'image:tag') or have 'image_tag' separately defined")
		}
	}

	parts := []string{cfg.Custom.BaseCommand}

	if cfg.Model != "" {
		modelFlag := "--model"
		if cfg.Custom.ModelFlag != nil {
			modelFlag = *cfg.Custom.ModelFlag
		}
		parts = append(parts, modelFlag, cfg.Model)
	}

	for _, arg := range cfg.Custom.Args {
		parts = append(parts, arg)
	}

	for key, value := range cfg.Custom.KVArgs {
		parts = append(parts, key, fmt.Sprint(value))
	}

	return parts, nil
}
