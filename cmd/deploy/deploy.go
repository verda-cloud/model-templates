package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/verda-cloud/model-templates/pkg/engine"
	"github.com/verda-cloud/model-templates/pkg/template"
	"github.com/verda-cloud/verdacloud-sdk-go/pkg/verda"
)

const ( // Default image versions when no tag has been provided
	VLLMDefaultImage   = "vllm/vllm-openai"
	SGLANGDefaultImage = "lmsysorg/sglang"

	VLLMDefaultTag   = "v0.13.0"
	SGLANGDefaultTag = "v0.5.6.post2-cu129-amd64"

	DefaultHost = "0.0.0.0"
)

var (
	// Authentication flags
	clientID     string
	clientSecret string

	// Deployment configuration flags
	deployName       string
	containerImage   string
	exposedPort      int
	isSpot           bool
	gpuTypeOverride  string
	gpuCountOverride int

	// Scaling options
	minReplicas int
	maxReplicas int

	// Additional flags
	envVars []string
	dryRun  bool
	timeout int
)

var rootCmd = &cobra.Command{
	Use:   "deploy <template-file> [template-files...]",
	Short: "Deploy model templates to Verda container deployments",
	Long: `Deploy one or more JSON template configuration files to Verda's container deployments.

Authentication can be provided via flags or environment variables:
  VERDA_CLIENT_ID and VERDA_CLIENT_SECRET`,
	RunE: runDeploy,
}

func init() {
	rootCmd.Flags().StringVar(&clientID, "client-id", "", "Verda OAuth client ID")
	rootCmd.Flags().StringVar(&clientSecret, "client-secret", "", "Verda OAuth client secret")

	rootCmd.Flags().StringVarP(&deployName, "name", "n", "", "Deployment name (auto-generated from template if not provided)")
	rootCmd.Flags().StringVarP(&containerImage, "image", "i", "", "Container image (uses engine default if not provided)")
	rootCmd.Flags().IntVarP(&exposedPort, "port", "p", 8000, "Exposed container port")
	rootCmd.Flags().BoolVar(&isSpot, "spot", false, "Use spot instances for deployment")
	rootCmd.Flags().StringVar(&gpuTypeOverride, "gpu-type", "", "Override GPU type from template (e.g., 'H100', 'A100')")
	rootCmd.Flags().IntVar(&gpuCountOverride, "gpu-count", 0, "Override GPU count (uses template TP value if not provided)")

	// Scaling options
	rootCmd.Flags().IntVar(&minReplicas, "min-replicas", 1, "Minimum number of replicas")
	rootCmd.Flags().IntVar(&maxReplicas, "max-replicas", 1, "Maximum number of replicas")

	// Additional flags
	rootCmd.Flags().StringArrayVarP(&envVars, "env", "e", []string{}, "Environment variables for container (format: KEY=VALUE, can be specified multiple times)")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Validate and show deployment plan without deploying")
	rootCmd.Flags().IntVar(&timeout, "timeout", 300, "Deployment timeout in seconds")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runDeploy(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no template files specified")
	}

	// Use environment variables if flags not provided
	if clientID == "" {
		clientID = os.Getenv("VERDA_CLIENT_ID")
	}
	if clientSecret == "" {
		clientSecret = os.Getenv("VERDA_CLIENT_SECRET")
	}

	var client *verda.Client
	var err error

	// Authentication is not needed for dry runs
	if !dryRun {
		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("authentication required: provide --client-id and --client-secret or set VERDA_CLIENT_ID and VERDA_CLIENT_SECRET environment variables")
		}

		client, err = createVerdaClient()
		if err != nil {
			return fmt.Errorf("failed to create Verda client: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	for _, file := range args {
		if err := deployTemplate(ctx, client, file); err != nil {
			fmt.Fprintf(os.Stderr, "Error deploying %s: %v\n", file, err)
			break
		}
	}

	return nil
}

func createVerdaClient() (*verda.Client, error) {
	var opts = []verda.ClientOption{
		verda.WithClientID(clientID),
		verda.WithClientSecret(clientSecret),
	}

	return verda.NewClient(opts...)
}

func parseEnvVars(envVars []string) ([]verda.ContainerEnvVar, error) {
	var envList []verda.ContainerEnvVar
	for _, env := range envVars {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid environment variable format: %s (expected KEY=VALUE)", env)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("environment variable key cannot be empty: %s", env)
		}
		envList = append(envList, verda.ContainerEnvVar{
			Type:                     "plain",
			Name:                     key,
			ValueOrReferenceToSecret: value,
		})
	}
	return envList, nil
}

func deployTemplate(ctx context.Context, client *verda.Client, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("unable to read file: %v", err)
	}

	var cfg template.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unable to parse file: %v", err)
	}

	// Host for containers should always be default
	host := DefaultHost
	cfg.Host = &host
	if exposedPort != 0 {
		cfg.Port = &exposedPort
	}

	command, err := engine.GenerateCommandFromTemplate(cfg)
	if err != nil {
		return fmt.Errorf("failed to generate command: %v", err)
	}

	name := deployName
	if name == "" {
		fileName := filepath.Base(filePath)
		name = strings.TrimSuffix(fileName, filepath.Ext(fileName))
		// Sanitize name for deployment
		name = strings.ReplaceAll(name, "_", "-")
	}

	gpuType := gpuTypeOverride
	if gpuType == "" && len(cfg.GPUTypes) > 0 {
		// Use last GPU type from template for now
		gpuType = strings.ToUpper(string(cfg.GPUTypes[len(cfg.GPUTypes)-1]))
	}
	if gpuType == "" {
		return fmt.Errorf("no GPU type specified and template doesn't specify gpu_types")
	}

	gpuCount := gpuCountOverride
	if gpuCount == 0 {
		// Try to infer from template parallelism settings
		if cfg.Engine == template.EngineSGLang && cfg.SGLang != nil && cfg.SGLang.TP != nil {
			gpuCount = *cfg.SGLang.TP
		} else if cfg.Engine == template.EngineVLLM && cfg.VLLM != nil && cfg.VLLM.TensorParallelSize != nil {
			gpuCount = *cfg.VLLM.TensorParallelSize
		} else {
			gpuCount = 1
		}
	}

	image := containerImage
	if image == "" {
		switch cfg.Engine {
		case template.EngineVLLM:
			tag := VLLMDefaultTag
			if cfg.ImageTag != nil && *cfg.ImageTag != "" {
				tag = *cfg.ImageTag
			}
			image = fmt.Sprintf("%s:%s", VLLMDefaultImage, tag)
		case template.EngineSGLang:
			tag := SGLANGDefaultTag
			if cfg.ImageTag != nil && *cfg.ImageTag != "" {
				tag = *cfg.ImageTag
			}
			image = fmt.Sprintf("%s:%s", SGLANGDefaultImage, tag)
		case template.EngineCustom:
			// For custom engine, use the image from custom config if provided
			if cfg.Custom != nil && cfg.Custom.Image != nil && *cfg.Custom.Image != "" {
				image = *cfg.Custom.Image
				// If image doesn't have a tag and image_tag is provided, append it
				if !strings.Contains(image, ":") && cfg.ImageTag != nil && *cfg.ImageTag != "" {
					image = fmt.Sprintf("%s:%s", image, *cfg.ImageTag)
				}
			} else {
				return fmt.Errorf("no image specified for custom engine")
			}
		default:
			return fmt.Errorf("no image specified and no default for engine %s", cfg.Engine)
		}
	}

	if len(command) == 0 {
		return fmt.Errorf("generated command is empty")
	}

	envList, err := parseEnvVars(envVars)
	if err != nil {
		return fmt.Errorf("failed to parse environment variables: %v", err)
	}

	fmt.Printf("Deploying %s...\n", filepath.Base(filePath))
	fmt.Printf("  Model: %s\n", cfg.Model)
	fmt.Printf("  Engine: %s\n", cfg.Engine)
	fmt.Printf("  GPU: %dx %s\n", gpuCount, gpuType)
	fmt.Printf("  Deployment Name: %s\n", name)
	if len(envList) > 0 {
		fmt.Printf("  Environment Variables: %d set\n", len(envList))
	}

	createReq := verda.CreateDeploymentRequest{
		Name:   name,
		IsSpot: isSpot,
		Compute: verda.ContainerCompute{
			Name: gpuType,
			Size: gpuCount,
		},
		ContainerRegistrySettings: verda.ContainerRegistrySettings{
			IsPrivate: false,
		},
		Scaling: verda.ContainerScalingOptions{
			MinReplicaCount:              minReplicas,
			MaxReplicaCount:              maxReplicas,
			ScaleDownPolicy:              &verda.ScalingPolicy{DelaySeconds: 300},
			ScaleUpPolicy:                &verda.ScalingPolicy{DelaySeconds: 300},
			QueueMessageTTLSeconds:       1,
			ConcurrentRequestsPerReplica: 1,
			ScalingTriggers: &verda.ScalingTriggers{
				QueueLoad: &verda.QueueLoadTrigger{Threshold: 2},
			},
		},
		Containers: []verda.CreateDeploymentContainer{
			{
				Image:       image,
				ExposedPort: exposedPort,
				Healthcheck: &verda.ContainerHealthcheck{
					Enabled: true,
					Port:    exposedPort,
					Path:    "/health", // TODO customise per engine and keep configurable
				},
				Env: envList,
				EntrypointOverrides: &verda.ContainerEntrypointOverrides{
					Enabled: true,
					Cmd:     command,
				},
			},
		},
	}

	if dryRun {
		jsonText, err := json.MarshalIndent(createReq, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal create deployment request: %v", err)
		}
		fmt.Printf("  DRY RUN! Would create deployment %s, but nah\nBody:\n\n%s\n", name, jsonText)
		os.Exit(0)
	}

	// Create the deployment
	deployment, err := client.ContainerDeployments.CreateDeployment(ctx, &createReq)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	fmt.Printf("✓ Deployed successfully\n")
	fmt.Printf("  Deployment Name: %s\n", deployment.Name)
	fmt.Printf("  Endpoint: %s\n", deployment.EndpointBaseURL)
	fmt.Printf("  Created At: %s\n", deployment.CreatedAt.Format(time.RFC3339))
	fmt.Println()

	return nil
}
