package model

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
