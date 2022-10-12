package p2p

import (
	"encoding/json"

	"github.com/pion/webrtc/v3"
	"github.com/szpnygo/gtc/log"
)

type P2PClient struct {
	peerConnection *webrtc.PeerConnection
	dataChannel    *webrtc.DataChannel
	isClosed       bool
	user_id        int64

	onClose     func(int64)
	onCandidate func(int64, string)
	onMessage   func([]byte)
}

func NewP2PClient(id int64) *P2PClient {
	return &P2PClient{
		user_id: id,
	}
}

func (p *P2PClient) OnClose(f func(int64)) {
	p.onClose = f
}

func (p *P2PClient) OnCandidate(f func(int64, string)) {
	p.onCandidate = f
}

func (p *P2PClient) OnMessage(f func([]byte)) {
	p.onMessage = f
}

func (p *P2PClient) Close() {
	if p.isClosed {
		return
	}
	if err := p.peerConnection.Close(); err == nil {
		p.isClosed = true
	}
}

// Create create webrtc peerconnection
func (p *P2PClient) Create() error {
	peerConnection, err := webrtc.NewPeerConnection(WebRTCConfig)
	if err != nil {
		log.GTCLog.Error(err)
		return err
	}

	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		p.dataChannel = dc
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			p.onMessage(msg.Data)
		})
	})

	peerConnection.OnConnectionStateChange(func(pcs webrtc.PeerConnectionState) {
		if pcs == webrtc.PeerConnectionStateFailed {
			p.onClose(p.user_id)
		}
	})

	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		p.onCandidate(p.user_id, i.ToJSON().Candidate)
	})

	p.peerConnection = peerConnection

	return nil
}

func (p *P2PClient) CreateOffer() ([]byte, error) {
	dataChannel, err := p.peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}
	p.dataChannel = dataChannel
	p.dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		p.onMessage(msg.Data)
	})

	offer, err := p.peerConnection.CreateOffer(nil)
	if err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}
	if err = p.peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}
	payload, err := json.Marshal(offer)
	if err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}

	return payload, nil
}

func (p *P2PClient) CreateAnswer(data []byte) ([]byte, error) {
	sdp := webrtc.SessionDescription{}
	if err := json.Unmarshal(data, &sdp); err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}

	if err := p.peerConnection.SetRemoteDescription(sdp); err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}

	answer, err := p.peerConnection.CreateAnswer(nil)
	if err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}
	if err = p.peerConnection.SetLocalDescription(answer); err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}

	payload, err := json.Marshal(answer)
	if err != nil {
		log.GTCLog.Error(err)
		return nil, err
	}

	return payload, nil
}

func (p *P2PClient) Answer(data []byte) error {
	sdp := webrtc.SessionDescription{}
	if err := json.Unmarshal(data, &sdp); err != nil {
		log.GTCLog.Error(err)
		return err
	}

	if err := p.peerConnection.SetRemoteDescription(sdp); err != nil {
		log.GTCLog.Error(err)
		return err
	}

	return nil
}

func (p *P2PClient) OnICECandidate(candidate []byte) error {
	if err := p.peerConnection.AddICECandidate(webrtc.ICECandidateInit{Candidate: string(candidate)}); err != nil {
		log.GTCLog.Error(err)
		return err
	}

	return nil
}

func (p *P2PClient) SendMessage(msg []byte) {
	_ = p.dataChannel.Send(msg)
}
