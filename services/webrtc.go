package services

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v3"
)

// WebRTCService manages WebRTC connections and data channels
type WebRTCService struct {
	mu       sync.Mutex
	peers    map[string]*webrtc.PeerConnection
	onTrack  func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)
	onSignal func(signalType string, data string)
}

// NewWebRTCService initializes a new WebRTCService
func NewWebRTCService() *WebRTCService {
	return &WebRTCService{
		peers: make(map[string]*webrtc.PeerConnection),
	}
}

// CreatePeerConnection creates a new WebRTC PeerConnection
func (w *WebRTCService) CreatePeerConnection(peerID string) (*webrtc.PeerConnection, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	// 오디오 트랜시버 추가 (opus 코덱 지원)
	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add audio transceiver: %v", err)
	}

	// 비디오 트랜시버 추가 (VP8 코덱 지원)
	_, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add video transceiver: %v", err)
	}

	// PeerConnection 저장
	w.mu.Lock()
	w.peers[peerID] = peerConnection
	w.mu.Unlock()

	return peerConnection, nil
}

// HandleSDP handles SDP offer/answer
func (w *WebRTCService) HandleSDP(peerID string, sdpType webrtc.SDPType, sdp string) error {
	w.mu.Lock()
	pc, exists := w.peers[peerID]
	w.mu.Unlock()
	if !exists {
		return fmt.Errorf("peer not found")
	}

	offer := webrtc.SessionDescription{
		Type: sdpType,
		SDP:  sdp,
	}

	// SDP 설정
	err := pc.SetRemoteDescription(offer)
	if err != nil {
		return fmt.Errorf("failed to set remote description: %v", err)
	}

	// SDP Answer 생성
	if sdpType == webrtc.SDPTypeOffer {
		answer, err := pc.CreateAnswer(nil)
		if err != nil {
			return fmt.Errorf("failed to create answer: %v", err)
		}
		err = pc.SetLocalDescription(answer)
		if err != nil {
			return fmt.Errorf("failed to set local description: %v", err)
		}
		fmt.Println("Generated SDP Answer:", answer.SDP)
	}

	return nil
}

// HandleCandidate processes an ICE candidate
func (w *WebRTCService) HandleCandidate(peerID string, candidate string) error {
	w.mu.Lock()
	pc, exists := w.peers[peerID]
	w.mu.Unlock()
	if !exists {
		return fmt.Errorf("peer not found")
	}

	iceCandidate := webrtc.ICECandidateInit{Candidate: candidate}
	err := pc.AddICECandidate(iceCandidate)
	if err != nil {
		return fmt.Errorf("failed to add ICE candidate: %v", err)
	}

	return nil
}

// ClosePeerConnection closes and removes a PeerConnection
func (w *WebRTCService) ClosePeerConnection(peerID string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if pc, exists := w.peers[peerID]; exists {
		pc.Close()
		delete(w.peers, peerID)
	}
}
