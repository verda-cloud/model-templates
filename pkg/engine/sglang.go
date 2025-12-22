package engine

import (
	"fmt"

	"github.com/verda-cloud/model-templates/pkg/template"
)

func generateSGLangCommand(cfg template.Config) []string {
	parts := []string{"python3", "-m", "sglang.launch_server"}

	if cfg.Model != "" {
		parts = append(parts, "--model", cfg.Model)
	}

	// SGLang-specific options
	if cfg.SGLang != nil {
		sg := cfg.SGLang

		if sg.TP != nil {
			parts = append(parts, "--tp", fmt.Sprint(*sg.TP))
		}
		if sg.DP != nil {
			parts = append(parts, "--dp", fmt.Sprint(*sg.DP))
		}
		if sg.EP != nil {
			parts = append(parts, "--ep", fmt.Sprint(*sg.EP))
		}

		if sg.MemFractionStatic != nil {
			parts = append(parts, "--mem_fraction_static", fmt.Sprint(*sg.MemFractionStatic))
		}
		if sg.KVCacheDtype != nil {
			parts = append(parts, "--kv_cache_type", fmt.Sprint(*sg.KVCacheDtype))
		}

		if sg.ContextLength != nil {
			parts = append(parts, "--context_length", fmt.Sprint(*sg.ContextLength))
		}
		if sg.MaxRunningRequests != nil {
			parts = append(parts, "--max_running_requests", fmt.Sprint(*sg.MaxRunningRequests))
		}

		if sg.EnableDPAttention != nil && *sg.EnableDPAttention {
			parts = append(parts, "--enable-dp-attention")
		}
		if sg.DisableSharedExpertsFusion != nil && *sg.DisableSharedExpertsFusion {
			parts = append(parts, "--disable-shared-experts-fusion")
		}

		if sg.ReasoningParser != nil {
			parts = append(parts, "--reasoning_parser", fmt.Sprint(*sg.ReasoningParser))
		}
		if sg.ToolCallParser != nil {
			parts = append(parts, "--tool_call_parser", fmt.Sprint(*sg.ToolCallParser))
		}
		if sg.ChatTemplate != nil {
			parts = append(parts, "--chat_template", fmt.Sprint(*sg.ChatTemplate))
		}

		if sg.SpeculativeAlgorithm != nil {
			parts = append(parts, "--speculative_algorithm", fmt.Sprint(*sg.SpeculativeAlgorithm))
		}
		if sg.SpeculativeNumSteps != nil {
			parts = append(parts, "--speculative_steps", fmt.Sprint(*sg.SpeculativeNumSteps))
		}
		if sg.SpeculativeEagleTopk != nil {
			parts = append(parts, "--speculative_ep", fmt.Sprint(*sg.SpeculativeEagleTopk))
		}
		if sg.SpeculativeNumDraftTokens != nil {
			parts = append(parts, "--speculative-num-draft-tokens", fmt.Sprint(*sg.SpeculativeNumDraftTokens))
		}
		if sg.SpeculativeDraftModelPath != nil {
			parts = append(parts, "--speculative-draft-model-path", *sg.SpeculativeDraftModelPath)
		}

		if sg.PipelineParallelSize != nil {
			parts = append(parts, "--pp", fmt.Sprint(*sg.PipelineParallelSize))
		}
		if sg.SchedulePolicy != nil {
			parts = append(parts, "--schedule-policy", *sg.SchedulePolicy)
		}

		if sg.DisableRadixCache != nil && *sg.DisableRadixCache {
			parts = append(parts, "--disable-radix-cache")
		}
		if sg.DisableCudaGraph != nil && *sg.DisableCudaGraph {
			parts = append(parts, "--disable-cuda-graph")
		}
	}

	// Add common options shared with vLLM
	parts = appendCommonOptions(parts, cfg, "--tokenizer-path", "--random-seed")

	return parts
}
