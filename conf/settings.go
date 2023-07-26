package conf

import (
	"fmt"
	"os"
	"strings"
	"time"

	machineryConfig "github.com/RichardKnop/machinery/v2/config"
	config "github.com/beego/beego/v2/core/config"
	"github.com/go-redis/redis/v8"
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
apigw_manager_maintainers = ${BK_APIGW_MANAGER_MAINTAINERS}

user_token_key_name = ${USER_TOKEN_KEY_NAME}
plugin_api_debug_username = ${PLUGIN_API_DEBUG_USERNAME}

gcs_mysql_name = ${GCS_MYSQL_NAME}
gcs_mysql_user = ${GCS_MYSQL_USER}
gcs_mysql_password = ${GCS_MYSQL_PASSWORD}
gcs_mysql_host = ${GCS_MYSQL_HOST}
gcs_mysql_port = ${GCS_MYSQL_PORT}

rabbitmq_vhost = ${RABBITMQ_VHOST}
rabbitmq_port = ${RABBITMQ_PORT}
rabbitmq_host = ${RABBITMQ_HOST}
rabbitmq_user = ${RABBITMQ_USER}
rabbitmq_password = ${RABBITMQ_PASSWORD}


store_backend = ${STORE_BACKEND}
`

var Settings config.Configer

var pluginName string
var pluginSecret string
var environment string
var port int
var apigwBackendHost string

var redisAddr string
var redisPassword string
var redisClient *redis.Client
var scheduleExpiration time.Duration
var finishedScheduleExpiration time.Duration
var workerConcurrency int

var mysqlConStr string

var apigwEndpoint string
var apigwApiName string
var apigwManagerMaintainers string

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

func initRedisPassword() {
	redisPassword = Settings.DefaultString("redis_password", "")
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
	workerConcurrency = Settings.DefaultInt("worker_concurrency", 0)
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

func initApigwManagerMaintainers() {
	apigwManagerMaintainers = Settings.DefaultString("apigw_manager_maintainers", "admin")

}
func ApigwManagerMaintainers() []string {
	return strings.Split(apigwManagerMaintainers, ",")
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

func ScheduleStoreMode() string {
	return Settings.DefaultString("store_backend", "mysql")
}

func MachineryCnf() *machineryConfig.Config {
	//rabbitmq_vhost = ${RABBITMQ_VHOST}
	//rabbitmq_port = ${RABBITMQ_PORT}
	//rabbitmq_host = ${RABBITMQ_HOST}
	//rabbitmq_user = ${RABBITMQ_USER}
	//rabbitmq_password = ${RABBITMQ_PASSWORD}

	rabbitmqUser := Settings.DefaultString("rabbitmq_user", "guest")
	rabbitmqPassword := Settings.DefaultString("rabbitmq_password", "guest")
	rabbitmqVhost := Settings.DefaultString("rabbitmq_vhost", "")
	rabbitmqPort := Settings.DefaultString("rabbitmq_port", "5672")
	rabbitmqHost := Settings.DefaultString("rabbitmq_host", "localhost")

	brokerUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		rabbitmqUser,
		rabbitmqPassword,
		rabbitmqHost,
		rabbitmqPort,
		rabbitmqVhost)

	cnf := &machineryConfig.Config{
		Broker:          brokerUrl,
		DefaultQueue:    "schedule",
		ResultsExpireIn: 3600,
		AMQP: &machineryConfig.AMQPConfig{
			Exchange:      "schedule_exchange",
			ExchangeType:  "direct",
			BindingKey:    "schedule_task",
			PrefetchCount: 3,
		},
	}
	return cnf
}

func initMysqlConAddr() {
	//gcs_mysql_name = ${GCS_MYSQL_NAME}
	//gcs_mysql_user = ${GCS_MYSQL_USER}
	//gcs_mysql_password = ${GCS_MYSQL_PASSWORD}
	//gcs_mysql_host = ${GCS_MYSQL_HOST}
	//gcs_mysql_port = ${GCS_MYSQL_PORT}

	mysqlName := Settings.DefaultString("gcs_mysql_name", "'")
	mysqlUsername := Settings.DefaultString("gcs_mysql_user", "root")
	mysqlPassword := Settings.DefaultString("gcs_mysql_password", "root")
	mysqlHost := Settings.DefaultString("gcs_mysql_host", "127.0.0.1")
	mysqlPort := Settings.DefaultInt("gcs_mysql_port", 3306)

	mysqlConStr = fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8",
		mysqlUsername,
		mysqlPassword,
		mysqlHost,
		mysqlPort,
		mysqlName,
	)
}

func MysqlConAddr() string {
	return mysqlConStr
}

func init() {
	var err error
	Settings, err = config.NewConfigData("ini", []byte(configData))
	if err != nil {
		fmt.Printf("runtime config load error: %v\n", err)
		os.Exit(2)
	}
	initMysqlConAddr()
	initPluginName()
	initPluginSecret()
	initEnvironment()
	initServerPort()
	initRedisAddr()
	initRedisPassword()
	initRedisClient()
	initScheduleExpiration()
	initWorkerConcurrency()
	initApigwEndpoint()
	initApigwApiName()
	initApigwManagerMaintainers()
	initApigwBackendHost()
	initUserTokenKeyName()
	initPluginApiDebugUsername()
	setupLog()

}
