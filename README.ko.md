# SpecVital

> 정적 코드 분석 기반 자동 테스트 인벤토리 생성

**[specvital.com](https://specvital.com)** | [English](README.md)

## SpecVital이란?

SpecVital은 AST 기반 정적 분석을 사용하여 GitHub 저장소에서 테스트 인벤토리를 자동 생성하는 플랫폼입니다. 저장소 URL만 입력하면 CI/CD 통합 없이 즉시 전체 테스트 현황을 확인할 수 있습니다.

### 핵심 가치

| 특징          | 설명                                        |
| ------------- | ------------------------------------------- |
| **즉시 분석** | 몇 시간이 아닌 몇 초 만에 결과 확인         |
| **정확성**    | Tree-sitter AST 파싱, 결정적 분석 (AI 아님) |
| **다중 언어** | 주요 테스트 프레임워크 지원                 |

## 작동 방식

```
GitHub URL → AST 파싱 (Tree-sitter) → 테스트 인벤토리
```

1. 공개 GitHub 저장소 URL 입력
2. SpecVital이 [specvital/core](https://github.com/specvital/core)를 사용해 테스트 파일 파싱
3. 대시보드에서 테스트 스위트, 케이스, 구조 확인

**참고**: 정적 분석 방식으로 테스트 정의를 파싱하며, 실행 결과는 제공하지 않습니다. 어떤 테스트가 존재하는지와 구조를 확인할 수 있으며, 통과/실패 여부는 알 수 없습니다.

## 지원 프레임워크

20개 이상의 테스트 프레임워크 지원:

| 언어                  | 프레임워크                               |
| --------------------- | ---------------------------------------- |
| JavaScript/TypeScript | Jest, Vitest, Playwright, Cypress, Mocha |
| Go                    | testing                                  |
| Python                | pytest, unittest                         |
| Java                  | JUnit 5, TestNG                          |
| Kotlin                | Kotest                                   |
| C#                    | NUnit, xUnit, MSTest                     |
| Ruby                  | RSpec, Minitest                          |
| PHP                   | PHPUnit                                  |
| Rust                  | cargo test                               |
| C++                   | Google Test                              |
| Swift                 | XCTest                                   |

자세한 내용은 [specvital/core](https://github.com/specvital/core)를 참조하세요.

## 현재 상태

SpecVital은 활발히 개발 중입니다. 현재 제공 기능:

- ✅ GitHub OAuth 인증
- ✅ URL 입력으로 저장소 분석
- ✅ 통계와 함께 테스트 트리 시각화
- ✅ 다중 프레임워크 감지
- ✅ 필터 및 검색 기능
- ✅ 다국어 지원 (영어, 한국어)

## 아키텍처

이 저장소는 웹 애플리케이션(프론트엔드 + 백엔드 API)을 포함합니다. 전체 시스템 구성:

| 저장소                                                  | 역할                               |
| ------------------------------------------------------- | ---------------------------------- |
| [specvital/web](https://github.com/specvital/web)       | 웹 대시보드 + REST API             |
| [specvital/core](https://github.com/specvital/core)     | 파서 라이브러리 (Tree-sitter 기반) |
| [specvital/worker](https://github.com/specvital/worker) | 분석 작업 백그라운드 워커          |
| [specvital/infra](https://github.com/specvital/infra)   | 데이터베이스 스키마 및 인프라      |

## 사용 사례

- **엔지니어링 매니저**: 팀의 테스트 커버리지 현황 파악
- **QA 리드**: 프로젝트 전반의 자동화 테스트 인벤토리 관리
- **스태프 엔지니어**: 레거시 리팩토링 전 기존 테스트 맵핑

## 피드백

- 질문 및 아이디어: [GitHub Discussions](https://github.com/orgs/specvital/discussions)
- 버그 리포트: [GitHub Issues](https://github.com/specvital/web/issues)

## 라이선스

[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

이 프로젝트는 Apache License 2.0 하에 라이선스됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

### 상표 고지

"SpecVital" 및 SpecVital 로고는 KubrickCode의 상표입니다. Fork 및 파생 저작물은 다른 이름과 브랜딩을 사용해야 합니다. 자세한 내용은 [NOTICE](NOTICE) 파일을 참조하세요.
