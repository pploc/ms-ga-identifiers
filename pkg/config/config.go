package config

import (
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Auth     AuthConfig     `yaml:"auth"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
	Env  string `yaml:"env"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type JWTConfig struct {
	Secret         string        `yaml:"secret"`
	ExpirationTime time.Duration `yaml:"expiration_time"`
	RefreshDuration time.Duration `yaml:"refresh_duration"`
}

type AuthConfig struct {
	ServiceURL string `yaml:"service_url"`
	Timeout    int    `yaml:"timeout"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

func Load() *Config {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.yaml"
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		// Return default config if file doesn't exist
		return defaultConfig()
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return defaultConfig()
	}

	// Override with environment variables if set
	overrideFromEnv(&cfg)

	return &cfg
}

func defaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Env:  "development",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "ms_ga_identifier",
			SSLMode:  "disable",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		JWT: JWTConfig{
			Secret:          "your-secret-key-change-in-production",
			ExpirationTime:  24 * time.Hour,
			RefreshDuration: 7 * 24 * time.Hour,
		},
		Auth: AuthConfig{
			ServiceURL: "http://localhost:8081",
			Timeout:    30,
		},
		Kafka: KafkaConfig{
			Brokers: []string{"localhost:9092"},
			Topic:   "identity-events",
		},
	}
}

func overrideFromEnv(cfg *Config) {
	if v := os.Getenv("SERVER_PORT"); v != "" {
		cfg.Server.Port = v
	}
	if v := os.Getenv("SERVER_ENV"); v != "" {
		cfg.Server.Env = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		cfg.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Database.Port = port
		}
	}
	if v := os.Getenv("DB_USER"); v != "" {
		cfg.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		cfg.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		cfg.Database.DBName = v
	}
	if v := os.Getenv("REDIS_HOST"); v != "" {
		cfg.Redis.Host = v
	}
	if v := os.Getenv("REDIS_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Redis.Port = port
		}
	}
	if v := os.Getenv("JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("AUTH_SERVICE_URL"); v != "" {
		cfg.Auth.ServiceURL = v
	}
}
