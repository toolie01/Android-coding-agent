---

## **main.go**
```go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// Safeguards represents the agent's rule engine
type Safeguards struct {
	Version         string                 `json:"version"`
	Metadata        map[string]interface{} `json:"metadata"`
	FileSystem      FileSystemRules        `json:"file_system"`
	CodeGeneration  CodeGenerationRules    `json:"code_generation"`
	GitOperations   GitOperationRules      `json:"git_operations"`
	Testing         TestingRules           `json:"testing"`
	APICallRules    APICallRules           `json:"api_calls"`
	LLMInteraction  LLMInteractionConfig   `json:"llm_interaction"`
	Logging         LoggingConfig          `json:"logging"`
	Security        SecurityRules          `json:"security"`
	ResourceLimits  ResourceLimits         `json:"resource_limits"`
	AndroidSpecific AndroidSpecificRules   `json:"android_specific"`
}

type FileSystemRules struct {
	AllowedPaths      []string `json:"allowed_paths"`
	ForbiddenPaths    []string `json:"forbidden_paths"`
	AllowedExtensions []string `json:"allowed_extensions"`
	MaxFileSizeMB     int      `json:"max_file_size_mb"`
	ReadOnlyPaths     []string `json:"read_only_paths"`
}

type CodeGenerationRules struct {
	Enabled           bool     `json:"enabled"`
	MaxFilesPerOp     int      `json:"max_files_per_operation"`
	MaxLinesPerFile   int      `json:"max_lines_per_file"`
	RequireApproval   bool     `json:"require_approval"`
	Languages         []string `json:"languages"`
	RestrictedImports []string `json:"restricted_imports"`
}

type GitOperationRules struct {
	Enabled               bool     `json:"enabled"`
	AllowedOperations     []string `json:"allowed_operations"`
	ForbiddenOperations   []string `json:"forbidden_operations"`
	RequireApprovalFor    []string `json:"require_approval_for"`
	CommitMessageTemplate string   `json:"commit_message_template"`
	AutoCommit            bool     `json:"auto_commit"`
	BranchProtection      bool     `json:"branch_protection"`
}

type TestingRules struct {
	Enabled             bool     `json:"enabled"`
	AutoRunTests        bool     `json:"auto_run_tests"`
	TestFramework       []string `json:"test_framework"`
	CoverageMinimum     int      `json:"coverage_minimum_percent"`
	RequirePassingTests bool     `json:"require_passing_tests"`
	AllowedTestCommands []string `json:"allowed_test_commands"`
}

type APICallRules struct {
	Enabled           bool                   `json:"enabled"`
	AllowedServices   []string               `json:"allowed_services"`
	ForbiddenServices []string               `json:"forbidden_services"`
	RateLimit         map[string]interface{} `json:"rate_limit"`
	TimeoutSeconds    int                    `json:"timeout_seconds"`
}

type LLMInteractionConfig struct {
	Enabled            bool     `json:"enabled"`
	ModelType          string   `json:"model_type"`
	MaxTokens          int      `json:"max_tokens"`
	Temperature        float64  `json:"temperature"`
	TopP               float64  `json:"top_p"`
	ContextWindow      int      `json:"context_window"`
	SystemPrompt       string   `json:"system_prompt"`
	CustomInstructions []string `json:"custom_instructions"`
}

type LoggingConfig struct {
	Enabled            bool   `json:"enabled"`
	Level              string `json:"level"`
	LogFile            string `json:"log_file"`
	AuditTrail         bool   `json:"audit_trail"`
	LogDecisions       bool   `json:"log_decisions"`
	LogAPICalls        bool   `json:"log_api_calls"`
	LogFileOperations  bool   `json:"log_file_operations"`
}

type SecurityRules struct {
	SandboxMode               bool     `json:"sandbox_mode"`
	ValidateChecksums         bool     `json:"validate_checksums"`
	ScanForMaliciousPatterns  bool     `json:"scan_for_malicious_patterns"`
	MaliciousPatterns         []string `json:"malicious_patterns"`
	RequireHumanReviewFor     []string `json:"require_human_review_for"`
}

type ResourceLimits struct {
	MaxMemoryMB           int `json:"max_memory_mb"`
	MaxCPUPercent         int `json:"max_cpu_percent"`
	TimeoutSeconds        int `json:"timeout_seconds"`
	MaxParallelOperations int `json:"max_parallel_operations"`
}

type AndroidSpecificRules struct {
	GradleModifications        bool     `json:"gradle_modifications"`
	ManifestModifications      bool     `json:"manifest_modifications"`
	ResourceGeneration         bool     `json:"resource_generation"`
	DependencyManagement       bool     `json:"dependency_management"`
	MinSDKVersion              int      `json:"min_sdk_version"`
	TargetSDKVersion           int      `json:"target_sdk_version"`
	RestrictDangerousPermissions []string `json:"restrict_dangerous_permissions"`
}

// Agent represents the main coding agent
type Agent struct {
	ID         string
	Safeguards *Safeguards
	Context    map[string]interface{}
}

// NewAgent creates a new agent instance with loaded safeguards
func NewAgent(safeguardsPath string) (*Agent, error) {
	data, err := os.ReadFile(safeguardsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read safeguards file: %w", err)
	}

	var safeguards Safeguards
	if err := json.Unmarshal(data, &safeguards); err != nil {
		return nil, fmt.Errorf("failed to parse safeguards JSON: %w", err)
	}

	agent := &Agent{
		ID:         uuid.New().String(),
		Safeguards: &safeguards,
		Context:    make(map[string]interface{}),
	}

	log.Printf("[AGENT %s] Initialized with safeguards version %s", agent.ID, safeguards.Version)
	return agent, nil
}

// ValidateFilePath checks if a file path is allowed by safeguards
func (a *Agent) ValidateFilePath(path string, operation string) (bool, string) {
	path = filepath.Clean(path)

	for _, forbidden := range a.Safeguards.FileSystem.ForbiddenPaths {
		if strings.Contains(path, forbidden) {
			return false, fmt.Sprintf("Path %s is forbidden", forbidden)
		}
	}

	allowed := false
	for _, allowedPath := range a.Safeguards.FileSystem.AllowedPaths {
		if strings.HasPrefix(path, allowedPath) {
			allowed = true
			break
		}
	}

	if !allowed {
		return false, fmt.Sprintf("Path %s not in allowed paths", path)
	}

	if operation == "write" {
		for _, readOnly := range a.Safeguards.FileSystem.ReadOnlyPaths {
			if strings.HasPrefix(path, readOnly) {
				return false, fmt.Sprintf("Path %s is read-only", path)
			}
		}
	}

	ext := filepath.Ext(path)
	extAllowed := false
	for _, allowedExt := range a.Safeguards.FileSystem.AllowedExtensions {
		if ext == allowedExt {
			extAllowed = true
			break
		}
	}

	if !extAllowed {
		return false, fmt.Sprintf("File extension %s not allowed", ext)
	}

	return true, ""
}

// ValidateOperation checks if an operation is allowed
func (a *Agent) ValidateOperation(opType string, opName string) (bool, string) {
	switch opType {
	case "git":
		for _, forbidden := range a.Safeguards.GitOperations.ForbiddenOperations {
			if forbidden == opName {
				return false, fmt.Sprintf("Git operation %s is forbidden", opName)
			}
		}
		allowed := false
		for _, op := range a.Safeguards.GitOperations.AllowedOperations {
			if op == opName {
				allowed = true
				break
			}
		}
		if !allowed {
			return false, fmt.Sprintf("Git operation %s is not in allowed list", opName)
		}

	case "code_generation":
		if !a.Safeguards.CodeGeneration.Enabled {
			return false, "Code generation is disabled"
		}
	}

	return true, ""
}

// ScanForMaliciousPatterns checks code for dangerous patterns
func (a *Agent) ScanForMaliciousPatterns(code string) (bool, []string) {
	var findings []string

	for _, pattern := range a.Safeguards.Security.MaliciousPatterns {
		if strings.Contains(code, pattern) {
			findings = append(findings, fmt.Sprintf("Found malicious pattern: %s", pattern))
		}
	}

	return len(findings) == 0, findings
}

// ProcessPrompt is where agent logic interacts with LLM
func (a *Agent) ProcessPrompt(ctx context.Context, prompt string) (string, error) {
	log.Printf("[AGENT %s] Processing prompt: %s", a.ID, prompt[:min(100, len(prompt))])

	if a.Safeguards.Security.ScanForMaliciousPatterns {
		if safe, findings := a.ScanForMaliciousPatterns(prompt); !safe {
			return "", fmt.Errorf("prompt contains malicious patterns: %v", findings)
		}
	}

	response := fmt.Sprintf("Agent %s processed: %s", a.ID, prompt)
	
	log.Printf("[AGENT %s] Response generated", a.ID)
	return response, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	configPath := flag.String("config", "safeguards.json", "Path to safeguards JSON file")
	prompt := flag.String("prompt", "", "Prompt for the agent")
	flag.Parse()

	if *prompt == "" {
		fmt.Println("Error: --prompt is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	agent, err := NewAgent(*configPath)
	if err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	ctx := context.Background()
	response, err := agent.ProcessPrompt(ctx, *prompt)
	if err != nil {
		log.Fatalf("Error processing prompt: %v", err)
	}

	fmt.Println(response)
}
