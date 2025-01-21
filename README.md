# 프로젝트 소개

이 저장소(Repository)는 **Google OAuth 2.0** 로그인을 **Go**로 구현한 예제와, **Go(백엔드)** 및 **React(프론트엔드)**로 제작된 **실시간 방송 애플리케이션** 예제를 함께 제공합니다.

---

## 목차

1. [Go를 활용한 Google OAuth 통합](#go를-활용한-google-oauth-통합)
   - [주요 기능](#주요-기능)
   - [폴더 구조](#폴더-구조)
2. [실시간 방송 애플리케이션](#실시간-방송-애플리케이션)
   - [주요 기능](#주요-기능-1)
   - [프로젝트 구조](#프로젝트-구조)
   - [설치 방법](#설치-방법)
     - [백엔드](#백엔드)
     - [프론트엔드](#프론트엔드)
   - [사용 방법](#사용-방법)
   - [주요 포인트](#주요-포인트)

---

## Go를 활용한 Google OAuth 통합

이 프로젝트 섹션은 **Google OAuth 2.0 로그인**을 Go로 구현하는 방법을 다룹니다.  
Google API를 호출하여 Access Token을 활용해 사용자 정보를 가져오는 과정을 포함하며, 사용자 인증 및 JWT 생성까지 처리합니다.

### 주요 기능

- **Resty**를 사용하여 Google API와 통신
- Google 사용자 정보(`/userinfo/v2/me`)를 가져오기
- `access_token`을 쿼리 파라미터로 전달하여 OAuth 인증 처리
- API 응답을 구조화된 Go 모델(`google_response.go`)로 파싱
- 데이터베이스에서 사용자 확인 및 자동 등록
- 기존 사용자에 대해 JWT 토큰 생성

### 폴더 구조

```plaintext
/toysgo
    ├── controllers
    │   └── rest
    │       ├── google.go            # GoogleController로 OAuth 로직 처리
    ├── services
    │   └── google_user_api.go       # Google API 호출을 위한 서비스
    ├── models
    │   └── google_response.go       # Google API 응답 모델
    ├── global
    │   └── auth.go                  # JWT 생성 로직
    └── main.go                      # 테스트를 위한 진입점

```

실시간 방송 애플리케이션

WebSocket을 사용한 실시간 방송 애플리케이션으로, **Go(백엔드)**와 **React(프론트엔드)**로 구현되었습니다.
STUN/TURN을 사용하여 NAT 및 방화벽을 넘어 원활한 P2P 연결을 지원합니다.

주요 기능
• 송출자 (Broadcaster): 비디오/오디오를 실시간으로 시청자에게 스트리밍
• 시청자 (Viewer): 실시간 방송을 시청
• WebSocket 통신: 방송 송출과 시청자 간 신호 처리 및 미디어 협상
• STUN/TURN 지원: NAT 및 방화벽을 넘어 P2P 연결 지원

프로젝트 구조
• 백엔드:
• Go로 작성
• WebSocket 연결, 신호 처리, 메시지 릴레이 담당
• STUN/TURN을 사용한 네트워크 트래버설 지원
• 프론트엔드:
• React + TypeScript
• Broadcast 컴포넌트: 실시간 스트리밍 송출
• Viewer 컴포넌트: 실시간 방송 시청

설치 방법

백엔드 1. 의존성 설치:

go mod tidy

    2.	서버 실행:

go run main.go

    3.	기본 서버 URL:

http://localhost:9000

프론트엔드 1. 프론트엔드 디렉토리로 이동:

cd frontend

    2.	의존성 설치:

npm install

    3.	React 앱 실행:

npm start

    4.	웹 앱 주소:

http://localhost:3000

사용 방법 1. 방송 송출
• 송출자 페이지 접속
• Start Broadcast 버튼으로 방송 시작 2. 방송 시청
• 시청자 페이지 접속
• 실시간으로 방송 시청

주요 포인트
• 백엔드와 프론트엔드가 모두 실행 중인지 확인
• 제한된 네트워크 환경에서는 TURN 서버를 설정해 안정적인 P2P 연결 지원
• WebSocket 로그를 확인하여 연결 문제 디버깅 가능

Author: 다양한 기능 테스트 및 학습 목적으로 작성된 예제 프로젝트.
