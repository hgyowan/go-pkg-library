package envs

import (
	"log"
	"os"
	"strconv"
)

const (
	PrdType             = "prd"
	StgType             = "stg"
	DevType             = "dev"
	UserTokenHeaderName = "x-user-token"
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

	RedisAddr       = getEnv("REDIS_ADDR", "")
	RedisPort       = getEnv("REDIS_PORT", "6379")
	RedisPassword   = getEnv("REDIS_PASSWORD", "")
	RedisMasterName = getEnv("REDIS_MASTER_NAME", "")

	SMTPAccount  = getEnv("SMTP_ACCOUNT", "")
	SMTPServer   = getEnv("SMTP_SERVER", "")
	SMTPPort     = getEnv("SMTP_PORT", "587")
	SMTPPassword = getEnv("SMTP_PASSWORD", "")
	SMTPSender   = getEnv("SMTP_SENDER", "")

	KakaoGranTType    = getEnv("KAKAO_GRANT_TYPE", "")
	KakaoClientID     = getEnv("KAKAO_CLIENT_ID", "")
	KakaoRedirectURI  = getEnv("KAKAO_REDIRECT_URI", "")
	KakaoClientSecret = getEnv("KAKAO_CLIENT_SECRET", "")

	JwtAccessSecret  = getEnv("JWT_ACCESS_SECRET", "")
	JwtRefreshSecret = getEnv("JWT_REFRESH_SECRET", "")
	JwtIssuer        = getEnv("JWT_ISSUER", "")

	SecretKey      = getEnv("SECRET_KEY", "")
	CBCSecretIVKey = getEnv("CBC_SECRET_IV_KEY", "")
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
