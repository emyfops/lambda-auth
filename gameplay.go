package main

import (
	"github.com/Tnze/go-mc/net"
	"github.com/Tnze/go-mc/yggdrasil/user"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
)

var successfulLogins = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "lambda_authentication_processed_logins_total",
	Help: "The total number of successful connections",
}, []string{"protocol"})

type NoGamePlay struct{}

// AcceptPlayer Disconnect the player and stop the player handling
func (n *NoGamePlay) AcceptPlayer(_ string, _ uuid.UUID, _ *user.PublicKey, _ []user.Property, protocol int32, _ *net.Conn) {
	successfulLogins.With(prometheus.Labels{"protocol": strconv.Itoa(int(protocol))}).Inc()
}
