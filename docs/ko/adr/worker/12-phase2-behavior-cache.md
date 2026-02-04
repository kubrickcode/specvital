---
title: Phase 2 행동 캐시
description: Phase 2 AI 생성 테스트 행동 설명 캐싱을 통한 API 비용 절감
---

# ADR-12: Phase 2 행동 캐시

> 🇺🇸 [English Version](/en/adr/worker/12-phase2-behavior-cache.md)

| 날짜       | 작성자       | 리포지토리    |
| ---------- | ------------ | ------------- |
| 2026-01-24 | @KubrickCode | worker, infra |

## Context

### Phase 2 비용 문제

AI 기반 스펙 문서 생성 파이프라인([ADR-14](/ko/adr/14-ai-spec-generation-pipeline.md))의 2단계 구조:

| Phase   | 모델                  | 비용     | 목적                      |
| ------- | --------------------- | -------- | ------------------------- |
| Phase 1 | gemini-2.5-flash      | $0.30/1M | 도메인별 테스트 분류      |
| Phase 2 | gemini-2.5-flash-lite | $0.10/1M | 테스트명 → 행동 설명 변환 |

Phase 1 결과는 `content_hash` 기반 문서 수준 캐싱 적용. Phase 2는 테스트별 AI 호출로 다음 상황에서 비용 증가:

- 동일 테스트 파일의 다중 커밋 분석
- 다른 리포지토리의 유사 테스트명
- 파서 버전 업데이트 또는 사용자 요청에 의한 재분석

### 캐싱 기회

테스트 행동 설명의 높은 캐시 재사용성:

| 시나리오               | 예시                                             |
| ---------------------- | ------------------------------------------------ |
| 동일 테스트, 다른 커밋 | 커밋 A와 B의 `test_user_login`은 동일 행동 생성  |
| 리포지토리 간 유사성   | `testAuthentication` 행동은 언어/프레임워크 무관 |
| 재분석                 | 파서 업그레이드가 행동 의미론에 영향 없음        |

### 캐시 키 설계 과제

**문제**: 테스트 행동의 고유 식별자 결정

| 접근법               | 장점          | 단점                                   |
| -------------------- | ------------- | -------------------------------------- |
| 테스트명만           | 최대 재사용   | 컨텍스트 무시 (동명이테스트 구분 불가) |
| 테스트명 + 파일 경로 | 컨텍스트 인식 | 경로 변경 시 캐시 무효화               |
| 테스트 콘텐츠 해시   | 정확한 매칭   | 유사 테스트 간 재사용 불가             |
| 의미론적 핑거프린트  | 의도 포착     | 복잡, 추가 AI 호출 필요                |

## Decision

**복합 키 `(test_name_hash, language, model_id)` 및 TTL 기반 만료를 사용하는 PostgreSQL 기반 행동 캐시 구현.**

### 테이블 스키마

```sql
CREATE TABLE behavior_caches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_name_hash TEXT NOT NULL,         -- 정규화된 테스트명의 SHA-256
    language VARCHAR(10) NOT NULL,        -- en, ko 등
    model_id VARCHAR(100) NOT NULL,       -- gemini-2.5-flash-lite
    behavior_description TEXT NOT NULL,   -- 캐시된 AI 출력
    confidence DECIMAL(3,2) NOT NULL,     -- 0.00-1.00
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,      -- TTL 만료
    hit_count INTEGER DEFAULT 0,          -- 사용 추적

    CONSTRAINT behavior_caches_unique
        UNIQUE (test_name_hash, language, model_id)
);

-- 조회용 인덱스
CREATE INDEX idx_behavior_caches_lookup
    ON behavior_caches (test_name_hash, language, model_id)
    WHERE expires_at > NOW();

-- 정리용 인덱스
CREATE INDEX idx_behavior_caches_expiry
    ON behavior_caches (expires_at);
```

### 캐시 키 전략

```
test_name_hash = SHA256(normalize(test_name))

normalize(test_name):
  1. 소문자 변환
  2. 공통 접두사 제거 (test_, it_, describe_, should_)
  3. 특수문자 및 숫자 제거
  4. 공백 정리
```

**예시**:

| 원본 테스트명                      | 정규화             | 해시 (축약) |
| ---------------------------------- | ------------------ | ----------- |
| `test_user_can_login`              | `user can login`   | `a3f2...`   |
| `TestUserCanLogin`                 | `user can login`   | `a3f2...`   |
| `it('should allow user to login')` | `allow user login` | `b7c1...`   |
| `describe('User Login')`           | `user login`       | `c9e4...`   |

### 캐시 조회 흐름

```
┌─────────────────────────────────────────────────────────────────┐
│                    Phase 2 처리                                  │
├─────────────────────────────────────────────────────────────────┤
│  Feature 내 각 테스트에 대해:                                    │
│                                                                  │
│  1. test_name_hash 계산                                          │
│                                                                  │
│  2. 캐시 조회:                                                   │
│     SELECT behavior_description, confidence                      │
│     FROM behavior_caches                                         │
│     WHERE test_name_hash = ? AND language = ? AND model_id = ?   │
│       AND expires_at > NOW()                                     │
│                                                                  │
│  3. HIT 시:                                                      │
│     - hit_count 증가                                             │
│     - 캐시된 행동 반환 (AI 호출 생략)                            │
│     - 쿼터 미소비                                                │
│                                                                  │
│  4. MISS 시:                                                     │
│     - Gemini API 호출                                            │
│     - 결과 TTL과 함께 저장                                       │
│     - 쿼터 소비                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### TTL 설정

| 티어       | TTL   | 근거                        |
| ---------- | ----- | --------------------------- |
| Free       | 7일   | 제한된 저장 공간, 높은 이탈 |
| Pro        | 30일  | 표준 비즈니스 보존          |
| Pro Plus   | 90일  | 확장 캐싱 가치              |
| Enterprise | 180일 | 최대 비용 최적화            |

### 쿼터 시스템 연동

[ADR-13](/ko/adr/13-billing-quota-architecture.md)에 따라 캐시 히트는 쿼터 미소비:

```go
type Phase2Result struct {
    Behavior   string
    Confidence float64
    FromCache  bool  // true 시 쿼터 미소비
}

// 사용량 추적
if !result.FromCache {
    quotaService.RecordUsage(ctx, userID, QuotaTypeSpecView)
}
```

## Options Considered

### Option A: PostgreSQL 테이블 (선택됨)

TTL 및 히트 추적 기능의 데이터베이스 기반 캐시.

**장점**:

- 기존 PostgreSQL 인프라와 통합
- 분석용 쿼리 가능 (히트율, 인기 테스트)
- 스케줄 작업을 통한 자동 정리
- 스펙 문서 쓰기와 트랜잭션 일관성

**단점**:

- 고볼륨 조회 시 데이터베이스 부하
- 대용량 캐시의 저장 비용

### Option B: Redis 캐시

자동 만료 기능의 인메모리 캐시.

**장점**:

- 밀리초 미만 조회
- 네이티브 TTL 지원
- 데이터베이스 부하 감소

**단점**:

- 추가 인프라 (현재 스택에 없음)
- Redis 재시작 시 캐시 손실
- 캐시 크기에 따른 메모리 비용 증가

### Option C: 문서 수준 캐싱만

기존 `spec_documents.content_hash` 캐싱에 의존.

**장점**:

- 새 인프라 불필요
- 이미 구현됨

**단점**:

- 리포지토리 간 유사 테스트 재사용 불가
- 테스트 변경 시 전체 Phase 2 재실행
- 문서 간 최적화 기회 상실

### Option D: 추가 캐싱 없음

Phase 2 비용을 운영 비용으로 수용.

**장점**:

- 가장 단순한 구현
- 캐시 무효화 복잡성 없음

**단점**:

- 스케일 시 높은 API 비용
- 반복 테스트의 느린 응답 시간
- 고볼륨 사용자의 비용 효율 저하

## Consequences

### Positive

| 영역          | 효과                                    |
| ------------- | --------------------------------------- |
| 비용 절감     | Phase 2 API 비용 40-60% 절감 (추정)     |
| 응답 시간     | 캐시 히트로 테스트당 1-2초 AI 지연 회피 |
| 쿼터 공정성   | 캐시 히트는 사용자 쿼터 미소비          |
| 분석          | 히트율 메트릭으로 캐시 튜닝 가능        |
| 리포지토리 간 | 유사 테스트의 캐시된 행동 공유          |

### Negative

| 영역          | 트레이드오프                               |
| ------------- | ------------------------------------------ |
| 저장 공간     | 고유 테스트 다양성에 따른 캐시 테이블 증가 |
| 오래된 데이터 | 캐시된 행동이 모델 업데이트 미반영 가능    |
| 정규화 오류   | 과도한 정규화로 잘못된 매칭 가능           |
| 정리 오버헤드 | TTL 적용을 위한 스케줄 작업 필요           |

### 기술 참고

- **캐시 워밍**: 미구현; 캐시는 자연적으로 축적
- **무효화**: 모델 버전 변경은 키의 `model_id`를 통해 무효화
- **충돌 해결**: 첫 번째 쓰기 우선; 동시 쓰기는 드묾
- **모니터링**: 분석당 캐시 히트/미스 비율 로깅

## Configuration

```
BEHAVIOR_CACHE_ENABLED=true
BEHAVIOR_CACHE_DEFAULT_TTL=30d
BEHAVIOR_CACHE_CLEANUP_SCHEDULE=0 4 * * *  # UTC 04:00 매일
BEHAVIOR_CACHE_CLEANUP_BATCH_SIZE=5000
```

## References

- [ADR-14: AI 기반 스펙 문서 생성 파이프라인](/ko/adr/14-ai-spec-generation-pipeline.md) - 상위 아키텍처
- [ADR-13: 빌링 및 쿼터 아키텍처](/ko/adr/13-billing-quota-architecture.md) - 쿼터 연동
- [ADR-18: GitHub API 캐시 테이블](/ko/adr/18-github-api-cache-tables.md) - 유사 캐싱 패턴
- 커밋: `8917156` (behavior_caches 테이블) - 2026-01-24
