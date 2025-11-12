package pkg

import "github.com/joho/godotenv"

const (
	BODY_LIMIT_MB = 500
)

var (
	FILE_FOLDER = GetEnv("FILE_FOLDER", "./app/fileDump")
	PORT        = GetEnv("PORT", "4000")
	CORS_ORIGIN = GetEnv("CORS_ORIGIN", "http://localhost:3000")

	// AMQP config
	AMPQ_USER  = GetEnv("AMQP_USER", "guest")
	AMPQ_PASS  = GetEnv("AMQP_PASS", "guest")
	AMPQ_HOST  = GetEnv("AMQP_HOST", "localhost")
	AMPQ_PORT  = GetEnv("AMQP_PORT", "5672")
	AMPQ_VHOST = GetEnv("AMQP_VHOST", "/")
)

var Env map[string]string

func GetEnv(key, def string) string {
	if val, ok := Env[key]; ok {
		return val
	}
	return def
}

func SetupEnvFile() {
	envFile := ".env"
	var err error
	_, err = godotenv.Read(envFile)
	if err != nil {
		panic(err)
	}
}
