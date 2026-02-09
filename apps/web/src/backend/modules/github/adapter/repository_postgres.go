package adapter

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/specvital/web/src/backend/internal/db"
	"github.com/specvital/web/src/backend/modules/github/domain/port"
)

var _ port.Repository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

func NewPostgresRepository(pool *pgxpool.Pool, queries *db.Queries) *PostgresRepository {
	return &PostgresRepository{pool: pool, queries: queries}
}

func (r *PostgresRepository) DeleteOrgRepositories(ctx context.Context, userID, orgID string) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}
	oid, err := parseUUID(orgID)
	if err != nil {
		return err
	}
	return r.queries.DeleteOrgGitHubRepositories(ctx, db.DeleteOrgGitHubRepositoriesParams{
		UserID: uid,
		OrgID:  oid,
	})
}

func (r *PostgresRepository) DeleteUserOrganizations(ctx context.Context, userID string) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserOrgMemberships(ctx, uid)
}

func (r *PostgresRepository) DeleteUserRepositories(ctx context.Context, userID string) error {
	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserGitHubRepositories(ctx, uid)
}

func (r *PostgresRepository) GetOrgIDByLogin(ctx context.Context, login string) (string, error) {
	row, err := r.queries.GetGitHubOrganizationByLogin(ctx, login)
	if err != nil {
		return "", err
	}
	return uuidToString(row.ID), nil
}

func (r *PostgresRepository) GetOrgRepositories(ctx context.Context, userID, orgID string) ([]port.RepositoryRecord, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}
	oid, err := parseUUID(orgID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetOrgGitHubRepositories(ctx, db.GetOrgGitHubRepositoriesParams{
		UserID: uid,
		OrgID:  oid,
	})
	if err != nil {
		return nil, err
	}

	repos := make([]port.RepositoryRecord, len(rows))
	for i, row := range rows {
		repos[i] = mapOrgRepoRowToRecord(row)
	}
	return repos, nil
}

func (r *PostgresRepository) GetUserOrganizations(ctx context.Context, userID string) ([]port.OrganizationRecord, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetUserGitHubOrganizations(ctx, uid)
	if err != nil {
		return nil, err
	}

	orgs := make([]port.OrganizationRecord, len(rows))
	for i, row := range rows {
		orgs[i] = mapOrgRowToRecord(row)
	}
	return orgs, nil
}

func (r *PostgresRepository) GetUserRepositories(ctx context.Context, userID string) ([]port.RepositoryRecord, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	rows, err := r.queries.GetUserGitHubRepositories(ctx, uid)
	if err != nil {
		return nil, err
	}

	repos := make([]port.RepositoryRecord, len(rows))
	for i, row := range rows {
		repos[i] = mapUserRepoRowToRecord(row)
	}
	return repos, nil
}

func (r *PostgresRepository) HasOrgRepositories(ctx context.Context, userID, orgID string) (bool, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return false, err
	}
	oid, err := parseUUID(orgID)
	if err != nil {
		return false, err
	}
	return r.queries.HasOrgGitHubRepositories(ctx, db.HasOrgGitHubRepositoriesParams{
		UserID: uid,
		OrgID:  oid,
	})
}

func (r *PostgresRepository) HasUserOrganizations(ctx context.Context, userID string) (bool, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return false, err
	}
	return r.queries.HasUserOrgMemberships(ctx, uid)
}

func (r *PostgresRepository) HasUserRepositories(ctx context.Context, userID string) (bool, error) {
	uid, err := parseUUID(userID)
	if err != nil {
		return false, err
	}
	return r.queries.HasUserGitHubRepositories(ctx, uid)
}

func (r *PostgresRepository) UpsertOrgRepositories(ctx context.Context, userID, orgID string, repos []port.RepositoryRecord) error {
	if len(repos) == 0 {
		return nil
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}
	oid, err := parseUUID(orgID)
	if err != nil {
		return err
	}

	return r.withTx(ctx, func(qtx *db.Queries) error {
		for _, repo := range repos {
			params := db.UpsertOrgGitHubRepositoryParams{
				Archived:        repo.Archived,
				DefaultBranch:   toPgText(repo.DefaultBranch),
				Description:     toPgText(repo.Description),
				Disabled:        repo.Disabled,
				Fork:            repo.Fork,
				FullName:        repo.FullName,
				GithubRepoID:    repo.ID,
				HtmlUrl:         repo.HTMLURL,
				IsPrivate:       repo.Private,
				Language:        toPgText(repo.Language),
				Name:            repo.Name,
				OrgID:           oid,
				PushedAt:        toPgTimestamptz(repo.PushedAt),
				StargazersCount: int32(repo.StarCount),
				UserID:          uid,
				Visibility:      repo.Visibility,
			}
			if err := qtx.UpsertOrgGitHubRepository(ctx, params); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostgresRepository) UpsertUserOrganizations(ctx context.Context, userID string, orgs []port.OrganizationRecord) error {
	if len(orgs) == 0 {
		return nil
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}

	return r.withTx(ctx, func(qtx *db.Queries) error {
		for _, org := range orgs {
			orgID, err := qtx.UpsertGitHubOrganization(ctx, db.UpsertGitHubOrganizationParams{
				AvatarUrl:   toPgText(org.AvatarURL),
				Description: toPgText(org.Description),
				GithubOrgID: org.ID,
				HtmlUrl:     toPgText(org.HTMLURL),
				Login:       org.Login,
			})
			if err != nil {
				return err
			}

			if err := qtx.UpsertUserOrgMembership(ctx, db.UpsertUserOrgMembershipParams{
				OrgID:  orgID,
				Role:   toPgText(org.Role),
				UserID: uid,
			}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostgresRepository) UpsertUserRepositories(ctx context.Context, userID string, repos []port.RepositoryRecord) error {
	if len(repos) == 0 {
		return nil
	}

	uid, err := parseUUID(userID)
	if err != nil {
		return err
	}

	return r.withTx(ctx, func(qtx *db.Queries) error {
		for _, repo := range repos {
			params := db.UpsertUserGitHubRepositoryParams{
				Archived:        repo.Archived,
				DefaultBranch:   toPgText(repo.DefaultBranch),
				Description:     toPgText(repo.Description),
				Disabled:        repo.Disabled,
				Fork:            repo.Fork,
				FullName:        repo.FullName,
				GithubRepoID:    repo.ID,
				HtmlUrl:         repo.HTMLURL,
				IsPrivate:       repo.Private,
				Language:        toPgText(repo.Language),
				Name:            repo.Name,
				PushedAt:        toPgTimestamptz(repo.PushedAt),
				StargazersCount: int32(repo.StarCount),
				UserID:          uid,
				Visibility:      repo.Visibility,
			}
			if err := qtx.UpsertUserGitHubRepository(ctx, params); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PostgresRepository) withTx(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(r.queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func extractOwnerFromFullName(fullName string) string {
	if parts := strings.SplitN(fullName, "/", 2); len(parts) == 2 {
		return parts[0]
	}
	return ""
}

func parseUUID(s string) (pgtype.UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}, err
	}
	return pgtype.UUID{Bytes: u, Valid: true}, nil
}

func toPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: s, Valid: true}
}

func toPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func uuidToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	return uuid.UUID(u.Bytes).String()
}

type repoRowFields struct {
	Archived        bool
	DefaultBranch   pgtype.Text
	Description     pgtype.Text
	Disabled        bool
	Fork            bool
	FullName        string
	GithubRepoID    int64
	HtmlUrl         string
	IsPrivate       bool
	Language        pgtype.Text
	Name            string
	PushedAt        pgtype.Timestamptz
	StargazersCount int32
	Visibility      string
}

func mapRepoRowFieldsToRecord(f repoRowFields) port.RepositoryRecord {
	rec := port.RepositoryRecord{
		Archived:   f.Archived,
		Disabled:   f.Disabled,
		Fork:       f.Fork,
		FullName:   f.FullName,
		HTMLURL:    f.HtmlUrl,
		ID:         f.GithubRepoID,
		Name:       f.Name,
		Owner:      extractOwnerFromFullName(f.FullName),
		Private:    f.IsPrivate,
		StarCount:  int(f.StargazersCount),
		Visibility: f.Visibility,
	}

	if f.DefaultBranch.Valid {
		rec.DefaultBranch = f.DefaultBranch.String
	}
	if f.Description.Valid {
		rec.Description = f.Description.String
	}
	if f.Language.Valid {
		rec.Language = f.Language.String
	}
	if f.PushedAt.Valid {
		t := f.PushedAt.Time
		rec.PushedAt = &t
	}

	return rec
}

func mapUserRepoRowToRecord(row db.GetUserGitHubRepositoriesRow) port.RepositoryRecord {
	return mapRepoRowFieldsToRecord(repoRowFields{
		Archived:        row.Archived,
		DefaultBranch:   row.DefaultBranch,
		Description:     row.Description,
		Disabled:        row.Disabled,
		Fork:            row.Fork,
		FullName:        row.FullName,
		GithubRepoID:    row.GithubRepoID,
		HtmlUrl:         row.HtmlUrl,
		IsPrivate:       row.IsPrivate,
		Language:        row.Language,
		Name:            row.Name,
		PushedAt:        row.PushedAt,
		StargazersCount: row.StargazersCount,
		Visibility:      row.Visibility,
	})
}

func mapOrgRepoRowToRecord(row db.GetOrgGitHubRepositoriesRow) port.RepositoryRecord {
	return mapRepoRowFieldsToRecord(repoRowFields{
		Archived:        row.Archived,
		DefaultBranch:   row.DefaultBranch,
		Description:     row.Description,
		Disabled:        row.Disabled,
		Fork:            row.Fork,
		FullName:        row.FullName,
		GithubRepoID:    row.GithubRepoID,
		HtmlUrl:         row.HtmlUrl,
		IsPrivate:       row.IsPrivate,
		Language:        row.Language,
		Name:            row.Name,
		PushedAt:        row.PushedAt,
		StargazersCount: row.StargazersCount,
		Visibility:      row.Visibility,
	})
}

func mapOrgRowToRecord(row db.GetUserGitHubOrganizationsRow) port.OrganizationRecord {
	rec := port.OrganizationRecord{
		ID:    row.GithubOrgID,
		Login: row.Login,
		OrgID: uuidToString(row.OrgID),
	}

	if row.AvatarUrl.Valid {
		rec.AvatarURL = row.AvatarUrl.String
	}
	if row.HtmlUrl.Valid {
		rec.HTMLURL = row.HtmlUrl.String
	}
	if row.Description.Valid {
		rec.Description = row.Description.String
	}
	if row.Role.Valid {
		rec.Role = row.Role.String
	}

	return rec
}
