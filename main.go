package main

import (
	"flag"
	"github.com/Tnze/go-mc/registry"
	"github.com/Tnze/go-mc/server"
	"go.uber.org/zap"
	"runtime/debug"
)

var isDebug = flag.Bool("debug", true, "Enable debug log output")

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
			OnlineMode:           true,
			EnforceSecureProfile: true,
			Threshold:            0, // compress all packets
			LoginChecker:         emptyList,
		},
		ConfigHandler: &server.Configurations{Registries: registry.NewNetworkCodec()},
		GamePlay:      &NoGamePlay{},
	}

	logger.Info("Server listening on :25565")
	err := s.Listen(":25565")
	if err != nil {
		logger.Fatal("Server listening error", zap.Error(err))
	}
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
