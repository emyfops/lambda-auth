package main

import (
	"github.com/Tnze/go-mc/chat"
	"github.com/Tnze/go-mc/data/packetid"
	"github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/Tnze/go-mc/yggdrasil/user"
	"github.com/google/uuid"
)

type NoGamePlay struct{}

// AcceptPlayer Disconnect the player and stop the player handling
func (n *NoGamePlay) AcceptPlayer(_ string, _ uuid.UUID, _ *user.PublicKey, _ []user.Property, _ int32, conn *net.Conn) {
	conn.WritePacket(pk.Marshal(
		packetid.ClientboundDisconnect, chat.Text("login complete")))
}
