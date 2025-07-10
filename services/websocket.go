package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/websocket/v2"
)

// 방송 정보 구조체
type BroadcastInfo struct {
	BroadcasterID   string    `json:"broadcaster_id"`
	BroadcasterName string    `json:"broadcaster_name"`
	StartTime       time.Time `json:"start_time"`
	ViewerCount     int       `json:"viewer_count"`
	IsLive          bool      `json:"is_live"`
}

// 연결 정보 구조체
type Connection struct {
	Conn   *websocket.Conn
	UserID string
	Name   string
	Role   string
}

// 시청자 정보 구조체
type ViewerInfo struct {
	Connection      *Connection
	BroadcasterID   string
	JoinTime        time.Time
}

// WebSocket 메시지 구조체
type Message struct {
	Type          string                 `json:"type"`
	Data          interface{}            `json:"data,omitempty"`
	BroadcasterID interface{}            `json:"broadcaster_id,omitempty"`
	ViewerID      interface{}            `json:"viewer_id,omitempty"`
	ViewerName    string                 `json:"viewer_name,omitempty"`
	Timestamp     string                 `json:"timestamp,omitempty"`
	Count         int                    `json:"count,omitempty"`
	Broadcasts    []BroadcastInfo        `json:"broadcasts,omitempty"`
	Broadcast     *BroadcastInfo         `json:"broadcast,omitempty"`
}

// ID를 문자열로 변환하는 헬퍼 함수
func toString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type WebSocketService struct {
	// 방송자 연결 관리 (broadcaster_id -> Connection)
	Broadcasters map[string]*Connection
	
	// 시청자 연결 관리 (viewer_id -> ViewerInfo)
	Viewers map[string]*ViewerInfo
	
	// 방송 목록 구독자 (connection -> user_id)
	ListSubscribers map[*websocket.Conn]string
	
	// 활성 방송 목록 (broadcaster_id -> BroadcastInfo)
	ActiveBroadcasts map[string]*BroadcastInfo
	
	// 대기 중인 Offer (viewer_id -> broadcaster_id)
	PendingOffers map[string]string
	
	Mutex       sync.RWMutex
	initialized bool
}

// 안전한 초기화
func NewWebSocketService() *WebSocketService {
	fmt.Println("🔧 WebSocketService 초기화 시작...")
	
	service := &WebSocketService{
		Broadcasters:     make(map[string]*Connection),
		Viewers:          make(map[string]*ViewerInfo),
		ListSubscribers:  make(map[*websocket.Conn]string),
		ActiveBroadcasts: make(map[string]*BroadcastInfo),
		PendingOffers:    make(map[string]string),
		initialized:      true,
	}
	
	fmt.Printf("✅ WebSocketService 초기화 완료:\n")
	fmt.Printf("  - Broadcasters: %v\n", service.Broadcasters != nil)
	fmt.Printf("  - Viewers: %v\n", service.Viewers != nil)
	fmt.Printf("  - ListSubscribers: %v\n", service.ListSubscribers != nil)
	fmt.Printf("  - ActiveBroadcasts: %v\n", service.ActiveBroadcasts != nil)
	fmt.Printf("  - PendingOffers: %v\n", service.PendingOffers != nil)
	
	return service
}

// 초기화 검증
func (wsService *WebSocketService) validateService() error {
	if !wsService.initialized {
		return fmt.Errorf("서비스가 초기화되지 않았습니다")
	}
	if wsService.ListSubscribers == nil {
		wsService.ListSubscribers = make(map[*websocket.Conn]string)
		fmt.Printf("⚠️ ListSubscribers 재초기화\n")
	}
	return nil
}

// 방송자 처리 (안전성 강화)
func (wsService *WebSocketService) HandleBroadcaster(conn *websocket.Conn) {
	if err := wsService.validateService(); err != nil {
		fmt.Printf("❌ 방송자 핸들러 초기화 오류: %v\n", err)
		return
	}

	userID := conn.Query("user_id")
	userName := conn.Query("user_name")
	
	fmt.Printf("🎥 방송자 핸들러 시작: userID=%s, userName=%s\n", userID, userName)
	
	if userID == "" {
		fmt.Println("❌ 방송자 user_id가 없습니다")
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","data":"user_id가 필요합니다"}`))
		return // Close() 제거하고 연결 유지
	}

	connection := &Connection{
		Conn:   conn,
		UserID: userID,
		Name:   userName,
		Role:   "broadcaster",
	}

	wsService.Mutex.Lock()
	wsService.Broadcasters[userID] = connection
	wsService.Mutex.Unlock()

	fmt.Printf("✅ 방송자 연결 등록: %s (%s)\n", userName, userID)

	defer func() {
		fmt.Printf("🧹 방송자 정리 시작: %s\n", userID)
		wsService.stopBroadcast(userID)
		wsService.Mutex.Lock()
		if wsService.Broadcasters != nil {
			delete(wsService.Broadcasters, userID)
		}
		wsService.Mutex.Unlock()
		fmt.Printf("🧹 방송자 정리 완료: %s\n", userID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("❌ 방송자 연결 오류 (%s): %v\n", userID, err)
			break
		}

		fmt.Printf("📨 방송자 메시지 수신 (%s): %s\n", userID, string(message))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("❌ JSON 파싱 오류: %v\n", err)
			continue
		}

		wsService.handleBroadcasterMessage(userID, &msg, message)
	}
}

// 시청자 처리 (안전성 강화)
func (wsService *WebSocketService) HandleViewer(conn *websocket.Conn) {
	if err := wsService.validateService(); err != nil {
		fmt.Printf("❌ 시청자 핸들러 초기화 오류: %v\n", err)
		return
	}

	userID := conn.Query("user_id")
	userName := conn.Query("user_name")
	broadcasterID := conn.Query("broadcaster_id")
	
	fmt.Printf("👀 시청자 핸들러 시작: userID=%s, userName=%s, broadcasterID=%s\n", userID, userName, broadcasterID)
	
	if userID == "" {
		fmt.Println("❌ 시청자 user_id가 없습니다")
		conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","data":"user_id가 필요합니다"}`))
		return
	}

	connection := &Connection{
		Conn:   conn,
		UserID: userID,
		Name:   userName,
		Role:   "viewer",
	}

	viewerInfo := &ViewerInfo{
		Connection:    connection,
		BroadcasterID: broadcasterID,
		JoinTime:      time.Now(),
	}

	wsService.Mutex.Lock()
	wsService.Viewers[userID] = viewerInfo
	wsService.Mutex.Unlock()

	fmt.Printf("✅ 시청자 연결 등록: %s (%s)\n", userName, userID)

	defer func() {
		fmt.Printf("🧹 시청자 정리 시작: %s\n", userID)
		wsService.removeViewer(userID)
		fmt.Printf("🧹 시청자 정리 완료: %s\n", userID)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("❌ 시청자 연결 오류 (%s): %v\n", userID, err)
			break
		}

		fmt.Printf("📨 시청자 메시지 수신 (%s): %s\n", userID, string(message))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("❌ JSON 파싱 오류: %v\n", err)
			continue
		}

		wsService.handleViewerMessage(userID, &msg, message)
	}
}

// 방송 목록 구독자 처리 (완전히 수정)
func (wsService *WebSocketService) HandleViewerList(conn *websocket.Conn) {
	if err := wsService.validateService(); err != nil {
		fmt.Printf("❌ 목록 구독자 핸들러 초기화 오류: %v\n", err)
		return
	}

	userID := conn.Query("user_id")
	fmt.Printf("📋 목록 구독자 핸들러 시작: userID=%s\n", userID)
	
	if userID == "" {
		fmt.Println("❌ 목록 구독자 user_id가 없습니다")
		// 익명 사용자 ID 생성
		userID = fmt.Sprintf("anonymous_%d", time.Now().UnixNano())
		fmt.Printf("⚠️ 익명 사용자 ID 할당: %s\n", userID)
		
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{
			"type": "warning",
			"data": "user_id가 없어 익명 사용자로 설정되었습니다: %s"
		}`, userID)))
	}

	// 안전한 등록
	fmt.Printf("🔐 뮤텍스 잠금 시도 (목록 구독자 등록)...\n")
	wsService.Mutex.Lock()
	
	// 재확인
	if wsService.ListSubscribers == nil {
		wsService.ListSubscribers = make(map[*websocket.Conn]string)
		fmt.Printf("⚠️ ListSubscribers 긴급 재초기화\n")
	}
	
	wsService.ListSubscribers[conn] = userID
	subscriberCount := len(wsService.ListSubscribers)
	wsService.Mutex.Unlock()
	fmt.Printf("🔓 뮤텍스 잠금 해제 완료\n")

	fmt.Printf("✅ 방송 목록 구독자 등록 완료: %s (총 %d명)\n", userID, subscriberCount)

	// 즉시 방송 목록 전송
	fmt.Printf("📤 즉시 방송 목록 전송...\n")
	wsService.sendBroadcastList(conn)

	defer func() {
		fmt.Printf("🧹 목록 구독자 정리 시작: %s\n", userID)
		wsService.Mutex.Lock()
		if wsService.ListSubscribers != nil {
			delete(wsService.ListSubscribers, conn)
		}
		wsService.Mutex.Unlock()
		fmt.Printf("🧹 목록 구독자 정리 완료: %s\n", userID)
	}()

	fmt.Printf("🔄 메시지 수신 루프 시작: %s\n", userID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("❌ 목록 구독자 연결 종료 (%s): %v\n", userID, err)
			break
		}

		fmt.Printf("📨 목록 구독자 메시지 수신 (%s): %s\n", userID, string(message))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Printf("❌ JSON 파싱 오류: %v\n", err)
			continue
		}

		switch msg.Type {
		case "get_broadcast_list":
			fmt.Printf("📋 방송 목록 요청 수신 from %s\n", userID)
			wsService.sendBroadcastList(conn)
		case "ping":
			fmt.Printf("🏓 Ping 수신 from %s\n", userID)
			pongMsg := fmt.Sprintf(`{"type": "pong", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
			conn.WriteMessage(websocket.TextMessage, []byte(pongMsg))
		default:
			fmt.Printf("❓ 알 수 없는 메시지 타입 (%s): %s\n", userID, msg.Type)
		}
	}

	fmt.Printf("🔌 HandleViewerList 종료: %s\n", userID)
}

// 방송자 메시지 처리
func (wsService *WebSocketService) handleBroadcasterMessage(broadcasterID string, msg *Message, rawMessage []byte) {
	fmt.Printf("🔍 방송자 메시지 처리: %s, 타입: %s\n", broadcasterID, msg.Type)
	
	switch msg.Type {
	case "start_broadcast":
		wsService.startBroadcast(broadcasterID, msg)
	case "stop_broadcast":
		wsService.stopBroadcast(broadcasterID)
	case "offer":
		wsService.forwardOfferToViewer(broadcasterID, msg, rawMessage)
	case "candidate":
		wsService.forwardCandidateToViewer(broadcasterID, msg, rawMessage)
	case "offer_request":
		wsService.handleOfferRequest(broadcasterID, msg)
	default:
		fmt.Printf("🔄 알 수 없는 방송자 메시지: %s\n", msg.Type)
	}
}

// 시청자 메시지 처리
func (wsService *WebSocketService) handleViewerMessage(viewerID string, msg *Message, rawMessage []byte) {
	fmt.Printf("🔍 시청자 메시지 처리: %s, 타입: %s\n", viewerID, msg.Type)
	
	switch msg.Type {
	case "request_stream":
		wsService.handleStreamRequest(viewerID, msg)
	case "answer":
		wsService.forwardAnswerToBroadcaster(viewerID, msg, rawMessage)
	case "candidate":
		wsService.forwardCandidateToBroadcaster(viewerID, msg, rawMessage)
	case "viewer_join":
		wsService.handleViewerJoin(viewerID, msg)
	case "viewer_leave":
		wsService.handleViewerLeave(viewerID, msg)
	default:
		fmt.Printf("🔄 알 수 없는 시청자 메시지: %s\n", msg.Type)
	}
}

// 방송 시작
func (wsService *WebSocketService) startBroadcast(broadcasterID string, msg *Message) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	broadcaster := wsService.Broadcasters[broadcasterID]
	if broadcaster == nil {
		fmt.Printf("❌ 방송자를 찾을 수 없음: %s\n", broadcasterID)
		return
	}

	broadcast := &BroadcastInfo{
		BroadcasterID:   broadcasterID,
		BroadcasterName: broadcaster.Name,
		StartTime:       time.Now(),
		ViewerCount:     0,
		IsLive:          true,
	}

	wsService.ActiveBroadcasts[broadcasterID] = broadcast
	fmt.Printf("🔴 방송 시작:\n")
	fmt.Printf("  - ID: %s\n", broadcast.BroadcasterID)
	fmt.Printf("  - 이름: %s\n", broadcast.BroadcasterName)
	fmt.Printf("  - 시작 시간: %s\n", broadcast.StartTime.Format("2006-01-02 15:04:05"))

	wsService.printCurrentState()

	// 모든 목록 구독자에게 알림
	notificationMsg := &Message{
		Type:      "broadcast_started",
		Broadcast: broadcast,
	}

	fmt.Printf("📢 방송 시작 알림을 %d명의 구독자에게 전송\n", len(wsService.ListSubscribers))
	wsService.broadcastToListSubscribers(notificationMsg)
}

// 방송 종료
func (wsService *WebSocketService) stopBroadcast(broadcasterID string) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	broadcast := wsService.ActiveBroadcasts[broadcasterID]
	if broadcast == nil {
		fmt.Printf("⚠️ 종료할 방송이 없음: %s\n", broadcasterID)
		return
	}

	delete(wsService.ActiveBroadcasts, broadcasterID)
	fmt.Printf("⚫ 방송 종료: %s (%s)\n", broadcast.BroadcasterName, broadcasterID)

	// 해당 방송의 모든 시청자에게 방송 종료 알림 전송
	broadcastEndMsg := &Message{
		Type:          "broadcast_ended",
		BroadcasterID: broadcasterID,
		Broadcast:     broadcast,
	}
	
	notifiedViewers := 0
	for _, viewer := range wsService.Viewers {
		if viewer.BroadcasterID == broadcasterID {
			// 방송 종료 알림 전송
			wsService.sendToConnection(viewer.Connection.Conn, broadcastEndMsg)
			notifiedViewers++
		}
	}
	fmt.Printf("📺 방송 %s 종료 알림을 %d명의 시청자에게 전송\n", broadcasterID, notifiedViewers)

	// 잠시 대기 후 연결 해제 (알림 전송 시간 확보)
	go func() {
		time.Sleep(100 * time.Millisecond)
		wsService.Mutex.Lock()
		defer wsService.Mutex.Unlock()
		
		// 해당 방송의 모든 시청자 연결 해제
		for viewerID, viewer := range wsService.Viewers {
			if viewer.BroadcasterID == broadcasterID {
				viewer.Connection.Conn.Close()
				delete(wsService.Viewers, viewerID)
			}
		}
	}()

	// 모든 목록 구독자에게 알림
	wsService.broadcastToListSubscribers(&Message{
		Type:          "broadcast_ended",
		BroadcasterID: broadcasterID,
	})
}

// 스트림 요청 처리
func (wsService *WebSocketService) handleStreamRequest(viewerID string, msg *Message) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	viewer := wsService.Viewers[viewerID]
	if viewer == nil {
		fmt.Printf("❌ 시청자를 찾을 수 없음: %s\n", viewerID)
		return
	}

	broadcasterID := toString(msg.BroadcasterID)
	if broadcasterID == "" {
		broadcasterID = viewer.BroadcasterID
	}

	fmt.Printf("🎯 스트림 요청: 시청자 %s -> 방송자 %s\n", viewerID, broadcasterID)

	broadcaster := wsService.Broadcasters[broadcasterID]
	broadcast := wsService.ActiveBroadcasts[broadcasterID]

	if broadcaster == nil || broadcast == nil {
		fmt.Printf("❌ 방송자 또는 방송을 찾을 수 없음: %s\n", broadcasterID)
		errorMsg := &Message{
			Type: "error",
			Data: "방송을 찾을 수 없습니다.",
		}
		wsService.sendToConnection(viewer.Connection.Conn, errorMsg)
		return
	}

	broadcast.ViewerCount++
	viewer.BroadcasterID = broadcasterID

	offerRequest := &Message{
		Type:       "offer_request",
		ViewerID:   viewerID,
		ViewerName: viewer.Connection.Name,
	}

	wsService.sendToConnection(broadcaster.Conn, offerRequest)
	wsService.PendingOffers[viewerID] = broadcasterID
	wsService.updateViewerCount(broadcasterID, broadcast.ViewerCount)

	fmt.Printf("✅ 스트림 요청 처리 완료: %s -> %s\n", viewerID, broadcasterID)
}

// Offer 전달
func (wsService *WebSocketService) forwardOfferToViewer(broadcasterID string, msg *Message, rawMessage []byte) {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	viewerID := toString(msg.ViewerID)
	if viewerID == "" {
		for vID, bID := range wsService.PendingOffers {
			if bID == broadcasterID {
				viewerID = vID
				delete(wsService.PendingOffers, vID)
				break
			}
		}
	}

	fmt.Printf("📤 Offer 전달 시도: %s -> %s\n", broadcasterID, viewerID)

	if viewerID != "" && wsService.Viewers[viewerID] != nil {
		viewer := wsService.Viewers[viewerID]
		var msgData map[string]interface{}
		json.Unmarshal(rawMessage, &msgData)
		msgData["broadcaster_id"] = broadcasterID

		modifiedMessage, _ := json.Marshal(msgData)
		if err := viewer.Connection.Conn.WriteMessage(websocket.TextMessage, modifiedMessage); err != nil {
			fmt.Printf("❌ Offer 전달 실패: %v\n", err)
		} else {
			fmt.Printf("✅ Offer 전달 성공: %s -> %s\n", broadcasterID, viewerID)
		}
	}
}

// Answer 전달
func (wsService *WebSocketService) forwardAnswerToBroadcaster(viewerID string, msg *Message, rawMessage []byte) {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	broadcasterID := toString(msg.BroadcasterID)
	if broadcasterID == "" && wsService.Viewers[viewerID] != nil {
		broadcasterID = wsService.Viewers[viewerID].BroadcasterID
	}

	fmt.Printf("📤 Answer 전달 시도: %s -> %s\n", viewerID, broadcasterID)

	if broadcasterID != "" && wsService.Broadcasters[broadcasterID] != nil {
		broadcaster := wsService.Broadcasters[broadcasterID]
		var msgData map[string]interface{}
		json.Unmarshal(rawMessage, &msgData)
		msgData["viewer_id"] = viewerID

		modifiedMessage, _ := json.Marshal(msgData)
		if err := broadcaster.Conn.WriteMessage(websocket.TextMessage, modifiedMessage); err != nil {
			fmt.Printf("❌ Answer 전달 실패: %v\n", err)
		} else {
			fmt.Printf("✅ Answer 전달 성공: %s -> %s\n", viewerID, broadcasterID)
		}
	}
}

// ICE Candidate 전달 (시청자에게)
func (wsService *WebSocketService) forwardCandidateToViewer(broadcasterID string, msg *Message, rawMessage []byte) {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	viewerID := toString(msg.ViewerID)
	if viewerID != "" && wsService.Viewers[viewerID] != nil {
		viewer := wsService.Viewers[viewerID]
		var msgData map[string]interface{}
		json.Unmarshal(rawMessage, &msgData)
		msgData["broadcaster_id"] = broadcasterID

		modifiedMessage, _ := json.Marshal(msgData)
		viewer.Connection.Conn.WriteMessage(websocket.TextMessage, modifiedMessage)
	}
}

// ICE Candidate 전달 (방송자에게)
func (wsService *WebSocketService) forwardCandidateToBroadcaster(viewerID string, msg *Message, rawMessage []byte) {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	broadcasterID := toString(msg.BroadcasterID)
	if broadcasterID == "" && wsService.Viewers[viewerID] != nil {
		broadcasterID = wsService.Viewers[viewerID].BroadcasterID
	}

	if broadcasterID != "" && wsService.Broadcasters[broadcasterID] != nil {
		broadcaster := wsService.Broadcasters[broadcasterID]
		var msgData map[string]interface{}
		json.Unmarshal(rawMessage, &msgData)
		msgData["viewer_id"] = viewerID

		modifiedMessage, _ := json.Marshal(msgData)
		broadcaster.Conn.WriteMessage(websocket.TextMessage, modifiedMessage)
	}
}

// 시청자 입장 처리
func (wsService *WebSocketService) handleViewerJoin(viewerID string, msg *Message) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	viewer := wsService.Viewers[viewerID]
	if viewer == nil {
		fmt.Printf("❌ 시청자를 찾을 수 없음: %s\n", viewerID)
		return
	}

	broadcasterID := toString(msg.BroadcasterID)
	if broadcasterID == "" {
		broadcasterID = viewer.BroadcasterID
	}

	broadcast := wsService.ActiveBroadcasts[broadcasterID]
	if broadcast == nil {
		fmt.Printf("❌ 방송을 찾을 수 없음: %s\n", broadcasterID)
		return
	}

	// 시청자 수 증가
	broadcast.ViewerCount++
	viewer.BroadcasterID = broadcasterID
	viewer.JoinTime = time.Now()

	fmt.Printf("👋 시청자 입장: %s (%s) -> 방송 %s, 총 시청자 수: %d명\n", 
		viewer.Connection.Name, viewerID, broadcasterID, broadcast.ViewerCount)

	// 방송자에게 시청자 입장 알림
	if broadcaster := wsService.Broadcasters[broadcasterID]; broadcaster != nil {
		joinMsg := &Message{
			Type:       "viewer_joined",
			ViewerID:   viewerID,
			ViewerName: viewer.Connection.Name,
			Count:      broadcast.ViewerCount,
		}
		wsService.sendToConnection(broadcaster.Conn, joinMsg)
	}

	// 시청자 수 업데이트
	wsService.updateViewerCount(broadcasterID, broadcast.ViewerCount)

	// 시청자에게 입장 확인 메시지
	confirmMsg := &Message{
		Type:          "join_confirmed",
		BroadcasterID: broadcasterID,
		Broadcast:     broadcast,
	}
	wsService.sendToConnection(viewer.Connection.Conn, confirmMsg)
}

// 시청자 퇴장 처리
func (wsService *WebSocketService) handleViewerLeave(viewerID string, msg *Message) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	viewer := wsService.Viewers[viewerID]
	if viewer == nil {
		fmt.Printf("❌ 시청자를 찾을 수 없음: %s\n", viewerID)
		return
	}

	broadcasterID := viewer.BroadcasterID
	if broadcasterID == "" {
		broadcasterID = toString(msg.BroadcasterID)
	}

	broadcast := wsService.ActiveBroadcasts[broadcasterID]
	if broadcast != nil && broadcast.ViewerCount > 0 {
		broadcast.ViewerCount--
		
		fmt.Printf("👋 시청자 퇴장: %s (%s) <- 방송 %s, 총 시청자 수: %d명\n", 
			viewer.Connection.Name, viewerID, broadcasterID, broadcast.ViewerCount)

		// 방송자에게 시청자 퇴장 알림
		if broadcaster := wsService.Broadcasters[broadcasterID]; broadcaster != nil {
			leaveMsg := &Message{
				Type:       "viewer_left",
				ViewerID:   viewerID,
				ViewerName: viewer.Connection.Name,
				Count:      broadcast.ViewerCount,
			}
			wsService.sendToConnection(broadcaster.Conn, leaveMsg)
		}

		// 시청자 수 업데이트
		wsService.updateViewerCount(broadcasterID, broadcast.ViewerCount)
	}

	// 시청자에게 퇴장 확인 메시지
	confirmMsg := &Message{
		Type:          "leave_confirmed",
		BroadcasterID: broadcasterID,
	}
	wsService.sendToConnection(viewer.Connection.Conn, confirmMsg)

	// 시청자 완전 제거 (중복 처리 방지)
	delete(wsService.Viewers, viewerID)
	delete(wsService.PendingOffers, viewerID)
}

// 시청자 수 업데이트
func (wsService *WebSocketService) updateViewerCount(broadcasterID string, count int) {
	fmt.Println("🔄 시청자 수 업데이트:", broadcasterID, "->", count)
	fmt.Println(broadcasterID, count)
	
	// 방송별 시청자 수 상세 로그
	fmt.Printf("📊 방송별 시청자 수 현황:\n")
	for id, broadcast := range wsService.ActiveBroadcasts {
		actualViewers := 0
		for _, viewer := range wsService.Viewers {
			if viewer.BroadcasterID == id {
				actualViewers++
			}
		}
		fmt.Printf("  방송 %s (%s): 기록된 시청자 %d명, 실제 연결된 시청자 %d명\n", 
			id, broadcast.BroadcasterName, broadcast.ViewerCount, actualViewers)
	}
	fmt.Printf("  총 연결된 시청자: %d명\n", len(wsService.Viewers))
	
	// 방송자에게 시청자 수 업데이트 전송
	if wsService.Broadcasters[broadcasterID] != nil {
		msg := &Message{
			Type:  "viewer_count_update",
			Count: count,
		}
		wsService.sendToConnection(wsService.Broadcasters[broadcasterID].Conn, msg)
	}

	// 해당 방송의 모든 시청자에게 시청자 수 업데이트 전송
	viewerUpdateMsg := &Message{
		Type:          "viewer_count_update",
		BroadcasterID: broadcasterID,
		Count:         count,
	}
	
	sentToViewers := 0
	for _, viewer := range wsService.Viewers {
		if viewer.BroadcasterID == broadcasterID {
			wsService.sendToConnection(viewer.Connection.Conn, viewerUpdateMsg)
			sentToViewers++
		}
	}
	fmt.Printf("📺 방송 %s의 시청자 %d명에게 시청자 수 업데이트 전송\n", broadcasterID, sentToViewers)

	// 방송목록 구독자들에게 전송
	wsService.broadcastToListSubscribers(&Message{
		Type:          "viewer_count_update",
		BroadcasterID: broadcasterID,
		Count:         count,
	})
}

// 시청자 제거
func (wsService *WebSocketService) removeViewer(viewerID string) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	viewer := wsService.Viewers[viewerID]
	if viewer != nil {
		broadcasterID := viewer.BroadcasterID
		broadcast := wsService.ActiveBroadcasts[broadcasterID]
		
		if broadcast != nil && broadcast.ViewerCount > 0 {
			broadcast.ViewerCount--
			wsService.updateViewerCount(broadcasterID, broadcast.ViewerCount)
		}

		delete(wsService.Viewers, viewerID)
		delete(wsService.PendingOffers, viewerID)
	}
}

// 현재 상태 출력
func (wsService *WebSocketService) printCurrentState() {
	fmt.Printf("\n📊 현재 서버 상태:\n")
	fmt.Printf("  방송자 수: %d\n", len(wsService.Broadcasters))
	fmt.Printf("  시청자 수: %d\n", len(wsService.Viewers))
	fmt.Printf("  목록 구독자 수: %d\n", len(wsService.ListSubscribers))
	fmt.Printf("  활성 방송 수: %d\n", len(wsService.ActiveBroadcasts))

	fmt.Printf("\n🔴 활성 방송 목록:\n")
	for id, broadcast := range wsService.ActiveBroadcasts {
		fmt.Printf("  - ID: %s, 이름: %s, 시청자: %d명\n", 
			id, broadcast.BroadcasterName, broadcast.ViewerCount)
	}
	fmt.Printf("\n")
}

// 방송 목록 전송
func (wsService *WebSocketService) sendBroadcastList(conn *websocket.Conn) {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	broadcasts := make([]BroadcastInfo, 0, len(wsService.ActiveBroadcasts))
	for _, broadcast := range wsService.ActiveBroadcasts {
		broadcasts = append(broadcasts, *broadcast)
	}

	msg := &Message{
		Type:       "broadcast_list",
		Broadcasts: broadcasts,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ JSON 마샬링 오류: %v\n", err)
		return
	}

	fmt.Printf("📋 방송 목록 전송: %d개 방송\n", len(broadcasts))
	for i, broadcast := range broadcasts {
		fmt.Printf("  [%d] %s (%s) - %d명 시청\n", 
			i+1, broadcast.BroadcasterName, broadcast.BroadcasterID, broadcast.ViewerCount)
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Printf("❌ 방송 목록 전송 실패: %v\n", err)
	} else {
		fmt.Printf("✅ 방송 목록 전송 성공\n")
	}
}

// 구독자들에게 브로드캐스트
func (wsService *WebSocketService) broadcastToListSubscribers(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ 구독자 알림 JSON 마샬링 오류: %v\n", err)
		return
	}

	fmt.Printf("📢 구독자들에게 브로드캐스트: %s\n", msg.Type)
	successCount := 0
	
	for conn, userID := range wsService.ListSubscribers {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Printf("❌ 구독자 %s에게 전송 실패: %v\n", userID, err)
		} else {
			successCount++
		}
	}
	fmt.Printf("📊 총 %d명 중 %d명에게 성공적으로 전송\n", len(wsService.ListSubscribers), successCount)
}

// 연결에 메시지 전송
func (wsService *WebSocketService) sendToConnection(conn *websocket.Conn, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("❌ JSON 마샬링 오류: %v\n", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Printf("❌ 메시지 전송 실패: %v\n", err)
	}
}

// Offer 요청 처리
func (wsService *WebSocketService) handleOfferRequest(broadcasterID string, msg *Message) {
	fmt.Printf("🔔 Offer 요청 처리: %s\n", broadcasterID)
}

// API용 활성 방송 목록 반환
func (wsService *WebSocketService) GetActiveBroadcasts() []BroadcastInfo {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	broadcasts := make([]BroadcastInfo, 0, len(wsService.ActiveBroadcasts))
	for _, broadcast := range wsService.ActiveBroadcasts {
		broadcasts = append(broadcasts, *broadcast)
	}
	return broadcasts
}

// API용 방송 통계 반환
func (wsService *WebSocketService) GetBroadcastStats(broadcasterID string) map[string]interface{} {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	broadcast := wsService.ActiveBroadcasts[broadcasterID]
	if broadcast == nil {
		return map[string]interface{}{
			"error": "방송을 찾을 수 없습니다",
		}
	}

	// 해당 방송의 시청자 목록
	viewers := make([]map[string]interface{}, 0)
	for _, viewer := range wsService.Viewers {
		if viewer.BroadcasterID == broadcasterID {
			viewers = append(viewers, map[string]interface{}{
				"viewer_id":   viewer.Connection.UserID,
				"viewer_name": viewer.Connection.Name,
				"join_time":   viewer.JoinTime,
			})
		}
	}

	return map[string]interface{}{
		"broadcast_info": broadcast,
		"viewers":        viewers,
		"viewer_count":   len(viewers),
	}
}

// 서버 상태 반환 (API용)
func (wsService *WebSocketService) GetServerStatus() map[string]interface{} {
	wsService.Mutex.RLock()
	defer wsService.Mutex.RUnlock()

	// 방송별 시청자 수 계산
	broadcasterViewers := make(map[string]int)
	for _, viewer := range wsService.Viewers {
		if viewer.BroadcasterID != "" {
			broadcasterViewers[viewer.BroadcasterID]++
		}
	}

	// 총 시청자 수 계산
	totalActiveViewers := 0
	for _, count := range broadcasterViewers {
		totalActiveViewers += count
	}

	averageViewers := 0.0
	if len(wsService.ActiveBroadcasts) > 0 {
		averageViewers = float64(totalActiveViewers) / float64(len(wsService.ActiveBroadcasts))
	}

	return map[string]interface{}{
		"server_info": map[string]interface{}{
			"uptime":  time.Now().Format("2006-01-02 15:04:05"),
			"version": "1.0.0",
		},
		"connections": map[string]interface{}{
			"broadcasters":     len(wsService.Broadcasters),
			"viewers":          len(wsService.Viewers),
			"list_subscribers": len(wsService.ListSubscribers),
			"pending_offers":   len(wsService.PendingOffers),
		},
		"broadcasts": map[string]interface{}{
			"active_count":        len(wsService.ActiveBroadcasts),
			"total_viewers":       totalActiveViewers,
			"average_viewers":     averageViewers,
			"broadcaster_viewers": broadcasterViewers,
		},
		"active_broadcasts": wsService.GetActiveBroadcasts(),
	}
}