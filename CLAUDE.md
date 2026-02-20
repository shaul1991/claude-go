# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 응답 언어

- 한글로 답변할 것

## 개요

Claude Code CLI용 Go 래퍼 라이브러리. `os/exec`를 통해 `claude -p`를 실행하고 결과를 Go 구조체로 파싱한다. Go 표준 라이브러리 외 외부 의존성 없음.

**사전 요구사항:** Claude Code CLI 설치 및 인증 필요 (`npm install -g @anthropic-ai/claude-code`)

## 명령어

```bash
# 테스트 실행
go test -v ./...

# 단일 테스트 실행
go test -v -run TestBuildArgsAllOptions ./...

# 예제 실행 (인증된 Claude CLI 필요)
go run examples/main.go
```

## 아키텍처

단일 패키지 라이브러리 (`package claude`), functional options 패턴 사용.

- **claude.go** — `Client` 구조체 및 모든 공개 메서드: `Ask` (텍스트), `AskJSON` (JSON 파싱), `AskWithSchema` (JSON 스키마 제약), `Resume` (세션 ID로 이어가기), `Continue` (최근 세션 이어가기), `Pipe` (stdin 전달). 내부 헬퍼 `buildArgs` (CLI 플래그 조립), `newCmd` (`exec.Cmd` 생성).
- **options.go** — Functional options (`Option` 타입 = `func(*Client)`): 모델, 프롬프트, 허용 도구, 최대 턴, 예산, 작업 디렉토리, CLI 경로 설정.
- **types.go** — `OutputFormat` 상수 (`text`, `json`, `stream-json`) 및 응답 타입: `Response`, `Cost`, `Usage`, `StreamEvent`.
- **stream.go** — `AskStream` 메서드: `stream-json` 포맷으로 CLI 실행, 고루틴에서 stdout 파이프의 NDJSON 라인을 읽어 `StreamEvent`를 채널로 전송. `(<-chan StreamEvent, <-chan error)` 반환.
- **claude_test.go** — 클라이언트 생성 및 `buildArgs` 플래그 조립 단위 테스트. 실제 CLI를 호출하지 않음.

### 핵심 패턴

- 모든 CLI 호출은 `buildArgs` → `newCmd` → `exec.CommandContext` 경로를 따른다. 새 CLI 플래그 추가 시: `Client` 필드 추가 → options.go에 `Option` 추가 → `buildArgs`에서 연결.
- 에러 래핑: `wrapExecError`가 `exec.ExitError`에서 stderr를 추출하여 에러 메시지를 개선.
- 스트리밍은 고루틴 + 두 개의 채널(이벤트 + 에러) 패턴 사용.
