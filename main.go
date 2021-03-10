package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"

	"gstreamer-transform/gst"
	"gstreamer-transform/signal"
)

func main() {

	// Everything below is the pion-WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	// Create a audio track
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion1")
	if err != nil {
		panic(err)
	}
	_, err = peerConnection.AddTrack(audioTrack)
	if err != nil {
		panic(err)
	}

	// Create a video track that we send video back to browser on
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion")
	if err != nil {
		panic(err)
	}

	// Add this newly created track to the PeerConnection
	rtpSender, err := peerConnection.AddTrack(videoTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are retuned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// Set a handler for when a new remote track starts, this handler creates a gstreamer pipeline
	// for the given codec
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				rtcpSendErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}})
				if rtcpSendErr != nil {
					fmt.Println(rtcpSendErr)
				}
			}
		}()

		if track.Kind().String() == "audio" {
			codecName := strings.Split(track.Codec().RTPCodecCapability.MimeType, "/")[1]
			pipeline := gst.CreatePipeline(codecName, []*webrtc.TrackLocalStaticSample{audioTrack})
			pipeline.Start()

			buf := make([]byte, 1400)
			for {
				i, _, readErr := track.Read(buf)
				if readErr != nil {
					panic(err)
				}

				pipeline.Push(buf[:i])
			}
		} else {
			codecName := strings.Split(track.Codec().RTPCodecCapability.MimeType, "/")[1]
			pipelineVideo := gst.CreatePipeline(strings.ToLower(codecName), []*webrtc.TrackLocalStaticSample{videoTrack})
			pipelineVideo.Start()

			buf := make([]byte, 1400)
			for {
				i, _, readErr := track.Read(buf)
				if readErr != nil {
					panic(err)
				}

				pipelineVideo.Push(buf[:i])
			}
		}

	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	signal.Decode(signal.MustReadStdin(), &offer)

	// Set the remote SessionDescription
	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	err = peerConnection.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))

	// Block forever
	select {}
}
