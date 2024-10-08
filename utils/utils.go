package utils

import (
	"log/slog"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func InitLogger(){
	opts := slog.HandlerOptions{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func FatalError(err error, msg string, keysAndValues ...interface{}) {
	
	if err != nil {
		slog.Error(msg, append(keysAndValues, "error", err)...)
		os.Exit(1)
	}
}

func Info(msg string) {
	slog.Info(msg)
}

func Warn(msg string){
	slog.Warn(msg)
}
