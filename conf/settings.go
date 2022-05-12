package conf

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/go-redis/redis/v8"
	"github.com/hibiken/asynq"
)

const configData string = `
plugin_name = ${BKPAAS_APP_ID}
plugin_secret = ${BKPAAS_APP_SECRET}
environment = ${BKPAAS_ENVIRONMENT}
log_file_prefix = ${BKPAAS_LOG_NAME_PREFIX}
process_type = ${BKPAAS_PROCESS_TYPE}
app_default_subdomains = ${BKPAAS_ENGINE_APP_DEFAULT_SUBDOMAINS}

redis_host = ${REDIS_HOST}
redis_port = ${REDIS_PORT}
redis_password = ${REDIS_PASSWORD}

schedule_expiration = ${SCHEDULE_EXPIRATION}
finished_schedule_expiration = ${FINISHED_SCHEDULE_EXPIRATION}
worker_concurrency = ${SCHEDULE_WORKER_CONCURRENCY}

apigw_api_name = ${BKPAAS_BK_PLUGIN_APIGW_NAME}
apigw_endpoint = ${BK_APIGW_MANAGER_URL_TEMPL}

user_token_key_name = ${USER_TOKEN_KEY_NAME}
plugin_api_debug_username = ${PLUGIN_API_DEBUG_USERNAME}
`

var Settings config.Configer

var pluginName string
var pluginSecret string
var environment string
var port int
var apigwBackendHost string

var redisAddr string
var redisPassword string
var asynqClient *asynq.Client
var redisClient *redis.Client
var scheduleExpiration time.Duration
var finishedScheduleExpiration time.Duration
var workerConcurrency int

var apigwEndpoint string
var apigwApiName string

var userTokenKeyName string
var pluginApiDebugUsername string

func IsDevMode() bool {
	return Settings.DefaultString("environment", "dev") == "dev"
}

func initPluginName() {
	pluginName = Settings.DefaultString("plugin_name", "")
}

func PluginName() string {
	return pluginName
}

func initPluginSecret() {
	pluginSecret = Settings.DefaultString("plugin_secret", "")
}

func PluginSecret() string {
	return pluginSecret
}

func initEnvironment() {
	environment = Settings.DefaultString("environment", "dev")
}

func Environment() string {
	return environment
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

func initApigwBackendHost() {
	subdomains := Settings.DefaultString("app_default_subdomains", "")
	apigwBackendHost = strings.Split(subdomains, ";")[0]
}

func ApigwBackendHost() string {
	return apigwBackendHost
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

func initRedisPassword() {
	redisPassword = Settings.DefaultString("redis_password", "")
}

func RedisPassword() string {
	return redisPassword
}

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
}

func RedisClient() *redis.Client {
	return redisClient
}

func initAsynqClient() {
	asynqClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
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

func initApigwEndpoint() {
	apigwEndpoint = Settings.DefaultString("apigw_endpoint", "")
}

func ApigwEndpoint() string {
	return apigwEndpoint
}

func initApigwApiName() {
	apigwApiName = Settings.DefaultString("apigw_api_name", "'")
}

func ApigwApiName() string {
	return apigwApiName
}
func initUserTokenKeyName() {
	var tokenDefaultKey string
	if IsDevMode() {
		tokenDefaultKey = "bk_token"
	} else {
		tokenDefaultKey = "jwt"
	}
	userTokenKeyName = Settings.DefaultString("user_token_key_name", tokenDefaultKey)
}

func UserTokenKeyName() string {
	return userTokenKeyName
}
func initPluginApiDebugUsername() {
	pluginApiDebugUsername = Settings.DefaultString("plugin_api_debug_username", "")
	if !IsDevMode() {
		pluginApiDebugUsername = ""
	}
}

func PluginApiDebugUsername() string {
	return pluginApiDebugUsername
}

func init() {
	var err error
	Settings, err = config.NewConfigData("ini", []byte(configData))
	if err != nil {
		fmt.Printf("runtime config load error: %v\n", err)
		os.Exit(2)
	}

	initPluginName()
	initPluginSecret()
	initEnvironment()
	initServerPort()
	initRedisAddr()
	initRedisPassword()
	initRedisClient()
	initAsynqClient()
	initScheduleExpiration()
	initWorkerConcurrency()
	initApigwEndpoint()
	initApigwApiName()
	initApigwBackendHost()
	setupLog()

	initUserTokenKeyName()
	initPluginApiDebugUsername()
}
