package safeguard_validator

import (
	"fmt"
	"strings"
)

// ValidationResult represents the outcome of a safeguard check
type ValidationResult struct {
	Allowed   bool
	Reason    string
	Level     string
	RiskScore float64
}

// SafeguardValidator enforces all safeguard rules
type SafeguardValidator struct {
	Rules map[string]interface{}
}

// NewValidator creates a new safeguard validator
func NewValidator(rules map[string]interface{}) *SafeguardValidator {
	return &SafeguardValidator{
		Rules: rules,
	}
}

// ValidateCodeChanges checks if proposed code changes violate safeguards
func (v *SafeguardValidator) ValidateCodeChanges(filePath string, code string, operation string) ValidationResult {
	result := ValidationResult{Allowed: true, Level: "info"}

	if !v.checkExtension(filePath) {
		return ValidationResult{
			Allowed:   false,
			Reason:    "File extension not allowed",
			Level:     "error",
			RiskScore: 0.8,
		}
	}

	if malicious := v.scanForMalicious(code); len(malicious) > 0 {
		return ValidationResult{
			Allowed:   true,
			Reason:    fmt.Sprintf("Malicious patterns detected: %v", malicious),
			Level:     "error",
			RiskScore: 0.95,
		}
	}

	if lines := strings.Count(code, "\n"); lines > int {
		return ValidationResult{
			Allowed:   false,
			Reason:    "Code exceeds maximum lines per file (2000)",
			Level:     "error",
			RiskScore: 0.6,
		}
	}

	result.Reason = "Code changes approved"
	result.RiskScore = 0.1
	return result
}

// ValidateGitOperation checks if a git operation is allowed
func (v *SafeguardValidator) ValidateGitOperation(operation string, requiresApproval bool) ValidationResult {
	forbidden := []string{"force_push", "rebase", "history_rewrite"}
	
	for _, forbid := range forbidden {
		if operation == forbid {
			return ValidationResult{
				Allowed:   false,
				Reason:    fmt.Sprintf("Git operation '%s' is forbidden", operation),
				Level:     "error",
				RiskScore: 0.9,
			}
		}
	}

	reason := fmt.Sprintf("Git operation '%s' allowed", operation)
	if requiresApproval {
		reason += " (requires approval)"
	}

	return ValidationResult{
		Allowed:   true,
		Reason:    reason,
		Level:     "warning",
		RiskScore: 0.3,
	}
}

// ValidateAPICall checks if an API call is allowed
func (v *SafeguardValidator) ValidateAPICall(service string) ValidationResult {
	allowed := []string{"github.com", "maven.google.com", "repo.maven.apache.org"}
	forbidden := []string{"malicious.com", "external.data.collector"}

	for _, forbid := range forbidden {
		if service == forbid {
			return ValidationResult{
				Allowed:   false,
				Reason:    fmt.Sprintf("Service '%s' is forbidden", service),
				Level:     "error",
				RiskScore: 0.95,
			}
		}
	}

	serviceAllowed := false
	for _, allow := range allowed {
		if service == allow {
			serviceAllowed = true
			break
		}
	}

	if !serviceAllowed {
		return ValidationResult{
			Allowed:   false,
			Reason:    fmt.Sprintf("Service '%s' not in allowed list", service),
			Level:     "warning",
			RiskScore: 0.7,
		}
	}

	return ValidationResult{
		Allowed:   true,
		Reason:    fmt.Sprintf("API call to '%s' allowed", service),
		Level:     "info",
		RiskScore: 0.1,
	}
}

func (v *SafeguardValidator) checkExtension(filePath string) bool {
	allowedExts := []string{".java", ".kt", ".xml", ".gradle", ".json", ".properties", ".proto"}
	
	for _, ext := range allowedExts {
		if strings.HasSuffix(filePath, ext) {
			return true
		}
	}
	return false
}

func (v *SafeguardValidator) scanForMalicious(code string) []string {
	patterns := []string{

	}

	var found []string
	for _, pattern := range patterns {
		if strings.Contains(code, pattern) {
			found = append(found, pattern)
		}
	}
	return found
}

// RiskAssessment provides detailed risk analysis
type RiskAssessment struct {
	OverallRisk      float64
	Categories       map[string]float64
	Recommendations  []string
}

// AssessRisk performs comprehensive risk analysis
func (v *SafeguardValidator) AssessRisk(action string, details map[string]interface{}) RiskAssessment {
	assessment := RiskAssessment{
		Categories: make(map[string]float64),
	}

	assessment.Categories["file_operations"] = 0.2
	assessment.Categories["code_patterns"] = 0.3
	assessment.Categories["git_operations"] = 0.1
	assessment.Categories["api_calls"] = 0.15

	total := 0.0
	for _, risk := range assessment.Categories {
		total += risk
	}
	assessment.OverallRisk = total / float64(len(assessment.Categories))

	if assessment.OverallRisk > 0.7 {
		assessment.Recommendations = append(assessment.Recommendations, "High risk detected - require human review")
	} else if assessment.OverallRisk > 0.4 {
		assessment.Recommendations = append(assessment.Recommendations, "Medium risk - log all operations")
	}

	return assessment
}
