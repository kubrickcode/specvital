package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"google.golang.org/genai"

	"github.com/specvital/web/src/backend/modules/spec-view/domain"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

const (
	// AI Model Configuration
	defaultModelID         = "gemini-2.5-flash-lite" // Cost-effective model optimized for high-volume, low-latency tasks
	defaultTemperature     = 0.1                     // Low temperature for consistent, deterministic outputs
	defaultMaxOutputTokens = 4096                    // Sufficient for ~200 test names per request

	// Retry Configuration
	defaultMaxRetries        = 3                      // Industry standard; balances reliability vs latency
	defaultInitialRetryDelay = 100 * time.Millisecond // Fast initial retry for transient failures
	defaultMaxRetryDelay     = 2 * time.Second        // Cap to prevent excessive wait times

	// Circuit Breaker Configuration
	cbFailureThreshold = 5                // Open circuit after 5 consecutive failures
	cbSuccessThreshold = 2                // Close circuit after 2 consecutive successes in half-open state
	cbTimeout          = 30 * time.Second // Wait before attempting recovery (half-open)
)

type GeminiConfig struct {
	APIKey  string
	ModelID string
}

type GeminiProvider struct {
	client         *genai.Client
	modelID        string
	circuitBreaker *CircuitBreaker
}

var _ port.AIProvider = (*GeminiProvider)(nil)

func NewGeminiProvider(ctx context.Context, cfg GeminiConfig) (*GeminiProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("gemini API key is required")
	}

	modelID := cfg.ModelID
	if modelID == "" {
		modelID = defaultModelID
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  cfg.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}

	circuitBreaker := NewCircuitBreaker(CircuitBreakerConfig{
		FailureThreshold: cbFailureThreshold,
		SuccessThreshold: cbSuccessThreshold,
		Timeout:          cbTimeout,
	})

	return &GeminiProvider{
		client:         client,
		modelID:        modelID,
		circuitBreaker: circuitBreaker,
	}, nil
}

func (p *GeminiProvider) ModelID() string {
	return p.modelID
}

func (p *GeminiProvider) ConvertTestNames(ctx context.Context, input port.ConvertInput) (map[string]string, error) {
	if !p.circuitBreaker.Allow() {
		return nil, domain.ErrAIProviderUnavailable
	}

	prompt := buildPrompt(input)
	config := p.buildConfig()

	resp, err := p.generateWithRetry(ctx, prompt, config)
	if err != nil {
		p.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("gemini convert: %w", err)
	}

	result, err := parseResponse(resp)
	if err != nil {
		p.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("parse response: %w", err)
	}

	p.circuitBreaker.RecordSuccess()

	return result, nil
}

func (p *GeminiProvider) buildConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		Temperature:      genai.Ptr(float32(defaultTemperature)),
		MaxOutputTokens:  defaultMaxOutputTokens,
		ResponseMIMEType: "application/json",
	}
}

func (p *GeminiProvider) generateWithRetry(ctx context.Context, prompt string, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	var lastErr error

	for attempt := range defaultMaxRetries {
		resp, err := p.client.Models.GenerateContent(ctx, p.modelID, genai.Text(prompt), config)
		if err == nil {
			return resp, nil
		}

		if !isRetryableError(err) {
			return nil, classifyError(err)
		}

		lastErr = err

		if attempt < defaultMaxRetries-1 {
			delay := calculateBackoff(attempt)

			select {
			case <-time.After(delay):
				continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, fmt.Errorf("max retries exceeded: %w", classifyError(lastErr))
}

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var apiErr genai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case 429, 500, 502, 503, 504:
			return true
		}
	}

	errStr := err.Error()
	return strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "temporarily unavailable")
}

func classifyError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr genai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case 429:
			return domain.ErrRateLimited
		case 401, 403:
			return fmt.Errorf("invalid API key: %w", err)
		case 400:
			return fmt.Errorf("invalid request: %w", err)
		default:
			return domain.ErrAIProviderUnavailable
		}
	}

	return domain.ErrConversionFailed
}

func calculateBackoff(attempt int) time.Duration {
	backoff := float64(defaultInitialRetryDelay) * math.Pow(2, float64(attempt))
	if backoff > float64(defaultMaxRetryDelay) {
		backoff = float64(defaultMaxRetryDelay)
	}

	jitter := 0.5 + rand.Float64()
	return time.Duration(backoff * jitter)
}

func buildPrompt(input port.ConvertInput) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`Convert test names to user-facing feature descriptions in %s.

<rules>
1. STRIP test syntax: should, must, it, test, spec, when, given, then, if, expect, assert, verify, returns, throws, describe, context
2. FOCUS on user value, not implementation:
   - "caches result" -> "Fast repeated access"
   - "validates input" -> "Input format verification"
   - "throws error" -> "Error notification"
   - "returns null" -> "Empty result handling"
3. USE hierarchy for context when test name alone is ambiguous:
   - hierarchy="EmailService > Validation", name="invalid" -> "Email format validation"
   - hierarchy="Cart > Checkout", name="empty" -> "Empty cart checkout prevention"
4. HANDLE special cases:
   - Parameterized tests (contains ":"): Extract concept, ignore specific values
   - Numeric/meaningless names: Derive meaning from hierarchy
   - Negation tests: Frame as protection or prevention
5. LENGTH: Max 50 characters. Prefer noun phrases.
6. LANGUAGE: Output in %s. Maintain technical terms that are universal (API, HTTP, JSON, etc.)
</rules>

<input>
`, languageDisplayName(input.Language), languageDisplayName(input.Language)))

	globalIndex := 1
	for _, suite := range input.Suites {
		sb.WriteString(fmt.Sprintf("## %s\n", suite.Hierarchy))
		for _, testName := range suite.Tests {
			sb.WriteString(fmt.Sprintf("%d|%s\n", globalIndex, testName))
			globalIndex++
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`</input>

<output_format>
Return JSON object where each key is the global_index (as string) and value is the converted name.
Example: {"1": "Session creation on login", "2": "Response caching", "3": "Email format validation"}
</output_format>`)

	return sb.String()
}

func languageDisplayName(lang entity.Language) string {
	names := map[entity.Language]string{
		entity.LanguageAr: "Arabic",
		entity.LanguageCs: "Czech",
		entity.LanguageDa: "Danish",
		entity.LanguageDe: "German",
		entity.LanguageEl: "Greek",
		entity.LanguageEn: "English",
		entity.LanguageEs: "Spanish",
		entity.LanguageFi: "Finnish",
		entity.LanguageFr: "French",
		entity.LanguageHi: "Hindi",
		entity.LanguageId: "Indonesian",
		entity.LanguageIt: "Italian",
		entity.LanguageJa: "Japanese",
		entity.LanguageKo: "Korean",
		entity.LanguageNl: "Dutch",
		entity.LanguagePl: "Polish",
		entity.LanguagePt: "Portuguese",
		entity.LanguageRu: "Russian",
		entity.LanguageSv: "Swedish",
		entity.LanguageTh: "Thai",
		entity.LanguageTr: "Turkish",
		entity.LanguageUk: "Ukrainian",
		entity.LanguageVi: "Vietnamese",
		entity.LanguageZh: "Chinese",
	}

	if name, ok := names[lang]; ok {
		return name
	}
	return "English"
}

func parseResponse(resp *genai.GenerateContentResponse) (map[string]string, error) {
	if resp == nil || len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("empty response from AI provider")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	var textContent string
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			textContent = part.Text
			break
		}
	}

	if textContent == "" {
		return nil, fmt.Errorf("no text content in response")
	}

	result := make(map[string]string)
	if err := json.Unmarshal([]byte(textContent), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	return result, nil
}
