# claude-go

[Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code)를 Go에서 프로그래밍적으로 사용할 수 있는 래퍼 라이브러리. `os/exec`로 `claude -p`를 호출하고, 결과를 Go 구조체로 파싱하여 반환합니다.

## 사전 요구사항

- Go 1.21+
- [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code) 설치 및 인증 완료

```bash
npm install -g @anthropic-ai/claude-code
claude  # 최초 실행 시 인증
```

## 설치

```bash
go get github.com/shaul1991/claude-go
```

## 빠른 시작

```go
package main

import (
    "context"
    "fmt"
    "log"

    claude "github.com/shaul1991/claude-go"
)

func main() {
    client := claude.NewClient(
        claude.WithMaxTurns(3),
    )

    answer, err := client.Ask(context.Background(), "What is 2+2?")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(answer) // 4
}
```

## API

### 클라이언트 생성

```go
client := claude.NewClient(opts ...Option) *Client
```

### 옵션

| 옵션 | CLI 플래그 | 설명 |
|------|-----------|------|
| `WithModel(model)` | `--model` | 사용할 모델 (예: `"sonnet"`, `"opus"`) |
| `WithSystemPrompt(prompt)` | `--system-prompt` | 시스템 프롬프트 |
| `WithAppendSystemPrompt(prompt)` | `--append-system-prompt` | 시스템 프롬프트에 추가 |
| `WithAllowedTools(tools...)` | `--allowedTools` | 허용할 도구 (예: `"bash"`, `"read"`) |
| `WithMaxTurns(n)` | `--max-turns` | 최대 에이전트 턴 수 |
| `WithMaxBudget(usd)` | `--max-budget-usd` | 최대 예산 (USD) |
| `WithWorkDir(dir)` | - | 프로세스 실행 디렉토리 |
| `WithCLIPath(path)` | - | claude 바이너리 경로 (기본값: `"claude"`) |

### 메서드

#### Ask - 텍스트 응답

```go
answer, err := client.Ask(ctx, "고루틴을 한 문장으로 설명해줘.")
// answer: "고루틴은 Go 런타임이 관리하는 경량 스레드입니다."
```

#### AskJSON - 구조화된 JSON 응답

```go
resp, err := client.AskJSON(ctx, "Go 언어의 장점 3가지를 알려줘.")
fmt.Println(resp.Result)    // 응답 텍스트
fmt.Println(resp.SessionID) // 세션 이어가기용 ID
fmt.Println(resp.Usage)     // {InputTokens: 5, OutputTokens: 42}
```

#### AskWithSchema - JSON 스키마 제약 응답

```go
schema := `{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}}}`
resp, err := client.AskWithSchema(ctx, "샘플 인물 데이터를 생성해줘.", schema)
```

#### Resume - 이전 세션 이어가기

```go
resp, _ := client.AskJSON(ctx, "Go란 무엇인가?")
resp2, _ := client.Resume(ctx, resp.SessionID, "동시성 모델에 대해 더 자세히 알려줘.")
```

#### Continue - 가장 최근 세션 이어가기

```go
resp, err := client.Continue(ctx, "아까 무슨 이야기 하고 있었지?")
```

#### Pipe - stdin으로 데이터 전달

```go
file, _ := os.Open("error.log")
defer file.Close()

answer, err := client.Pipe(ctx, file, "이 에러들을 분석해줘.")
```

#### AskStream - 스트리밍 응답 (채널 기반)

```go
events, errc := client.AskStream(ctx, "1부터 5까지 세어줘.")

for ev := range events {
    fmt.Printf("[%s] %s\n", ev.Type, string(ev.Event))
}
if err := <-errc; err != nil {
    log.Fatal(err)
}
```

## 응답 타입

```go
type Response struct {
    SessionID string `json:"session_id"`
    Result    string `json:"result"`
    Cost      Cost   `json:"cost_usd"`
    Usage     Usage  `json:"usage"`
    Model     string `json:"model"`
    Duration  int    `json:"duration_ms"`
}

type StreamEvent struct {
    Type  string          `json:"type"`
    Event json.RawMessage `json:"event,omitempty"`
}
```

## 예제 실행

```bash
git clone https://github.com/shaul1991/claude-go.git
cd claude-go
go run examples/main.go
```

## 테스트 실행

```bash
go test -v ./...
```

## 프로젝트 구조

```
claude-go/
├── go.mod          # 모듈 정의
├── claude.go       # Client 구조체 및 핵심 메서드
├── options.go      # Functional options
├── types.go        # 요청/응답 타입 정의
├── stream.go       # 스트리밍 응답 처리
├── claude_test.go  # 테스트
└── examples/
    └── main.go     # 사용 예제
```

## 라이선스

MIT
