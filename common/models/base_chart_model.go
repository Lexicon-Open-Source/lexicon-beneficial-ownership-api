package models

import "gopkg.in/guregu/null.v4"

type BaseChartModel struct {
	Name  null.String `json:"name"`
	Value int64       `json:"value"`
}

type BaseChartModelFloatValue struct {
	Name  null.String `json:"name"`
	Value float64     `json:"value"`
}
