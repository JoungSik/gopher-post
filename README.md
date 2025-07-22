# Gopher Post

고퍼가 배달해주는 우편, 매일 관심사를 메일로 받아보는 서비스

## 개요

Gopher Post는 RSS 피드에서 최신 글을 자동으로 수집하여 매일 이메일로 전송하는 Go 기반 서비스입니다. Docker 컨테이너로 실행되어 `feeds.yml`에 설정된 RSS 피드들을 파싱하고, 새로운
글들을 구독자들에게 이메일로 배달합니다.

## 주요 기능

- **RSS 피드 자동 수집**: feeds.yml에 설정된 다양한 RSS/Atom 피드에서 최신 글 파싱
- **이메일 자동 배달**: 새로운 글들을 HTML 템플릿으로 이메일 전송
- **Docker 컨테이너화**: 쉬운 배포와 확장성
- **일정 기반 실행**: 매일 자동으로 실행되는 스케줄링
- **한글 지원**: 한국어 인터페이스와 Noto Sans KR 폰트 적용
- **.env 파일 지원**: 환경변수 자동 로드

## 빠른 시작

### 필요 조건

- Go 1.24.0 이상
- Docker

### 개발 환경 설정

```bash
# 저장소 클론
git clone https://github.com/JoungSik/gopher-post.git
cd gopher-post

# 의존성 설치
go mod download

# 개발 모드로 실행
go run cmd/main.go
```

### Docker로 실행

```bash
# Docker 이미지 빌드
docker build -t gopher-post .

# 환경 변수와 함께 컨테이너 실행
docker run -e SMTP_HOST=smtp.gmail.com \
           -e SMTP_PORT=587 \
           -e SMTP_USERNAME=your-email@gmail.com \
           -e SMTP_PASSWORD=your-app-password \
           -e FROM_EMAIL=gopher-post@example.com \
           gopher-post

# 또는 환경 변수 파일을 사용
docker run --env-file .env gopher-post
```

## 설정

### feeds.yml 설정

```yaml
# 예시 설정
feeds:
  - name: "김정식"
    url: "https://joungsik.github.io"
    rss: "https://joungsik.github.io/post/index.xml"
  - name: "기술 블로그"
    url: "https://example.com/tech"
    rss: "https://example.com/tech/feed.xml"
```

### recipients.yml 설정

```yaml
# 이메일 수신자 설정
recipients:
  - email: "user@example.com"
  - email: "another@example.com"
```

### 환경 변수

```bash
# SMTP 설정
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# 발송자 정보
FROM_EMAIL=gopher-post@example.com
```

## 아키텍처

```
┌─────────────────┐    ┌──────────────┐    ┌─────────────────┐
│   feeds.yml     │───▶│  Feed Parser │───▶│  Email Service  │
│   (RSS 설정)    │    │   (RSS 파싱) │    │   (메일 발송)   │
└─────────────────┘    └──────────────┘    └─────────────────┘
                              │
                       ┌──────────────┐
                       │   Template   │
                       │   Engine     │
                       │  (이메일 템플릿)│
                       └──────────────┘
```

## 개발

### 일반적인 명령어

```bash
# 코드 포맷팅
go fmt ./...

# 정적 분석
go vet ./...

# 테스트 실행
go test ./...

# 의존성 정리
go mod tidy
```

### 프로젝트 구조

```
gopher-post/
├── .github/workflows/
│   ├── build.yml       # Docker 이미지 빌드 및 푸시
│   └── mailer.yml      # 매일 7시 자동 메일 발송
├── cmd/
│   └── main.go         # 메인 애플리케이션
├── internal/
│   ├── config/         # 설정 관리 (feeds.yml, recipients.yml 파싱)
│   ├── feed/           # RSS 피드 파싱
│   ├── email/          # 이메일 발송 (SMTP)
│   └── template/       # 이메일 템플릿 처리
├── templates/          # 이메일 템플릿
├── feeds.yml           # RSS 피드 설정 파일
├── recipients.yml      # 이메일 수신자 설정 파일
├── Dockerfile          # Docker 컨테이너 설정
├── go.mod              # Go 모듈 정의
├── CLAUDE.md           # Claude Code 가이드
└── README.md
```

## GitHub Actions 자동화

### 필요한 GitHub Secrets 설정

Repository Settings > Secrets and variables > Actions에서 다음 secrets를 설정하세요:

```
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=your-email@gmail.com
```

### 자동화 워크플로우

- **build.yml**: main 브랜치에 push시 자동으로 Docker 이미지 빌드 및 GitHub Container Registry에 푸시
- **mailer.yml**: 매일 오전 7시(KST)에 자동으로 뉴스레터 발송

### 수동 실행

GitHub Actions 탭에서 "Daily Newsletter Mailer" 워크플로우를 수동으로 실행할 수 있습니다.

## 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다.

## 문제 신고 및 기능 요청

GitHub Issues를 통해 버그 신고나 기능 요청을 해주세요.
