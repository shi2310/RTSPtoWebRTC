package main

import (
	"github.com/deepch/vdk/av"
)

// 流结构体
type StreamST struct {
	URL    string `json:"url"`
	Codecs []av.CodecData
	Cl     viwer //包通道
}

type viwer struct {
	c chan av.Packet
}

// 投放数据到某流的通道
func (element *StreamST) cast(pck av.Packet) {
	if len(element.Cl.c) < cap(element.Cl.c) {
		element.Cl.c <- pck
	}
}

// 设置某流的编译码数据
func (element *StreamST) coAd(codecs []av.CodecData) {
	element.Codecs = codecs
}

// 获取某流的编译码数据
func (element *StreamST) coGe() []av.CodecData {
	return element.Codecs
}

// 某流声明包通道
func (element *StreamST) clAd() chan av.Packet {
	ch := make(chan av.Packet, 100)
	element.Cl = viwer{c: ch}
	return ch
}
