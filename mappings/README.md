# Template Mappings

This directory contains JSON mapping files that translate model configuration fields to inference engine CLI arguments. These mappings are the single source of truth used across all programming languages (Go, Python, etc.) to generate engine commands.

## Files

- `common.json` - Shared parameters across all engines (model, host, port, quantization, etc.)
- `vllm.json` - vLLM-specific parameters (tensor parallelism, memory management, etc.)
- `sglang.json` - SGLang-specific parameters (TP/DP/EP, memory, speculative decoding, etc.)
- `custom.json` - Template for custom inference engines
- `schema.json` - JSON Schema for validating mapping files

## Adding a New Field

When adding a new parameter mapping, you need to edit **two files**:

1. **The appropriate mapping file** (`common.json`, `vllm.json`, or `sglang.json`)
2. **The Go struct** in `pkg/template/config.go`

### Example: Adding a new vLLM parameter

**1. Add to `vllm.json`:**

```json
{
  "mappings": {
    "vllm.new_parameter": {
      "json_path": "vllm.new_parameter",
      "cli_flag": "--new-parameter",
      "type": "int",
      "description": "Description of what this parameter does"
    }
  }
}
```

**2. Add to `pkg/template/config.go`:**

```go
type VLLMConfig struct {
    // ... existing fields ...
    NewParameter *int `json:"new_parameter,omitempty"`
}
```

The mapping types (`string`, `int`, `float`, `bool`, `bool_flag`, `json_object`, `array_of_strings`, `object`) are defined in `schema.json`. Parameters with engine-specific CLI flags can use an object for `cli_flag` (see `seed` in `common.json` for an example).

## Mapping Type Examples

### `string`
Simple string value passed directly to CLI.
```json
"vllm.dtype": {
  "json_path": "vllm.dtype",
  "cli_flag": "--dtype",
  "type": "string"
}
```
Input: `{"vllm": {"dtype": "bfloat16"}}` → Output: `--dtype bfloat16`

### `int`
Integer value.
```json
"port": {
  "json_path": "port",
  "cli_flag": "--port",
  "type": "int"
}
```
Input: `{"port": 8000}` → Output: `--port 8000`

### `float`
Floating-point value with optional format string.
```json
"vllm.gpu_memory_utilization": {
  "json_path": "vllm.gpu_memory_utilization",
  "cli_flag": "--gpu-memory-utilization",
  "type": "float",
  "format": "%.2f"
}
```
Input: `{"vllm": {"gpu_memory_utilization": 0.9}}` → Output: `--gpu-memory-utilization 0.90`

### `bool`
Boolean value that outputs "true" or "false" string.
```json
"some_feature": {
  "json_path": "some_feature",
  "cli_flag": "--some-feature",
  "type": "bool"
}
```
Input: `{"some_feature": true}` → Output: `--some-feature true`

### `bool_flag`
Boolean flag only included when true (no value).
```json
"trust_remote_code": {
  "json_path": "trust_remote_code",
  "cli_flag": "--trust-remote-code",
  "type": "bool_flag"
}
```
Input: `{"trust_remote_code": true}` → Output: `--trust-remote-code`
Input: `{"trust_remote_code": false}` → Output: *(nothing)*

### `json_object`
Object marshaled to JSON string.
```json
"model_loader_extra_config": {
  "json_path": "model_loader_extra_config",
  "cli_flag": "--model-loader-extra-config",
  "type": "json_object",
  "transform": "json_marshal"
}
```
Input: `{"model_loader_extra_config": {"key": "value"}}` → Output: `--model-loader-extra-config '{"key":"value"}'`

### `array_of_strings`
Array where each element is passed as a separate argument with the same flag.
```json
"allowed_models": {
  "json_path": "allowed_models",
  "cli_flag": "--allowed-model",
  "type": "array_of_strings"
}
```
Input: `{"allowed_models": ["model1", "model2"]}` → Output: `--allowed-model model1 --allowed-model model2`

### `object`
Generic object/map.
```json
"custom_config": {
  "json_path": "custom_config",
  "cli_flag": "--config",
  "type": "object"
}
```

### Engine-specific CLI flags
Use an object for `cli_flag` when different engines use different flag names.
```json
"seed": {
  "json_path": "seed",
  "cli_flag": {
    "vllm": "--seed",
    "sglang": "--random-seed"
  },
  "type": "int"
}
```
For vLLM: `{"seed": 42}` → `--seed 42`
For SGLang: `{"seed": 42}` → `--random-seed 42`
