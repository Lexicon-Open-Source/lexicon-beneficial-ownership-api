package bo_v1_services

import (
	"context"
	"lexicon/bo-api/beneficiary_ownership"
	models "lexicon/bo-api/beneficiary_ownership/v1/models"
)

func GetLkppChartData(ctx context.Context) (models.LkppChartsModel, error) {
	tx, err := beneficiary_ownership.Pool.Begin(ctx)

	if err != nil {
		return models.LkppChartsModel{}, err
	}

	list, err := models.LkppChartData(ctx, tx)

	if err != nil {
		return models.LkppChartsModel{}, err
	}

	tx.Commit(ctx)

	return list, nil
}
