package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

const configData string = `
pluginname = ${BKPAAS_APP_ID}
runmode = ${BKPAAS_ENVIRONMENT}

redis_host = ${REDIS_HOST}
redis_port = ${REDIS_PORT}
redis_password = ${REDIS_PASSWORD}

schedule_expiration = ${SCHEDULE_EXPIRATION}
finished_schedule_expiration = ${FINISHED_SCHEDULE_EXPIRATION}

worker_concurrency = ${SCHEDULE_WORKER_CONCURRENCY}
`

var Settings config.Configer
var pluginName string
var port int
var workerNum int
var redisAddr string
var asynqClient *asynq.Client
var redisClient *redis.Client
var scheduleExpiration time.Duration
var finishedScheduleExpiration time.Duration
var workerConcurrency int

func IsDevMode() bool {
	return Settings.DefaultString("runmode", "dev") == "dev"
}

func initPluginName() {
	pluginName = Settings.DefaultString("pluginname", "")
}

func PluginName() string {
	return pluginName
}

func initServerPort() {
	if IsDevMode() {
		port = 8000
	} else {
		port = 5000
	}
}

func Port() int {
	return port
}

func initRedisAddr() {
	redisAddr = fmt.Sprintf(
		"%v:%v",
		Settings.DefaultString("redis_host", "127.0.0.1"),
		Settings.DefaultString("redis_port", "6379"),
	)
}

func RedisAddr() string {
	return redisAddr
}

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: Settings.DefaultString("redis_password", ""),
		DB:       0,
	})
}

func RedisClient() *redis.Client {
	return redisClient
}

func initAsynqClient() {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
}

func AsynqClient() *asynq.Client {
	return asynqClient
}

func initScheduleExpiration() {
	scheduleExpiration = time.Duration(Settings.DefaultInt("schedule_expiration", 30)) * 24 * time.Hour
	finishedScheduleExpiration = time.Duration(Settings.DefaultInt("schedule_expiration", 7)) * 24 * time.Hour
}

func ScheduleExpiration() time.Duration {
	return scheduleExpiration
}

func FinishedScheduleExpiration() time.Duration {
	return finishedScheduleExpiration
}

func initWorkerConcurrency() {
	workerConcurrency = Settings.DefaultInt("worker_concurrency", 20)
}

func WorkerConcurrency() int {
	return workerConcurrency
}

func init() {
	var err error
	Settings, err = config.NewConfigData("ini", []byte(configData))
	if err != nil {
		fmt.Printf("runtime config load error: %v\n", err)
		os.Exit(2)
	}

	initServerPort()
	initRedisAddr()
	initRedisClient()
	initAsynqClient()
	initScheduleExpiration()
	initWorkerConcurrency()
}
