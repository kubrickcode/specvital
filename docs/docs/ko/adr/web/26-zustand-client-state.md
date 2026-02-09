---
title: Zustand 클라이언트 상태 관리
description: 웹 플랫폼의 클라이언트 사이드 상태 관리를 위한 Zustand 도입 ADR
---

# ADR-26: Zustand 클라이언트 상태 관리

> 🇺🇸 [English Version](/en/adr/web/26-zustand-client-state)

| 날짜       | 작성자       | 리포지토리 |
| ---------- | ------------ | ---------- |
| 2026-01-23 | @KubrickCode | web        |

## Context

Specvital 웹 플랫폼에 기존 상태 관리 솔루션과 다른 **임시 클라이언트 사이드 상태** 필요:

| 상태 카테고리       | 관리 도구      | 예시                                        |
| ------------------- | -------------- | ------------------------------------------- |
| **서버 상태**       | TanStack Query | 리포지토리 데이터, 분석 결과, 사용자 프로필 |
| **URL 상태**        | nuqs           | 대시보드 필터, 검색 쿼리, 뷰 모드           |
| **클라이언트 상태** | _필요_         | 백그라운드 작업 추적, 세션 임시 데이터      |

### 문제

모달 기반 UI의 백그라운드 작업(분석, 스펙 생성)에 치명적 한계 존재:

1. 모달 닫을 시 상태 손실
2. 페이지 네비게이션 시 폴링 중단
3. 실행 중인 작업의 페이지 간 가시성 부재

Issue #240 사용자 의견: _"모달을 닫거나 다른 페이지로 이동하면 스펙 생성 진행 상황 추적 불가."_

### 제약 조건

- **App Router RSC 경계**: 상태 관리는 클라이언트/서버 컴포넌트 경계 준수 필요
- **Hydration 불일치 위험**: SSR + 클라이언트 상태 조합 시 신중한 처리 필요
- **번들 크기 민감성**: 대시보드 성능 중요
- **기존 스택 통합**: TanStack Query, nuqs와 공존 필요

## Decision

클라이언트 사이드 전역 상태 관리를 위해 **Zustand** 도입.

### 근거

1. 상태 관리 3종 세트 완성: TanStack Query (서버), nuqs (URL), Zustand (클라이언트)
2. 최소 번들 영향 (~3 KB gzipped)
3. Provider 래퍼 불필요
4. 내장 미들웨어: `persist`, `devtools`
5. React 18+ 최적화를 위한 `useSyncExternalStore` 네이티브 사용
6. 완전한 타입 추론의 TypeScript 우선 설계

## Options Considered

### Option A: Zustand (선택됨)

**특성**:

- ~3 KB gzipped 번들 크기
- Provider 래퍼 불필요
- 선택적 구독이 가능한 Hooks 기반 API

**장점**:

- 보일러플레이트 제로, 즉각적 생산성
- 셋업 없이 컴포넌트 트리 어디서나 동작
- Redux DevTools 호환
- 스토리지 동기화를 위한 내장 persist 미들웨어

**단점**:

- 대규모 코드베이스에서 Redux보다 구조화 부족
- 팀 컨벤션 수립 필요 (강제 패턴 없음)

### Option B: Redux Toolkit

**특성**:

- ~10-15 KB gzipped
- Provider 래퍼 필수
- 서버 상태용 RTK Query

**평가**: 기각. RTK Query가 기존 TanStack Query와 중복. 남은 클라이언트 상태 요구사항에 비해 추가 보일러플레이트와 Provider 요구사항 정당화 불가.

### Option C: React Context + useReducer

**특성**:

- 0 KB (내장)
- Provider 래퍼 필수
- 수동 최적화 필요

**평가**: 기각. 자주 변경되는 상태에서 리렌더링 문제 발생. Context는 저빈도 업데이트용 설계, 작업 추적 같은 동적 상태에 부적합.

### Option D: Jotai (Atomic State)

**특성**:

- ~3.5 KB gzipped
- Atom 기반 구성
- Provider 선택적 (권장)

**평가**: 적합하나 미선택. Atom 기반 모델은 다른 사고 모델 필요. Zustand의 스토어 기반 접근 방식이 단순한 작업 추적 요구에 더 직관적.

## Implementation

### Store 구조

```typescript
type TaskStoreState = {
  tasks: Map<string, BackgroundTask>;
};

type TaskStoreActions = {
  addTask: (task: Omit<BackgroundTask, "createdAt">) => void;
  updateTask: (id: string, updates: Partial<BackgroundTask>) => void;
  removeTask: (id: string) => void;
  clearCompletedTasks: () => void;
};
```

### 주요 파일

| 파일                                             | 용도                   |
| ------------------------------------------------ | ---------------------- |
| `lib/background-tasks/task-store.ts`             | 메인 Zustand 스토어    |
| `lib/background-tasks/hooks.ts`                  | 소비자용 커스텀 훅     |
| `lib/background-tasks/task-store.spec.ts`        | 유닛 테스트            |
| `lib/background-tasks/components/task-badge.tsx` | 활성 작업 수 표시 뱃지 |

### 사용 패턴

| 패턴            | 구현                                               |
| --------------- | -------------------------------------------------- |
| 싱글톤 Store    | `create()`로 단일 `useTaskStore` 생성              |
| 파생 Selectors  | `useActiveTasks`로 필터링된 결과 계산              |
| Shallow 비교    | `useShallow`로 불필요한 리렌더링 방지              |
| React 외부 접근 | 비컴포넌트에서 `useTaskStore.getState()` 사용      |
| 수동 영속화     | persist 미들웨어 대신 커스텀 sessionStorage 동기화 |

### 영속화

Zustand의 `persist` 미들웨어 대신 커스텀 sessionStorage 구현:

- 스토리지 키: `specvital:background-tasks`
- `typeof window === "undefined"` 체크로 SSR 안전
- 손상된 스토리지나 쿼터 초과 시 조용한 오류 처리

## Consequences

### 긍정적

- **최소 통합 마찰**: Provider 불필요, 즉시 동작
- **개발자 경험**: 상태 + 액션이 동일 파일에 위치한 단일 파일 스토어
- **성능**: 선택적 구독으로 불필요한 리렌더링 방지
- **디버깅**: Redux DevTools 호환

### 부정적

- **팀 일관성 리스크**: 강제 패턴 없음
  - **완화책**: 스토어 구성 컨벤션 문서화, 코드 리뷰로 강제
- **제한된 미들웨어 생태계**: Redux보다 적은 커뮤니티 미들웨어
  - **완화책**: 핵심 요구사항(persist, devtools) 충족됨
- **Next.js SSR 고려사항**: 스토어가 모듈 레벨
  - **완화책**: 공식 Zustand Next.js 가이드 준수

## References

### 내부 문서

- [ADR-04: TanStack Query 선택](/ko/adr/web/04-tanstack-query-selection) - "전역 상태 라이브러리 불필요" 명시
- [ADR-16: nuqs URL 상태 관리](/ko/adr/web/16-nuqs-url-state-management) - URL/공유 가능 상태

### 외부 문서

- [Zustand 공식 문서](https://zustand.docs.pmnd.rs/)
- [Zustand GitHub 리포지토리](https://github.com/pmndrs/zustand)
- [Zustand Next.js 셋업 가이드](https://zustand.docs.pmnd.rs/guides/nextjs)

### 관련 커밋

- `8664cbc`: feat(background-tasks): add global task store with persistence
- `ecb4434`: feat(background-tasks): integrate Account Badge and Tasks Dropdown
- `fc99ce5`: feat(background-tasks): add Dashboard active tasks section
- `2e1e6df`: refactor(dashboard): migrate use-reanalyze hook to TanStack Query polling
