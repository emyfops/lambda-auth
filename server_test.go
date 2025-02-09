package main

import (
	"context"
	"flag"
	"github.com/Tnze/go-mc/bot"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Set("online", "false")
	go main()

	time.Sleep(2 * time.Second)

	os.Exit(m.Run())
}

func TestEmptyPingInfo(t *testing.T) {
	resp, _, err := bot.PingAndListTimeout(":25565", time.Second)
	if err != nil {
		t.Fatal(err)
	}

	if len(resp) != 109 {
		t.Fatalf("expect resp len %d, got %d", 109, len(resp))
	}
}

func TestServerConnection(t *testing.T) {
	c := bot.NewClient()
	err := c.JoinServer(":25565")
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	go func() {
		c.HandleGame()
		cancel()
	}()

	select {
	case <-ctx.Done():
		log.Fatal("the server did not kick the bot in time", ctx.Err())
	default:
	}
}
