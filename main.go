package main

import (
	"context"
	bo "lexicon/bo-api/beneficiary_ownership"
	"lexicon/bo-api/common/utils"
	"net/http"
	"time"

	"github.com/golang-module/carbon/v2"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// INITIATE CONFIGURATION
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg(".env file not found, using environment variables from system")
	}
	cfg := defaultConfig()
	cfg.loadFromEnv()

	ctx := context.Background()

	carbon.SetDefault(carbon.Default{
		Layout:       carbon.ISO8601Layout,
		Timezone:     carbon.UTC,
		WeekStartsAt: carbon.Monday,
		Locale:       "en",
	})

	// INITIATE DATABASES
	// PGSQL
	pgsqlClient, err := pgxpool.New(ctx, cfg.PgSql.ConnStr())

	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to PGSQL Database")
	}
	defer pgsqlClient.Close()

	bo.SetDatabase(pgsqlClient)

	// init httpClient
	httpClient := http.Client{
		Timeout: time.Minute * 5,
	}
	utils.SetClient(&httpClient)
	// INITIATE SERVER
	server, err := NewLexiconBOServer(cfg)

	if err != nil {
		log.Error().Err(err).Msg("Failed to start the server")
	}

	server.setupRoute()
	server.start()

}
