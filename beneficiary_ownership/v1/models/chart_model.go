package bo_v1_models

import (
	"context"

	common_models "lexicon/bo-api/common/models"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type ChartsModel struct {
	Countries    []common_models.BaseChartModel `json:"countries"`
	SubjectTypes []common_models.BaseChartModel `json:"subjet_types"`
	CaseTypes    []common_models.BaseChartModel `json:"case_types"`
}

var emptyChartModel ChartsModel

func ChartData(ctx context.Context, tx pgx.Tx) (ChartsModel, error) {

	var chartsResult ChartsModel

	// Get Countries Chart Data
	countriesQuery := `
		SELECT 
			c.nation name,
			count(c.nation) value
		FROM 
			cases c 
		GROUP BY
			c.nation
	`
	log.Info().Msg("Executing query: " + countriesQuery)

	countries, err := tx.Query(ctx, countriesQuery)

	log.Info().Msg("Finish countries query")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyChartModel, err
	}

	defer countries.Close()

	for countries.Next() {
		var chartResult common_models.BaseChartModel

		err = countries.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyChartModel, err
		}

		chartsResult.Countries = append(chartsResult.Countries, chartResult)
	}

	// Get Subject Types Chart Data
	subjectTypesQuery := `
		SELECT 
			CASE 
				WHEN c.subject_type = 1 THEN 'Individual'
				WHEN c.subject_type = 2 THEN 'Company'
				WHEN c.subject_type = 3 THEN 'Organization'
			END name,
			count(c.subject_type) value
		FROM 
			cases c
		GROUP BY
			c.subject_type
	`
	log.Info().Msg("Executing query: " + subjectTypesQuery)

	subjectTypes, err := tx.Query(ctx, subjectTypesQuery)

	log.Info().Msg("Subject Types query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyChartModel, err
	}

	defer subjectTypes.Close()

	for subjectTypes.Next() {
		var chartResult common_models.BaseChartModel

		err = subjectTypes.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyChartModel, err
		}

		chartsResult.SubjectTypes = append(chartsResult.SubjectTypes, chartResult)
	}

	// Get Case Types Chart Data
	caseTypesQuery := `
		SELECT
			CASE 
				WHEN c.case_type = 1 THEN 'Verdict'
				WHEN c.case_type = 2 THEN 'Blacklist'
				WHEN c.case_type = 3 THEN 'Sanction'
			END name,
			count(c.case_type) value
		FROM 
			cases c 
		GROUP BY
			c.case_type
	`
	log.Info().Msg("Executing query: " + caseTypesQuery)

	caseTypes, err := tx.Query(ctx, caseTypesQuery)

	log.Info().Msg("Case Types query executed")

	if err != nil {
		log.Error().Err(err).Msg("Error querying database")
		return emptyChartModel, err
	}

	defer caseTypes.Close()

	for caseTypes.Next() {
		var chartResult common_models.BaseChartModel

		err = caseTypes.Scan(&chartResult.Name, &chartResult.Value)

		if err != nil {
			return emptyChartModel, err
		}

		chartsResult.CaseTypes = append(chartsResult.CaseTypes, chartResult)
	}

	return chartsResult, nil
}
