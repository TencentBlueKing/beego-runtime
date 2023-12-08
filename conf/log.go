package conf

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

type loggersHook struct{}

func (h *loggersHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *loggersHook) Fire(entry *log.Entry) error {
	pc, file, _, _ := runtime.Caller(9)
	funcName := runtime.FuncForPC(pc).Name()
	entry.Data["funcName"] = funcName
	entry.Data["pathname"] = file
	return nil
}

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

	log.AddHook(&loggersHook{})
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999999Z07:00",
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
