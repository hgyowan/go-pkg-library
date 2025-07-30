package envs

import (
	"log"
	"os"
	"strconv"
)

const (
	PrdType                      = "prd"
	StgType                      = "stg"
	DevType                      = "dev"
	UserTokenHeaderName          = "x-user-auth-token"
	UserLangCodeHeaderName       = "langCode"
	DefaultUploadPartSize  int64 = 5 << 20 //5MiB
	APP_ID                       = "PTL"
	OLIM_EMAIL_DOMAIN            = "olimplanet.com"
)

var (
	ServerName  = getEnv("SERVER_NAME", "")
	ServiceType = getEnv("SERVICE_TYPE", "dev") // prd / stg / dev
	ServerPort  = getEnv("SERVER_PORT", ":50051")

	LogLevel = getEnv("LOG_LEVEL", "debug") //debug | info | warn | error |

	DBHost     = getEnv("DB_HOST", "")
	DBPort     = getEnv("DB_PORT", "")
	DBName     = getEnv("DB_NAME", "")
	DBUser     = getEnv("DB_USER", "")
	DBPassword = getEnv("DB_PASSWORD", "")

	CFMAPIHost = getEnv("CFM_API_HOST", "")

	RedisAddr     = getEnv("REDIS_ADDR", "")
	RedisPort     = getEnv("REDIS_PORT", "6379")
	RedisPassword = getEnv("REDIS_PASSWORD", "")

	SMTPAccount  = getEnv("SMTP_ACCOUNT", "")
	SMTPServer   = getEnv("SMTP_SERVER", "")
	SMTPPort     = getEnv("SMTP_PORT", "587")
	SMTPPassword = getEnv("SMTP_PASSWORD", "")
	SMTPSender   = getEnv("SMTP_SENDER", "")
)

func getEnv(envName, defaultVal string) string {
	envVal := os.Getenv(envName)
	if envVal == "" {
		envVal = defaultVal
	}
	return envVal
}

func mustAtoi(val string) int {
	intVal, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	return intVal
}
