package pkg

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

// init は、パッケージの初期状態を構成します。
func init() {
	zerolog.TimeFieldFormat = "15:04:05"
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.SetFlags(0)
}

// Logging は、コンテキストのリクエストIDと指定値をログへ出力します。
func Logging(
	ctx context.Context,
	data interface{},
) {
	// テスト中であればロギングしない
	if flag.Lookup("test.v") != nil {
		log.Println("run under go test")
		return
	}

	_, fullPath, line, _ := runtime.Caller(1)
	trimPath := fullPath
	if idx := strings.Index(fullPath, "internal/"); idx != -1 {
		trimPath = fullPath[idx:]
	}

	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("file", fmt.Sprintf("%s:%d", trimPath, line))

	event := logger.Logger()

	switch v := data.(type) {
	case error:
		event.Error().Msg(v.Error())
	default:
		event.Info().Interface("data", data).Msg("")
	}
}
