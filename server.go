package main

import (
	"flag"
	"github.com/Tnze/go-mc/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

var isDebug = flag.Bool("debug", true, "Enable debug log output")
var isOnline = flag.Bool("online", true, "Enable online-mode")

func main() {
	flag.Parse()
	// initialize log library
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
		ConfigHandler: &EmptyConfig{},
		GamePlay:      &NoGamePlay{},
	}

	go startPrometheus(logger)

	logger.Info("Server listening on :25565")
	err := s.Listen(":25565")
	if err != nil {
		logger.Fatal("Server listening error", zap.Error(err))
	}
}

func startPrometheus(logger *zap.Logger) {
	logger.Info("Starting prometheus metrics server on :9100")

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}

// printBuildInfo reading compile information of the binary program with runtime/debug package，and print it to log
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
