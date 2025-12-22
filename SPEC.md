# Configuration Specification

This document describes the configuration fields available for creating LLM inference server templates. Configurations are written in JSON format and support three inference engines: **SGLang**, **vLLM**, and **Custom**.

## Configuration Structure

All configurations have a common top-level structure:

```json
{
  "engine": "sglang|vllm|custom",
  "model": "model-identifier",
  "trust_remote_code": true,
  "quantization": "fp8",
  "gpu_types": ["h100", "h200"],
  "load_format": "auto",
  "model_loader_extra_config": {
    "enable_multithread_load": true,
    "num_threads": 8
  },
  "seed": 42,
  "tokenizer": "path/to/tokenizer",
  "image_tag": "v0.3.5",

  "sglang": { /* SGLang-specific options */ },
  "vllm": { /* vLLM-specific options */ },
  "custom": { /* Custom engine options */ }
}
```

---

## Field Reference

### Common Fields (All Engines)

These fields are available at the top level of all configurations.

#### `engine` (string, required)
Specifies which inference engine to use.
- **Valid values:** `"sglang"`, `"vllm"`, `"custom"`
- **Example:** `"engine": "sglang"`

#### `model` (string, required)
The model identifier or path to load.
- **Format:** HuggingFace path (e.g., `"meta-llama/Llama-3.1-70B-Instruct"`) or local path
- **Example:** `"model": "deepseek-ai/DeepSeek-V3.2"`

#### `explanation` (string, optional)
Detailed explanation of the configuration's purpose, use case, and key features.
- **Use case:** Document the configuration for reference and understanding
- **Example:** `"explanation": "Llama 3.1 70B with basic configuration using 4-way tensor parallelism. This is the baseline setup with 32k context length."`

#### `short_explanation` (string, optional)
Brief one-line description of the configuration.
- **Use case:** Quick reference for configuration summary
- **Example:** `"short_explanation": "Llama 3.1 70B basic configuration"`

#### `trust_remote_code` (boolean, optional)
Whether to trust and execute remote code from model repositories. Used by SGLang and vLLM.
- **Default:** `false`
- **Warning:** Only enable for trusted models
- **Example:** `"trust_remote_code": true`

#### `quantization` (string, optional)
Quantization method to reduce model size and memory usage. Used by SGLang and vLLM.
- **Valid values:** `"awq"`, `"fp8"`, `"fp4"`, `"int8"`, `"int4"`, etc.
- **Example:** `"quantization": "fp8"`

#### `gpu_types` (array of strings, optional)
List of compatible GPU types that can run this model configuration.
- **Valid values:** `"h100"`, `"h200"`, `"b200"`, `"a100"`, `"l40s"`, `"v100"`, `"rtx6000-ada"`, `"rtx6000-pro"`
- **Purpose:** Documents hardware requirements for deployment planning
- **Example:** `"gpu_types": ["h100", "h200", "b200"]`

#### `load_format` (string, optional)
Format for loading model weights. Used by SGLang and vLLM.
- **Valid values:**
    - `"auto"` - Auto-detect format (default)
    - `"pt"` - PyTorch checkpoint
    - `"safetensors"` - SafeTensors format
    - `"runai_streamer"` - Run:ai model streaming
    - `"runai_streamer_sharded"` - Sharded Run:ai streaming
    - `"tensorizer"` - Tensorized format
    - `"gguf"` - GGUF format
    - `"bitsandbytes"` - BitsAndBytes quantized format
    - Other values: `"npcache"`, `"dummy"`, `"sharded_state"`, `"mistral"`
- **Example:** `"load_format": "runai_streamer"`

#### `model_loader_extra_config` (object, optional)
Additional configuration for the model loader. Used by SGLang and vLLM.
- **Common options:**
    - `enable_multithread_load`: Enable multi-threaded model loading
    - `num_threads`: Number of threads for loading (SGLang)
    - `concurrency`: Number of concurrent download streams (vLLM with Run:ai)
    - `memory_limit`: Memory limit in bytes (vLLM with Run:ai)
- **Example for multi-threading:**
  ```json
  "model_loader_extra_config": {
    "enable_multithread_load": true,
    "num_threads": 64
  }
  ```
- **Example for Run:ai streaming:**
  ```json
  "model_loader_extra_config": {
    "concurrency": 16,
    "memory_limit": 5368709120
  }
  ```

#### `seed` (integer, optional)
Random seed for reproducibility. Used by SGLang and vLLM.
- **Use case:** Ensure deterministic generation across runs
- **Example:** `"seed": 42`

#### `tokenizer` (string, optional)
Override the default tokenizer. Used by SGLang and vLLM.
- **Format:** HuggingFace path or local path
- **Use case:** Use a different tokenizer than the model's default
- **SGLang flag:** `--tokenizer-path`
- **vLLM flag:** `--tokenizer`
- **Example:** `"tokenizer": "meta-llama/Llama-3.1-70B-Instruct"`

#### `image_tag` (string, optional)
Image tag to use for the container image, overriding the default tag.
- **Use case:** Specify a particular version or variant of the inference engine image
- **Example:** `"image_tag": "v0.3.5"`
- **Note:** For custom engines with the `image` field, either the image must include a tag (e.g., `"image:tag"`) or this field must be provided

### SGLang Fields

SGLang-specific configuration options are nested under the `"sglang"` key.

#### Parallelism Options

##### `tp` (integer, optional)
Tensor parallelism: splits the model across multiple GPUs horizontally.
- **Use when:** Model is too large for a single GPU
- **Typical values:** 1, 2, 4, 8
- **Example:** `"tp": 8` (split across 8 GPUs)

##### `dp` (integer, optional)
Data parallelism: creates multiple replicas to process different requests in parallel.
- **Use when:** Need higher throughput
- **Requires:** `enable_dp_attention` set to `true`
- **Example:** `"dp": 4`

##### `ep` (integer, optional)
Expert parallelism: distributes experts in MoE (Mixture of Experts) models across GPUs.
- **Use when:** Running MoE models like DeepSeek, Qwen3, GLM
- **Example:** `"ep": 8`

#### Memory Management

##### `mem_fraction_static` (float, optional)
Fraction of GPU memory to allocate statically for model weights and KV cache.
- **Range:** 0.0 to 1.0
- **Default:** Engine default (typically ~0.85)
- **Use case:** Increase for throughput, decrease if running out of memory
- **Example:** `"mem_fraction_static": 0.90`

##### `kv_cache_dtype` (string, optional)
Data type for the key-value cache in attention layers.
- **Valid values:** `"fp8_e4m3"`, `"bf16"`, `"fp16"`
- **Trade-off:** `fp8_e4m3` uses less memory but may reduce quality slightly
- **Example:** `"kv_cache_dtype": "fp8_e4m3"`

#### Model Configuration

##### `context_length` (integer, optional)
Maximum context window length in tokens.
- **Use when:** Need to override model's default context length
- **Example:** `"context_length": 32768`

##### `max_running_requests` (integer, optional)
Maximum number of requests the server will process concurrently.
- **Use case:** Limit concurrency to control memory usage
- **Example:** `"max_running_requests": 64`

#### Features

##### `enable_dp_attention` (boolean, optional)
Enable distributed attention for data parallelism.
- **Required when:** Using `dp` > 1
- **Example:** `"enable_dp_attention": true`

##### `reasoning_parser` (string, optional)
Parser for extracting reasoning/thinking tokens from model outputs.
- **Valid values:** `"deepseek-v3"`, `"deepseek-r1"`, `"qwen3"`, `"nano_v3"`, `"glm45"`, `"kimi"`
- **Use with:** Models that support chain-of-thought reasoning
- **Example:** `"reasoning_parser": "deepseek-v3"`

##### `tool_call_parser` (string, optional)
Parser for extracting tool/function calls from model outputs.
- **Valid values:** `"deepseekv32"`, `"qwen3_coder"`, `"glm45"`, `"kimi"`
- **Use with:** Models that support function calling
- **Example:** `"tool_call_parser": "deepseekv32"`

##### `chat_template` (string, optional)
Path to a custom Jinja chat template file.
- **Use when:** Need to override the model's default chat format
- **Example:** `"chat_template": "/path/to/template.jinja"`

##### `disable_shared_experts_fusion` (boolean, optional)
Disable optimization for shared experts in MoE models.
- **Use when:** Experiencing issues with expert fusion
- **Example:** `"disable_shared_experts_fusion": true`

#### Speculative Decoding (Multi-Token Prediction)

##### `speculative_algorithm` (string, optional)
Algorithm for speculative decoding to improve latency.
- **Valid values:** `"EAGLE"`
- **Example:** `"speculative_algorithm": "EAGLE"`

##### `speculative_num_steps` (integer, optional)
Number of speculative decoding steps per forward pass.
- **Typical value:** 3
- **Example:** `"speculative_num_steps": 3`

##### `speculative_eagle_topk` (integer, optional)
Top-k parameter for EAGLE speculative decoding.
- **Example:** `"speculative_eagle_topk": 4`

##### `speculative_num_draft_tokens` (integer, optional)
Number of draft tokens to generate in speculative decoding.
- **Example:** `"speculative_num_draft_tokens": 4`

##### `speculative_draft_model_path` (string, optional)
Path or identifier for the draft model used in speculative decoding.
- **Use case:** Specify a smaller, faster model to generate draft tokens
- **Example:** `"speculative_draft_model_path": "lmsys/sglang-EAGLE3-Llama-4-Scout-17B-16E-Instruct-v1"`

#### Pipeline & Scheduling

##### `pipeline_parallel_size` (integer, optional)
Pipeline parallelism: splits model layers vertically across GPUs.
- **Use when:** Need additional parallelism beyond tensor parallelism
- **SGLang flag:** `--pp`
- **Example:** `"pipeline_parallel_size": 2`

##### `schedule_policy` (string, optional)
Scheduling strategy for request processing.
- **Valid values:** `"fcfs"` (first come first served), `"lpm"` (longest prompt first), `"dfs-weight"`
- **Default:** `"fcfs"`
- **Example:** `"schedule_policy": "lpm"`

#### Cache & Optimization

##### `disable_radix_cache` (boolean, optional)
Disable radix cache (prefix caching) optimization.
- **Use when:** Troubleshooting caching issues
- **Default:** `false` (caching enabled)
- **Example:** `"disable_radix_cache": true`

##### `disable_cuda_graph` (boolean, optional)
Disable CUDA graph optimization.
- **Use when:** Debugging or when CUDA graphs cause issues
- **Default:** `false` (CUDA graphs enabled)
- **Example:** `"disable_cuda_graph": true`

### vLLM Fields

vLLM-specific configuration options are nested under the `"vllm"` key.

#### Parallelism Options

##### `tensor_parallel_size` (integer, optional)
Tensor parallelism: splits the model across multiple GPUs horizontally.
- **Use when:** Model is too large for a single GPU
- **Typical values:** 1, 2, 4, 8
- **Example:** `"tensor_parallel_size": 4`

##### `pipeline_parallel_size` (integer, optional)
Pipeline parallelism: splits the model across multiple GPUs vertically (by layers).
- **Use when:** Need additional parallelism beyond tensor parallelism
- **Less common:** Most deployments use only tensor parallelism
- **Example:** `"pipeline_parallel_size": 2`

#### Model Configuration

##### `dtype` (string, optional)
Data type for model weights and computation.
- **Valid values:** `"auto"`, `"bfloat16"`, `"float16"`, `"float32"`
- **Default:** `"auto"` (infers from model)
- **Trade-offs:** `bfloat16` has better numerical stability than `float16`
- **Example:** `"dtype": "bfloat16"`

##### `max_model_len` (integer, optional)
Maximum context length in tokens.
- **Use when:** Need to override model's default max length or limit memory usage
- **Example:** `"max_model_len": 32768`

#### Memory Management

##### `gpu_memory_utilization` (float, optional)
Fraction of GPU memory to allocate for model execution.
- **Range:** 0.0 to 1.0
- **Default:** 0.90
- **Trade-off:** Higher values allow more concurrent requests but risk OOM
- **Example:** `"gpu_memory_utilization": 0.95`

##### `max_num_seqs` (integer, optional)
Maximum number of sequences (requests) to process in parallel.
- **Use case:** Control concurrency and memory usage
- **Example:** `"max_num_seqs": 256`

#### Features

##### `enable_prefix_caching` (boolean, optional)
Cache and reuse computation for common prompt prefixes.
- **Use when:** Many requests share the same system prompt or prefix
- **Benefit:** Reduces latency and improves throughput
- **Example:** `"enable_prefix_caching": true`

##### `enable_chunked_prefill` (boolean, optional)
Process long prompts in chunks to reduce latency for first token.
- **Use when:** Serving requests with very long prompts
- **Benefit:** Better time-to-first-token (TTFT)
- **Example:** `"enable_chunked_prefill": true`

##### `enforce_eager` (boolean, optional)
Force eager execution mode instead of using CUDA graphs.
- **Use when:** Debugging or when CUDA graphs cause issues
- **Trade-off:** Slower but more flexible
- **Example:** `"enforce_eager": true`

##### `served_model_name` (string, optional)
Override the model name returned by the API.
- **Use case:** Expose a different model name to clients
- **Example:** `"served_model_name": "my-custom-model"`

#### Parallelism & Distribution

##### `data_parallel_size` (integer, optional)
Number of data parallel replicas.
- **Use when:** Need to scale throughput with multiple model replicas
- **Example:** `"data_parallel_size": 2`

##### `distributed_executor_backend` (string, optional)
Backend for distributed execution.
- **Valid values:** `"ray"`, `"mp"` (multiprocessing), `"uni"`, `"external_launcher"`
- **Use case:** Choose distributed computing framework
- **Example:** `"distributed_executor_backend": "ray"`

#### Memory & Performance

##### `swap_space` (integer, optional)
CPU swap space per GPU in GiB.
- **Use when:** Need additional memory beyond GPU VRAM
- **Default:** 4
- **Example:** `"swap_space": 8`

##### `max_num_batched_tokens` (integer, optional)
Maximum number of tokens processed per iteration.
- **Use case:** Control batch size and memory usage
- **Example:** `"max_num_batched_tokens": 8192`

#### Scheduling

##### `scheduling_policy` (string, optional)
Request scheduling policy.
- **Valid values:** `"fcfs"` (first come first served), `"priority"`
- **Default:** `"fcfs"`
- **Example:** `"scheduling_policy": "priority"`

### Custom Engine Fields

Custom engine configuration is for running inference with non-standard or proprietary engines. Options are nested under the `"custom"` key.

##### `base_command` (string, required)
The base command to execute your custom inference engine.
- **Required:** Must be specified when using `"engine": "custom"`
- **Example:** `"base_command": "python -m my_inference_engine.serve"`

##### `model_flag` (string, optional)
The command-line flag used to specify the model path.
- **Default:** `"--model"`
- **Use when:** Your engine uses a different flag for the model argument
- **Example:** `"model_flag": "--model-path"`

##### `image` (string, optional)
The full container image to use for the custom inference engine.
- **Format:** Image name with optional registry and tag (e.g., `"myregistry.io/my-engine:v1.0"`)
- **Use when:** Specifying the complete container image for your custom engine
- **Validation:** Either the image must include a tag (e.g., `"image:tag"`) or the top-level `image_tag` field must be provided
- **Example:** `"image": "myregistry.io/custom-inference:latest"`

##### `args` (array of strings, optional)
List of additional command-line arguments to pass to your engine.
- **Format:** Each argument as a string (can include both flag and value in one string)
- **Example:**
  ```json
  "args": [
    "--workers 4",
    "--batch-size 32",
    "--max-tokens 2048"
  ]
  ```

##### `kv_args` (object, optional)
Key-value pairs for command-line arguments.
- **Format:** Key is the flag, value is the argument value
- **Use when:** Prefer structured key-value format over string arrays
- **Example:**
  ```json
  "kv_args": {
    "--timeout": 60,
    "--max-tokens": 2048,
    "--temperature": 0.7
  }
  ```

---

## JSON Examples

### SGLang - Minimal
```json
{
  "engine": "sglang",
  "model": "meta-llama/Llama-3.1-8B-Instruct",
  "sglang": {
    "tp": 1
  }
}
```

### SGLang - Full Featured
```json
{
  "engine": "sglang",
  "model": "deepseek-ai/DeepSeek-V3.2",
  "gpu_types": ["h200", "b200"],
  "model_loader_extra_config": {
    "enable_multithread_load": true,
    "num_threads": 64
  },
  "sglang": {
    "tp": 8,
    "dp": 4,
    "ep": 2,
    "mem_fraction_static": 0.90,
    "kv_cache_dtype": "fp8_e4m3",
    "enable_dp_attention": true,
    "reasoning_parser": "deepseek-v3",
    "tool_call_parser": "deepseekv32"
  }
}
```

### vLLM - Minimal
```json
{
  "engine": "vllm",
  "model": "mistralai/Mistral-7B-Instruct-v0.3",
  "vllm": {
    "tensor_parallel_size": 1
  }
}
```

### vLLM - Full Featured
```json
{
  "engine": "vllm",
  "model": "meta-llama/Llama-3.1-70B-Instruct",
  "quantization": "awq",
  "gpu_types": ["h100", "h200", "b200"],
  "load_format": "runai_streamer",
  "model_loader_extra_config": {
    "concurrency": 8
  },
  "vllm": {
    "tensor_parallel_size": 4,
    "dtype": "bfloat16",
    "max_model_len": 32768,
    "gpu_memory_utilization": 0.95,
    "max_num_seqs": 512,
    "enable_prefix_caching": true,
    "enable_chunked_prefill": true
  }
}
```

### SGLang - With Multithread Model Loading
```json
{
  "engine": "sglang",
  "model": "deepseek-ai/DeepSeek-V3",
  "trust_remote_code": true,
  "gpu_types": ["h200", "b200"],
  "model_loader_extra_config": {
    "enable_multithread_load": true,
    "num_threads": 64
  },
  "sglang": {
    "tp": 8
  }
}
```

### vLLM - With RunAI Model Streamer
```json
{
  "engine": "vllm",
  "model": "meta-llama/Llama-3.1-70B-Instruct",
  "gpu_types": ["h100", "h200"],
  "load_format": "runai_streamer",
  "model_loader_extra_config": {
    "concurrency": 16,
    "memory_limit": 5368709120
  },
  "vllm": {
    "tensor_parallel_size": 4
  }
}
```

### vLLM - With Multithread Model Loading
```json
{
  "engine": "vllm",
  "model": "mistralai/Mistral-7B-Instruct-v0.3",
  "gpu_types": ["h100", "h200", "b200", "a100", "l40s", "v100"],
  "model_loader_extra_config": {
    "enable_multithread_load": true
  },
  "vllm": {
    "tensor_parallel_size": 1
  }
}
```

### Custom Engine
```json
{
  "engine": "custom",
  "model": "my-org/my-model",
  "gpu_types": ["h100", "h200", "b200"],
  "image_tag": "v1.0",
  "custom": {
    "base_command": "python -m my_engine.serve",
    "model_flag": "--model-path",
    "image": "myregistry.io/custom-inference",
    "args": [
      "--workers 4",
      "--batch-size 32"
    ],
    "kv_args": {
      "--timeout": 60,
      "--max-tokens": 2048
    }
  }
}
```

### KimiK2 - Reasoning Model
```json
{
  "engine": "sglang",
  "model": "Kimi-K2-Instruct",
  "gpu_types": ["h200", "b200"],
  "model_loader_extra_config": {
    "enable_multithread_load": true,
    "num_threads": 32
  },
  "sglang": {
    "tp": 8,
    "dp": 4,
    "ep": 4,
    "reasoning_parser": "kimi",
    "tool_call_parser": "kimi"
  }
}
```

---

## Usage

### Generating Commands

Use `cmd-gen.go` to generate inference server startup commands from your configuration:

```bash
go run cmd-gen.go examples/deepseek-sglang.json
```

This will output the complete command-line invocation for starting the inference server with all configured options.

### Configuration Best Practices

1. **Start with examples:** Use the provided example configurations as templates
2. **Match GPU types:** Ensure your `gpu_types` field matches your available hardware
3. **Test incrementally:** Start with minimal configs, then add optimizations
4. **Memory management:** Adjust `mem_fraction_static` or `gpu_memory_utilization` based on your memory constraints
5. **Parallelism:** Use `tp` (tensor parallelism) for model size, `dp` for throughput
6. **Model-specific parsers:** Enable `reasoning_parser` and `tool_call_parser` only for models that support them
