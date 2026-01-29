package batch

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/genai"

	"github.com/specvital/worker/internal/domain/specview"
)

// Provider handles Batch API operations for Gemini.
type Provider struct {
	client *genai.Client
	config BatchConfig
}

// NewProvider creates a new Batch API provider.
func NewProvider(ctx context.Context, config BatchConfig) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  config.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("batch: failed to create Gemini client: %w", err)
	}

	return &Provider{
		client: client,
		config: config,
	}, nil
}

// CreateJob creates a batch job for Phase 1 classification.
func (p *Provider) CreateJob(ctx context.Context, req BatchRequest) (*BatchResult, error) {
	if len(req.Requests) == 0 {
		return nil, errors.New("batch: no requests provided")
	}

	slog.InfoContext(ctx, "creating batch job",
		"model", req.Model,
		"request_count", len(req.Requests),
		"analysis_id", req.AnalysisID,
	)

	src := &genai.BatchJobSource{
		InlinedRequests: req.Requests,
	}

	config := &genai.CreateBatchJobConfig{
		DisplayName: fmt.Sprintf("phase1-classification-%s", req.AnalysisID),
	}

	job, err := p.client.Batches.Create(ctx, req.Model, src, config)
	if err != nil {
		return nil, fmt.Errorf("batch: failed to create job: %w", err)
	}

	slog.InfoContext(ctx, "batch job created",
		"job_name", job.Name,
		"state", job.State,
	)

	return &BatchResult{
		JobName: job.Name,
		State:   mapJobState(job.State),
	}, nil
}

// GetJobStatus retrieves the current status of a batch job.
func (p *Provider) GetJobStatus(ctx context.Context, jobName string) (*BatchResult, error) {
	if jobName == "" {
		return nil, errors.New("batch: job name is required")
	}

	job, err := p.client.Batches.Get(ctx, jobName, nil)
	if err != nil {
		return nil, fmt.Errorf("batch: failed to get job status: %w", err)
	}

	result := &BatchResult{
		JobName: job.Name,
		State:   mapJobState(job.State),
	}

	// If job failed, capture error
	if job.State == genai.JobStateFailed && job.Error != nil {
		result.Error = fmt.Errorf("batch job failed: %s", job.Error.Message)
	}

	// If job succeeded, extract responses
	if job.State == genai.JobStateSucceeded {
		// Extract inline responses
		if job.Dest != nil && len(job.Dest.InlinedResponses) > 0 {
			result.Responses = job.Dest.InlinedResponses

			// Aggregate token usage from all responses
			totalPromptTokens := int32(0)
			totalCandidatesTokens := int32(0)

			for _, resp := range job.Dest.InlinedResponses {
				if resp.Response != nil && resp.Response.UsageMetadata != nil {
					totalPromptTokens += resp.Response.UsageMetadata.PromptTokenCount
					totalCandidatesTokens += resp.Response.UsageMetadata.CandidatesTokenCount
				}
			}

			result.TokenUsage = &specview.TokenUsage{
				Model:            job.Model,
				PromptTokens:     totalPromptTokens,
				CandidatesTokens: totalCandidatesTokens,
				TotalTokens:      totalPromptTokens + totalCandidatesTokens,
			}
		}
	}

	return result, nil
}

// CancelJob cancels a running batch job.
func (p *Provider) CancelJob(ctx context.Context, jobName string) error {
	if jobName == "" {
		return errors.New("batch: job name is required")
	}

	err := p.client.Batches.Cancel(ctx, jobName, nil)
	if err != nil {
		return fmt.Errorf("batch: failed to cancel job: %w", err)
	}

	slog.InfoContext(ctx, "batch job cancelled", "job_name", jobName)
	return nil
}

// Close releases resources held by the provider.
func (p *Provider) Close() error {
	// genai.Client doesn't require explicit close
	return nil
}
