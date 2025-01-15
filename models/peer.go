package models

type Peer struct {
	ID           string `json:"id"`
	SDP          string `json:"sdp"`
	ICECandidate string `json:"ice_candidate"`
}
