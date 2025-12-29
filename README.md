# Model templates
Model format for defining AI models and their suggested engine options, especially for vLLM and SGLang. 

These templates are also used by Verda command line client (TODO: add link when available) that allows easily running the template.

## Supported Engines

- **sglang** - SGLang inference engine
- **vllm** - vLLM inference engine
- **custom** - Any custom inference engine

## Configuration Format

Structured JSON with common fields at the top level and engine-specific options nested. More
information about the configuration format is available in a separate [spec file](SPEC.md).

The JSON-to-CLI parameter mappings are defined in the [`mappings/`](mappings/) directory.


## Acknowledgements

This project was inspired by and partially uses model data from:
- **SGL Cookbook**: https://github.com/sgl-project/sgl-cookbook - Many example configurations are adapted from the awesome SGL Cookbook
- **InferenceMAX**: https://github.com/InferenceMAX/InferenceMAX - GPT-OSS example configurations are based on InferenceMAX benchmarks

## Resources

- **SGLang**: https://github.com/sgl-project/sglang
- **vLLM**: https://github.com/vllm-project/vllm
- **SGL Cookbook**: https://github.com/sgl-project/sgl-cookbook
- **InferenceMAX**: https://inferencemax.ai
- 