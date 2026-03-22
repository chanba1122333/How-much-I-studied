# `main.go` 코드 설명 (한 줄·한 블록씩)

이 문서는 `study-timer` 프로젝트의 `main.go`가 무엇을 하는지, 각 줄이 어떤 의미인지 정리한 것입니다.

---

## 1. 패키지와 import (1~13행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 1 | `package main` | 실행 가능한 Go 프로그램의 진입점이 되는 패키지 이름입니다. `main` 패키지에 `func main()`이 있어야 `go run`으로 실행됩니다. |
| 3~13 | `import (...)` | 이 파일에서 쓰는 표준 라이브러리를 가져옵니다. |
| 4 | `"bufio"` | 표준 입력을 **한 줄씩** 읽을 때 쓰는 `Scanner` 등이 들어 있습니다. |
| 5 | `"encoding/json"` | JSON 문자열 ↔ Go 구조체/슬라이스 변환(`Marshal`, `Unmarshal`)에 사용합니다. |
| 6 | `"fmt"` | `Printf`, `Println`, `Sprintf` 같은 **포맷 출력** 함수입니다. |
| 7 | `"os"` | 현재 작업 디렉터리, 파일 읽기/쓰기, 종료 코드(`Exit`) 등 **운영체제와의 인터페이스**입니다. |
| 8 | `"os/exec"` | **외부 명령**(`git` 등)을 자식 프로세스로 실행할 때 사용합니다. |
| 9 | `"os/signal"` | Ctrl+C 같은 **시그널**을 Go 채널로 받기 위해 사용합니다. |
| 10 | `"path/filepath"` | OS에 맞는 경로 구분자로 경로를 이을 때 `filepath.Join`을 씁니다. |
| 11 | `"strings"` | 문자열 자르기, 소문자 변환, 합치기 등 문자열 유틸입니다. |
| 12 | `"time"` | 시각(`Now`), 간격(`Since`), 타이머(`Ticker`), 포맷(`Format`)을 다룹니다. |

---

## 2. 상수와 데이터 구조 (15~21행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 15 | `const gitLogPath = "data/study_sessions.json"` | Git에 올릴 **학습 기록 파일의 저장소 루트 기준 상대 경로**입니다. `git add` 등에 그대로 넘깁니다. |
| 17~21 | `type session struct { ... }` | JSON 한 건(한 번의 학습 세션)에 대응하는 **구조체**입니다. |
| 18 | `StartedAt string \`json:"started_at"\`` | 세션 **시작 시각**을 문자열로 저장합니다. 태그는 JSON 필드 이름이 `started_at`이 되도록 합니다. |
| 19 | `EndedAt string \`json:"ended_at"\`` | 세션 **종료 시각**입니다. JSON 키는 `ended_at`. |
| 20 | `DurationSeconds int \`json:"duration_seconds"\`` | **초 단위** 학습 시간입니다. JSON 키는 `duration_seconds`. |

---

## 3. `projectRoot` (23~25행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 23 | `func projectRoot() (string, error)` | 프로그램이 실행되는 **현재 작업 디렉터리**(프로젝트 루트로 가정)를 돌려줍니다. |
| 24 | `return os.Getwd()` | OS에 물어본 **현재 디렉터리 절대/상대 경로**와, 실패 시 에러를 반환합니다. |

---

## 4. `logPath` (27~29행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 27 | `func logPath(root string) string` | 루트 경로와 파일 이름을 합쳐 **학습 기록 JSON의 전체 경로**를 만듭니다. |
| 28 | `return filepath.Join(root, "data", "study_sessions.json")` | Windows는 `\`, Unix는 `/`로 알아서 이어 줍니다. |

---

## 5. `ensureData` (31~41행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 31 | `func ensureData(root string) error` | `data` 폴더와 빈 JSON 파일이 없으면 **만들어 두는** 준비 단계입니다. |
| 32 | `dir := filepath.Join(root, "data")` | `data` 디렉터리 경로입니다. |
| 33~35 | `if err := os.MkdirAll(dir, 0o755); err != nil` | 없으면 **중첩 포함해 전부 생성**합니다. `0o755`는 Unix 권한(소유자 읽기/쓰기/실행, 그 외 읽기/실행)입니다. |
| 36 | `p := logPath(root)` | JSON 파일 전체 경로입니다. |
| 37 | `if _, err := os.Stat(p); os.IsNotExist(err)` | 파일이 **없으면** `IsNotExist`가 참입니다. |
| 38 | `return os.WriteFile(p, []byte("[]\n"), 0o644)` | **빈 JSON 배열** `[]`과 줄바꿈을 써서 파일을 만듭니다. `0o644`는 일반 파일 권한입니다. |
| 40 | `return nil` | 이미 있으면 아무 것도 안 하고 성공입니다. |

---

## 6. `loadSessions` (43~57행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 43 | `func loadSessions(path string) ([]session, error)` | JSON 파일을 읽어 **세션 슬라이스**로 파싱합니다. |
| 44~47 | `raw, err := os.ReadFile(path)` | 파일 전체를 바이트로 읽습니다. 실패하면 에러 반환. |
| 48 | `s := strings.TrimSpace(string(raw))` | 앞뒤 공백·줄바꿈을 제거한 문자열입니다. |
| 49~51 | `if s == "" { return []session{}, nil }` | 내용이 비어 있으면 **빈 목록**으로 간주합니다. |
| 52 | `var out []session` | 파싱 결과를 담을 슬라이스입니다. |
| 53~55 | `json.Unmarshal([]byte(s), &out)` | JSON 배열을 Go의 `[]session`으로 채웁니다. 실패 시 한국어 메시지와 함께 `fmt.Errorf`로 감쌉니다. |
| 56 | `return out, nil` | 성공 시 기록 목록을 돌려줍니다. |

---

## 7. `saveSessions` (59~65행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 59 | `func saveSessions(path string, sessions []session) error` | 세션 목록을 JSON 파일로 **덮어씁니다**. |
| 60 | `json.MarshalIndent(sessions, "", "  ")` | 들여쓰기(스페이스 2칸) 있는 **읽기 쉬운 JSON**으로 직렬화합니다. |
| 64 | `os.WriteFile(path, append(b, '\n'), 0o644)` | 파일 끝에 **줄바꿈 한 번**을 붙여 POSIX 관례에 맞춥니다. |

---

## 8. `isGitRepo` (67~71행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 67 | `func isGitRepo(root string) bool` | 이 폴더가 **Git 저장소 안**인지 확인합니다. |
| 68 | `exec.Command("git", "rev-parse", "--git-dir")` | `.git` 디렉터리 위치를 묻는 Git 명령입니다. 성공하면 저장소로 인정됩니다. |
| 69 | `cmd.Dir = root` | 명령을 **프로젝트 루트**에서 실행합니다. |
| 70 | `return cmd.Run() == nil` | 에러 없이 끝나면 `true`, 아니면 `false`입니다. |

---

## 9. `runGit` (73~86행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 73 | `func runGit(root string, args ...string) (string, int)` | `git` 뒤에 붙일 인자들을 받아 실행하고, **출력 문자열**과 **종료 코드**를 돌려줍니다. |
| 74 | `exec.Command("git", args...)` | 예: `runGit(root, "add", "file")` → `git add file`. |
| 75 | `cmd.Dir = root` | 항상 프로젝트 루트에서 Git을 돌립니다. |
| 76 | `out, err := cmd.CombinedOutput()` | 표준 출력과 표준 에러를 **합친** 바이트를 가져옵니다. |
| 77 | `code := 0` | 기본은 성공(0)입니다. |
| 78~84 | `if err != nil { ... }` | 실패하면 `*exec.ExitError`이면 **실제 exit code**를 쓰고, 그 외는 1로 둡니다. |
| 85 | `return string(out), code` | 사람이 읽을 문자열과 코드를 반환합니다. |

---

## 10. `formatHMS` (88~93행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 88 | `func formatHMS(total int) string` | 총 **초**를 `시:분:초` 형태로 만듭니다. |
| 89~91 | `h`, `m`, `s` | 각각 시간·분·초로 나눈 몫과 나머지입니다. |
| 92 | `fmt.Sprintf("%02d:%02d:%02d", h, m, s)` | 각 항목을 **두 자리 0 패딩**합니다 (예: `01:05:09`). |

---

## 11. `formatDurationKo` (95~117행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 95 | `func formatDurationKo(total int) string` | 한국어로 **“N시간 M분”** 같은 짧은 설명 문자열을 만듭니다. |
| 96~98 | `if total < 60` | 1분 미만이면 **“N초”**만 씁니다. |
| 99~102 | `h`, `rem`, `m`, `sec` | 시간·분·초를 정수로 쪼갭니다. |
| 103 | `var parts []string` | 문장 조각을 모을 슬라이스입니다. |
| 104~106 | `if h > 0` | 시간이 있으면 `"%d시간"`을 붙입니다. |
| 107~109 | `if m > 0` | 분이 있으면 `"%d분"`을 붙입니다. |
| 110~112 | `if sec > 0 && h == 0` | **시간이 0일 때만** 초를 넣습니다 (1시간 넘는 세션에서는 초 생략으로 문장을 짧게). |
| 113~115 | `if len(parts) == 0` | 위 조건으로 아무 것도 안 붙었으면(예: 경계 케이스) 다시 `"%d초"`로 처리합니다. |
| 116 | `strings.Join(parts, " ")` | 조각들을 공백으로 이어 한 문장으로 만듭니다. |

---

## 12. `pushStudyLog` (119~152행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 119 | `func pushStudyLog(root string, durationSec int, endedAt time.Time)` | 학습이 끝난 뒤 **Git add → (필요 시) commit → push**까지 수행합니다. |
| 120~123 | `git add -- data/study_sessions.json` | `--` 뒤는 경로만 해석해 실수로 옵션으로 읽히지 않게 합니다. 실패 시 stderr에 출력하고 `Exit(1)`. |
| 124~128 | `git status --porcelain -- 파일` | 해당 파일에 **스테이징된 변경이 있는지** 한 줄 형식으로 확인합니다. |
| 129~136 | 변경이 없으면 | 이미 커밋된 내용과 같다는 뜻일 수 있어 메시지를 찍고, **`git push`만** 시도한 뒤 return합니다. |
| 138~142 | `msg := fmt.Sprintf(...)` | 커밋 메시지에 **한국어 기간**과 **종료 시각(로컬)**을 넣습니다. `2006-01-02 15:04`는 Go의 레이아웃 참조 시각입니다. |
| 143~146 | `git commit -m msg` | 스테이징된 변경이 있을 때만 의미 있는 커밋이 됩니다. |
| 147~150 | `git push` | 원격으로 올립니다. |
| 151 | 성공 메시지 출력 | |

---

## 13. `main` — 시작과 준비 (154~168행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 154 | `func main()` | 프로그램 **진입점**입니다. |
| 155~159 | `projectRoot()` | 현재 디렉터리를 루트로 쓰고, 실패하면 종료합니다. |
| 160~163 | `ensureData(root)` | `data` 폴더와 빈 `study_sessions.json`을 보장합니다. |
| 164~168 | `!isGitRepo(root)` | Git 저장소가 아니면 안내 문구를 stderr에 쓰고 종료합니다. |

---

## 14. `main` — Ctrl+C 처리 (170~176행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 170 | `interrupt := make(chan os.Signal, 1)` | 시그널을 **버퍼 1**인 채널로 받습니다 (블로킹 없이 한 번 받을 공간). |
| 171 | `signal.Notify(interrupt, os.Interrupt)` | **Ctrl+C**(일반적으로 `SIGINT`)를 이 채널로 보냅니다. |
| 172~176 | `go func() { <-interrupt; ... os.Exit(130) }()` | 별도 고루틴에서 Ctrl+C를 기다립니다. **130**은 관례적인 “사용자 중단” 종료 코드입니다. 이 경로에서는 **JSON 저장·Git 푸시를 하지 않습니다**. |

---

## 15. `main` — 타이머 UI (178~195행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 178 | 안내 `Println` | 사용자에게 **`end` 입력으로 종료**한다고 알립니다. |
| 180 | `startWall := time.Now().UTC()` | 세션 **시작 시각**(UTC)입니다. 기록과 길이 계산의 기준입니다. |
| 181 | `quitTick := make(chan struct{})` | **빈 struct 채널**으로 “타이머 고루틴 그만” 신호를 보냅니다. |
| 183~195 | 고루틴 + `time.NewTicker(time.Second)` | **1초마다** 깨어나 `select`로 `quitTick`과 비교합니다. |
| 188~189 | `<-quitTick` | 닫히거나 값이 오면 **루프 종료**하고 고루틴이 return합니다. |
| 190~192 | `<-ticker.C` | 1초마다 경과 초를 `Since(startWall)`으로 구해 **`\r`로 같은 줄 덮어쓰기** 하며 `경과: HH:MM:SS`를 표시합니다. |
| 185 | `defer ticker.Stop()` | 고루틴이 끝날 때 타이머 리소스를 정리합니다. |

---

## 16. `main` — 표준 입력 루프 (197~211행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 197 | `bufio.NewScanner(os.Stdin)` | **한 줄씩** 읽는 스캐너입니다. |
| 198 | `for scanner.Scan()` | 줄이 있을 때마다 반복합니다. EOF면 false가 됩니다. |
| 199 | `line := strings.TrimSpace(strings.ToLower(scanner.Text()))` | 앞뒤 공백 제거, **소문자**로 통일해 `end`/`End` 모두 인식합니다. |
| 200~201 | `if line == "end" { break }` | 종료 키워드면 입력 루프를 빠져나갑니다. |
| 203~205 | 그 외 비어 있지 않은 줄 | `end`만 허용한다고 다시 알려 줍니다. |
| 207~211 | `scanner.Err()` | 읽기 오류면 `quitTick`을 닫아 타이머를 멈추고, 에러 출력 후 종료합니다. |

---

## 17. `main` — 종료 후 정리·저장·푸시 (213~240행)

| 줄 | 코드 | 의미 |
|---|------|------|
| 213 | `close(quitTick)` | 타이머 고루틴이 **select에서 quit 쪽**을 받고 종료되도록 채널을 닫습니다. |
| 214 | `time.Sleep(50 * time.Millisecond)` | 마지막 `\r` 출력이 터미널에 반영될 **짧은 여유**입니다. |
| 215 | `fmt.Println()` | 줄바꿈으로 프롬프트를 다음 줄로 깔끔히 내립니다. |
| 217 | `endWall := time.Now().UTC()` | 세션 **종료 시각**(UTC). |
| 218 | `durationSec := int(endWall.Sub(startWall).Seconds())` | 시작~종료 **실제 경과**를 초 단위 정수로 버림합니다. |
| 220 | `path := logPath(root)` | JSON 파일 경로입니다. |
| 221~225 | `loadSessions` | 기존 기록을 읽습니다. 실패 시 종료. |
| 226~230 | `append(..., session{...})` | 새 세션 한 건을 **RFC3339 나노초** 문자열과 초 단위 길이와 함께 슬라이스 끝에 붙입니다. |
| 231~234 | `saveSessions` | 전체 목록을 파일에 다시 씁니다. |
| 236~237 | `Printf` / `Println` | 이번 세션 요약과 푸시 안내를 출력합니다. |
| 239 | `pushStudyLog(root, durationSec, endWall)` | Git add / commit / push 흐름을 실행합니다. |

---

## 전체 흐름 요약

1. 현재 폴더를 루트로 두고 `data/study_sessions.json`을 준비합니다.  
2. Git 저장소인지 확인합니다.  
3. Ctrl+C는 별도 처리해 **기록 없이** 종료할 수 있습니다.  
4. `end`를 입력할 때까지 1초마다 경과 시간을 같은 줄에 갱신합니다.  
5. `end` 후 경과 시간을 JSON 배열에 append하고 저장한 뒤, 해당 파일을 커밋·푸시합니다.  

`go.mod`의 `module study-timer`는 모듈 이름이며, `go 1.21`은 이 프로젝트가 사용하는 **최소 Go 버전**입니다.
