# Examples

This directory contains example programs demonstrating how to use the model template mapping system in different languages.

## Python: translate_config.py

Python implementation that reads the mapping files from `mappings/` to translate model configurations into CLI commands.

### Usage

```bash
python3 examples/translate_config.py <config-file.json>
```

### Examples

**SGLang:**
```bash
python3 examples/translate_config.py templates/qwen3-sglang.json
```

**vLLM:**
```bash
python3 examples/translate_config.py templates/llama-vllm.json
```

**Custom:**
```bash
python3 examples/translate_config.py templates/custom-engine.json
```

## Go: mapping_usage.go

Go implementation that uses the `pkg/mapping` package with reflection on `template.Config` structs.

### Usage

```bash
go run examples/mapping_usage.go <config-file.json>
```

### Examples

**SGLang:**
```bash
go run examples/mapping_usage.go templates/qwen3-sglang.json
```

**vLLM:**
```bash
go run examples/mapping_usage.go templates/llama-vllm.json
```

**Custom:**
```bash
go run examples/mapping_usage.go templates/custom-engine.json
```

Both implementations produce identical output by reading the same mapping files.

## See Also

- [MAPPING_SPEC.md](../MAPPING_SPEC.md) - List of all available mappings and their parameters
- [SPEC.md](../SPEC.md) - Model configuration JSON specification
- [mappings/](../mappings/) - Declarative mapping files (single source of truth)
