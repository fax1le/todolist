package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Addr          string
	PGHost        string
	PGUser        string
	PGPassword    string
	DBName        string
	RedisHost     string
	RedisPassword string
	RedisDb       int
	RedisProtocol int
	LogPath       string
	LogLevel      string
}

func Load() Config {
	return Config{
		Addr:          ":" + getStringEnv("ADDR"),
		PGHost:        getStringEnv("PG_HOST"),
		PGUser:        getStringEnv("PG_USER"),
		PGPassword:    getStringEnv("PG_PASSWORD"),
		DBName:        getStringEnv("DB_NAME"),
		RedisHost:     getStringEnv("REDIS_HOST"),
		RedisPassword: getStringEnv("REDIS_PASSWORD"),
		RedisDb:       getIntEnv("REDIS_DB"),
		RedisProtocol: getIntEnv("REDIS_PROTOCOL"),
		LogPath:       getStringEnv("LOG_PATH"),
		LogLevel:      getStringEnv("LOG_LEVEL"),
	}
}

var notRequiredVars = map[string]string{
	"REDIS_PASSWORD": "REDIS_PASSWORD",
	"LOG_PATH":       "LOG_PATH",
}

func getStringEnv(key string) string {
	env_var := os.Getenv(key)
	_, ok := notRequiredVars[key]

	if env_var == "" && !ok {
		log.Fatal("failed to load config, env variable missing:", key)
	}

	return env_var
}

func getIntEnv(key string) int {
	env_var := os.Getenv(key)

	if env_var == "" {
		log.Fatal("failed to load config, env variable missing:", key)
	}

	val, err := strconv.Atoi(env_var)

	if err != nil {
		log.Fatal("failed to load config, invalid format:", key, err)
	}

	return val
}
