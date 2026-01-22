-- name: GetSpecDocumentByAnalysisID :one
SELECT
    sd.id,
    sd.analysis_id,
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
SELECT
    sd.id,
    sd.analysis_id,
    sd.language,
    sd.version,
    sd.executive_summary,
    sd.model_id,
    sd.created_at
FROM spec_documents sd
WHERE sd.analysis_id = $1 AND sd.language = $2
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
-- Returns latest generation status for an analysis (any language)
SELECT
    rj.state,
    rj.created_at,
    rj.finalized_at,
    rj.errors
FROM river_job rj
WHERE rj.kind = 'specview:generate'
  AND rj.args->>'analysis_id' = $1
ORDER BY rj.created_at DESC
LIMIT 1;

-- name: GetSpecGenerationStatusByLanguage :one
-- Returns generation status for a specific analysis + language combination
SELECT
    rj.state,
    rj.created_at,
    rj.finalized_at,
    rj.errors
FROM river_job rj
WHERE rj.kind = 'specview:generate'
  AND rj.args->>'analysis_id' = @analysis_id
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

-- name: GetAvailableLanguagesByAnalysisID :many
-- Returns all available languages for an analysis with their latest version info
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

