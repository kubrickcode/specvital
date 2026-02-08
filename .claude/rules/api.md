# API Rules

## Common

### Field Naming

- Boolean: `is/has/can` prefix
- Date: `~At` suffix (ISO 8601 UTC)
- Consistent terminology (unify on "create" or "add", etc.)

### Pagination (Cursor-Based)

- REST: `?cursor=xyz&limit=20` → `{ data, nextCursor, hasNext }`
- GraphQL: Relay Connection (`first`, `after`, `PageInfo`)

### Sorting

- Parameters: `sortBy`, `sortOrder` (REST) or `orderBy` array (GraphQL)
- Support multiple criteria
- Specify defaults clearly

### Filtering

- Range: `{ min, max }` or `{ gte, lte }`
- Complex conditions: nested objects

## REST

- Nested resources: max 2 levels
- Verbs only when not expressible as resource (`/users/:id/activate`)
- List response: `data` + pagination info
- Creation: 201 + resource (exclude sensitive info)
- Error: RFC 7807 ProblemDetail (`type`, `title`, `status`, `detail`, `instance`)
- Batch: `/batch` suffix with success/failure counts

## GraphQL

### Type Naming

- Input: `{Verb}{Type}Input`
- Connection: `{Type}Connection`
- Edge: `{Type}Edge`

### Input Design

- Separate create/update (required vs optional fields)
- Avoid nesting - use IDs only

### Error Handling

- Default: `code`, `field` in `errors[].extensions`
- Type-safe: Union types (`User | ValidationError`)

### Performance

- N+1: DataLoader mandatory

### Documentation

- `"""description"""` required for all types
- State input constraints explicitly
- Deprecation: `@deprecated(reason: "...")`, never delete types
