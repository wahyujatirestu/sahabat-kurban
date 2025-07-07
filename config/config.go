package config

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host 		string
	Port 		string
	Username 	string
	Password 	string
	DBName		string
	Driver		string
}

type ApiConfig struct {
	ApiPort		string
}

type TokenConfig struct {
	AppName				string
	JwtSignatureKey		[]byte
	JwtSigningMethod	*jwt.SigningMethodHMAC
	AccessTokenLifetime	time.Duration
}

type Config struct {
	DBConfig
	ApiConfig
	TokenConfig
}

func (c *Config) ReadConfig() error {
	_ = godotenv.Load()

	c.DBConfig = DBConfig{
		Host: 		os.Getenv("DB_HOST"),
		Port: 		os.Getenv("DB_PORT"),
		Username: 	os.Getenv("DB_USERNAME"),
		Password: 	os.Getenv("DB_PASSWORD"),
		DBName: 	os.Getenv("DB_NAME"),
		Driver: 	os.Getenv("DB_DRIVER"),
	}

	c.ApiConfig = ApiConfig{
		ApiPort: 	os.Getenv("API_PORT"),
	}

	accessTokenLifetime := time.Duration(20) * time.Minute

	c.TokenConfig = TokenConfig{
		AppName: 				"Sahabat-Kurban",
		JwtSignatureKey: 		[]byte(os.Getenv("ACCESS_TOKEN")),
		JwtSigningMethod: 		jwt.SigningMethodHS256,
		AccessTokenLifetime: 	accessTokenLifetime,
	}

	if c.Host == "" || c.Port == "" || c.Username == "" || c.Password == "" || c.DBName == "" {
		return  errors.New("Some config is empty")
	}

	return nil
}

func NewConfig() (*Config, error) {
	config := &Config{}

	if err := config.ReadConfig(); err != nil {
		return nil, err
	}

	return config, nil
}