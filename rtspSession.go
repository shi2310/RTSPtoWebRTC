package main

import (
	"log"
	"time"

	"github.com/deepch/vdk/format/rtsp"
)

// 建立rtsp会话
func rtspSession(stream *StreamST) (*rtsp.Client, error) {
	rtsp.DebugRtsp = true
	session, err := rtsp.Dial(stream.URL)
	if err != nil {
		log.Println("rtsp连接异常", err)
		return nil, err
	}
	session.RtpKeepAliveTimeout = time.Duration(10 * time.Second)
	codec, err := session.Streams()
	if err != nil {
		log.Println("获取代码数据异常", err)
		return nil, err
	}
	stream.coAd(codec)

	// 开辟协程从会话中循环读取packet投入到缓冲通道
	go func() {
		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				log.Println("ReadPacket error", err)
				break
			}
			stream.cast(pkt)
		}
	}()

	return session, nil
}
