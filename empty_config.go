package main

import "github.com/Tnze/go-mc/net"

type EmptyConfig struct{}

func (c *EmptyConfig) AcceptConfig(*net.Conn) error { return nil }
