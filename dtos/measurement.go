package dtos

type Measurement struct {
	Unit   string `json:"unit"`
	Values []interface{}
}

type MeasurementList struct {
	Measurements []string
}
