---
title: ADR (한국어)
---

# 아키텍처 의사결정 기록 (ADR)

> 🇺🇸 [English Version](/en/adr/)

Specvital 프로젝트의 아키텍처 의사결정 문서화

## ADR이란?

ADR(Architecture Decision Record)은 중요한 아키텍처 결정을 그 배경 및 결과와 함께 기록하는 문서입니다. 멀티-리포지토리 마이크로서비스 환경에서 의사결정 히스토리를 유지하는 데 도움이 됩니다.

## ADR 작성 시점

| 카테고리     | 예시                                              |
| ------------ | ------------------------------------------------- |
| 기술 스택    | 프레임워크 선택, 라이브러리 도입, 버전 업그레이드 |
| 아키텍처     | 서비스 경계, 통신 패턴, 데이터 흐름               |
| API 설계     | 엔드포인트 구조, 버저닝 전략, 에러 처리           |
| 데이터베이스 | 스키마 설계, 마이그레이션 전략, 인덱싱 방식       |
| 인프라       | 배포 플랫폼, 스케일링 전략, 모니터링              |
| 공통 관심사  | 보안, 성능 최적화, 옵저버빌리티                   |

## 템플릿

| 템플릿                       | 용도                                  |
| ---------------------------- | ------------------------------------- |
| [template.md](./template.md) | 대부분의 의사결정에 사용하는 표준 ADR |

## 명명 규칙

```
XX-brief-decision-title.md
```

- `XX`: 두 자리 순차 번호 (01, 02, ...)
- 소문자와 하이픈 사용
- 간결하고 명확한 제목

## 기술 영역

| 영역           | 영향받는 리포지토리 |
| -------------- | ------------------- |
| Parser         | core                |
| API            | web                 |
| Worker         | collector           |
| Database       | infra               |
| Infrastructure | infra               |
| Cross-cutting  | 복수                |

## ADR 목록

### 공통 (전체 리포지토리)

| #   | 제목                                                              | 영역           | 날짜       |
| --- | ----------------------------------------------------------------- | -------------- | ---------- |
| 01  | [정적 분석 기반 즉시 분석](./01-static-analysis-approach.md)      | Cross-cutting  | 2024-12-17 |
| 02  | [경쟁 차별화 전략](./02-competitive-differentiation.md)           | Cross-cutting  | 2024-12-17 |
| 03  | [API와 Worker 서비스 분리](./04-api-worker-service-separation.md) | Architecture   | 2024-12-17 |
| 04  | [큐 기반 비동기 처리](./05-queue-based-async-processing.md)       | Architecture   | 2024-12-17 |
| 05  | [Polyrepo 리포지토리 전략](./06-repository-strategy.md)           | Architecture   | 2024-12-17 |
| 06  | [PaaS 우선 인프라 전략](./07-paas-first-infrastructure.md)        | Infrastructure | 2024-12-17 |
| 07  | [공유 인프라 전략](./08-shared-infrastructure.md)                 | Infrastructure | 2024-12-17 |

### Core 리포지토리

| #   | 제목                                                         | 영역 | 날짜       |
| --- | ------------------------------------------------------------ | ---- | ---------- |
| 01  | [코어 라이브러리 분리](./core/01-core-library-separation.md) | Core | 2024-12-17 |

### Collector 리포지토리

| #   | 제목                                                                    | 영역         | 날짜       |
| --- | ----------------------------------------------------------------------- | ------------ | ---------- |
| 01  | [스케줄 기반 재수집 아키텍처](./collector/01-scheduled-recollection.md) | Architecture | 2024-12-18 |

## 프로세스

1. **생성**: [template.md](./template.md) 복사 → `XX-title.md`
2. **작성**: 확정된 의사결정 내용으로 모든 섹션 작성
3. **현지화**: `adr/`에 영어 버전 생성
4. **리뷰**: 팀 리뷰를 위해 PR 제출
5. **병합**: 승인 후 목록에 추가

## 관련 리포지토리

- [specvital/core](https://github.com/specvital/core) - 파서 엔진
- [specvital/web](https://github.com/specvital/web) - 웹 플랫폼
- [specvital/collector](https://github.com/specvital/collector) - 워커 서비스
- [specvital/infra](https://github.com/specvital/infra) - 인프라
