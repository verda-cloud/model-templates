package engine

import (
	"fmt"

	"github.com/verda-cloud/model-templates/pkg/template"
)

func generateVLLMCommand(cfg template.Config) []string {
	parts := []string{""} // vLLM container runs base command by default

	if cfg.Model != "" {
		parts = append(parts, "--model", cfg.Model)
	}

	// vLLM-specific options
	if cfg.VLLM != nil {
		vllm := cfg.VLLM

		if vllm.TensorParallelSize != nil {
			parts = append(parts, "--tensor-parallel-size", fmt.Sprint(*vllm.TensorParallelSize))
		}
		if vllm.PipelineParallelSize != nil {
			parts = append(parts, "--pipeline-parallel-size", fmt.Sprint(*vllm.PipelineParallelSize))
		}

		if vllm.Dtype != nil {
			parts = append(parts, "--dtype", *vllm.Dtype)
		}
		if vllm.MaxModelLen != nil {
			parts = append(parts, "--max-model-len", fmt.Sprint(*vllm.MaxModelLen))
		}

		if vllm.GPUMemoryUtilization != nil {
			parts = append(parts, "--gpu-memory-utilization", fmt.Sprintf("%.2f", *vllm.GPUMemoryUtilization))
		}
		if vllm.MaxNumSeqs != nil {
			parts = append(parts, "--max-num-seqs", fmt.Sprint(*vllm.MaxNumSeqs))
		}

		if vllm.EnablePrefixCaching != nil && *vllm.EnablePrefixCaching {
			parts = append(parts, "--enable-prefix-caching")
		}
		if vllm.EnableChunkedPrefill != nil && *vllm.EnableChunkedPrefill {
			parts = append(parts, "--enable-chunked-prefill")
		}
		if vllm.EnforceEager != nil && *vllm.EnforceEager {
			parts = append(parts, "--enforce-eager")
		}

		if vllm.ServedModelName != nil {
			parts = append(parts, "--served-model-name", *vllm.ServedModelName)
		}

		if vllm.DataParallelSize != nil {
			parts = append(parts, "--data-parallel-size", fmt.Sprint(*vllm.DataParallelSize))
		}
		if vllm.DistributedExecutorBackend != nil {
			parts = append(parts, "--distributed-executor-backend", *vllm.DistributedExecutorBackend)
		}

		if vllm.SwapSpace != nil {
			parts = append(parts, "--swap-space", fmt.Sprint(*vllm.SwapSpace))
		}
		if vllm.MaxNumBatchedTokens != nil {
			parts = append(parts, "--max-num-batched-tokens", fmt.Sprint(*vllm.MaxNumBatchedTokens))
		}

		if vllm.SchedulingPolicy != nil {
			parts = append(parts, "--scheduling-policy", *vllm.SchedulingPolicy)
		}
	}

	// Add common options shared with SGLang
	parts = appendCommonOptions(parts, cfg, "--tokenizer", "--seed")

	return parts
}
