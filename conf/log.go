package conf

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func setupLog() {
	if IsDevMode() {
		return
	}

	rand.Seed(time.Now().UnixNano())

	log.SetFormatter(&log.JSONFormatter{
		FieldMap: log.FieldMap{
			log.FieldKeyLevel: "levelname",
			log.FieldKeyMsg:   "message",
			log.FieldKeyFile:  "pathname",
			log.FieldKeyFunc:  "funcName",
		},
	})
	log.SetOutput(&lumberjack.Logger{
		Filename: fmt.Sprintf(
			"/app/v3logs/%v-%v-%v.log",
			Settings.DefaultString("log_file_prefix", ""),
			randSeq(4),
			Settings.DefaultString("process_type", ""),
		),
		MaxSize: 500,
	})
}
