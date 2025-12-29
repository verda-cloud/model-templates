package engine

import (
	"fmt"
	"strings"

	"github.com/verda-cloud/model-templates/pkg/mapping"
	"github.com/verda-cloud/model-templates/pkg/template"
)

func GenerateCommandFromTemplate(cfg template.Config) ([]string, error) {
	engine := strings.ToLower(string(cfg.Engine))

	switch template.EngineOption(engine) {
	case template.EngineVLLM, template.EngineSGLang, template.EngineCustom:
		// Valid engine
	default:
		return nil, fmt.Errorf("unknown engine option: %s", string(cfg.Engine))
	}

	builder, err := mapping.NewCommandBuilder(engine)
	if err != nil {
		return nil, fmt.Errorf("failed to create command builder: %w", err)
	}

	command, err := builder.Build(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build command: %w", err)
	}

	return command, nil
}
