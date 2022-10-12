package p2p

import "github.com/pion/webrtc/v3"

var WebRTCConfig = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
		{
			URLs: []string{"stun:129.211.18.139:31478"},
		},
		{
			URLs:       []string{"turn:129.211.18.139:31478"},
			Username:   "szpnygo",
			Credential: "szpnygo",
		},
	},
}
