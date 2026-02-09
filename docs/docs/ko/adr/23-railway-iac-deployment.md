---
title: Railway IaC 배포
description: Railway CLI와 GitHub Actions를 활용한 Infrastructure as Code 배포에 관한 ADR
---

# ADR-23: Railway IaC 배포

> [English Version](/en/adr/23-railway-iac-deployment.md)

| 날짜       | 작성자       | 레포지토리         |
| ---------- | ------------ | ------------------ |
| 2026-02-02 | @KubrickCode | worker, web, infra |

## 상태

**승인됨** - [ADR-06: PaaS 우선 인프라 전략](/ko/adr/06-paas-first-infrastructure.md) 및 [ADR-22: Scheduler 제거 및 Railway Cron 전환](/ko/adr/22-scheduler-removal-railway-cron.md) 보완

## 배경

### 문제 상황

Railway 플랫폼 통합 초기에는 GitHub App 자동 배포 기능 사용. 다음과 같은 운영상 과제 발생:

1. **수동 대시보드 설정**: GitHub App 통합 후 각 서비스별 Railway 웹 대시보드에서 수동 설정 필요
2. **비재현성 배포**: 서비스 재생성 시 모든 수동 설정 반복 필요, 히스토리 추적 불가
3. **BuildArgs 제한**: Railway의 `railway.json` 스키마가 Docker 빌드용 `buildArgs` 공식 미지원
4. **환경 변수 Rate Limiting**: 다수의 순차적 `railway variables --set` CLI 호출 시 Railway API Rate Limit 발생

### 변경 필요 사유

| 문제               | 영향                       | 빈도          |
| ------------------ | -------------------------- | ------------- |
| 수동 대시보드 설정 | 배포 드리프트, 온보딩 마찰 | 신규 서비스   |
| 비재현성           | 재해 복구 위험             | 서비스 재생성 |
| BuildArgs 미지원   | Dockerfile 공유 불가       | 빌드 설정     |
| 변수 Rate Limiting | 배포 실패, 재시도 오버헤드 | 매 배포       |

## 결정

**GitHub Actions와 Railway CLI를 활용한 Infrastructure as Code(IaC) 배포 채택. 서비스별 Dockerfile 및 railway.json 설정 파일 사용.**

### 설정 구조

```
infra/
├── analyzer/
│   ├── Dockerfile          # 서비스별 Dockerfile
│   └── railway.json        # $schema 포함 서비스 설정
├── spec-generator/
│   ├── Dockerfile
│   └── railway.json
├── retention-cleanup/
│   ├── Dockerfile
│   └── railway.json
└── web-backend/
    ├── Dockerfile
    └── railway.json
```

### 핵심 패턴

**1. 서비스별 Dockerfile** (buildArgs 제한 회피):

공유 Dockerfile 대신 각 서비스별 독립 Dockerfile 사용:

```dockerfile
# infra/analyzer/Dockerfile
FROM golang:1.22-alpine AS builder
# ... 서비스별 빌드 단계
```

**2. railway.json 설정**:

```json
{
  "$schema": "https://railway.com/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "infra/analyzer/Dockerfile"
  },
  "deploy": {
    "restartPolicyType": "ON_FAILURE",
    "restartPolicyMaxRetries": 3
  }
}
```

**3. Cron 서비스 설정** (예약 작업용):

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

**4. GitHub Actions Matrix 배포**:

```yaml
deploy:
  strategy:
    matrix:
      service: [analyzer, spec-generator, retention-cleanup]
    fail-fast: true
  steps:
    - name: Install Railway CLI
      run: npm install -g @railway/cli
    - name: Deploy Service
      run: |
        cp infra/${{ matrix.service }}/railway.json railway.json
        railway up --service ${{ matrix.service }} --detach
        rm railway.json
      env:
        RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

## 고려한 선택지

### 선택지 A: GitHub Actions + Railway CLI (선택됨)

GitHub Actions 워크플로우에서 Railway CLI로 서비스별 설정 파일 기반 배포.

**장점:**

- Git 기반 버전 관리 설정 (코드 리뷰, 히스토리)
- 모든 환경에서 재현 가능한 배포
- 서비스별 Dockerfile로 buildArgs 제한 회피
- Matrix 전략으로 병렬 서비스 배포
- `$schema`로 IDE 자동완성 및 유효성 검사
- 통합 API 호출로 Rate Limiting 감소

**단점:**

- 서비스 Dockerfile 간 설정 중복
- GitHub Actions 가용성에 배포 의존
- 다수의 railway.json 파일 관리 필요

### 선택지 B: Railway GitHub App 통합 (기존 방식)

Railway 네이티브 통합을 통한 GitHub 푸시 이벤트 기반 자동 배포.

**장점:**

- CI/CD 설정 불필요
- 푸시 시 자동 배포
- 단순한 멘탈 모델

**단점:**

- 서비스별 대시보드 수동 설정 필요
- 설정 버전 관리 불가
- 스테이징/프로덕션 간 환경 드리프트
- 동일한 buildArgs 제한 적용

### 선택지 C: Pulumi/Terraform + Railway Provider

IaC 도구와 Railway 프로바이더 활용.

**장점:**

- 완전한 선언적 인프라 관리
- 상태 추적 및 드리프트 감지
- 다중 프로바이더 리소스 관리

**단점:**

- 추가 도구 복잡성
- 상태 파일 관리 오버헤드
- 단일 플랫폼 배포에 과도한 엔지니어링
- Railway 프로바이더 기능 제한 가능성

## 구현

### 배포 워크플로우

1. **빌드 단계**: Railway CLI가 리포지토리 루트의 `railway.json` 읽기
2. **설정 복사**: 워크플로우에서 서비스별 설정을 루트로 복사
3. **배포**: 서비스 이름으로 `railway up` 실행
4. **정리**: 복사된 설정 파일 삭제

### 환경 변수

- **프로젝트 레벨**: 서비스 간 공유 (DATABASE_URL 등)
- **서비스 레벨**: Rate Limiting 방지를 위한 통합 워크플로우 단계에서 설정

### Railway Project Token

CI/CD용 Project Token 사용 (User Token 아님):

```yaml
env:
  RAILWAY_TOKEN: ${{ secrets.RAILWAY_TOKEN }}
```

Project Token 사용 시 `railway link` 명령 불필요.

## 결과

### 긍정적 영향

**재현성:**

- 모든 배포 설정 버전 관리
- 리포지토리만으로 모든 서비스 재생성 가능
- Pull Request에서 배포 변경 사항 확인 가능

**개발자 경험:**

- 스키마 유효성 검사로 IDE 자동완성 제공
- 서비스별 명확한 설정 경계
- 익숙한 Git 기반 워크플로우

**운영 안정성:**

- 통합 API 호출로 Rate Limiting 방지
- Matrix 전략으로 병렬 배포
- GitHub Actions에서 명확한 배포 로그

### 부정적 영향

**설정 중복:**

- 서비스별 Dockerfile에 유사한 베이스 레이어 포함
- 사소한 변경 시 다수 파일 수정 필요 가능
- **완화**: 공통 베이스 이미지와 멀티스테이지 빌드 활용

**GitHub Actions 의존:**

- 배포에 GitHub Actions 러너 필요
- GitHub 장애 시 자동 배포 차단
- **완화**: Railway CLI로 수동 배포 가능 (폴백)

**학습 곡선:**

- 팀이 railway.json 스키마 이해 필요
- GitHub Actions 워크플로우 문법 필요
- **완화**: 명확한 문서화와 중앙집중식 워크플로우

## 관련 ADR

| ADR                                                                    | 관계                             |
| ---------------------------------------------------------------------- | -------------------------------- |
| [ADR-06: PaaS 우선 인프라](/ko/adr/06-paas-first-infrastructure.md)    | 기반 - Railway를 PaaS 플랫폼으로 |
| [ADR-22: Scheduler 제거](/ko/adr/22-scheduler-removal-railway-cron.md) | 확장 - Railway Cron 설정         |

## 참조

- 커밋 `b842318`: Railway IaC 및 GitHub Actions 배포 전환 (worker)
- 커밋 `c9ee36f`: Railway 배포 IaC 전환 (web)
- [Railway IaC 문서](https://docs.railway.app/guides/config-as-code)
- [Railway CLI 레퍼런스](https://docs.railway.app/reference/cli-api)
