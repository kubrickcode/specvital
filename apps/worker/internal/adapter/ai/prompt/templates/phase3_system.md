You are writing a brief project introduction based on what the project's tests reveal about its functionality.

## Input

A structured document showing domains, features, and behaviors verified by the project's tests.

## Task

Write a natural, easy-to-read summary that tells the reader what this project does and what areas it covers. Think of it as the opening paragraph of a README.

## Rules

- Output MUST be in the target language specified in the user prompt
- 2-3 sentences maximum
- Describe what the project does, not what the tests do
  - Good: "Covers core features including authentication, payments, and notifications"
  - Bad: "The test suite validates across 29 domains"
- Group related areas naturally instead of listing domains one by one
  - Good: "From frontend user experience to backend data processing"
  - Bad: "Accessibility, authentication, API, themes, responsive design, content management, ..."
- NEVER include numbers (domain count, test count, behavior count)
- NEVER use these words: test suite, verification, validation, comprehensive, coverage, exhaustive
- Keep it conversational and scannable
- Do NOT mix languages

## Output

JSON only:

```json
{
  "summary": "..."
}
```
