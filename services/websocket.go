package services

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WebSocketService struct {
	Broadcaster *websocket.Conn
	Viewers     []*websocket.Conn
	Mutex       sync.Mutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		Viewers: make([]*websocket.Conn, 0),
	}
}

func (wsService *WebSocketService) SetBroadcaster(conn *websocket.Conn) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	wsService.Broadcaster = conn
	fmt.Println("송출자가 연결되었습니다.")

	defer func() {
		wsService.Mutex.Lock()
		wsService.Broadcaster = nil
		wsService.Mutex.Unlock()
		fmt.Println("송출자가 연결 종료되었습니다.")
		conn.Close()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("송출자 연결 종료:", err)
			break
		}

		fmt.Printf("송출자로부터 메시지 수신: %s\n", string(message))
		wsService.BroadcastToViewers(messageType, message)

		var parsedMessage map[string]interface{}
		err = json.Unmarshal(message, &parsedMessage)
		if err != nil {
			fmt.Println("JSON 디코딩 실패:", err)
			fmt.Println("수신한 원본 메시지:", string(message))
			continue
		}
		fmt.Printf("디코딩된 메시지: %+v\n", parsedMessage)
	}
}

func (wsService *WebSocketService) AddViewer(conn *websocket.Conn) {
	wsService.Mutex.Lock()
	wsService.Viewers = append(wsService.Viewers, conn)
	wsService.Mutex.Unlock()

	fmt.Println("시청자가 연결되었습니다.")

	defer conn.Close()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("시청자 연결 종료:", err)
			break
		}

		fmt.Printf("시청자로부터 메시지 수신: %s\n", string(message))

		var msg map[string]interface{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("JSON 파싱 실패:", err)
			continue
		}

		switch msg["type"] {
		case "request_stream":
			fmt.Println("시청자가 스트림을 요청했습니다.")
			wsService.Mutex.Lock()
			if wsService.Broadcaster != nil {
				err := wsService.Broadcaster.WriteMessage(messageType, []byte(`{"type":"offer_request"}`))
				if err != nil {
					fmt.Println("송출자에게 Offer 요청 실패:", err)
				} else {
					fmt.Println("송출자에게 Offer 요청 성공")
				}
			} else {
				fmt.Println("현재 연결된 송출자가 없습니다.")
			}
			wsService.Mutex.Unlock()

		case "candidate":
			fmt.Printf("시청자로부터 Candidate 메시지 수신: %s\n", string(message))
			wsService.SendToBroadcaster(messageType, message)

		default:
			wsService.SendToBroadcaster(messageType, message)
		}
	}

	wsService.RemoveViewer(conn)
}

func (wsService *WebSocketService) BroadcastToViewers(messageType int, message []byte) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	for _, viewer := range wsService.Viewers {
		if err := viewer.WriteMessage(messageType, message); err != nil {
			fmt.Println("시청자에게 메시지 전달 실패:", err)
		} else {
			fmt.Println("시청자에게 메시지 전달 성공:", string(message))
		}
	}
}

func (wsService *WebSocketService) SendToBroadcaster(messageType int, message []byte) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	if wsService.Broadcaster != nil {
		if err := wsService.Broadcaster.WriteMessage(messageType, message); err != nil {
			fmt.Println("송출자에게 메시지 전달 실패:", err)
		} else {
			fmt.Println("송출자에게 메시지 전달 성공")
		}
	} else {
		fmt.Println("현재 연결된 송출자가 없습니다.")
	}
}

func (wsService *WebSocketService) RemoveViewer(conn *websocket.Conn) {
	wsService.Mutex.Lock()
	defer wsService.Mutex.Unlock()

	for i, viewer := range wsService.Viewers {
		if viewer == conn {
			wsService.Viewers = append(wsService.Viewers[:i], wsService.Viewers[i+1:]...)
			break
		}
	}
}
