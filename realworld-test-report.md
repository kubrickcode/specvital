# Parser Validation Report: nunit

**Date**: 2025-12-26
**Repository**: [nunit/nunit](https://github.com/nunit/nunit)
**Framework**: NUnit

---

## 📊 비교 결과

| Source                                | Test Count    |
| ------------------------------------- | ------------- |
| **Ground Truth** (AI Manual Analysis) | 4,455         |
| **SpecVital Parser**                  | 3,575         |
| **Delta**                             | -880 (-19.8%) |

**Status**: ❌ FAIL

---

## 🔍 Ground Truth 세부 정보

**Method**: AI Manual Analysis

**.NET SDK 10.0이 필요하나 devcontainer에는 9.0만 설치되어 CLI 실행 불가:**

```bash
# CLI 시도 시 에러
error NETSDK1045: The current .NET SDK does not support targeting .NET 10.0.
```

**수동 분석 접근법**:

테스트 프로젝트 디렉토리에서 NUnit 어트리뷰트별로 카운트:

| Attribute                   | Count |
| --------------------------- | ----- |
| `[Test]` (standalone)       | 2,972 |
| `[Test(...)]` (with params) | 24    |
| `[TestCase(...)]`           | 1,298 |
| `[TestCaseSource(...)]`     | 147   |
| `[Theory]`                  | 14    |

**분석 범위**:

- `src/NUnitFramework/tests/` (312 files)
- `src/NUnitFramework/nunit.framework.legacy.tests/` (25 files)
- `src/NUnitFramework/nunitlite.tests/` (15 files)
- `src/NUnitFramework/windows-tests/` (2 files)
- `src/NUnitFramework/slow-tests/` (1 file)
- `CakeScripts.Tests/` (1 file)

**제외**: `testdata/`, `mock-assembly/` (fixture data)

---

## 📈 Parser 결과

| Metric         | Value |
| -------------- | ----- |
| Files Scanned  | 361   |
| Files Matched  | 361   |
| Tests Detected | 3,575 |
| Duration       | 190ms |

### Framework 분포

| Framework | Files | Tests |
| --------- | ----- | ----- |
| NUnit     | 361   | 3,575 |

---

## 🐛 불일치 분석

### 핵심 버그: [TestCase] 어트리뷰트 카운트 누락

**현재 파서 동작** (`pkg/parser/strategies/nunit/definition.go:223-276`):

```go
func parseTestMethod(...) *domain.Test {
    for _, attr := range attributes {
        switch name {
        case "TestCase", "TestCaseAttribute":
            isTest = true  // ← 여러 [TestCase]가 있어도 true 한 번만 설정
        }
    }
    // 메서드당 Test 1개만 반환
    return &domain.Test{...}
}
```

**문제점**:

- `[TestCase]` 어트리뷰트가 여러 개 있어도 메서드당 1개로만 카운트
- NUnit에서 `[TestCase]`는 parameterized test로, 각 어트리뷰트가 독립적인 테스트 케이스임

### 예시: WarningTests.cs

| Method                         | [TestCase] Count | Parser Result | Expected |
| ------------------------------ | ---------------- | ------------- | -------- |
| `WarningPasses`                | 32               | 1             | 32       |
| `WarningFails`                 | 32               | 1             | 32       |
| `WarningUsedInSetUpOrTearDown` | 12               | 1             | 12       |

**파일 전체**:

- Ground Truth: ~86-89 테스트
- Parser Result: 8 테스트
- Delta: -78 테스트 (1개 파일에서만)

### 누락된 패턴

1. **[TestCase] 다중 어트리뷰트**
   - 영향: 약 1,298개 [TestCase] 중 대부분 누락

2. **[TestCaseSource] 처리**
   - 현재: 1개로 카운트 (정확)
   - 런타임 확장 불가능하므로 현재 동작 적절

3. **[Theory] 처리**
   - 현재: 감지 안됨
   - 영향: 14개 테스트 누락

### Root Cause

`parseTestMethod` 함수가 메서드 단위로 Test 객체 1개만 반환하도록 설계됨.
NUnit의 `[TestCase]` parameterized test 패턴을 지원하려면 각 어트리뷰트당 1개의 Test를 생성해야 함.

---

## 📋 결론

> Parser에 심각한 정확도 문제가 있음. repos.yaml에 추가하지 마세요.

**수정 필요 사항**:

1. **[TestCase] 다중 어트리뷰트 지원**
   - 각 `[TestCase]` 어트리뷰트를 개별 테스트로 카운트
   - 예상 영향: +1,200 테스트 이상

2. **[Theory] 어트리뷰트 지원**
   - `[Theory]` 인식 및 카운트 추가
   - 예상 영향: +14 테스트

3. **testdata 폴더 필터링 검토**
   - 현재 `testdata/` 일부 파일이 테스트로 감지됨 (147개 테스트)
   - Ground Truth에서는 제외되어야 함

---

## 📝 다음 단계

- [ ] `parseTestMethod` 함수 수정: `[TestCase]` 어트리뷰트 개수만큼 Test 생성
- [ ] `[Theory]` 어트리뷰트 지원 추가
- [ ] testdata 경로 필터링 로직 검토
- [ ] 수정 후 재검증 실행
