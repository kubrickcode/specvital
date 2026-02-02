-- name: GetSpecDocumentByAnalysisID :one
-- (legacy - no user filter) Returns most recent spec document for an analysis
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = $1
ORDER BY sd.version DESC
LIMIT 1;

-- name: GetSpecDocumentByAnalysisIDAndLanguage :one
-- (legacy - no user filter) Returns most recent spec document for an analysis and language
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = $1 AND sd.language = $2
ORDER BY sd.version DESC
LIMIT 1;

-- name: GetSpecDocumentByUserAndAnalysis :one
-- Returns spec document for a specific user and analysis (for access control)
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = @analysis_id AND sd.user_id = @user_id
ORDER BY sd.version DESC
LIMIT 1;

-- name: GetSpecDocumentByUserAndAnalysisAndLanguage :one
-- Returns spec document for a specific user, analysis, and language (for access control)
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = @analysis_id AND sd.user_id = @user_id AND sd.language = @language
ORDER BY sd.version DESC
LIMIT 1;

-- name: GetSpecDomainsByDocumentID :many
SELECT
    d.id,
    d.name,
    d.description,
    d.sort_order,
    d.classification_confidence
FROM spec_domains d
WHERE d.document_id = $1
ORDER BY d.sort_order, d.name;

-- name: GetSpecFeaturesByDomainIDs :many
SELECT
    f.id,
    f.domain_id,
    f.name,
    f.description,
    f.sort_order
FROM spec_features f
WHERE f.domain_id = ANY($1::uuid[])
ORDER BY f.domain_id, f.sort_order, f.name;

-- name: GetSpecBehaviorsByFeatureIDs :many
SELECT
    b.id,
    b.feature_id,
    b.original_name,
    b.converted_description,
    b.sort_order,
    b.source_test_case_id
FROM spec_behaviors b
WHERE b.feature_id = ANY($1::uuid[])
ORDER BY b.feature_id, b.sort_order;

-- name: GetSpecBehaviorSourceInfo :many
SELECT
    b.id AS behavior_id,
    tc.line_number,
    tc.status::text,
    tf.file_path,
    tf.framework
FROM spec_behaviors b
JOIN test_cases tc ON tc.id = b.source_test_case_id
JOIN test_suites ts ON ts.id = tc.suite_id
JOIN test_files tf ON tf.id = ts.file_id
WHERE b.id = ANY($1::uuid[]);

-- name: GetSpecGenerationStatus :one
-- Returns latest generation status for a specific user and analysis (any language)
SELECT
    rj.state,
    rj.created_at,
    rj.finalized_at,
    rj.errors
FROM river_job rj
WHERE rj.kind = 'specview:generate'
  AND rj.args->>'analysis_id' = @analysis_id
  AND rj.args->>'user_id' = @user_id
ORDER BY rj.created_at DESC
LIMIT 1;

-- name: GetSpecGenerationStatusByLanguage :one
-- Returns generation status for a specific user, analysis, and language combination
SELECT
    rj.state,
    rj.created_at,
    rj.finalized_at,
    rj.errors
FROM river_job rj
WHERE rj.kind = 'specview:generate'
  AND rj.args->>'analysis_id' = @analysis_id
  AND rj.args->>'user_id' = @user_id
  AND rj.args->>'language' = @language
ORDER BY rj.created_at DESC
LIMIT 1;

-- name: CheckSpecDocumentExistsByLanguage :one
SELECT EXISTS(
    SELECT 1 FROM spec_documents WHERE analysis_id = $1 AND language = $2
) AS exists;

-- name: CheckAnalysisExists :one
SELECT EXISTS(
    SELECT 1 FROM analyses WHERE id = $1 AND status = 'completed'
) AS exists;

-- name: GetAnalysisTestCount :one
-- Returns the total test count for a completed analysis.
-- Used for quota calculation before spec generation.
SELECT total_tests FROM analyses WHERE id = $1 AND status = 'completed';

-- name: GetAvailableLanguagesByAnalysisID :many
-- Returns all available languages for an analysis with their latest version info (legacy - no user filter)
SELECT
    sd.language,
    sd.version AS latest_version,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = $1
  AND sd.version = (
      SELECT MAX(sd2.version)
      FROM spec_documents sd2
      WHERE sd2.analysis_id = sd.analysis_id
        AND sd2.language = sd.language
  )
ORDER BY sd.language;

-- name: GetAvailableLanguagesByUserAndAnalysis :many
-- Returns all available languages for a specific user and analysis with their latest version info
SELECT
    sd.language,
    sd.version AS latest_version,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = @analysis_id
  AND sd.user_id = @user_id
  AND sd.version = (
      SELECT MAX(sd2.version)
      FROM spec_documents sd2
      WHERE sd2.analysis_id = sd.analysis_id
        AND sd2.user_id = sd.user_id
        AND sd2.language = sd.language
  )
ORDER BY sd.language;

-- name: GetVersionsByLanguage :many
-- Returns all versions for a specific analysis and language, ordered by version descending (legacy - no user filter)
SELECT
    sd.version,
    sd.created_at,
    sd.model_id
FROM spec_documents sd
WHERE sd.analysis_id = $1 AND sd.language = $2
ORDER BY sd.version DESC;

-- name: GetVersionsByUserAndLanguage :many
-- Returns all versions for a specific user, analysis and language, ordered by version descending
SELECT
    sd.version,
    sd.created_at,
    sd.model_id
FROM spec_documents sd
WHERE sd.analysis_id = @analysis_id AND sd.user_id = @user_id AND sd.language = @language
ORDER BY sd.version DESC;

-- name: GetSpecDocumentByVersion :one
-- Returns a specific version of a spec document (legacy - no user filter)
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = $1 AND sd.language = $2 AND sd.version = $3;

-- name: GetSpecDocumentByUserAndVersion :one
-- Returns a specific version of a spec document for a specific user
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = @analysis_id AND sd.user_id = @user_id AND sd.language = @language AND sd.version = @version;

-- name: GetAiSpecSummariesByAnalysisIDs :many
-- Returns AI Spec summary aggregation for multiple analysis IDs
-- Used for Dashboard RepositoryCard to show [AI Spec] badge
SELECT
    sd.analysis_id,
    sd.user_id,
    COUNT(DISTINCT sd.language) AS language_count,
    MAX(sd.created_at) AS latest_generated_at
FROM spec_documents sd
WHERE sd.analysis_id = ANY(@analysis_ids::uuid[])
  AND sd.user_id = @user_id
GROUP BY sd.analysis_id, sd.user_id;

-- name: HasPreviousSpecByLanguage :one
-- Checks if user has generated a spec document for the same codebase but different analysis
-- Used to determine if behavior cache might exist for the selected language
SELECT EXISTS(
    SELECT 1 FROM spec_documents sd
    JOIN analyses a ON a.id = sd.analysis_id
    WHERE a.codebase_id = (SELECT a2.codebase_id FROM analyses a2 WHERE a2.id = @current_analysis_id)
      AND sd.user_id = @user_id
      AND sd.language = @language
      AND sd.analysis_id != @current_analysis_id
) AS has_previous_spec;

-- name: GetLanguagesWithPreviousSpec :many
-- Returns all languages where user has a spec document for the same codebase but different analysis
-- Batch version of HasPreviousSpecByLanguage to avoid N+1 queries
SELECT DISTINCT sd.language
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
WHERE a.codebase_id = (SELECT a2.codebase_id FROM analyses a2 WHERE a2.id = @current_analysis_id)
  AND sd.user_id = @user_id
  AND sd.analysis_id != @current_analysis_id;

-- name: GetSpecDocumentByRepository :one
-- Returns latest spec document for a repository (owner/repo) via JOIN
-- For repository-centric API that provides cross-analysis version access
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at,
    a.commit_sha
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
JOIN codebases c ON c.id = a.codebase_id
WHERE c.owner = @owner AND c.name = @repo
  AND sd.user_id = @user_id
  AND sd.language = @language
ORDER BY sd.created_at DESC, sd.version DESC
LIMIT 1;

-- name: GetSpecDocumentByRepositoryAndVersion :one
-- Returns specific version of spec document for a repository
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at,
    a.commit_sha
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
JOIN codebases c ON c.id = a.codebase_id
WHERE c.owner = @owner AND c.name = @repo
  AND sd.user_id = @user_id
  AND sd.language = @language
  AND sd.version = @version;

-- name: GetSpecDocumentByRepositoryAndDocumentId :one
-- Returns specific spec document by document ID for a repository
SELECT
    sd.id,
    sd.analysis_id,
    sd.user_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at,
    a.commit_sha
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
JOIN codebases c ON c.id = a.codebase_id
WHERE c.owner = @owner AND c.name = @repo
  AND sd.user_id = @user_id
  AND sd.id = @document_id;

-- name: GetVersionHistoryByRepository :many
-- Returns spec versions for a repository (across analyses) with commit SHA
-- Ordered by creation date descending to show most recent first
-- Limited to 100 versions to prevent unbounded result sets
SELECT
    sd.id,
    sd.analysis_id,
    sd.version,
    sd.language,
    sd.model_id,
    sd.created_at,
    a.commit_sha
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
JOIN codebases c ON c.id = a.codebase_id
WHERE c.owner = @owner AND c.name = @repo
  AND sd.user_id = @user_id
  AND sd.language = @language
ORDER BY sd.created_at DESC
LIMIT 100;

-- name: GetAvailableLanguagesByRepository :many
-- Returns all available languages for a repository with their latest version info
SELECT
    sd.language,
    MAX(sd.version)::int AS latest_version,
    MAX(sd.created_at)::timestamptz AS latest_created_at
FROM spec_documents sd
JOIN analyses a ON a.id = sd.analysis_id
JOIN codebases c ON c.id = a.codebase_id
WHERE c.owner = @owner AND c.name = @repo
  AND sd.user_id = @user_id
GROUP BY sd.language
ORDER BY sd.language;

-- name: CheckCodebaseExists :one
-- Checks if codebase exists by owner/name
SELECT EXISTS(
    SELECT 1 FROM codebases WHERE owner = @owner AND name = @name
) AS exists;
