package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/entity"
	"github.com/specvital/web/src/backend/modules/spec-view/domain/port"
)

type PostgresRepository struct {
	queries *db.Queries
}

var _ port.SpecViewRepository = (*PostgresRepository)(nil)

func NewPostgresRepository(queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{queries: queries}
}

func (r *PostgresRepository) CheckAnalysisExists(ctx context.Context, analysisID string) (bool, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return false, err
	}

	return r.queries.CheckAnalysisExists(ctx, uid)
}

func (r *PostgresRepository) CheckSpecDocumentExists(ctx context.Context, analysisID string) (bool, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return false, err
	}

	return r.queries.CheckSpecDocumentExists(ctx, uid)
}

func (r *PostgresRepository) GetSpecDocument(ctx context.Context, analysisID string) (*entity.SpecDocument, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	docRow, err := r.queries.GetSpecDocumentByAnalysisID(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	domains, err := r.queries.GetSpecDomainsByDocumentID(ctx, docRow.ID)
	if err != nil {
		return nil, err
	}

	domainIDs := make([]pgtype.UUID, len(domains))
	for i, d := range domains {
		domainIDs[i] = d.ID
	}

	features, err := r.queries.GetSpecFeaturesByDomainIDs(ctx, domainIDs)
	if err != nil {
		return nil, err
	}

	featureIDs := make([]pgtype.UUID, len(features))
	for i, f := range features {
		featureIDs[i] = f.ID
	}

	behaviors, err := r.queries.GetSpecBehaviorsByFeatureIDs(ctx, featureIDs)
	if err != nil {
		return nil, err
	}

	behaviorIDs := make([]pgtype.UUID, 0, len(behaviors))
	for _, b := range behaviors {
		if b.SourceTestCaseID.Valid {
			behaviorIDs = append(behaviorIDs, b.ID)
		}
	}

	sourceInfoMap := make(map[string]*entity.BehaviorSourceInfo)
	if len(behaviorIDs) > 0 {
		sourceInfoRows, err := r.queries.GetSpecBehaviorSourceInfo(ctx, behaviorIDs)
		if err != nil {
			return nil, err
		}
		for _, si := range sourceInfoRows {
			sourceInfoMap[uuidToString(si.BehaviorID)] = &entity.BehaviorSourceInfo{
				FilePath:   si.FilePath,
				Framework:  si.Framework.String,
				LineNumber: int(si.LineNumber.Int32),
				Status:     si.TcStatus,
			}
		}
	}

	featureMap := make(map[string][]entity.SpecBehavior)
	for _, b := range behaviors {
		fid := uuidToString(b.FeatureID)
		bid := uuidToString(b.ID)
		behavior := entity.SpecBehavior{
			ConvertedDescription: b.ConvertedDescription,
			ID:                   bid,
			OriginalName:         b.OriginalName,
			SortOrder:            int(b.SortOrder),
			SourceInfo:           sourceInfoMap[bid],
		}
		if b.SourceTestCaseID.Valid {
			tcID := uuidToString(b.SourceTestCaseID)
			behavior.SourceTestCaseID = &tcID
		}
		featureMap[fid] = append(featureMap[fid], behavior)
	}

	domainMap := make(map[string][]entity.SpecFeature)
	for _, f := range features {
		did := uuidToString(f.DomainID)
		fid := uuidToString(f.ID)
		feature := entity.SpecFeature{
			Behaviors: featureMap[fid],
			ID:        fid,
			Name:      f.Name,
			SortOrder: int(f.SortOrder),
		}
		if f.Description.Valid {
			feature.Description = &f.Description.String
		}
		domainMap[did] = append(domainMap[did], feature)
	}

	entityDomains := make([]entity.SpecDomain, len(domains))
	for i, d := range domains {
		did := uuidToString(d.ID)
		entityDomains[i] = entity.SpecDomain{
			Features:  domainMap[did],
			ID:        did,
			Name:      d.Name,
			SortOrder: int(d.SortOrder),
		}
		if d.Description.Valid {
			entityDomains[i].Description = &d.Description.String
		}
		if d.ClassificationConfidence.Valid {
			f, _ := d.ClassificationConfidence.Float64Value()
			entityDomains[i].ClassificationConfidence = &f.Float64
		}
	}

	doc := &entity.SpecDocument{
		AnalysisID: uuidToString(docRow.AnalysisID),
		CreatedAt:  docRow.CreatedAt.Time,
		Domains:    entityDomains,
		ID:         uuidToString(docRow.ID),
		Language:   docRow.Language,
		ModelID:    docRow.ModelID,
	}
	if docRow.ExecutiveSummary.Valid {
		doc.ExecutiveSummary = &docRow.ExecutiveSummary.String
	}

	return doc, nil
}

func (r *PostgresRepository) GetGenerationStatus(ctx context.Context, analysisID string) (*entity.SpecGenerationStatus, error) {
	if _, err := uuid.Parse(analysisID); err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecGenerationStatus(ctx, []byte(analysisID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	status := &entity.SpecGenerationStatus{
		AnalysisID: analysisID,
		Status:     mapRiverJobState(string(row.State)),
	}

	if row.CreatedAt.Valid {
		status.StartedAt = &row.CreatedAt.Time
	}
	if row.FinalizedAt.Valid {
		status.CompletedAt = &row.FinalizedAt.Time
	}

	return status, nil
}

func parseUUID(s string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: parsed, Valid: true}, nil
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

func mapRiverJobState(state string) entity.GenerationStatus {
	switch state {
	case "available", "pending", "scheduled", "retryable":
		return entity.StatusPending
	case "running":
		return entity.StatusRunning
	case "completed":
		return entity.StatusCompleted
	case "discarded", "cancelled":
		return entity.StatusFailed
	default:
		return entity.StatusNotFound
	}
}
