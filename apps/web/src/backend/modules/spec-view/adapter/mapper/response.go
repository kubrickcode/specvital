package mapper

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/specvital/web/src/backend/internal/api"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
)

func ToSpecDocumentResponse(doc *entity.SpecDocument) ([]byte, error) {
	analysisUID, err := uuid.Parse(doc.AnalysisID)
	if err != nil {
		return nil, fmt.Errorf("invalid analysis ID %q: %w", doc.AnalysisID, err)
	}
	docUID, err := uuid.Parse(doc.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid document ID %q: %w", doc.ID, err)
	}

	domains, err := toAPIDomains(doc.Domains)
	if err != nil {
		return nil, err
	}

	var availableLanguages *[]api.AvailableLanguageInfo
	if len(doc.AvailableLanguages) > 0 {
		langs := toAPIAvailableLanguages(doc.AvailableLanguages)
		availableLanguages = &langs
	}

	specDoc := api.SpecDocument{
		AnalysisID:         analysisUID,
		AvailableLanguages: availableLanguages,
		CreatedAt:          doc.CreatedAt,
		Domains:            domains,
		ExecutiveSummary:   doc.ExecutiveSummary,
		ID:                 docUID,
		Language:           api.SpecLanguage(doc.Language),
		ModelID:            &doc.ModelID,
		Version:            doc.Version,
	}

	return marshalCompleted(specDoc)
}

func toAPIAvailableLanguages(langs []entity.AvailableLanguageInfo) []api.AvailableLanguageInfo {
	result := make([]api.AvailableLanguageInfo, len(langs))
	for i, l := range langs {
		result[i] = api.AvailableLanguageInfo{
			CreatedAt:       l.CreatedAt,
			HasPreviousSpec: l.HasPreviousSpec,
			Language:        api.SpecLanguage(l.Language),
			LatestVersion:   l.LatestVersion,
		}
	}
	return result
}

func ToGeneratingResponse(status *entity.SpecGenerationStatus) ([]byte, error) {
	genStatus := api.SpecGenerationStatus{
		Status: toAPIStatusEnum(status.Status),
	}
	if status.StartedAt != nil {
		genStatus.StartedAt = status.StartedAt
	}
	if status.CompletedAt != nil {
		genStatus.CompletedAt = status.CompletedAt
	}
	if status.ErrorMessage != nil {
		genStatus.ErrorMessage = status.ErrorMessage
	}

	return marshalGenerating(genStatus)
}

func ToRequestGenerationResponse(output *entity.SpecGenerationStatus) (api.RequestSpecGenerationResponse, error) {
	analysisUID, err := uuid.Parse(output.AnalysisID)
	if err != nil {
		return api.RequestSpecGenerationResponse{}, fmt.Errorf("invalid analysis ID %q: %w", output.AnalysisID, err)
	}
	return api.RequestSpecGenerationResponse{
		AnalysisID: analysisUID,
		Status:     toAPIStatusEnum(output.Status),
	}, nil
}

func ToGenerationStatusResponse(status *entity.SpecGenerationStatus) (api.SpecGenerationStatusResponse, error) {
	analysisUID, err := uuid.Parse(status.AnalysisID)
	if err != nil {
		return api.SpecGenerationStatusResponse{}, fmt.Errorf("invalid analysis ID %q: %w", status.AnalysisID, err)
	}
	resp := api.SpecGenerationStatusResponse{
		AnalysisID: analysisUID,
		Status:     toAPIStatusEnum(status.Status),
	}
	if status.StartedAt != nil {
		resp.StartedAt = status.StartedAt
	}
	if status.CompletedAt != nil {
		resp.CompletedAt = status.CompletedAt
	}
	if status.ErrorMessage != nil {
		resp.ErrorMessage = status.ErrorMessage
	}
	return resp, nil
}

func toAPIDomains(domains []entity.SpecDomain) ([]api.SpecDomain, error) {
	result := make([]api.SpecDomain, len(domains))
	for i, d := range domains {
		domainUID, err := uuid.Parse(d.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid domain ID %q: %w", d.ID, err)
		}
		features, err := toAPIFeatures(d.Features)
		if err != nil {
			return nil, err
		}
		result[i] = api.SpecDomain{
			Description: d.Description,
			Features:    features,
			ID:          domainUID,
			Name:        d.Name,
			SortOrder:   d.SortOrder,
		}
		if d.ClassificationConfidence != nil {
			f32 := float32(*d.ClassificationConfidence)
			result[i].ClassificationConfidence = &f32
		}
	}
	return result, nil
}

func toAPIFeatures(features []entity.SpecFeature) ([]api.SpecFeature, error) {
	result := make([]api.SpecFeature, len(features))
	for i, f := range features {
		featureUID, err := uuid.Parse(f.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid feature ID %q: %w", f.ID, err)
		}
		behaviors, err := toAPIBehaviors(f.Behaviors)
		if err != nil {
			return nil, err
		}
		result[i] = api.SpecFeature{
			Behaviors:   behaviors,
			Description: f.Description,
			ID:          featureUID,
			Name:        f.Name,
			SortOrder:   f.SortOrder,
		}
	}
	return result, nil
}

func toAPIBehaviors(behaviors []entity.SpecBehavior) ([]api.SpecBehavior, error) {
	result := make([]api.SpecBehavior, len(behaviors))
	for i, b := range behaviors {
		behaviorUID, err := uuid.Parse(b.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid behavior ID %q: %w", b.ID, err)
		}
		result[i] = api.SpecBehavior{
			ConvertedDescription: b.ConvertedDescription,
			ID:                   behaviorUID,
			OriginalName:         b.OriginalName,
			SortOrder:            b.SortOrder,
		}
		if b.SourceTestCaseID != nil {
			tcUID, err := uuid.Parse(*b.SourceTestCaseID)
			if err != nil {
				return nil, fmt.Errorf("invalid test case ID %q: %w", *b.SourceTestCaseID, err)
			}
			result[i].SourceTestCaseID = &tcUID
		}
		if b.SourceInfo != nil {
			result[i].SourceInfo = &api.SpecBehaviorSourceInfo{
				FilePath:   b.SourceInfo.FilePath,
				Framework:  api.Framework(b.SourceInfo.Framework),
				LineNumber: b.SourceInfo.LineNumber,
				Status:     api.TestStatus(b.SourceInfo.Status),
			}
		}
	}
	return result, nil
}

func toAPIStatusEnum(status entity.GenerationStatus) api.SpecGenerationStatusEnum {
	switch status {
	case entity.StatusPending:
		return api.SpecGenerationStatusEnumPending
	case entity.StatusRunning:
		return api.SpecGenerationStatusEnumRunning
	case entity.StatusCompleted:
		return api.SpecGenerationStatusEnumCompleted
	case entity.StatusFailed:
		return api.SpecGenerationStatusEnumFailed
	default:
		return api.SpecGenerationStatusEnumNotFound
	}
}

func marshalCompleted(doc api.SpecDocument) ([]byte, error) {
	completed := api.SpecDocumentCompleted{
		Data:   doc,
		Status: "completed",
	}
	return json.Marshal(completed)
}

func marshalGenerating(status api.SpecGenerationStatus) ([]byte, error) {
	generating := api.SpecDocumentGenerating{
		GenerationStatus: status,
		Status:           "generating",
	}
	return json.Marshal(generating)
}

func ToVersionHistoryResponse(output *VersionHistoryInput) api.VersionHistoryResponse {
	data := make([]api.VersionInfo, len(output.Versions))
	for i, v := range output.Versions {
		data[i] = api.VersionInfo{
			CreatedAt: v.CreatedAt,
			Version:   v.Version,
		}
		if v.ModelID != "" {
			data[i].ModelID = &v.ModelID
		}
	}

	return api.VersionHistoryResponse{
		Data:          data,
		Language:      api.SpecLanguage(output.Language),
		LatestVersion: output.LatestVersion,
	}
}

type VersionHistoryInput struct {
	Language      string
	LatestVersion int
	Versions      []entity.VersionInfo
}

// ToRepoSpecDocumentResponse converts RepoSpecDocument to API response
func ToRepoSpecDocumentResponse(doc *entity.RepoSpecDocument) (api.RepoSpecDocumentResponse, error) {
	analysisUID, err := uuid.Parse(doc.AnalysisID)
	if err != nil {
		return api.RepoSpecDocumentResponse{}, fmt.Errorf("invalid analysis ID %q: %w", doc.AnalysisID, err)
	}
	docUID, err := uuid.Parse(doc.ID)
	if err != nil {
		return api.RepoSpecDocumentResponse{}, fmt.Errorf("invalid document ID %q: %w", doc.ID, err)
	}

	domains, err := toAPIDomains(doc.Domains)
	if err != nil {
		return api.RepoSpecDocumentResponse{}, err
	}

	var availableLanguages *[]api.AvailableLanguageInfo
	if len(doc.AvailableLanguages) > 0 {
		langs := toAPIAvailableLanguages(doc.AvailableLanguages)
		availableLanguages = &langs
	}

	repoDoc := api.RepoSpecDocument{
		AnalysisID:         analysisUID,
		AvailableLanguages: availableLanguages,
		CommitSHA:          doc.CommitSHA,
		CreatedAt:          doc.CreatedAt,
		Domains:            domains,
		ExecutiveSummary:   doc.ExecutiveSummary,
		ID:                 docUID,
		Language:           api.SpecLanguage(doc.Language),
		ModelID:            &doc.ModelID,
		Version:            doc.Version,
	}

	resp := api.RepoSpecDocumentResponse{}
	if err := resp.FromRepoSpecDocumentCompleted(api.RepoSpecDocumentCompleted{
		Data:   repoDoc,
		Status: "completed",
	}); err != nil {
		return api.RepoSpecDocumentResponse{}, fmt.Errorf("failed to create response: %w", err)
	}

	return resp, nil
}

// ToRepoVersionHistoryResponse converts version history to API response
func ToRepoVersionHistoryResponse(language string, versions []entity.RepoVersionInfo) (api.RepoVersionHistoryResponse, error) {
	data := make([]api.RepoVersionInfo, len(versions))
	for i, v := range versions {
		analysisUID, err := uuid.Parse(v.AnalysisID)
		if err != nil {
			return api.RepoVersionHistoryResponse{}, fmt.Errorf("invalid analysis ID %q: %w", v.AnalysisID, err)
		}
		docUID, err := uuid.Parse(v.ID)
		if err != nil {
			return api.RepoVersionHistoryResponse{}, fmt.Errorf("invalid document ID %q: %w", v.ID, err)
		}

		lang := api.SpecLanguage(v.Language)
		data[i] = api.RepoVersionInfo{
			AnalysisID: analysisUID,
			CommitSHA:  v.CommitSHA,
			CreatedAt:  v.CreatedAt,
			ID:         docUID,
			Language:   &lang,
			Version:    v.Version,
		}
		if v.ModelID != "" {
			data[i].ModelID = &v.ModelID
		}
	}

	return api.RepoVersionHistoryResponse{
		Data:     data,
		Language: api.SpecLanguage(language),
	}, nil
}
