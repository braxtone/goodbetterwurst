package main

import (
	"time"
)

type PercentageChartIngredient struct {
	IngredientType string
	Ingredient     string
	Amount         int
	Note           string
}

type PercentageChart struct {
	ID          int32
	Name        string
	Source      string
	DateCreated time.Time
	Notes       string
	Ingredients []PercentageChartIngredient
}

type PercentageChartInstance struct {
	SourceChart int32
	PercentageChart
}
