# API Documentation

Base URL: `http://localhost:8000` (default)

## Endpoints

### Health Check

```
GET /health
```

**Response** `200 OK`:

```json
{
  "status": "ok"
}
```

---

### Analyze Repository

```
GET /api/analyze/{owner}/{repo}
```

Analyzes test files in a public GitHub repository.

**Parameters**:

| Name  | Type   | Description      |
| ----- | ------ | ---------------- |
| owner | string | Repository owner |
| repo  | string | Repository name  |

**Response** `200 OK`:

```json
{
  "analyzedAt": "2024-01-15T10:30:00Z",
  "commitSha": "abc123def456",
  "owner": "facebook",
  "repo": "react",
  "suites": [
    {
      "filePath": "packages/react/src/__tests__/React-test.js",
      "framework": "jest",
      "tests": [
        {
          "filePath": "packages/react/src/__tests__/React-test.js",
          "framework": "jest",
          "line": 15,
          "name": "should render correctly",
          "status": "active"
        }
      ]
    }
  ],
  "summary": {
    "active": 150,
    "skipped": 10,
    "todo": 5,
    "total": 165,
    "frameworks": [
      {
        "framework": "jest",
        "active": 100,
        "skipped": 5,
        "todo": 3,
        "total": 108
      }
    ]
  }
}
```

---

## Types

### Framework

```typescript
type Framework = "jest" | "vitest" | "playwright" | "go";
```

### TestStatus

```typescript
type TestStatus = "active" | "skipped" | "todo";
```

### TestCase

```typescript
type TestCase = {
  filePath: string;
  framework: Framework;
  line: number;
  name: string;
  status: TestStatus;
};
```

### TestSuite

```typescript
type TestSuite = {
  filePath: string;
  framework: Framework;
  tests: TestCase[];
};
```

### AnalysisResult

```typescript
type AnalysisResult = {
  analyzedAt: string;
  commitSha: string;
  owner: string;
  repo: string;
  suites: TestSuite[];
  summary: Summary;
};
```

---

## Error Responses

All errors follow [RFC 7807](https://tools.ietf.org/html/rfc7807) Problem Details format.

**Content-Type**: `application/problem+json`

### Error Response Format

```json
{
  "type": "about:blank",
  "title": "Not Found",
  "status": 404,
  "detail": "Repository not found",
  "instance": "/api/analyze/owner/repo",
  "rateLimit": {
    "limit": 60,
    "remaining": 0,
    "resetAt": 1705312200
  }
}
```

### Error Codes

| Status | Title                 | Description                       |
| ------ | --------------------- | --------------------------------- |
| 400    | Bad Request           | Missing or invalid parameters     |
| 403    | Forbidden             | Private repository or no access   |
| 404    | Not Found             | Repository does not exist         |
| 422    | Repository Too Large  | Repository exceeds analysis limit |
| 429    | Too Many Requests     | GitHub API rate limit exceeded    |
| 500    | Internal Server Error | Unexpected server error           |

### Rate Limit Extension

Error responses include `rateLimit` field when available:

| Field     | Type   | Description                      |
| --------- | ------ | -------------------------------- |
| limit     | number | Maximum requests per hour        |
| remaining | number | Remaining requests               |
| resetAt   | number | Unix timestamp when limit resets |

---

## Rate Limiting

GitHub API rate limits apply:

- **Without token**: 60 requests/hour
- **With token**: 5,000 requests/hour

Set `GITHUB_TOKEN` environment variable to increase the limit.

When rate limit is exhausted, the API returns `429 Too Many Requests`.

---

## Examples

### cURL

```bash
# Analyze repository
curl http://localhost:8000/api/analyze/facebook/react

# Health check
curl http://localhost:8000/health
```

### JavaScript

```typescript
const response = await fetch("http://localhost:8000/api/analyze/facebook/react");
const result = await response.json();
console.log(result.summary.total); // Total test count
```
