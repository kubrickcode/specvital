---
title: Scheduler 제거 및 Railway Cron 전환
description: Scheduler 서비스 제거 및 Railway Cron + 개별 바이너리 전환 ADR
---

# ADR-22: Scheduler 제거 및 Railway Cron 전환

> 🇺🇸 [English Version](/en/adr/22-scheduler-removal-railway-cron.md)

| Date       | Author       | Repos         |
| ---------- | ------------ | ------------- |
| 2026-02-02 | @KubrickCode | worker, infra |

## 상태

**Accepted** - [ADR-01: Scheduled Re-collection](/ko/adr/worker/01-scheduled-recollection.md) 대체

## Context

### 기존 아키텍처 문제

Scheduler 서비스([ADR-01](/ko/adr/worker/01-scheduled-recollection.md))는 사용자에게 즉각적인 응답 제공을 위한 사전 분석 목적으로 설계. 그러나 운영 데이터 분석 결과 근본적인 비용-효과 문제 발견:

| 지표           | 예상      | 실제                            |
| -------------- | --------- | ------------------------------- |
| 분석 시간      | 30초 이상 | ~5초                            |
| 사전 계산 가치 | 높음      | 낮음 (5초 대기 허용 가능)       |
| 데이터 최신성  | 유지 가능 | 활발한 레포지토리는 유지 불가능 |
| 스토리지 증가  | 제어 가능 | 급격한 증가 (미조회 결과 누적)  |
| 24/7 운영 비용 | 정당화됨  | 실제 유틸리티 대비 과도함       |

### 사전 계산 실패 원인

1. **낮은 가치**: 5초 분석 시간은 사용자에게 허용 가능한 수준
2. **최신성 유지 불가**: 활발한 레포지토리는 커밋이 빈번하여 사전 계산 결과 즉시 구식화
3. **데이터베이스 비대화**: 미조회 분석 결과 급격히 누적
4. **비용 비효율**: 24/7 Scheduler 운영 비용이 실제 유틸리티 초과

### Scheduler 아키텍처 오버헤드

Scheduler 서비스가 도입한 복잡성:

- **분산 락**: 단일 인스턴스 보장을 위한 PostgreSQL 기반 락
- **go-cron 내부 스케줄링**: 프로세스 내 cron 작업 관리
- **24/7 운영 비용**: 최소한의 실제 작업에 항상 켜진 서비스
- **장애 결합**: Scheduler 장애 시 모든 스케줄 작업 영향

## Decision

**Scheduler 서비스 완전 제거. Railway Cron 트리거 + 개별 단일 목적 바이너리로 전환.**

### 신규 아키텍처

```
┌─────────────────────────────────────────────────────────────┐
│                    Before (Scheduler Service)                │
├─────────────────────────────────────────────────────────────┤
│  cmd/scheduler/    → 24/7 실행, 내부 go-cron                │
│                      ├── Auto-refresh cron job              │
│                      ├── 분산 락 (PostgreSQL)               │
│                      └── Cleanup 작업 (내장)                │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                    After (Railway Cron)                      │
├─────────────────────────────────────────────────────────────┤
│  cmd/analyzer/         → 큐 소비자 (River, ON_FAILURE)      │
│  cmd/spec-generator/   → 큐 소비자 (River, ON_FAILURE)      │
│  cmd/retention-cleanup/→ Cron 바이너리 (Railway, "0 3 * *") │
│  cmd/enqueue/          → 수동 유틸리티                      │
└─────────────────────────────────────────────────────────────┘
```

### Railway Cron 설정

각 cron 작업은 별도의 Railway 서비스로 개별 설정:

**retention-cleanup/railway.json**:

```json
{
  "$schema": "https://railway.com/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "infra/retention-cleanup/Dockerfile"
  },
  "deploy": {
    "cronSchedule": "0 3 * * *",
    "restartPolicyType": "NEVER"
  }
}
```

**주요 설정**:

- `cronSchedule`: 스케줄링용 표준 cron 표현식
- `restartPolicyType: NEVER`: 바이너리 완료 후 종료
- 서비스별 별도 Dockerfile로 빌드 격리

## Options Considered

### Option A: Railway Cron + 개별 바이너리 (선택)

**동작 방식:**

1. 각 주기 작업을 독립 Go 바이너리로 구현
2. Railway가 cron 표현식으로 바이너리 실행 트리거
3. 바이너리 완료 후 종료 (24/7 프로세스 없음)
4. 분산 락 불필요 (Railway가 단일 실행 관리)

**장점:**

- cron 스케줄링을 위한 24/7 운영 비용 제거
- 분산 락 복잡성 제거 (플랫폼이 관리)
- Railway 대시보드에서 작업별 비용 가시성 확보
- Railway가 스케줄링 안정성 및 재시도 처리
- 간단한 배포 (완료 후 종료하는 바이너리)
- 작업별 독립적 스케일링 및 설정

**단점:**

- 실행 시 콜드 스타트 지연
- Railway 플랫폼 의존성
- 각 cron 작업마다 Railway IaC 필요

### Option B: Scheduler 유지 (범위 축소)

**설명:**

Scheduler 유지하되 auto-refresh 제거; cleanup 작업만 유지.

**장점:**

- 최소한의 코드 변경
- 기존 모니터링 및 알림 유지
- 익숙한 운영 모델

**단점:**

- 드문 작업에도 24/7 프로세스 필요
- 분산 락 복잡성 유지
- 실제 작업 대비 비용 비례하지 않음

### Option C: 외부 Cron 서비스 (GitHub Actions, CloudWatch)

**설명:**

API 엔드포인트 호출 또는 큐 작업 등록용 외부 cron 트리거 사용.

**장점:**

- 무료 티어 이용 가능 (GitHub Actions)
- 플랫폼 독립적 접근
- Railway 전용 설정 불필요

**단점:**

- 추가 보안 표면 (노출된 엔드포인트)
- 서비스 간 조율 복잡성
- 레이트 리미팅 및 재시도 로직 필요
- 모니터링 분산

## Implementation

### 제거된 컴포넌트

worker 레포지토리에서 삭제된 파일:

```
src/cmd/scheduler/main.go
src/internal/app/bootstrap/scheduler.go
src/internal/app/container_scheduler.go
src/internal/domain/analysis/autorefresh.go
src/internal/domain/analysis/decay.go
src/internal/handler/scheduler/autorefresh.go
src/internal/infra/scheduler/cron.go
src/internal/infra/scheduler/lock.go
src/internal/usecase/autorefresh/refresh.go
```

### 신규 바이너리 구조

```
cmd/
├── analyzer/           # 큐 소비자 (River)
├── spec-generator/     # 큐 소비자 (River)
├── retention-cleanup/  # Cron 바이너리 (Railway)
└── enqueue/            # 수동 유틸리티
```

### 인프라 설정

```
infra/
├── analyzer/
│   ├── Dockerfile
│   └── railway.json
├── spec-generator/
│   ├── Dockerfile
│   └── railway.json
└── retention-cleanup/
    ├── Dockerfile
    └── railway.json
```

### 배포 비교

| 항목      | Before (Scheduler)   | After (Railway Cron)        |
| --------- | -------------------- | --------------------------- |
| 운영 비용 | 24/7 (유휴 상태에도) | 실행 시간만                 |
| 분산 락   | 필요 (PostgreSQL)    | 불필요 (Railway 관리)       |
| 스케일링  | 고정 단일 인스턴스   | 작업별 독립                 |
| 장애 격리 | 모든 작업 동시 실패  | 작업별 격리                 |
| 설정      | 환경 변수 + 코드     | 서비스별 Railway IaC        |
| 모니터링  | 단일 서비스 메트릭   | Railway에서 서비스별 메트릭 |

## Consequences

### Positive

**비용 최적화:**

- 드문 cron 작업에 대한 24/7 운영 비용 제거
- 실제 실행 시간에 대해서만 비용 발생
- 작업별 비용 가시성으로 최적화 가능

**운영 단순화:**

- 분산 락 관리 및 디버깅 불필요
- 각 cron 작업은 단순한 완료 후 종료 바이너리
- Railway가 스케줄링, 재시도, 단일 실행 처리

**배포 독립성:**

- 각 cron 작업을 독립적으로 배포 가능
- 다른 스케줄도 코드 변경 불필요
- IaC 기반 설정 (Infrastructure as Code)

**장애 격리:**

- 하나의 cron 작업 실패가 다른 작업에 영향 없음
- 명확한 작업별 로그 및 메트릭
- 독립적 재시도 정책

### Negative

**플랫폼 의존성:**

- Railway의 cron 구현에 종속
- 마이그레이션 시 모든 cron 작업 재설정 필요
- Railway 전용 IaC 형식

**콜드 스타트 지연:**

- 각 실행마다 새 컨테이너 시작
- 1분 미만 간격에는 부적합
- 초기 연결 설정 오버헤드

**설정 분산:**

- 유지 관리할 railway.json 파일 다수
- Dockerfile과 railway.json 간 동기화 필요
- infra 레포지토리 파일 증가

### 대체된 ADR

| ADR                                                               | 상태      | 비고                                                       |
| ----------------------------------------------------------------- | --------- | ---------------------------------------------------------- |
| [Worker ADR-01](/ko/adr/worker/01-scheduled-recollection.md)      | 대체됨    | Auto-refresh scheduler 제거                                |
| [Worker ADR-05](/ko/adr/worker/05-worker-scheduler-separation.md) | 부분 대체 | Scheduler 더 이상 존재하지 않음; 바이너리 분리 패턴은 유효 |

### 관련 업데이트 필요

| 문서                                                      | 필요한 업데이트                  |
| --------------------------------------------------------- | -------------------------------- |
| [ADR-04](/ko/adr/04-queue-based-async-processing.md)      | Railway Cron 대안 관련 노트 추가 |
| [ADR-12](/ko/adr/12-worker-centric-analysis-lifecycle.md) | Scheduler 참조 제거              |

## References

- Commit `c163239`: Scheduler 서비스 제거
- Commit `f3fae45`: worker 바이너리 분리
- Commit `6e03a7f`: retention-cleanup bootstrap 추가
- [Railway Cron Documentation](https://docs.railway.app/reference/cron-jobs)
