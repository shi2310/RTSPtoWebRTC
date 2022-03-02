package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/deepch/vdk/codec/h264parser"
	"github.com/deepch/vdk/format/rtsp"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

func webrtcAnswer(offerSdp string, stream *StreamST, session *rtsp.Client) (string, error) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	videoTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", pseudoUUID()+"_video")
	if err != nil {
		return "", nil
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		return "", nil
	}
	sdp, err := base64.StdEncoding.DecodeString(offerSdp)
	if err != nil {
		return "", nil
	}
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  string(sdp),
	}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		return "", nil
	}
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		return "", nil
	}
	peerConnection.SetLocalDescription(answer)

	// 获取代码数据
	codecs := stream.coGe()
	if codecs == nil {
		return "", errors.New("没有获取到解码数据")
	}
	//sps, pps := []byte{}, []byte{}
	sps := codecs[0].(h264parser.CodecData).SPS()
	pps := codecs[0].(h264parser.CodecData).PPS()

	go func() {
		close := make(chan bool, 1)
		conected := make(chan bool, 1)
		defer peerConnection.Close()
		defer session.Close()

		peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			log.Println("Connection State has changed:", connectionState.String())
			if connectionState != webrtc.ICEConnectionStateConnected {
				log.Println("Client Close Exit")
				close <- true
				return
			}
			if connectionState == webrtc.ICEConnectionStateConnected {
				conected <- true
			}
		})
		// 等待连接
		<-conected
		// webrtc连接成功以后建立缓冲通道
		ch := stream.clAd()
		var pre uint32
		var start bool
		for {
			select {
			case <-close:
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					start = true
				}
				if !start {
					continue
				}
				if pck.IsKeyFrame {
					pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)
				} else {
					pck.Data = pck.Data[4:]
				}
				var ts uint32
				if pre != 0 {
					ts = uint32(timeToTs(pck.Time)) - pre
				}
				err := videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: uint32(ts)})
				pre = uint32(timeToTs(pck.Time))
				if err != nil {
					return
				}
			}
		}
	}()

	return base64.StdEncoding.EncodeToString([]byte(answer.SDP)), nil
}

func timeToTs(tm time.Duration) int64 {
	return int64(tm * time.Duration(90000) / time.Second)
}

func pseudoUUID() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
