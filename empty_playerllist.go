package main

import (
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/server"
	"github.com/google/uuid"
)

// EmptyPlayerList Implements part of [server.ListPingHandler] with "empty" values and implements [server.LoginChecker]
type EmptyPlayerList struct{}

func NewPlayerList() *EmptyPlayerList {
	return &EmptyPlayerList{}
}

// Protocol Returns the input client protocol to allow all clients to connect
func (pl *EmptyPlayerList) Protocol(client int32) int {
	return int(client)
}

func (pl *EmptyPlayerList) MaxPlayer() int    { return 0 }
func (pl *EmptyPlayerList) OnlinePlayer() int { return 0 }

func (pl *EmptyPlayerList) PlayerSamples() []server.PlayerSample       { return nil }
func (pl *EmptyPlayerList) OnlinePlayerSamples() []server.PlayerSample { return nil }

func (pl *EmptyPlayerList) CheckPlayer(string, uuid.UUID, int32) (bool, chat.Message) {
	return true, chat.Message{}
}
