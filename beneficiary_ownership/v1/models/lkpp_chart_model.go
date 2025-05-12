package bo_v1_models

import (
	"context"

	common_models "lexicon/bo-api/common/models"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type LkppChartsModel struct {
	BlacklistProvince    []common_models.BaseChartModel `json:"blacklist_province"`
	CeilingDistribution  []common_models.BaseChartModel `json:"ceiling_distribution"`
	TopTenReporter       []common_models.BaseChartModel `json:"top_ten_reporter"`
	ScenarioDistribution []common_models.BaseChartModel `json:"scenario_distribution"`
}

var emptyLkppChartModel LkppChartsModel

func LkppChartData(ctx context.Context, tx pgx.Tx) (LkppChartsModel, error) {

	var lkppChartsResult LkppChartsModel

	// Get Blacklist by Province Chart Data
	blacklistProvinceQuery := `
		SELECT
			c.extra_data -> 0 -> 'data' ->> 'province' dimension,
			count(c.extra_data -> 0 -> 'data' ->> 'province') total
		FROM 
			public.cases c
		WHERE
			c.extra_data -> 0 ->> 'type' = 'LKPP'
			AND c.extra_data -> 0 -> 'data' ->> 'province' <> ''
			AND c.extra_data -> 0 -> 'data' ->> 'province' <> '-'
		GROUP BY
			1
		ORDER BY 
			2 desc;
	`
	log.Info().Msg("Executing query: " + blacklistProvinceQuery)

	blacklistProvinces, err := tx.Query(ctx, blacklistProvinceQuery)

	log.Info().Msg("Finish countries query")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyLkppChartModel, err
	}

	defer blacklistProvinces.Close()

	for blacklistProvinces.Next() {
		var chartResult common_models.BaseChartModel

		err = blacklistProvinces.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.BlacklistProvince = append(lkppChartsResult.BlacklistProvince, chartResult)
	}

	// Get Subject Types Chart Data
	ceilingDistributionQuery := `
		SELECT
			a.*
		FROM (
			SELECT
			CASE
				WHEN (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint >= 0 AND (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint < 2500000000 THEN '0 - 2.5 B'
				WHEN (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint >= 2500000000 AND (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint < 15000000000 THEN '2.5 B - 15 B'
				WHEN (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint >= 15000000000 AND (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint < 50000000000 THEN '15 B - 50 B'
				WHEN (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint >= 50000000000 AND (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint <= 100000000000 THEN '50 B - 100 B'
				WHEN (c.extra_data -> 0 -> 'data' ->> 'ceiling')::bigint > 100000000000 THEN '> 100 B'
			END AS dimension,
			COUNT(*) AS value
			FROM 
				public.cases c
			WHERE
				c.extra_data -> 0 ->> 'type' = 'LKPP'
			GROUP BY dimension
		) a
		ORDER BY
			CASE
				WHEN a.dimension = '0 - 2.5 B' THEN 1
				WHEN a.dimension = '2.5 B - 15 B' THEN 2
				WHEN a.dimension = '15 B - 50 B' THEN 3
				WHEN a.dimension = '50 B - 100 B' THEN 4
				WHEN a.dimension = '> 100 B' THEN 5
			END ASC;
	`
	log.Info().Msg("Executing query: " + ceilingDistributionQuery)

	ceilingDistributions, err := tx.Query(ctx, ceilingDistributionQuery)

	log.Info().Msg("Subject Types query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyLkppChartModel, err
	}

	defer ceilingDistributions.Close()

	for ceilingDistributions.Next() {
		var chartResult common_models.BaseChartModel

		err = ceilingDistributions.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.CeilingDistribution = append(lkppChartsResult.CeilingDistribution, chartResult)
	}

	// Get Case Types Chart Data
	topTenReporterQuery := `
		SELECT
			c.extra_data -> 0 -> 'data' ->> 'institution_area' AS dimension,
			count(*) total_report
		FROM
			public.cases c
		WHERE
			c.extra_data -> 0 ->> 'type' = 'LKPP'
		GROUP BY
			1
		ORDER BY
			2 DESC
		LIMIT 10;
	`
	log.Info().Msg("Executing query: " + topTenReporterQuery)

	topTenReporters, err := tx.Query(ctx, topTenReporterQuery)

	log.Info().Msg("Case Types query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyLkppChartModel, err
	}

	defer topTenReporters.Close()

	for topTenReporters.Next() {
		var chartResult common_models.BaseChartModel

		err = topTenReporters.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.TopTenReporter = append(lkppChartsResult.TopTenReporter, chartResult)
	}

	return lkppChartsResult, nil
}
