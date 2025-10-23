package bo_v1_services

import (
	"context"
	"fmt"
	"lexicon/bo-api/beneficiary_ownership"
	models "lexicon/bo-api/beneficiary_ownership/v1/models"
	bo_v1_repositories "lexicon/bo-api/beneficiary_ownership/v1/repositories"
	"os"

	"github.com/rs/zerolog/log"
)

func GetDetail(ctx context.Context, id string) (models.DetailResultModel, error) {
	tx, err := beneficiary_ownership.Pool.Begin(ctx)

	if err != nil {
		return models.DetailResultModel{}, err
	}

	detail, err := models.GetDetailById(ctx, tx, id)

	if err != nil {
		return models.DetailResultModel{}, err
	}

	tx.Commit(ctx)

	return detail, nil
}

func GetUrlByCaseNumber(ctx context.Context, caseNumber []string) ([]string, error) {
	tx, err := beneficiary_ownership.Pool.Begin(ctx)

	if err != nil {
		log.Err(err).Msg("Error starting transaction")
		return nil, err
	}

	id, err := bo_v1_repositories.GetIdsByCaseNumbers(ctx, tx, caseNumber)
	if err != nil {
		log.Err(err).Msg("Error getting ids by case numbers")
		tx.Rollback(ctx)
		return nil, err
	}

	// Get base URL from environment variable
	baseURL := os.Getenv("BASE_URL")

	var urls []string

	for _, v := range id {
		url := fmt.Sprintf("%s/data/%s", baseURL, v)
		urls = append(urls, url)
	}

	tx.Commit(ctx)

	return urls, nil
}
