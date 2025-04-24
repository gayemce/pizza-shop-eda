package config

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/gayemce/pizza-shop-eda/order-service/logger"
	"github.com/joho/godotenv"
)

var env ConfigDto

type ConfigDto struct {
	port                string
	database_url        string
	database_name       string
	kafka_host          string
	kafka_port          string
	kafka_default_topic string
	kafka_group_id      string
}

func init() {
	if env.port == "" {
		ConfigEnv()
	}
}

func ConfigEnv() {
	LoadEnvVariable()
	env = ConfigDto{
		port:                os.Getenv("PORT"),
		database_url:        os.Getenv("MONGO_DB_URL"),
		database_name:       os.Getenv("MONGO_DB_NAME"),
		kafka_host:          os.Getenv("KAFKA_HOST"),
		kafka_port:          os.Getenv("KAFKA_PORT"),
		kafka_default_topic: os.Getenv("KAFKA_DEFAULT_TOPIC"),
		kafka_group_id:      os.Getenv("KAFKA_GROUP_ID"),
	}
}

func LoadEnvVariable() {
	if _, err := os.Stat(".env"); err != nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("error loading .env file: %v", err)
		}
	} else {
		logger.Log("No .env file found using system environment variable")
	}
}

func accessField(key string) (string, error) {
	v := reflect.ValueOf(env)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct got %v", t)
	}

	_, ok := t.FieldByName(key)
	if !ok {
		return "", fmt.Errorf("key %v does not exists", key)
	}

	fv := v.FieldByName(key)
	return fv.String(), nil
}

func GetEnvProperty(key string) string {
	if env.port == "" {
		ConfigEnv()
	}
	val, err := accessField(key)
	if err != nil {
		logger.Log(fmt.Sprintf("error accessing field: %v", err))
	}
	return val
}