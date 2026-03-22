# 학습 체크 타이머 (터미널, Go)

학습을 시작한 뒤 터미널에서 **`end`** 를 입력하면, 그동안의 학습 시간을 기록하고 **Git 커밋 후 `git push`** 로 GitHub에 올립니다.

## 필요한 것

- **Go 1.21+**
- **Git** 설치 및 PATH 설정
- **미리 만든 GitHub 저장소** (빈 저장소여도 됨)

## 사전 준비 (GitHub 연동)

1. GitHub에서 새 저장소를 만듭니다.
2. 이 프로젝트 폴더에서 Git을 초기화하고 원격을 연결합니다.

```bash
cd 프로젝트_폴더
git init
git remote add origin https://github.com/사용자이름/저장소이름.git
```

3. 첫 커밋을 한 번 올려 두면 이후 푸시가 수월합니다.

```bash
git add .
git commit -m "초기 설정"
git branch -M main
git push -u origin main
```

> **인증**: HTTPS는 Personal Access Token, SSH는 키 등 GitHub 권장 방식으로 로그인해 두어야 `git push`가 됩니다.

## 빌드 및 실행

**프로젝트 루트(이 README가 있는 폴더)에서** 실행해야 `data/study_sessions.json` 경로와 Git 동작이 맞습니다.

```bash
go run .
```

실행 파일로 쓰려면:

```bash
go build -o study-timer .
./study-timer
```

Windows PowerShell 예:

```powershell
go build -o study-timer.exe .
.\study-timer.exe
```

## 사용 방법

- 실행되면 **경과 시간**이 같은 줄에 갱신됩니다.
- 학습을 마치면 **`end`** 를 입력하고 Enter 합니다.
- `data/study_sessions.json`에 세션이 추가되고, 해당 파일만 스테이징한 뒤 커밋·푸시합니다.

### 기록 형식 (`data/study_sessions.json`)

각 항목은 대략 다음과 같습니다.

```json
{
  "started_at": "2025-03-23T12:00:00.000000000Z",
  "ended_at": "2025-03-23T12:45:00.000000000Z",
  "duration_seconds": 2700
}
```

## 동작 요약

| 단계 | 내용 |
|------|------|
| 시작 | 프로그램 실행 시각을 세션 시작으로 기록 |
| 종료 | `end` 입력 시 종료 시각까지의 초 단위 시간 저장 |
| Git | `git add data/study_sessions.json` → `git commit` → `git push` |

변경이 없으면(이미 커밋된 상태와 동일하면) 커밋은 생략하고 `git push`만 시도할 수 있습니다.

## 주의

- **Ctrl+C** 로 끊으면 기록·푸시는 하지 않고 종료합니다.
- 이 폴더가 Git 저장소가 아니면 시작 시 안내 메시지 후 종료합니다.
- `origin`에 푸시 권한이 있어야 합니다.

## 앞으로 확장하기 (아이디어)

- 과목·프로그램 이름 필드 추가 (`subject`, `tag` 등)
- 하루 합계·주간 통계 출력
- 웹 UI 또는 데스크톱 알림

현재 버전은 **“학습 시간만 기록 + 종료 시 GitHub 반영”** 에 집중했습니다.

## 라이선스

개인 학습용으로 자유롭게 수정해 사용하시면 됩니다.
