package main

import "github.com/Tnze/go-mc/chat"

type EmptyPingInfo struct {
	description chat.Message
}

func NewPingInfo() *EmptyPingInfo {
	text := chat.Text("")
	s := EmptyPingInfo{text}

	return &s
}

func (ep *EmptyPingInfo) Name() string               { return "" }
func (ep *EmptyPingInfo) FavIcon() string            { return "" }
func (ep *EmptyPingInfo) Description() *chat.Message { return &ep.description }
