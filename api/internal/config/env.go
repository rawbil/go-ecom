package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	ParseTime  bool
}

type JwtConfig struct {
	JwtSecret string
	JwtExpire int
}

func InitConfig() Config {
	return Config{
		Port:       getEnv("PORT", "8080"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "rawbil"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_ADDRESS", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "go1"),
		ParseTime:  getEnv("PARSE_TIME", "true") == "true",
	}
}

func GetServerAddr() string {
	return getEnv("SERVER_ADDR", ":8080")
}

func GetJwtConfig() *JwtConfig {
	return &JwtConfig{
		JwtSecret: getEnv("JWT_SECRET", "myverysecretkE4"),
		JwtExpire: int(getIntEnv("JWT_EXPIRE", 3600)),
	}
}

// If ok is false, return fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
