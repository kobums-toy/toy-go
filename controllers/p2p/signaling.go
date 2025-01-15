package p2p

import (
	"encoding/json"
	"fmt"
	"toysgo/services"

	"github.com/gofiber/websocket/v2"
	"github.com/pion/webrtc/v3"
)

// Create a WebRTC service instance
var webrtcService = services.NewWebRTCService()

// Message defines the structure of WebSocket messages
type Message struct {
	Type string `json:"type"` // "offer", "answer", or "candidate"
	Data string `json:"data"` // SDP or ICE candidate
}

func WebSocketHandler(c *websocket.Conn) {
	defer c.Close()

	peerID := c.Query("peer_id") // Peer ID from client
	_, err := webrtcService.CreatePeerConnection(peerID)
	if err != nil {
		fmt.Println("Failed to create PeerConnection:", err)
		return
	}
	defer webrtcService.ClosePeerConnection(peerID)

	for {
		// Read message from client
		_, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket error:", err)
			break
		}

		// Parse the message
		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Invalid message format:", err)
			continue
		}

		// Handle WebRTC signaling messages
		switch msg.Type {
		case "offer":
			fmt.Printf("Received SDP offer:\n%s\n", msg.Data)
			err := webrtcService.HandleSDP(peerID, webrtc.SDPTypeOffer, msg.Data)
			if err != nil {
				fmt.Println("Error handling offer:", err)
			}
		case "answer":
			fmt.Println("Received SDP answer")
			err := webrtcService.HandleSDP(peerID, webrtc.SDPTypeAnswer, msg.Data)
			if err != nil {
				fmt.Println("Error handling answer:", err)
			}
		case "candidate":
			fmt.Println("Received ICE candidate")
			err := webrtcService.HandleCandidate(peerID, msg.Data)
			if err != nil {
				fmt.Println("Error handling candidate:", err)
			}
		default:
			fmt.Println("Unknown message type:", msg.Type)
		}
	}
}
