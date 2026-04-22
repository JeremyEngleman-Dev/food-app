package dto

// Requests
type CreateIngredient struct {
	Name                   string `json:"name"`
	Category               int64  `json:"category"`
	DefaultMeasurementType int64  `json:"defaultMeasurementType"`
	Description            string `json:"description"`
	CreatedBy              int64  `json:"createdBy"`
}

// Responses
