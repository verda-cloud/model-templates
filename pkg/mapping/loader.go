package mapping

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func getMappingsDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "mappings"
	}

	// Go up to the project root and then to mappings/
	pkgDir := filepath.Dir(filename)                  // pkg/mapping
	projectRoot := filepath.Dir(filepath.Dir(pkgDir)) // project root
	return filepath.Join(projectRoot, "mappings")
}

// LoadMappings loads mapping files for a specific engine
// Returns both engine-specific and common mappings
func LoadMappings(engine string) (*MappingFile, *MappingFile, error) {
	mappingsDir := getMappingsDir()

	// Engine-specific mappings
	enginePath := filepath.Join(mappingsDir, fmt.Sprintf("%s.json", engine))
	engineMapping, err := LoadMappingFile(enginePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load %s mappings: %w", engine, err)
	}

	// Common mappings
	commonPath := filepath.Join(mappingsDir, "common.json")
	commonMapping, err := LoadMappingFile(commonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load common mappings: %w", err)
	}

	return engineMapping, commonMapping, nil
}

func LoadMappingFile(path string) (*MappingFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read mapping file %s: %w", path, err)
	}

	var mapping MappingFile
	if err := json.Unmarshal(data, &mapping); err != nil {
		return nil, fmt.Errorf("failed to parse mapping file %s: %w", path, err)
	}

	return &mapping, nil
}
