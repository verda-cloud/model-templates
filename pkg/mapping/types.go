package mapping

type MappingFile struct {
	Version     string                  `json:"version"`
	Engine      string                  `json:"engine,omitempty"`
	Description string                  `json:"description,omitempty"`
	BaseCommand []string                `json:"base_command,omitempty"`
	Mappings    map[string]MappingEntry `json:"mappings"`
}

type MappingEntry struct {
	JSONPath    string      `json:"json_path"`
	CLIFlag     interface{} `json:"cli_flag"` // Can be string or map[string]string for engine-specific flags
	Type        string      `json:"type"`
	Format      string      `json:"format,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Notes       string      `json:"notes,omitempty"`
	Transform   string      `json:"transform,omitempty"`
	Example     interface{} `json:"example,omitempty"`
}

// GetCLIFlag returns the CLI flag for a given engine
// If the flag is engine-specific, it returns the flag for the specified engine
// Otherwise, it returns the single flag string
func (m *MappingEntry) GetCLIFlag(engine string) string {
	switch v := m.CLIFlag.(type) {
	case string:
		return v
	case map[string]interface{}:
		if flag, ok := v[engine]; ok {
			if flagStr, ok := flag.(string); ok {
				return flagStr
			}
		}
		return ""
	default:
		return ""
	}
}

// SupportedTypes lists all supported parameter types
const (
	TypeString         = "string"
	TypeInt            = "int"
	TypeFloat          = "float"
	TypeBool           = "bool"
	TypeBoolFlag       = "bool_flag"
	TypeJSONObject     = "json_object"
	TypeArrayOfStrings = "array_of_strings"
	TypeObject         = "object"
)

const (
	TransformJSONMarshal = "json_marshal"
)
