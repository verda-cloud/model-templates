package template

type EngineOption string
type GPUType string

const (
	EngineSGLang EngineOption = "sglang"
	EngineVLLM   EngineOption = "vllm"
	EngineCustom EngineOption = "custom"

	GPUTypeA100       GPUType = "a100"
	GPUTypeB200       GPUType = "b200"
	GPUTypeB300       GPUType = "b300"
	GPUTypeH100       GPUType = "h100"
	GPUTypeH200       GPUType = "h200"
	GPUTypeL40S       GPUType = "l40s"
	GPUTypeRTX6000    GPUType = "rtx6000-ada"
	GPUTypeRTXPro6000 GPUType = "rtx6000-pro"
	GPUTypeV100       GPUType = "v100"
)

// Config represents the container configuration with common fields
type Config struct {
	Engine           EngineOption `json:"engine"`
	Model            string       `json:"model"`
	Explanation      string       `json:"explanation,omitempty"`
	ShortExplanation string       `json:"short_explanation,omitempty"`

	// Common fields (supported by multiple engines)
	TrustRemoteCode        *bool          `json:"trust_remote_code,omitempty"`
	Quantization           *string        `json:"quantization,omitempty"`
	GPUTypes               []GPUType      `json:"gpu_types,omitempty"`
	LoadFormat             *string        `json:"load_format,omitempty"`
	ModelLoaderExtraConfig map[string]any `json:"model_loader_extra_config,omitempty"`
	Seed                   *int           `json:"seed,omitempty"`
	Tokenizer              *string        `json:"tokenizer,omitempty"`
	ImageTag               *string        `json:"image_tag,omitempty"`

	// Container configuration fields (not part of template JSON, used for deployment)
	// Host specifies the network interface to bind to (default: 0.0.0.0)
	// Port specifies the container port to expose (set by deployment tool)
	Host *string
	Port *int

	// Engine-specific configurations
	SGLang *SGLangConfig `json:"sglang,omitempty"`
	VLLM   *VLLMConfig   `json:"vllm,omitempty"`
	Custom *CustomConfig `json:"custom,omitempty"`
}

// SGLangConfig contains SGLang-specific configuration options
type SGLangConfig struct {
	TP *int `json:"tp,omitempty"`
	DP *int `json:"dp,omitempty"`
	EP *int `json:"ep,omitempty"`

	MemFractionStatic *float64 `json:"mem_fraction_static,omitempty"`
	KVCacheDtype      *string  `json:"kv_cache_dtype,omitempty"`

	ContextLength      *int `json:"context_length,omitempty"`
	MaxRunningRequests *int `json:"max_running_requests,omitempty"`

	EnableDPAttention          *bool   `json:"enable_dp_attention,omitempty"`
	ReasoningParser            *string `json:"reasoning_parser,omitempty"`
	ToolCallParser             *string `json:"tool_call_parser,omitempty"`
	ChatTemplate               *string `json:"chat_template,omitempty"`
	DisableSharedExpertsFusion *bool   `json:"disable_shared_experts_fusion,omitempty"`

	SpeculativeAlgorithm      *string `json:"speculative_algorithm,omitempty"`
	SpeculativeNumSteps       *int    `json:"speculative_num_steps,omitempty"`
	SpeculativeEagleTopk      *int    `json:"speculative_eagle_topk,omitempty"`
	SpeculativeNumDraftTokens *int    `json:"speculative_num_draft_tokens,omitempty"`
	SpeculativeDraftModelPath *string `json:"speculative_draft_model_path,omitempty"`

	PipelineParallelSize *int    `json:"pipeline_parallel_size,omitempty"`
	SchedulePolicy       *string `json:"schedule_policy,omitempty"`

	DisableRadixCache *bool `json:"disable_radix_cache,omitempty"`
	DisableCudaGraph  *bool `json:"disable_cuda_graph,omitempty"`
}

// VLLMConfig contains vLLM-specific configuration options
type VLLMConfig struct {
	TensorParallelSize   *int `json:"tensor_parallel_size,omitempty"`
	PipelineParallelSize *int `json:"pipeline_parallel_size,omitempty"`

	Dtype       *string `json:"dtype,omitempty"`
	MaxModelLen *int    `json:"max_model_len,omitempty"`

	GPUMemoryUtilization *float64 `json:"gpu_memory_utilization,omitempty"`
	MaxNumSeqs           *int     `json:"max_num_seqs,omitempty"`

	EnablePrefixCaching  *bool   `json:"enable_prefix_caching,omitempty"`
	EnableChunkedPrefill *bool   `json:"enable_chunked_prefill,omitempty"`
	EnforceEager         *bool   `json:"enforce_eager,omitempty"`
	ServedModelName      *string `json:"served_model_name,omitempty"`

	DataParallelSize           *int    `json:"data_parallel_size,omitempty"`
	DistributedExecutorBackend *string `json:"distributed_executor_backend,omitempty"`

	SwapSpace           *int `json:"swap_space,omitempty"`
	MaxNumBatchedTokens *int `json:"max_num_batched_tokens,omitempty"`

	SchedulingPolicy *string `json:"scheduling_policy,omitempty"`
}

// CustomConfig contains custom engine configuration options
type CustomConfig struct {
	BaseCommand string         `json:"base_command"`
	ModelFlag   *string        `json:"model_flag,omitempty"`
	Image       *string        `json:"image,omitempty"`
	Args        []string       `json:"args,omitempty"`
	KVArgs      map[string]any `json:"kv_args,omitempty"`
}
