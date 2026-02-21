# claude-go

[Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code)를 Go에서 프로그래밍적으로 사용할 수 있는 래퍼 라이브러리. `os/exec`로 `claude -p`를 호출하고, 결과를 Go 구조체로 파싱하여 반환합니다.

Anthropic Messages API(`POST /v1/messages`) 호환 HTTP 서버도 내장하고 있어, 기존 Anthropic API 클라이언트와 호환되는 프록시 서버로 활용할 수 있습니다.

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

## Go 라이브러리 API

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

## HTTP 서버 (Anthropic Messages API 호환)

내장 HTTP 서버는 Anthropic Messages API(`POST /v1/messages`)와 동일한 인터페이스를 제공합니다. 기존 Anthropic API 클라이언트에서 엔드포인트만 변경하면 바로 사용할 수 있습니다.

### 서버 실행

```bash
# 빌드
go build -o bin/claude-server ./cmd/server

# 기본 실행 (포트 8080, 인증 없음)
./bin/claude-server

# 옵션 지정
./bin/claude-server \
  -port 8080 \
  -api-key "my-secret-key" \
  -cli-path /usr/local/bin/claude \
  -work-dir /path/to/project \
  -max-budget 1.0 \
  -max-turns 10
```

### 서버 옵션

| 플래그 | 환경변수 | 설명 |
|--------|---------|------|
| `-port` | `PORT` | 서버 포트 (기본값: `8080`) |
| `-host` | - | 서버 호스트 (기본값: `0.0.0.0`) |
| `-api-key` | `API_KEY` | API 키 (빈 값이면 인증 건너뜀) |
| `-cli-path` | `CLAUDE_CLI_PATH` | claude CLI 바이너리 경로 |
| `-work-dir` | `CLAUDE_WORK_DIR` | claude CLI 실행 디렉토리 |
| `-max-budget` | - | 요청당 최대 예산 (USD) |
| `-max-turns` | - | 요청당 최대 턴 수 |

### API 엔드포인트

#### `GET /health`

헬스 체크.

```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

#### `POST /v1/messages`

Anthropic Messages API 호환 엔드포인트. 비스트리밍 및 스트리밍 모두 지원.

**인증:** `x-api-key` 헤더 또는 `Authorization: Bearer <key>` (서버에 `-api-key` 설정 시).

**비스트리밍 요청:**

```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "content-type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "sonnet",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "What is 2+2?"}]
  }'
```

**스트리밍 요청:**

```bash
curl -X POST http://localhost:8080/v1/messages \
  -H "content-type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "sonnet",
    "max_tokens": 1024,
    "stream": true,
    "messages": [{"role": "user", "content": "Count 1 to 5"}]
  }'
```

**시스템 프롬프트 (문자열 또는 배열):**

```json
{
  "model": "sonnet",
  "max_tokens": 1024,
  "system": "You are a helpful assistant.",
  "messages": [{"role": "user", "content": "Hello"}]
}
```

## 응답 타입

```go
// Go 라이브러리 응답
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
├── go.mod              # 모듈 정의
├── claude.go           # Client 구조체 및 핵심 메서드
├── options.go          # Functional options
├── types.go            # 요청/응답 타입 정의
├── stream.go           # 스트리밍 응답 처리
├── claude_test.go      # 테스트
├── examples/
│   └── main.go         # 사용 예제
├── cmd/
│   └── server/
│       └── main.go     # HTTP 서버 엔트리포인트
└── internal/
    └── server/
        ├── server.go     # 서버 설정, 라우팅, 미들웨어 체인
        ├── handler.go    # Messages API 핸들러 (스트리밍/비스트리밍)
        ├── middleware.go  # 로깅, CORS, 인증 미들웨어
        └── types.go      # Anthropic Messages API 타입 정의
```

## 라이선스

MIT
