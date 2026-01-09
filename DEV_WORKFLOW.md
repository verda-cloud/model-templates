A good workflow would be: 
1. [Provision a GPU instance on Verda cloud](https://docs.verda.com/cpu-and-gpu-instances/set-up-a-gpu-instance)
2. [SSH into the instance](https://docs.verda.com/cpu-and-gpu-instances/connecting-to-your-datacrunch.io-server)
3. Running the container with flags generated from the template inside the instance

> [!IMPORTANT]
> If the template has tensor parallel size set to 2 ensure the provisioned instance has the required number of GPUs

Example: 
For the [`llama-vllm.json`](./templates/llama-vllm.json) template first get the correct flags on your local machine by running:

```bash
go run examples/mapping_usage.go templates/llama-vllm.json
```

which should output

```bash
Loading mappings for engine: vllm

Generated Command:
--model meta-llama/Llama-3.1-8B-Instruct --gpu-memory-utilization 0.95 --max-model-len 8192 --tensor-parallel-size 2 --dtype bfloat16 --max-num-seqs 256 --enable-prefix-caching --enable-chunked-prefill

Formatted Command:
--model meta-llama/Llama-3.1-8B-Instruct \
  --gpu-memory-utilization 0.95 \
  --max-model-len 8192 \
  --tensor-parallel-size 2 \
  --dtype bfloat16 \
  --max-num-seqs 256 \
  --enable-prefix-caching \
  --enable-chunked-prefill
```

> [!TIP]
> For vllm we use `"vllm/vllm-openai:v0.13.0"` docker image (see [`deploy.go`](./cmd/deploy/deploy.go)), 

Then, in the SSH instance run:
```bash
docker run --runtime nvidia --gpus all -p 8000:8000     --ipc=host     vllm/vllm-openai:v0.13.0 --model meta-llama/Llama-3.1-8B-Instruct --gpu-memory-utilization 0.95 --max-model-len 8192 --tensor-parallel-size 2 --dtype bfloat16 --max-num-seqs 256 --enable-prefix-caching --enable-chunked-prefill
```

> [!TIP]
> To use HuggingFace gates models provide a `$HF_TOKEN` via `--env "HF_TOKEN=$HF_TOKEN"`

If all goes well this should launch a vLLM server.

Once the template works correctly you can deploy to the Serverless Containers service.