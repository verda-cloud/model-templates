package mapping

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/verda-cloud/model-templates/pkg/template"
)

// CommandBuilder builds CLI commands from configuration using mapping files
type CommandBuilder struct {
	engineMapping *MappingFile
	commonMapping *MappingFile
	engine        string
}

// NewCommandBuilder creates a new command builder for the specified engine
func NewCommandBuilder(engine string) (*CommandBuilder, error) {
	engineMapping, commonMapping, err := LoadMappings(engine)
	if err != nil {
		return nil, err
	}

	return &CommandBuilder{
		engineMapping: engineMapping,
		commonMapping: commonMapping,
		engine:        engine,
	}, nil
}

// Build constructs a command from the configuration struct
func (b *CommandBuilder) Build(config template.Config) ([]string, error) {
	if b.engine == "custom" {
		return b.buildCustomCommand(config)
	}

	var parts []string

	if len(b.engineMapping.BaseCommand) > 0 {
		parts = append(parts, b.engineMapping.BaseCommand...)
	}

	if config.Host != nil && *config.Host != "" {
		parts = append(parts, "--host", *config.Host)
	}
	if config.Port != nil {
		parts = append(parts, "--port", fmt.Sprint(*config.Port))
	}
	if config.Model != "" {
		parts = append(parts, "--model", config.Model)
	}

	engineParts, err := b.processMappings(config, b.engineMapping.Mappings)
	if err != nil {
		return nil, fmt.Errorf("error processing engine mappings: %w", err)
	}
	parts = append(parts, engineParts...)

	commonParts, err := b.processMappings(config, b.commonMapping.Mappings)
	if err != nil {
		return nil, fmt.Errorf("error processing common mappings: %w", err)
	}
	parts = append(parts, commonParts...)

	return parts, nil
}

// buildCustomCommand handles custom engine configuration
func (b *CommandBuilder) buildCustomCommand(config template.Config) ([]string, error) {
	if config.Custom == nil {
		return nil, fmt.Errorf("custom engine requires 'custom' configuration object")
	}

	customConfig := config.Custom

	// Get base command
	if customConfig.BaseCommand == "" {
		return nil, fmt.Errorf("custom engine requires 'base_command'")
	}

	// Split base command into parts
	parts := strings.Fields(customConfig.BaseCommand)

	// Add model with model flag
	if config.Model != "" {
		modelFlag := "--model"
		if customConfig.ModelFlag != nil && *customConfig.ModelFlag != "" {
			modelFlag = *customConfig.ModelFlag
		}
		parts = append(parts, modelFlag, config.Model)
	}

	// Add args array
	for _, arg := range customConfig.Args {
		// Split each arg string (e.g., "--workers 4" becomes "--workers", "4")
		parts = append(parts, strings.Fields(arg)...)
	}

	// Add kv_args
	for key, value := range customConfig.KVArgs {
		parts = append(parts, key, fmt.Sprint(value))
	}

	return parts, nil
}

// processMappings processes a set of mappings and returns command parts
func (b *CommandBuilder) processMappings(config template.Config, mappings map[string]MappingEntry) ([]string, error) {
	var parts []string

	for _, entry := range mappings {
		// This is handled previously so skipped here
		if entry.JSONPath == "model" {
			continue
		}

		value := b.getValueByJSONPath(config, entry.JSONPath)
		if value == nil {
			if entry.Required {
				return nil, fmt.Errorf("required parameter %s is missing", entry.JSONPath)
			}
			continue
		}

		cliFlag := entry.GetCLIFlag(b.engine)
		if cliFlag == "" {
			continue
		}

		formatted, skip, err := b.formatValue(value, entry)
		if err != nil {
			return nil, fmt.Errorf("error formatting %s: %w", entry.JSONPath, err)
		}
		if skip {
			continue
		}

		switch entry.Type {
		case TypeBoolFlag:
			// Boolean flags are added without a value
			parts = append(parts, cliFlag)
		default:
			// Other types have flag + value
			parts = append(parts, cliFlag, formatted)
		}
	}

	return parts, nil
}

// getValueByJSONPath extracts a value from the Config struct using JSON path notation
// Example: "vllm.tensor_parallel_size" -> config.VLLM.TensorParallelSize
func (b *CommandBuilder) getValueByJSONPath(config template.Config, jsonPath string) interface{} {
	parts := strings.Split(jsonPath, ".")

	value := reflect.ValueOf(config)

	for _, part := range parts {
		field := b.findFieldByJSONTag(value, part)
		if !field.IsValid() {
			return nil
		}

		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return nil
			}
			field = field.Elem()
		}

		value = field
	}

	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil
		}
		return value.Elem().Interface()
	}

	return value.Interface()
}

func (b *CommandBuilder) findFieldByJSONTag(value reflect.Value, jsonTag string) reflect.Value {
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return reflect.Value{}
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)

		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}

		// Remove options like ",omitempty"
		tagName := strings.Split(tag, ",")[0]

		if tagName == jsonTag {
			return value.Field(i)
		}
	}

	return reflect.Value{}
}

func (b *CommandBuilder) formatValue(value interface{}, entry MappingEntry) (string, bool, error) {
	if value == nil {
		return "", true, nil
	}

	switch entry.Type {
	case TypeString:
		if v, ok := value.(string); ok {
			return v, false, nil
		}
		return fmt.Sprint(value), false, nil

	case TypeInt:
		switch v := value.(type) {
		case int:
			return fmt.Sprint(v), false, nil
		case int64:
			return fmt.Sprint(v), false, nil
		case float64:
			return fmt.Sprint(int(v)), false, nil
		default:
			return "", true, fmt.Errorf("expected int, got %T", value)
		}

	case TypeFloat:
		var floatVal float64
		switch v := value.(type) {
		case float64:
			floatVal = v
		case float32:
			floatVal = float64(v)
		case int:
			floatVal = float64(v)
		default:
			return "", true, fmt.Errorf("expected float, got %T", value)
		}

		if entry.Format != "" {
			return fmt.Sprintf(entry.Format, floatVal), false, nil
		}
		return fmt.Sprint(floatVal), false, nil

	case TypeBool:
		if v, ok := value.(bool); ok {
			return fmt.Sprint(v), false, nil
		}
		return "", true, fmt.Errorf("expected bool, got %T", value)

	case TypeBoolFlag:
		// Only include the flag if the value is true
		if v, ok := value.(bool); ok && v {
			return "", false, nil
		}
		return "", true, nil

	case TypeJSONObject:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Map && v.Len() == 0 {
			return "", true, nil
		}

		if entry.Transform == TransformJSONMarshal {
			jsonBytes, err := json.Marshal(value)
			if err != nil {
				return "", true, fmt.Errorf("failed to marshal JSON: %w", err)
			}
			return string(jsonBytes), false, nil
		}
		return fmt.Sprint(value), false, nil

	case TypeArrayOfStrings:
		if reflect.TypeOf(value).Kind() == reflect.Slice {
			slice := reflect.ValueOf(value)
			var strSlice []string
			for i := 0; i < slice.Len(); i++ {
				strSlice = append(strSlice, fmt.Sprint(slice.Index(i).Interface()))
			}
			return strings.Join(strSlice, " "), false, nil
		}
		return fmt.Sprint(value), false, nil

	case TypeObject:
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return "", true, fmt.Errorf("failed to marshal object: %w", err)
		}
		return string(jsonBytes), false, nil

	default:
		return fmt.Sprint(value), false, nil
	}
}
