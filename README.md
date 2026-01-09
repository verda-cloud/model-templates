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

## ⚠️ Disclaimer

The templates in the [`templates/`](templates/) directory have **not been fully verified**. While they are based on recommended configurations from various sources, they may require adjustments for your specific hardware, software versions, or use cases. Please test and debug serverless deployments thoroughly.

Take a look at [DEV_WORKFLOW.md](./DEV_WORKFLOW.md) for a way to test templates on a GPU instance before deploying to a Serverless Container.


## Acknowledgements

This project was inspired by and partially uses model data from:
- **SGL Cookbook**: https://github.com/sgl-project/sgl-cookbook - Many example configurations are adapted from the awesome SGL Cookbook
- **InferenceMAX**: https://github.com/InferenceMAX/InferenceMAX - GPT-OSS example configurations are based on InferenceMAX benchmarks

## Resources

### LLM

- **SGLang**: https://github.com/sgl-project/sglang
- **vLLM**: https://github.com/vllm-project/vllm
- **SGL Cookbook**: https://github.com/sgl-project/sgl-cookbook
- **InferenceMAX**: https://inferencemax.ai

### Verda

- **Containers Docs**: https://docs.verda.com/containers/overview
- **API Docs**: https://docs.verda.com
- **GO SDK**: https://github.com/verda-cloud/verdacloud-sdk-go
- **Python SDK**: https://github.com/verda-cloud/sdk-python