package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
)

const gitLogPath = "data/study_sessions.json"

type session struct {
	StartedAt       string `json:"started_at"`
	EndedAt         string `json:"ended_at"`
	DurationSeconds int    `json:"duration_seconds"`
}

func projectRoot() (string, error) {
	return os.Getwd()
}

func logPath(root string) string {
	return filepath.Join(root, "data", "study_sessions.json")
}

func ensureData(root string) error {
	dir := filepath.Join(root, "data")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	p := logPath(root)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return os.WriteFile(p, []byte("[]\n"), 0o644)
	}
	return nil
}

func loadSessions(path string) ([]session, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := strings.TrimSpace(string(raw))
	if s == "" {
		return []session{}, nil
	}
	var out []session
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return nil, fmt.Errorf("study_sessions.json 형식 오류: %w", err)
	}
	return out, nil
}

func saveSessions(path string, sessions []session) error {
	b, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func isGitRepo(root string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = root
	return cmd.Run() == nil
}

func runGit(root string, args ...string) (string, int) {
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if x, ok := err.(*exec.ExitError); ok {
			code = x.ExitCode()
		} else {
			code = 1
		}
	}
	return string(out), code
}

func formatHMS(total int) string {
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatDurationKo(total int) string {
	if total < 60 {
		return fmt.Sprintf("%d초", total)
	}
	h := total / 3600
	rem := total % 3600
	m := rem / 60
	sec := rem % 60
	var parts []string
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%d시간", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%d분", m))
	}
	if sec > 0 && h == 0 {
		parts = append(parts, fmt.Sprintf("%d초", sec))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%d초", total)
	}
	return strings.Join(parts, " ")
}

func pushStudyLog(root string, durationSec int, endedAt time.Time) {
	if out, code := runGit(root, "add", "--", gitLogPath); code != 0 {
		fmt.Fprintf(os.Stderr, "git add 실패: %s\n", strings.TrimSpace(out))
		os.Exit(1)
	}
	statusOut, statusCode := runGit(root, "status", "--porcelain", "--", gitLogPath)
	if statusCode != 0 {
		fmt.Fprintf(os.Stderr, "git status 실패: %s\n", strings.TrimSpace(statusOut))
		os.Exit(1)
	}
	if strings.TrimSpace(statusOut) == "" {
		fmt.Println("커밋할 변경이 없습니다. (이미 동일한 기록이 커밋되어 있을 수 있습니다.)")
		if out, code := runGit(root, "push"); code != 0 {
			fmt.Fprintf(os.Stderr, "git push 실패: %s\n", strings.TrimSpace(out))
			os.Exit(1)
		}
		fmt.Println("git push 완료.")
		return
	}
	msg := fmt.Sprintf(
		"학습 기록: %s (%s)",
		formatDurationKo(durationSec),
		endedAt.Local().Format("2006-01-02 15:04"),
	)
	if out, code := runGit(root, "commit", "-m", msg); code != 0 {
		fmt.Fprintf(os.Stderr, "git commit 실패: %s\n", strings.TrimSpace(out))
		os.Exit(1)
	}
	if out, code := runGit(root, "push"); code != 0 {
		fmt.Fprintf(os.Stderr, "git push 실패: %s\n", strings.TrimSpace(out))
		os.Exit(1)
	}
	fmt.Println("GitHub에 푸시했습니다.")
}

func main() {
	root, err := projectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "작업 폴더를 알 수 없습니다: %v\n", err)
		os.Exit(1)
	}
	if err := ensureData(root); err != nil {
		fmt.Fprintf(os.Stderr, "data 폴더 준비 실패: %v\n", err)
		os.Exit(1)
	}
	if !isGitRepo(root) {
		fmt.Fprintln(os.Stderr, "이 폴더가 Git 저장소가 아닙니다.\n"+
			"GitHub에서 저장소를 만든 뒤 `git init`, `git remote add origin ...` 후 다시 실행하세요.")
		os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		fmt.Fprintln(os.Stderr, "\n중단되었습니다. 기록 및 푸시는 수행하지 않습니다.")
		os.Exit(130)
	}()

	fmt.Println("학습 타이머를 시작합니다. 종료하려면 end 를 입력하세요.\n")

	startWall := time.Now().UTC()
	quitTick := make(chan struct{})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-quitTick:
				return
			case <-ticker.C:
				elapsed := int(time.Since(startWall).Seconds())
				fmt.Printf("\r경과: %s   ", formatHMS(elapsed))
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if line == "end" {
			break
		}
		if line != "" {
			fmt.Println("종료만 하려면 end 를 입력하세요.")
		}
	}
	if err := scanner.Err(); err != nil {
		close(quitTick)
		fmt.Fprintf(os.Stderr, "입력 오류: %v\n", err)
		os.Exit(1)
	}

	close(quitTick)
	time.Sleep(50 * time.Millisecond)
	fmt.Println()

	endWall := time.Now().UTC()
	durationSec := int(endWall.Sub(startWall).Seconds())

	path := logPath(root)
	sessions, err := loadSessions(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "기록 읽기 실패: %v\n", err)
		os.Exit(1)
	}
	sessions = append(sessions, session{
		StartedAt:       startWall.Format(time.RFC3339Nano),
		EndedAt:         endWall.Format(time.RFC3339Nano),
		DurationSeconds: durationSec,
	})
	if err := saveSessions(path, sessions); err != nil {
		fmt.Fprintf(os.Stderr, "기록 저장 실패: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("이번 세션: %s (%s)\n", formatDurationKo(durationSec), formatHMS(durationSec))
	fmt.Println("저장 후 원격 저장소로 푸시합니다...\n")

	pushStudyLog(root, durationSec, endWall)
}
