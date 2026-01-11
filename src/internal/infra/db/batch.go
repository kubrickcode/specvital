package db

const InsertTestSuiteBatch = `
INSERT INTO test_suites (file_id, parent_id, name, line_number, depth)
VALUES ($1, $2, $3, $4, $5)
RETURNING id`

var TestCaseCopyColumns = []string{"suite_id", "name", "line_number", "status", "tags", "modifier"}
