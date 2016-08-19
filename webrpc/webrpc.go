package lambda

import (
	"encoding/json"
	"fmt"
	"github.com/keroserene/go-webrtc"
)

var pc *webrtc.PeerConnection
var dc *webrtc.DataChannel
var Mode ModeType
var err error
var OnOpen func()
var OnData func(msg []byte)

// Janky state machine.
type ModeType int

const (
	ModeInit ModeType = iota
	ModeConnect
	ModeChat
)

func SetLoggingVerbosity(level int) {
	webrtc.SetLoggingVerbosity(level)
}

func generateOffer() {
	fmt.Println("Generating offer...")
	offer, err := pc.CreateOffer() // blocking
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(offer)
}

func generateAnswer() {
	fmt.Println("Generating answer...")
	answer, err := pc.CreateAnswer() // blocking
	if err != nil {
		fmt.Println(err)
		return
	}
	pc.SetLocalDescription(answer)
}

func receiveDescription(sdp *webrtc.SessionDescription) {
	err = pc.SetRemoteDescription(sdp)
	if nil != err {
		fmt.Println("ERROR", err)
		return
	}
	fmt.Println("SDP " + sdp.Type + " successfully received.")
	if "offer" == sdp.Type {
		go generateAnswer()
	}
}

// Manual "copy-paste" signaling channel.
func signalSend(msg string) {
	fmt.Println("\n ---- Please copy below to peer ---- \n")
	fmt.Println(msg + "\n")
}

func SignalReceive(msg string) {
	var parsed map[string]interface{}
	err = json.Unmarshal([]byte(msg), &parsed)
	if nil != err {
		// fmt.Println(err, ", try again.")
		return
	}

	// If this is a valid signal and no PeerConnection has been instantiated,
	// start as the "answerer."
	if nil == pc {
		Start(false)
	}

	if nil != parsed["sdp"] {
		sdp := webrtc.DeserializeSessionDescription(msg)
		if nil == sdp {
			fmt.Println("Invalid SDP.")
			return
		}
		receiveDescription(sdp)
	}

	// Allow individual ICE candidate messages, but this won't be necessary if
	// the remote peer also doesn't use trickle ICE.
	if nil != parsed["candidate"] {
		ice := webrtc.DeserializeIceCandidate(msg)
		if nil == ice {
			fmt.Println("Invalid ICE candidate.")
			return
		}
		pc.AddIceCandidate(*ice)
		fmt.Println("ICE candidate successfully received.")
	}
}

// Attach callbacks to a newly created data channel.
// In this demo, only one data channel is expected, and is only used for chat.
// But it is possible to send any sort of bytes over a data channel, for many
// more interesting purposes.
func prepareDataChannel(channel *webrtc.DataChannel) {
	channel.OnOpen = func() {
		fmt.Println("Data Channel established")
		startChat()
	}
	channel.OnClose = func() {
		fmt.Println("Data Channel closed")
		endChat()
	}
	channel.OnMessage = OnData /* func(msg []byte) {
		receiveChat(string(msg))
	}*/
}

func startChat() {
	Mode = ModeChat
	if OnOpen != nil {
		OnOpen()
	}
}

func endChat() {
	Mode = ModeInit
	fmt.Println("------- chat disabled -------")
}

func SendData(msg string) {
	dc.Send([]byte(msg))
}

//func receiveChat(msg string) {
//	fmt.Println("\n" + string(msg))
//}

// Create a PeerConnection.
// If |instigator| is true, create local data channel which causes a
// negotiation-needed, leading to preparing an SDP offer to be sent to the
// remote peer. Otherwise, await an SDP offer from the remote peer, and send an
// answer back.
func Start(instigator bool) {
	Mode = ModeConnect
	fmt.Println("Starting up PeerConnection...")
	// TODO: Try with TURN servers.
	config := webrtc.NewConfiguration(
		webrtc.OptionIceServer("stun:stun.l.google.com:19302"))

	pc, err = webrtc.NewPeerConnection(config)
	if nil != err {
		fmt.Println("Failed to create PeerConnection.")
		return
	}

	// OnNegotiationNeeded is triggered when something important has occurred in
	// the state of PeerConnection (such as creating a new data channel), in which
	// case a new SDP offer must be prepared and sent to the remote peer.
	pc.OnNegotiationNeeded = func() {
		go generateOffer()
	}
	// Once all ICE candidates are prepared, they need to be sent to the remote
	// peer which will attempt reaching the local peer through NATs.
	pc.OnIceComplete = func() {
		fmt.Println("Finished gathering ICE candidates.")
		sdp := pc.LocalDescription().Serialize()
		signalSend(sdp)
	}
	/*
		pc.OnIceGatheringStateChange = func(state webrtc.IceGatheringState) {
			fmt.Println("Ice Gathering State:", webrtc.IceGatheringStateString[state])
			if webrtc.IceGatheringStateComplete == state {
				// send local description.
			}
		}
	*/
	// A DataChannel is generated through this callback only when the remote peer
	// has initiated the creation of the data channel.
	pc.OnDataChannel = func(channel *webrtc.DataChannel) {
		fmt.Println("Datachannel established by remote... ", channel.Label())
		dc = channel
		prepareDataChannel(channel)
	}

	if instigator {
		// Attempting to create the first datachannel triggers ICE.
		fmt.Println("Initializing datachannel....")
		dc, err = pc.CreateDataChannel("test", webrtc.Init{})
		if nil != err {
			fmt.Println("Unexpected failure creating Channel.")
			return
		}
		prepareDataChannel(dc)
	}
}
