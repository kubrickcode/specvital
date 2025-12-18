---
title: Asynq slog 어댑터
---

# ADR-06: Asynq 로깅 slog 어댑터

> 🇺🇸 [English Version](/en/adr/collector/06-asynq-slog-adapter.md)

| 날짜       | 작성자       | 리포지토리 |
| ---------- | ------------ | ---------- |
| 2024-12-18 | @KubrickCode | collector  |

## 배경

### 문제: 클라우드 플랫폼에서 stderr vs stdout

Asynq는 기본적으로 Go 표준 `log` 패키지를 사용하며, 이는 stderr로 출력함. 많은 클라우드 플랫폼과 로그 수집 서비스는 실제 로그 심각도와 관계없이 stderr 출력을 ERROR 레벨로 분류함.

**영향:**

- Asynq의 INFO 레벨 메시지가 모니터링 대시보드에서 에러로 표시
- 정상 운영 로그로 인해 알림이 발생할 수 있음
- 실제 에러와 일반 메시지 구분이 어려움
- 로그 필터링의 신뢰성 저하

### 구조화된 로깅 요구사항

현대 관측성(Observability) 관행은 다음을 위해 구조화된(JSON) 로깅을 선호한다:

- 기계 파싱 가능한 로그 항목
- 일관된 필드 추출
- 서비스 간 상관관계 분석
- 효율적인 쿼리 및 필터링

표준 `log` 패키지는 비구조화된 텍스트를 출력하여, 구조화된 로깅 파이프라인과 호환되지 않음.

## 결정

**`slog.Logger`를 래핑하여 `asynq.Logger` 인터페이스를 충족하는 `SlogAdapter` 구현**

핵심 설계 선택:

1. **slog 사용**: Go 1.21+ 표준 라이브러리 구조화 로거
2. **stdout 출력**: stderr 대신 stdout으로 JSON 로그 출력
3. **Fatal에서 os.Exit 미호출**: severity 속성으로 로깅, 종료는 호출자에게 위임

## 검토한 옵션

### 옵션 A: SlogAdapter (선택)

**구현:**

```go
type SlogAdapter struct {
    logger *slog.Logger
}

func (l *SlogAdapter) Info(args ...any) {
    l.logger.Info(fmt.Sprint(args...))
}

func (l *SlogAdapter) Fatal(args ...any) {
    l.logger.Error(fmt.Sprint(args...), "severity", "fatal")
}
```

**장점:**

- Go 표준 라이브러리 사용 (외부 의존성 없음)
- 플랫폼 호환을 위한 stdout JSON 출력
- 애플리케이션 전체 로깅 전략과 일관성
- Fatal이 갑작스럽게 종료하지 않음; 우아한 처리 가능

**단점:**

- 어댑터 구현 및 유지보수 필요
- 예상과 약간 다른 Fatal 동작

### 옵션 B: Asynq 기본 Logger

- stderr 출력을 사용하는 Go `log.Logger`
- 코드 변경 불필요
- 클라우드 환경에서 플랫폼 호환성 문제

### 옵션 C: 서드파티 로깅 라이브러리

- zap, zerolog, 또는 logrus를 Asynq 어댑터와 사용
- 추가 의존성 발생
- 기존 로깅 선택과 충돌 가능성

## 구현 노트

### asynq.Logger 인터페이스

```go
type Logger interface {
    Debug(args ...any)
    Info(args ...any)
    Warn(args ...any)
    Error(args ...any)
    Fatal(args ...any)
}
```

모든 메서드는 가변 인수를 받으며, 어댑터는 `fmt.Sprint`로 이를 결합함.

### Fatal 동작

어댑터는 의도적으로 `os.Exit` 호출을 피함. 대신 `"severity": "fatal"` 속성과 함께 Error 레벨로 로깅함. 이를 통해:

- 지연된 정리(deferred cleanup)가 실행됨
- 우아한 종료 처리가 가능
- 갑작스러운 종료 없이 로그에서 fatal 레벨 가시성 제공

### 설정

```go
srv := asynq.NewServer(
    redisOpt,
    asynq.Config{
        Logger: queue.NewSlogAdapter(slog.Default()),
        // 기타 설정...
    },
)
```

## 결과

### 긍정적

**플랫폼 호환성:**

- 클라우드 플랫폼에서 심각도에 따른 올바른 로그 분류
- 정상 운영에서 잘못된 에러 알림 없음
- 신뢰할 수 있는 로그 레벨 필터링

**구조화된 로깅:**

- 기계 파싱이 가능한 JSON 형식
- 애플리케이션 로깅과 일관성
- 구조화된 필드로 더 쉬운 디버깅

**운영 명확성:**

- 실제 에러와 일반 메시지의 명확한 분리
- 모니터링 대시보드 정확도 향상
- 알림 피로도 감소

### 부정적

**유지보수 오버헤드:**

- 유지해야 할 커스텀 어댑터 코드
- asynq.Logger 인터페이스 변경 추적 필요

**의미론적 차이:**

- Fatal이 프로세스를 종료하지 않음
- 즉시 종료를 기대하는 호출자가 놀랄 수 있음

## 참고 자료

- [Go slog 패키지 문서](https://pkg.go.dev/log/slog)
- [Asynq Logger 인터페이스](https://pkg.go.dev/github.com/hibiken/asynq#Logger)
