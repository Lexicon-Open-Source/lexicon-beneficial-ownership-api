package main

import (
	"fmt"
	"os"
	"strconv"
)

func loadEnvString(key string, result *string) {
	s, ok := os.LookupEnv(key)

	if !ok {
		return
	}
	*result = s
}

func loadEnvUint(key string, result *uint) {
	s, ok := os.LookupEnv(key)

	if !ok {
		return
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return
	}
	*result = uint(n)
}

/* Configuration */

/* PgSQL Configuration */
type pgSqlConfig struct {
	Host     string `json:"host"`
	Port     uint   `json:"port"`
	Database string `json:"database"`
	SslMode  string `json:"ssl_mode"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (p pgSqlConfig) ConnStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=%s", p.Host, p.Port, p.User, p.Password, p.Database, p.SslMode)
}

func defaultPgSql() pgSqlConfig {
	return pgSqlConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "lexicon_beneficiary_ownership",
		User:     "",
		Password: "",
		SslMode:  "disable",
	}
}

func (p *pgSqlConfig) loadFromEnv() {
	loadEnvString("POSTGRES_HOST", &p.Host)
	loadEnvUint("POSTGRES_PORT", &p.Port)
	loadEnvString("POSTGRES_DB_NAME", &p.Database)
	loadEnvString("APP_POSTGRES_SSLMODE", &p.SslMode)
	loadEnvString("POSTGRES_USERNAME", &p.User)
	loadEnvString("POSTGRES_PASSWORD", &p.Password)
}

/* Listen Configuration */

type listenConfig struct {
	Host string `json:"host"`
	Port uint   `json:"port"`
}

func (l listenConfig) Addr() string {
	return fmt.Sprintf("%s:%d", l.Host, l.Port)
}

func defaultListenConfig() listenConfig {
	return listenConfig{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

func (l *listenConfig) loadFromEnv() {
	loadEnvString("APP_LISTEN_HOST", &l.Host)
	loadEnvUint("APP_LISTEN_PORT", &l.Port)
}

type config struct {
	Listen         listenConfig `json:"listen"`
	PgSql          pgSqlConfig  `json:"pgsql"`
	BackendApiKey  string       `json:"api_key"`
	ServerSalt     string       `json:"salt"`
	ChatbotBaseURL string       `json:"chatbot_base_url"`
}

func (c *config) loadFromEnv() {
	c.Listen.loadFromEnv()
	c.PgSql.loadFromEnv()
	loadEnvString("API_KEY", &c.BackendApiKey)
	loadEnvString("SALT", &c.ServerSalt)
	loadEnvString("CHATBOT_BASE_URL", &c.ChatbotBaseURL)
}

func defaultConfig() config {
	return config{
		Listen:         defaultListenConfig(),
		PgSql:          defaultPgSql(),
		BackendApiKey:  "", // TODO: Default value ​​need to be defined
		ServerSalt:     "", // TODO: Default value ​​need to be defined
		ChatbotBaseURL: "",
	}
}
