# Commit Message

## 한국어

```
feat(domain-hints): Python 테스트 파일 DomainHints 추출 구현

Python 테스트 파일에서 import 문과 function call 정보를 tree-sitter 쿼리로 추출하는 PythonExtractor 구현.

- `import x`, `from x import y` 구문 파싱 지원
- 상대 경로 import 지원 (`.models`, `..services`)
- function call 2-segment normalization 적용
- pytest/unittest/mock 프레임워크 call 필터링
- pytest 빌트인 fixtures 필터링 (raises, monkeypatch, caplog 등)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

## English

```
feat(domain-hints): implement DomainHints extraction for Python test files

Implement PythonExtractor to extract import statements and function calls from Python test files using tree-sitter queries.

- Support `import x`, `from x import y` syntax parsing
- Support relative imports (`.models`, `..services`)
- Apply 2-segment normalization for function calls
- Filter pytest/unittest/mock framework calls
- Filter pytest built-in fixtures (raises, monkeypatch, caplog, etc.)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```
