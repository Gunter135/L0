package utils

import (
	"log/slog"
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandStringBytes(n int) string {
	//один генератор на всех, потенциальное бутылочное горлышко
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
// SetDefaultLogger в каждом методе на случай если вдруг вызовут один раз только одного из них
func FatalError(err error, msg string, keysAndValues ...interface{}) {
	SetDefaultLogger()
	if err != nil {
		slog.Error(msg, append(keysAndValues, "error", err)...)
		os.Exit(1)
	}
}

func Info(msg string) {
	SetDefaultLogger()
	slog.Info(msg)
}

func Warn(msg string){
	SetDefaultLogger()
	slog.Warn(msg)
}
//do once?
func SetDefaultLogger(){
	opts := slog.HandlerOptions{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}