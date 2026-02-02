package adapter

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/spec-view/domain"
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

func (r *PostgresRepository) GetAnalysisTestCount(ctx context.Context, analysisID string) (int, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.GetAnalysisTestCount(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrAnalysisNotCompleted
		}
		return 0, err
	}

	return int(count), nil
}

func (r *PostgresRepository) CheckSpecDocumentExistsByLanguage(ctx context.Context, analysisID string, language string) (bool, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return false, err
	}

	return r.queries.CheckSpecDocumentExistsByLanguage(ctx, db.CheckSpecDocumentExistsByLanguageParams{
		AnalysisID: uid,
		Language:   language,
	})
}

func (r *PostgresRepository) GetAvailableLanguages(ctx context.Context, analysisID string) ([]entity.AvailableLanguageInfo, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetAvailableLanguagesByAnalysisID(ctx, uid)
	if err != nil {
		return nil, err
	}

	result := make([]entity.AvailableLanguageInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.AvailableLanguageInfo{
			CreatedAt:     row.CreatedAt.Time,
			Language:      row.Language,
			LatestVersion: int(row.LatestVersion),
		}
	}
	return result, nil
}

func (r *PostgresRepository) GetAvailableLanguagesByUser(ctx context.Context, userID string, analysisID string) ([]entity.AvailableLanguageInfo, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetAvailableLanguagesByUserAndAnalysis(ctx, db.GetAvailableLanguagesByUserAndAnalysisParams{
		UserID:     uid,
		AnalysisID: aid,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.AvailableLanguageInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.AvailableLanguageInfo{
			CreatedAt:     row.CreatedAt.Time,
			Language:      row.Language,
			LatestVersion: int(row.LatestVersion),
		}
	}
	return result, nil
}

func (r *PostgresRepository) GetSpecDocumentByLanguage(ctx context.Context, analysisID string, language string) (*entity.SpecDocument, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	var docRow db.GetSpecDocumentByAnalysisIDRow
	if language == "" {
		docRow, err = r.queries.GetSpecDocumentByAnalysisID(ctx, uid)
	} else {
		langRow, langErr := r.queries.GetSpecDocumentByAnalysisIDAndLanguage(ctx, db.GetSpecDocumentByAnalysisIDAndLanguageParams{
			AnalysisID: uid,
			Language:   language,
		})
		if langErr != nil {
			err = langErr
		} else {
			docRow = db.GetSpecDocumentByAnalysisIDRow{
				ID:               langRow.ID,
				AnalysisID:       langRow.AnalysisID,
				UserID:           langRow.UserID,
				Language:         langRow.Language,
				Version:          langRow.Version,
				ExecutiveSummary: langRow.ExecutiveSummary,
				ModelID:          langRow.ModelID,
				CreatedAt:        langRow.CreatedAt,
			}
		}
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildSpecDocument(ctx, docRow)
}

func (r *PostgresRepository) GetSpecDocumentByUser(ctx context.Context, userID string, analysisID string, language string) (*entity.SpecDocument, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	var docRow db.GetSpecDocumentByAnalysisIDRow
	if language == "" {
		row, queryErr := r.queries.GetSpecDocumentByUserAndAnalysis(ctx, db.GetSpecDocumentByUserAndAnalysisParams{
			UserID:     uid,
			AnalysisID: aid,
		})
		if queryErr != nil {
			err = queryErr
		} else {
			docRow = db.GetSpecDocumentByAnalysisIDRow{
				ID:               row.ID,
				AnalysisID:       row.AnalysisID,
				UserID:           row.UserID,
				Language:         row.Language,
				Version:          row.Version,
				ExecutiveSummary: row.ExecutiveSummary,
				ModelID:          row.ModelID,
				CreatedAt:        row.CreatedAt,
			}
		}
	} else {
		langRow, langErr := r.queries.GetSpecDocumentByUserAndAnalysisAndLanguage(ctx, db.GetSpecDocumentByUserAndAnalysisAndLanguageParams{
			UserID:     uid,
			AnalysisID: aid,
			Language:   language,
		})
		if langErr != nil {
			err = langErr
		} else {
			docRow = db.GetSpecDocumentByAnalysisIDRow{
				ID:               langRow.ID,
				AnalysisID:       langRow.AnalysisID,
				UserID:           langRow.UserID,
				Language:         langRow.Language,
				Version:          langRow.Version,
				ExecutiveSummary: langRow.ExecutiveSummary,
				ModelID:          langRow.ModelID,
				CreatedAt:        langRow.CreatedAt,
			}
		}
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildSpecDocument(ctx, docRow)
}

func (r *PostgresRepository) buildSpecDocument(ctx context.Context, docRow db.GetSpecDocumentByAnalysisIDRow) (*entity.SpecDocument, error) {
	domains, err := r.buildDomainHierarchy(ctx, docRow.ID)
	if err != nil {
		return nil, err
	}

	doc := &entity.SpecDocument{
		AnalysisID: uuidToString(docRow.AnalysisID),
		CreatedAt:  docRow.CreatedAt.Time,
		Domains:    domains,
		ID:         uuidToString(docRow.ID),
		Language:   docRow.Language,
		ModelID:    docRow.ModelID,
		Version:    int(docRow.Version),
	}
	if docRow.ExecutiveSummary.Valid {
		doc.ExecutiveSummary = &docRow.ExecutiveSummary.String
	}

	return doc, nil
}

func (r *PostgresRepository) GetGenerationStatus(ctx context.Context, userID string, analysisID string) (*entity.SpecGenerationStatus, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(analysisID); err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecGenerationStatus(ctx, db.GetSpecGenerationStatusParams{
		UserID:     []byte(userID),
		AnalysisID: []byte(analysisID),
	})
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

func (r *PostgresRepository) GetGenerationStatusByLanguage(ctx context.Context, userID string, analysisID string, language string) (*entity.SpecGenerationStatus, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return nil, err
	}
	if _, err := uuid.Parse(analysisID); err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecGenerationStatusByLanguage(ctx, db.GetSpecGenerationStatusByLanguageParams{
		UserID:     []byte(userID),
		AnalysisID: []byte(analysisID),
		Language:   []byte(language),
	})
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

// buildDomainHierarchy fetches and builds the complete domain hierarchy for a spec document.
// This is the common logic shared by buildSpecDocument and buildRepoSpecDocument.
func (r *PostgresRepository) buildDomainHierarchy(ctx context.Context, documentID pgtype.UUID) ([]entity.SpecDomain, error) {
	domains, err := r.queries.GetSpecDomainsByDocumentID(ctx, documentID)
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

	return entityDomains, nil
}

func (r *PostgresRepository) GetSpecDocumentByVersion(ctx context.Context, analysisID string, language string, version int) (*entity.SpecDocument, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	docRow, err := r.queries.GetSpecDocumentByVersion(ctx, db.GetSpecDocumentByVersionParams{
		AnalysisID: uid,
		Language:   language,
		Version:    int32(version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildSpecDocument(ctx, db.GetSpecDocumentByAnalysisIDRow{
		ID:               docRow.ID,
		AnalysisID:       docRow.AnalysisID,
		UserID:           docRow.UserID,
		Language:         docRow.Language,
		Version:          docRow.Version,
		ExecutiveSummary: docRow.ExecutiveSummary,
		ModelID:          docRow.ModelID,
		CreatedAt:        docRow.CreatedAt,
	})
}

func (r *PostgresRepository) GetSpecDocumentByUserAndVersion(ctx context.Context, userID string, analysisID string, language string, version int) (*entity.SpecDocument, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	docRow, err := r.queries.GetSpecDocumentByUserAndVersion(ctx, db.GetSpecDocumentByUserAndVersionParams{
		UserID:     uid,
		AnalysisID: aid,
		Language:   language,
		Version:    int32(version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildSpecDocument(ctx, db.GetSpecDocumentByAnalysisIDRow{
		ID:               docRow.ID,
		AnalysisID:       docRow.AnalysisID,
		UserID:           docRow.UserID,
		Language:         docRow.Language,
		Version:          docRow.Version,
		ExecutiveSummary: docRow.ExecutiveSummary,
		ModelID:          docRow.ModelID,
		CreatedAt:        docRow.CreatedAt,
	})
}

func (r *PostgresRepository) GetVersionsByLanguage(ctx context.Context, analysisID string, language string) ([]entity.VersionInfo, error) {
	uid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetVersionsByLanguage(ctx, db.GetVersionsByLanguageParams{
		AnalysisID: uid,
		Language:   language,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.VersionInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.VersionInfo{
			CreatedAt: row.CreatedAt.Time,
			ModelID:   row.ModelID,
			Version:   int(row.Version),
		}
	}
	return result, nil
}

func (r *PostgresRepository) GetVersionsByUser(ctx context.Context, userID string, analysisID string, language string) ([]entity.VersionInfo, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetVersionsByUserAndLanguage(ctx, db.GetVersionsByUserAndLanguageParams{
		UserID:     uid,
		AnalysisID: aid,
		Language:   language,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.VersionInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.VersionInfo{
			CreatedAt: row.CreatedAt.Time,
			ModelID:   row.ModelID,
			Version:   int(row.Version),
		}
	}
	return result, nil
}

func (r *PostgresRepository) HasPreviousSpecByLanguage(ctx context.Context, userID, analysisID, language string) (bool, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return false, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return false, err
	}

	return r.queries.HasPreviousSpecByLanguage(ctx, db.HasPreviousSpecByLanguageParams{
		CurrentAnalysisID: aid,
		UserID:            uid,
		Language:          language,
	})
}

func (r *PostgresRepository) GetLanguagesWithPreviousSpec(ctx context.Context, userID, analysisID string) ([]string, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	aid, err := parseUUID(analysisID)
	if err != nil {
		return nil, err
	}

	return r.queries.GetLanguagesWithPreviousSpec(ctx, db.GetLanguagesWithPreviousSpecParams{
		CurrentAnalysisID: aid,
		UserID:            uid,
	})
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

// Repository-based queries (cross-analysis access)

func (r *PostgresRepository) CheckCodebaseExists(ctx context.Context, owner, name string) (bool, error) {
	return r.queries.CheckCodebaseExists(ctx, db.CheckCodebaseExistsParams{
		Owner: owner,
		Name:  name,
	})
}

func (r *PostgresRepository) GetSpecDocumentByRepository(ctx context.Context, userID, owner, name, language string) (*entity.RepoSpecDocument, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecDocumentByRepository(ctx, db.GetSpecDocumentByRepositoryParams{
		Owner:    owner,
		Repo:     name,
		UserID:   uid,
		Language: language,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildRepoSpecDocument(ctx, row)
}

func (r *PostgresRepository) GetSpecDocumentByRepositoryAndVersion(ctx context.Context, userID, owner, name, language string, version int) (*entity.RepoSpecDocument, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecDocumentByRepositoryAndVersion(ctx, db.GetSpecDocumentByRepositoryAndVersionParams{
		Owner:    owner,
		Repo:     name,
		UserID:   uid,
		Language: language,
		Version:  int32(version),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildRepoSpecDocument(ctx, db.GetSpecDocumentByRepositoryRow{
		ID:               row.ID,
		AnalysisID:       row.AnalysisID,
		UserID:           row.UserID,
		Language:         row.Language,
		Version:          row.Version,
		ExecutiveSummary: row.ExecutiveSummary,
		ModelID:          row.ModelID,
		CreatedAt:        row.CreatedAt,
		CommitSha:        row.CommitSha,
	})
}

func (r *PostgresRepository) GetSpecDocumentByRepositoryAndDocumentId(ctx context.Context, userID, owner, name, documentID string) (*entity.RepoSpecDocument, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	docID, err := parseUUID(documentID)
	if err != nil {
		return nil, err
	}

	row, err := r.queries.GetSpecDocumentByRepositoryAndDocumentId(ctx, db.GetSpecDocumentByRepositoryAndDocumentIdParams{
		Owner:      owner,
		Repo:       name,
		UserID:     uid,
		DocumentID: docID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.buildRepoSpecDocument(ctx, db.GetSpecDocumentByRepositoryRow{
		ID:               row.ID,
		AnalysisID:       row.AnalysisID,
		UserID:           row.UserID,
		Language:         row.Language,
		Version:          row.Version,
		ExecutiveSummary: row.ExecutiveSummary,
		ModelID:          row.ModelID,
		CreatedAt:        row.CreatedAt,
		CommitSha:        row.CommitSha,
	})
}

func (r *PostgresRepository) buildRepoSpecDocument(ctx context.Context, row db.GetSpecDocumentByRepositoryRow) (*entity.RepoSpecDocument, error) {
	domains, err := r.buildDomainHierarchy(ctx, row.ID)
	if err != nil {
		return nil, err
	}

	doc := &entity.RepoSpecDocument{
		AnalysisID: uuidToString(row.AnalysisID),
		CommitSHA:  row.CommitSha,
		CreatedAt:  row.CreatedAt.Time,
		Domains:    domains,
		ID:         uuidToString(row.ID),
		Language:   row.Language,
		ModelID:    row.ModelID,
		Version:    int(row.Version),
	}
	if row.ExecutiveSummary.Valid {
		doc.ExecutiveSummary = &row.ExecutiveSummary.String
	}

	return doc, nil
}

func (r *PostgresRepository) GetVersionHistoryByRepository(ctx context.Context, userID, owner, name, language string) ([]entity.RepoVersionInfo, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetVersionHistoryByRepository(ctx, db.GetVersionHistoryByRepositoryParams{
		Owner:    owner,
		Repo:     name,
		UserID:   uid,
		Language: language,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.RepoVersionInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.RepoVersionInfo{
			AnalysisID: uuidToString(row.AnalysisID),
			CommitSHA:  row.CommitSha,
			CreatedAt:  row.CreatedAt.Time,
			ID:         uuidToString(row.ID),
			Language:   row.Language,
			ModelID:    row.ModelID,
			Version:    int(row.Version),
		}
	}
	return result, nil
}

func (r *PostgresRepository) GetAvailableLanguagesByRepository(ctx context.Context, userID, owner, name string) ([]entity.AvailableLanguageInfo, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetAvailableLanguagesByRepository(ctx, db.GetAvailableLanguagesByRepositoryParams{
		Owner:  owner,
		Repo:   name,
		UserID: uid,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.AvailableLanguageInfo, len(rows))
	for i, row := range rows {
		result[i] = entity.AvailableLanguageInfo{
			CreatedAt:     row.LatestCreatedAt.Time,
			Language:      row.Language,
			LatestVersion: int(row.LatestVersion),
		}
	}
	return result, nil
}
