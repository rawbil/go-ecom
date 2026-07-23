package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
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
	JwtSecret          string
	JwtExpire          int
	RefreshTokenExpire int
}

type ResendConfig struct {
	ApiKey    string
	EmailFrom string
}

func InitConfig() Config {
	return Config{
		Port:       getEnv("PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "rawbil"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_ADDRESS", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:     getEnv("DB_NAME", "go1"),
		ParseTime:  getEnv("PARSE_TIME", "true") == "true",
	}
}

var (
	loadEnvOnce sync.Once
	loadEnvErr  error
)

func LoadEnv() error {
	loadEnvOnce.Do(func() {
		loadEnvErr = godotenv.Load(envFilePath())
	})

	return loadEnvErr
}

func GetServerAddr() string {
	return getEnv("SERVER_ADDR", ":8080")
}

func GetJwtConfig() JwtConfig {
	return JwtConfig{
		JwtSecret:          getEnv("JWT_SECRET", ""),
		JwtExpire:          int(getIntEnv("JWT_EXPIRE", 3600)),
		RefreshTokenExpire: int(getIntEnv("REFRESH_TOKEN_EXPIRE", 7)),
	}
}

func GetResendConfig() *ResendConfig {
	return &ResendConfig{
		ApiKey:    getEnv("RESEND_API_KEY", ""),
		EmailFrom: getEnv("RESEND_EMAIL", ""),
	}
}

// If ok is false, return fallback
func getEnv(key, fallback string) string {
	_ = LoadEnv()

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func envFilePath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ".env"
	}

	for {
		path := filepath.Join(dir, ".env")
		if _, err := os.Stat(path); err == nil {
			return path
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ".env"
		}

		dir = parent
	}
}

func getIntEnv(key string, fallback int64) int64 {
	_ = LoadEnv()
	
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
