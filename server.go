package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Tnze/go-mc/registry"
	"github.com/Tnze/go-mc/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	isOnline = flag.Bool("online", true, "Enable online-mode")
	isDebug  = flag.Bool("debug", true, "Enable debug log output")
	port     = flag.Int("port", 25565, "Server port")
	promPort = flag.Int("prom_port", 11900, "Prometheus port")
)

func main() {
	flag.Parse()

	var logger *zap.Logger
	if *isDebug {
		logger = unwrap(zap.NewDevelopment())
	} else {
		logger = unwrap(zap.NewProduction())
	}
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}(logger)

	printBuildInfo(logger)
	defer logger.Info("Stopping the server")

	emptyList := NewPlayerList()
	emptyPing := NewPingInfo()

	s := server.Server{
		Logger: zap.NewStdLog(logger),
		ListPingHandler: struct {
			*EmptyPlayerList
			*EmptyPingInfo
		}{emptyList, emptyPing},
		LoginHandler: &server.MojangLoginHandler{
			OnlineMode:           *isOnline,
			EnforceSecureProfile: *isOnline,
			Threshold:            0, // compress all packets
			LoginChecker:         emptyList,
		},
		ConfigHandler: &server.Configurations{Registries: registry.NewNetworkCodec()},
		GamePlay:      &NoGamePlay{},
	}

	go startPrometheus(logger)

	logger.Info("Server started", zap.Int("port", *port))
	err := s.Listen(fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("Server listening error", zap.Error(err))
	}
}

func startPrometheus(logger *zap.Logger) {
	logger.Info("Starting prometheus metrics server", zap.Int("port", *promPort))

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%d", *promPort), nil)
}

func printBuildInfo(logger *zap.Logger) {
	binaryInfo, _ := debug.ReadBuildInfo()
	settings := make(map[string]string)
	for _, v := range binaryInfo.Settings {
		settings[v.Key] = v.Value
	}
	logger.Debug("Build info", zap.Any("settings", settings))
}

func unwrap[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
