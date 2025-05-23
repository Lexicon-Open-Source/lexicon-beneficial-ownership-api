package bo_v1_repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/guregu/null.v4"
)

type DetailResultModel struct {
	ID                   ulid.ULID   `json:"id"`
	Subject              string      `json:"subject"`
	SubjectType          string      `json:"subject_type"`
	PersonInCharge       null.String `json:"person_in_charge"`
	BenificiaryOwnership null.String `json:"benificiary_ownership"`
	Date                 null.Time   `json:"date"`
	DecisionNumber       null.String `json:"decision_number"`
	Source               string      `json:"source"`
	Link                 string      `json:"link"`
	Nation               string      `json:"nation"`
	PunishmentDuration   null.String `json:"punishment_duration"`
	Type                 string      `json:"type"`
	Year                 string      `json:"year"`
	Summary              string      `json:"summary"`
}

var emptyDetail DetailResultModel

func GetDetailById(ctx context.Context, tx pgx.Tx, id string) (DetailResultModel, error) {
	var result DetailResultModel

	log.Info().Msg("Start getting detail by id: " + id)
	query := `
	SELECT id, subject, subject_type, person_in_charge, benificiary_ownership, date, decision_number, source, link, nation, punishment_duration, type, year, summary
	FROM cases
	WHERE id = $1
	LIMIT 1
	`

	log.Info().Msg("Executing query: " + query)
	row := tx.QueryRow(ctx, query, id)
	err := row.Scan(&result.ID, &result.Subject, &result.SubjectType, &result.PersonInCharge, &result.BenificiaryOwnership, &result.Date, &result.DecisionNumber, &result.Source, &result.Link, &result.Nation, &result.PunishmentDuration, &result.Type, &result.Year, &result.Summary)

	if err != nil {
		log.Info().Msg("Data Not Found")

		return emptyDetail, err
	}

	log.Info().Msg("Finish getting detail by id: " + id)
	log.Info().Msg("Data Found")
	return result, nil
}

func GetIdsByCaseNumbers(ctx context.Context, tx pgx.Tx, caseNumbers []string) ([]string, error) {
	var ids []string

	query := `
	SELECT id
	FROM cases
	WHERE decision_number = ANY($1)
	`

	row, err := tx.Query(ctx, query, caseNumbers)
	if err != nil {
		log.Err(err).Msg("Error executing query")
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		var id string
		err := row.Scan(&id)
		if err != nil {
			log.Err(err).Msg("Error scanning row")
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}
