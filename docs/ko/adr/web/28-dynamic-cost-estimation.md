---
title: 동적 비용 예측
description: 스펙 생성 모달의 캐시 인식 비용 예측 ADR
---

# ADR-28: 동적 비용 예측

> 🇺🇸 [English Version](/en/adr/web/28-dynamic-cost-estimation.md)

| 날짜       | 작성자       | 리포지토리 |
| ---------- | ------------ | ---------- |
| 2026-02-03 | @KubrickCode | web        |

## Context

### 문제 상황

스펙 생성 모달에서 AI 기반 테스트 행동 생성 전 예상 사용량 표시. [ADR-13: 빌링 및 쿼터 아키텍처](/ko/adr/13-billing-quota-architecture.md)에 따라 사용자는 월간 쿼터 제한이 있으며, 캐시 히트는 쿼터 미소모.

그러나 캐시 옵션 선택과 무관하게 예상 사용량이 동일하게 표시:

| 사용자 선택 | 예상 동작                 | 실제 동작        |
| ----------- | ------------------------- | ---------------- |
| 캐시 사용   | 감소된 예상치 (캐시 히트) | 새로 분석과 동일 |
| 새로 분석   | 전체 테스트 수            | 캐시 사용과 동일 |

발생 원인:

1. 예상치가 `current_test_count`로 하드코딩
2. 캐시 히트율 예측 메커니즘 부재
3. Worker의 캐시 조회 알고리즘 Web 레이어에서 접근 불가

### 기술적 과제

Worker 서비스는 실제 캐시 키 조회에 `SHA256(normalize(test_name))` 사용 ([Worker ADR-12](/ko/adr/worker/12-phase2-behavior-cache.md)). Web에서 복제 시 필요 사항:

- 리포지토리 간 정규화 로직 공유
- 프론트엔드 SHA256 해시 계산
- Worker 알고리즘 변경 시 동기화 유지

## Decision

**이전 스펙 행동과의 이름+경로 매칭을 통한 클라이언트 측 캐시 히트 예측 구현.**

### 알고리즘

```
현재 분석의 각 테스트에 대해:
  1. 이전 spec_document.behaviors에서 매칭 테스트 탐색
     매칭 기준: test.name == behavior.test_name
              AND test.file_path == behavior.file_path
  2. 매칭 발견 시: 캐시 히트 예측
  3. 매칭 미발견 시: 캐시 미스 예측

예상_비용 = 현재_테스트 - 예측_캐시_히트
```

### 설계 선택

| 선택               | 근거                                    |
| ------------------ | --------------------------------------- |
| 근사 매칭          | 단순 구현으로 ~99% 정확도 달성          |
| 클라이언트 측 예측 | API 왕복 불필요; 즉시 모달 렌더링       |
| 고지 툴팁          | "실제 사용량은 예상과 다를 수 있음"     |
| 우아한 성능 저하   | 이전 분석 없을 시 → 전체 테스트 수 표시 |

## Options Considered

### Option A: 이름+경로 예측 (채택)

이전 스펙 문서 행동과 `name + file_path` 비교를 통한 테스트 매칭.

**장점:**

- 제로 레이턴시: 기존 데이터로 클라이언트 측 계산
- API 의존성 없음: 오프라인 동작
- 단순 구현: 문자열 비교만
- 유지보수 용이: Worker 정규화 변경과 독립
- 높은 정확도: 실제 캐시 히트와 ~99% 상관

**단점:**

- 예측 불일치: 이름+경로 일치하나 정규화 해시 불일치 엣지 케이스
- 거짓 음성: 경로 변경된 동일 논리 테스트가 미스로 표시
- 교차 리포지토리 미지원: 다른 리포지토리 캐시 히트 예측 불가

### Option B: Worker의 SHA256 해시 알고리즘 복제

Web 프론트엔드에 동일한 정규화 및 해싱 구현.

**장점:**

- 실제 캐시 동작과 100% 정확도
- 교차 리포지토리 예측 이론적 가능

**단점:**

- 코드 중복: Worker(Go)와 Web(TypeScript)에 정규화 로직
- 동기화 유지: 알고리즘 변경 시 조율된 배포 필요
- 번들 크기: 프론트엔드에 SHA256 라이브러리 추가
- 복잡성: 정규식 기반 정규화 드리프트 발생 가능
- 테스트 부담: 언어 간 패리티 테스트 필요

**기각**: 커플링 비용이 정확도 이점 초과.

### Option C: 항상 최악 케이스 예상치 표시

최대 가능 비용 표시 (모든 테스트 = 캐시 미스).

**장점:**

- 예측 리스크 제로
- 가장 단순한 구현

**단점:**

- 불량 UX: 캐시 옵션이 무가치해 보임
- 전환율 영향: 유익한 캐시 모드 건너뜀
- 신뢰 저하: 캐시 시나리오에서 예상치 항상 오류

**기각**: 캐시 투명성 목적 훼손.

### Option D: Worker 서비스 호출 예측

쿼터 소모 없이 실제 캐시 조회 실행하는 API 엔드포인트 추가.

**장점:**

- 100% 정확도
- 단일 진실 원천

**단점:**

- 추가 레이턴시: 네트워크가 모달 렌더링 지연
- API 복잡성: 유지할 새 엔드포인트
- 실패 모드: Worker 불가 시 모달 작동 불가
- Rate Limiting: 캐시 상태 탐색 남용 가능

**기각**: 예측 기능에 과잉 엔지니어링.

## Implementation Details

### API 엔드포인트

```
GET /api/spec-view/{analysisId}/cache-prediction?language={language}
```

### 데이터 구조

```typescript
interface CachePrediction {
  totalBehaviors: number;
  cacheableCount: number;
  estimatedCost: number;
}
```

### SQL 쿼리 구조

```sql
WITH
  previous_spec AS (
    -- 동일 코드베이스/언어의 최근 스펙 탐색
    SELECT id FROM spec_documents
    WHERE ... ORDER BY created_at DESC LIMIT 1
  ),
  previous_behaviors AS (
    -- 테스트명과 파일 경로 포함 행동 추출
    SELECT test_name, file_path FROM spec_behaviors
    WHERE spec_document_id = (SELECT id FROM previous_spec)
  ),
  current_tests AS (
    -- 현재 분석의 모든 테스트 케이스 조회
    SELECT name, file_path FROM test_cases
    WHERE analysis_id = $1
  ),
  matched_behaviors AS (
    -- 이름 AND 파일경로 일치 카운트
    SELECT COUNT(*) FROM current_tests c
    JOIN previous_behaviors p
      ON c.name = p.test_name AND c.file_path = p.file_path
  )
SELECT total_behaviors, cacheable_count;
```

### 프론트엔드 통합

```typescript
const dynamicEstimatedCost = forceRegenerate
  ? estimatedCost
  : (cachePrediction?.estimatedCost ?? estimatedCost);
```

## Consequences

### Positive

- **사용자 경험**: 예상치에 캐시 이점 반영, 정보 기반 의사결정 가능
- **비용 투명성**: 쿼터 커밋 전 대략적 절감액 확인
- **구현 속도**: 교차 리포지토리 조율 없이 단순 접근법 배포
- **복원력**: 외부 의존성 없음; Worker 불가 시에도 동작
- **유지보수성**: Worker의 캐시 키 진화와 Web 레이어 격리

### Negative

- **예측 정확도**: ~1% 케이스에서 부정확한 예상치 가능
- **알고리즘 분기**: Worker의 실제 알고리즘과 의도적 불일치
- **제한된 범위**: 교차 리포지토리 캐시 히트 예측 불가
- **사용자 기대**: 예상 != 실제 관리 위한 고지 필요

### 엣지 케이스

| 시나리오                 | 예측      | 실제           | 영향             |
| ------------------------ | --------- | -------------- | ---------------- |
| 테스트명 변경, 동일 경로 | 캐시 미스 | 캐시 미스      | 정확             |
| 테스트 이동, 동일 이름   | 캐시 미스 | 캐시 히트 가능 | 절감액 과소 예측 |
| 테스트 새 파일에 복사    | 캐시 미스 | 캐시 히트      | 절감액 과소 예측 |
| 공백만 이름 변경         | 캐시 히트 | 캐시 히트      | 정확             |

엣지 케이스는 절감액 과소 예측 (더 높은 비용 표시) 결과, 빌링 컨텍스트에서 안전한 방향.

## References

### Internal

- [ADR-13: 빌링 및 쿼터 아키텍처](/ko/adr/13-billing-quota-architecture.md) - 이벤트 기반 사용량 추적
- [ADR-21: 할당량 예약](/ko/adr/21-quota-reservation.md) - 원자적 쿼터 예약
- [Worker ADR-12: Phase 2 행동 캐시](/ko/adr/worker/12-phase2-behavior-cache.md) - 실제 캐시 알고리즘
- [Web ADR-26: Zustand 클라이언트 상태 관리](/ko/adr/web/26-zustand-client-state.md) - 상태 관리 패턴

### External

- [Commit 3645d52](https://github.com/specvital/web/commit/3645d52) - 초기 구현
- [Commit e5a7a53](https://github.com/specvital/web/commit/e5a7a53) - 캐시 재사용 선택 UI
