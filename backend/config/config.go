package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	ExpireTime int
}

var AppConfig *Config

func InitConfig() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AutomaticEnv()

	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("jwt.secret", "JWT_SECRET")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Warning: Config file not found, using environment variables: %v\n", err)
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getString("server.port", "8080"),
			Mode: getString("server.mode", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getString("database.host", "localhost"),
			Port:     getString("database.port", "3306"),
			User:     getString("database.user", "root"),
			Password: getString("database.password", "root123456"),
			DBName:   getString("database.dbname", "canteen_system"),
		},
		Redis: RedisConfig{
			Host:     getString("redis.host", "localhost"),
			Port:     getString("redis.port", "6379"),
			Password: getString("redis.password", ""),
			DB:       getInt("redis.db", 0),
		},
		JWT: JWTConfig{
			Secret:     getString("jwt.secret", "your_jwt_secret_key_here_development"),
			ExpireTime: getInt("jwt.expire_time", 24),
		},
	}
}

func getString(key string, defaultValue string) string {
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	return defaultValue
}

func getInt(key string, defaultValue int) int {
	if viper.IsSet(key) {
		return viper.GetInt(key)
	}
	return defaultValue
}

func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func (c *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
