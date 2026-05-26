package models

import "time"

type Ingredient struct {
	Id                     int64
	Name                   string
	Category               int64
	DefaultMeasurementType int64
	Description            string
	CreatedAt              time.Time
	CreatedBy              int64
	ModifiedAt             time.Time
	ModifiedBy             int64
}

// Requests
type CreateIngredient struct {
	Name                   string `json:"name"`
	Category               int64  `json:"category"`
	DefaultMeasurementType int64  `json:"defaultMeasurementType"`
	Description            string `json:"description"`
	CreatedBy              int64  `json:"createdBy"`
}

// Responses
