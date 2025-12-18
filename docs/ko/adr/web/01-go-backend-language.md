---
title: Go 백엔드 언어
---

# ADR-01: 백엔드 언어로 Go 선택

> :us: [English Version](/en/adr/web/01-go-backend-language.md)

| 날짜       | 작성자       | 리포지토리 |
| ---------- | ------------ | ---------- |
| 2024-12-18 | @KubrickCode | web        |

## Context

### 언어 선택 문제

웹 플랫폼은 백엔드 언어 선택이 필요했다. 두 가지 주요 후보가 있었다:

1. **Go**: 기존 collector와 core 서비스와 일치
2. **NestJS (TypeScript)**: Next.js 프론트엔드와 일치

### 기존 아키텍처

시스템은 이미 다음에서 Go를 사용 중이었다:

- **Core 라이브러리**: 파서 엔진, crypto 유틸리티, 도메인 모델
- **Collector 서비스**: asynq를 통한 분석 작업 처리 백그라운드 워커
- **공유 인프라**: Redis 기반 태스크 큐 (asynq)

### 핵심 고려사항

웹 백엔드는 다음이 필요했다:

- collector 처리를 위한 분석 태스크 enqueue
- collector와의 암호화 작업 공유 (OAuth 토큰 암호화)
- 필요시 core 라이브러리 기능 접근

## Decision

**기존 서비스와의 기술 스택 통합을 극대화하기 위해 Go를 웹 백엔드 언어로 채택한다.**

핵심 원칙:

1. **단일 큐 시스템**: collector와 asynq 공유 (별도 BullMQ 불필요)
2. **직접 라이브러리 접근**: RPC 오버헤드 없이 core 라이브러리 import
3. **통합 도구**: 단일 언어 CI/CD, 모니터링, 배포
4. **암호화 공유**: OAuth 토큰에 대해 동일한 암호화/복호화

## Options Considered

### Option A: Go 백엔드 (선택됨)

**작동 방식:**

- Go HTTP 서버 (Chi/Gin/Echo)가 REST API 제공
- `github.com/specvital/core` 패키지 직접 import
- 태스크 enqueue를 위한 공유 asynq 클라이언트
- collector와 동일한 Redis 인스턴스

**장점:**

- **통합 오버헤드 없음**: core 라이브러리 직접 import
- **공유 큐 인프라**: 단일 Redis 인스턴스, 통합 asynq 프로토콜
- **일관된 암호화**: 서비스 간 동일한 암호화
- **운영 단순성**: 관리할 언어 런타임이 하나
- **리소스 효율성**: Node.js보다 낮은 메모리 사용량

**단점:**

- 프론트엔드 개발자가 Go 기초를 배워야 함
- npm에 비해 작은 생태계
- TypeScript 프론트엔드와 타입 공유 불가

### Option B: NestJS + TypeScript

**작동 방식:**

- TypeScript 기반 NestJS 프레임워크
- 태스크 큐로 BullMQ (asynq와 별도)
- gRPC 래퍼 또는 TypeScript 재작성을 통한 core 라이브러리 접근

**장점:**

- Next.js 프론트엔드와 언어 공유
- 더 큰 패키지 생태계 (npm)
- 프론트엔드와 백엔드 간 타입 공유

**단점:**

- **Core 라이브러리 비호환**: Go core를 직접 사용 불가
- **이중 큐 시스템**: BullMQ (web) + asynq (collector) = 2배 Redis 또는 복잡한 브릿징
- **암호화 재구현**: NaCl 암호화를 TypeScript로 재작성 필요
- **운영 복잡성**: 두 언어 런타임, 별도 CI/CD

### Option C: NestJS BFF + Go Core API

**작동 방식:**

- NestJS를 Backend-for-Frontend 레이어로
- core 라이브러리를 래핑하는 Go 서비스
- 레이어 간 gRPC 통신

**평가:**

- 3단계 레이턴시 오버헤드
- 복잡한 배포 토폴로지
- 현재 규모에 과도한 최적화
- **기각**: 요구사항 대비 과잉

## Implementation Considerations

### Core 라이브러리 통합

```
Web Service (Go)                    Collector (Go)
      │                                   │
      └─── import core/pkg/crypto ────────┘
      │                                   │
      └─── import core/pkg/domain ────────┘
```

공유되는 핵심 패키지:

- `core/pkg/crypto`: OAuth 토큰용 NaCl SecretBox 암호화
- `core/pkg/domain`: 타입 안전 도메인 모델

### 큐 아키텍처

```
Web (Producer)          Redis (Upstash)         Collector (Consumer)
      │                       │                        │
      ├─ asynq.Enqueue() ───→ task_queue ────→ asynq.HandleFunc()
      │                       │                        │
      └──────────── shared asynq protocol ─────────────┘
```

이점:

- 단일 Redis 인스턴스 (비용 절감)
- 타입 안전 태스크 페이로드
- 내장 재시도, 스케줄링, dead-letter queue

### 암호화 공유

OAuth 플로우 요구사항:

1. Web이 DB 저장 전 GitHub 토큰 암호화
2. Collector가 GitHub API 접근 시 토큰 복호화

Go 사용 시:

- 동일한 `crypto.Encryptor` 인터페이스
- 동일한 암호화 키 (환경 변수)
- 호환성 보장

NestJS 사용 시:

- TypeScript로 NaCl SecretBox 재구현 필요
- 미묘한 암호화 비호환 위험
- 추가 테스트 부담

## Consequences

### Positive

**인프라 효율성:**

- 단일 Redis 인스턴스가 web과 collector 모두 지원
- 통합 모니터링 및 알림
- 공유 배포 패턴

**개발 속도:**

- web과 core 간 직렬화 레이어 없음
- 서비스 간 컴파일 타임 타입 안전성
- 일관된 에러 처리 패턴

**운영 단순성:**

- 백엔드 서비스에 단일 언어
- 단일 CI/CD 파이프라인 패턴
- 통합 의존성 관리 (go.mod)

### Negative

**학습 곡선:**

- 프론트엔드 개발자가 Go에 익숙해져야 함
- **완화책**: 집중 교육, 페어 프로그래밍

**타입 공유:**

- 프론트엔드용 자동 타입 생성 불가
- **완화책**: TypeScript 클라이언트용 OpenAPI codegen

**생태계:**

- npm보다 적은 패키지
- **완화책**: 기능 시작 전 대안 평가

## References

- [ADR-05: 큐 기반 비동기 처리](/ko/adr/05-queue-based-async-processing.md)
- [ADR-08: 공유 인프라 전략](/ko/adr/08-shared-infrastructure.md)
- [Core ADR-01: 코어 라이브러리 분리](/ko/adr/core/01-core-library-separation.md)
