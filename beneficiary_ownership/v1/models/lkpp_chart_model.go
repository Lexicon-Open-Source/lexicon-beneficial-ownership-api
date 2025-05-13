package bo_v1_models

import (
	"context"

	common_models "lexicon/bo-api/common/models"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	// "gopkg.in/guregu/null.v4"
)

// type common_models.BaseChartModel struct {
// 	Name  null.String `json:"name"`
// 	Value float64     `json:"value"`
// }

type LkppChartsModel struct {
	BlacklistProvinces    []common_models.BaseChartModel           `json:"blacklist_province"`
	CeilingDistribution   []common_models.BaseChartModel           `json:"ceiling_distribution"`
	TopTenReporters       []common_models.BaseChartModel           `json:"top_ten_reporter"`
	ScenarioDistribution  []common_models.BaseChartModelFloatValue `json:"scenario_distribution"`
	ViolationDistribution []common_models.BaseChartModelFloatValue `json:"violation_distribution"`
}

var emptyLkppChartModel LkppChartsModel

func LkppChartData(ctx context.Context, tx pgx.Tx) (LkppChartsModel, error) {

	var lkppChartsResult LkppChartsModel

	// Get Blacklist by Province Chart Data
	blacklistProvincesQuery := `
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
	log.Info().Msg("Executing query: " + blacklistProvincesQuery)

	blacklistProvinces, err := tx.Query(ctx, blacklistProvincesQuery)

	log.Info().Msg("Blacklist by Province Chart Data query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying Blacklist by Province Chart Data")
		return emptyLkppChartModel, err
	}

	defer blacklistProvinces.Close()

	for blacklistProvinces.Next() {
		var chartResult common_models.BaseChartModel

		err = blacklistProvinces.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.BlacklistProvinces = append(lkppChartsResult.BlacklistProvinces, chartResult)
	}

	// Get Ceiling Distribution Chart Data
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

	log.Info().Msg("Ceiling Distribution Chart Data query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying Ceiling Distribution Chart Data")
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

	// Get Top Ten Reporters Chart Data
	topTenReportersQuery := `
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
	log.Info().Msg("Executing query: " + topTenReportersQuery)

	topTenReporters, err := tx.Query(ctx, topTenReportersQuery)

	log.Info().Msg("Top Ten Reporters Chart Data query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying Top Ten Reporters Chart Data")
		return emptyLkppChartModel, err
	}

	defer topTenReporters.Close()

	for topTenReporters.Next() {
		var chartResult common_models.BaseChartModel

		err = topTenReporters.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.TopTenReporters = append(lkppChartsResult.TopTenReporters, chartResult)
	}

	// Get Blacklist Distribution by Scenario Chart Data
	scenarioDistributionQuery := `
		SELECT
			c.extra_data -> 0 -> 'data' ->> 'scenario' AS dimension,
			round(count(*)::decimal / sum(count(*)) OVER (), 3)*100 AS percentage
		FROM
			public.cases c
		WHERE
			c.extra_data -> 0 ->> 'type' = 'LKPP'
		GROUP BY
			1
		ORDER BY 
			2 DESC;
	`
	log.Info().Msg("Executing query: " + scenarioDistributionQuery)

	scenarioDistribution, err := tx.Query(ctx, scenarioDistributionQuery)

	log.Info().Msg("Blacklist Distribution by Scenario Chart Data query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying Blacklist Distribution by Scenario Chart Data")
		return emptyLkppChartModel, err
	}

	defer scenarioDistribution.Close()

	for scenarioDistribution.Next() {
		var chartResult common_models.BaseChartModelFloatValue

		err = scenarioDistribution.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.ScenarioDistribution = append(lkppChartsResult.ScenarioDistribution, chartResult)
	}

	// Get Distribution of Violation Chart Data
	violationDistributionQuery := `
		WITH ranked AS (
			SELECT
				c.extra_data -> 0 -> 'data' ->> 'rule' AS dimension,
				count(*) AS total,
				rank() OVER (ORDER BY count(*) desc) AS rnk
			FROM
				public.cases c 
			WHERE
				c.extra_data -> 0 ->> 'type' = 'LKPP'
			GROUP BY
				1
		),
		dimensions AS (
			SELECT
				CASE
					WHEN rnk <= 5 THEN dimension
					ELSE 'Other'
				END AS dimension,
				sum(total) AS total
			FROM 
				ranked
			GROUP BY
				1
		),
		final_data AS (
			SELECT
				dimension,
				round(total::decimal / sum(total) OVER (), 3) AS percentage
			FROM
				dimensions
		)
		SELECT 
			dimension,
			percentage * 100 as percentage
		FROM 
			final_data
		ORDER BY
			CASE
				WHEN dimension = 'Other' THEN 99
				ELSE 1
			END,
			percentage desc;
	`
	log.Info().Msg("Executing query: " + violationDistributionQuery)

	violationDistribution, err := tx.Query(ctx, violationDistributionQuery)

	log.Info().Msg("Distribution of Violation Chart Data query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying Distribution of Violation Chart Data")
		return emptyLkppChartModel, err
	}

	defer violationDistribution.Close()

	for violationDistribution.Next() {
		var chartResult common_models.BaseChartModelFloatValue

		err = violationDistribution.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyLkppChartModel, err
		}

		lkppChartsResult.ViolationDistribution = append(lkppChartsResult.ViolationDistribution, chartResult)
	}

	return lkppChartsResult, nil
}
