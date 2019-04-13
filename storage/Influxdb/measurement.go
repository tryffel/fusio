package Influxdb

import "time"

// Single measurement
type Point struct {
	Value     float32
	Timestamp time.Time
}

// Series of points
type Series []Point

// Multiple series
type Batch map[string][]Point

//
//type Measurements map[string]Point

// Measurements: key: measurement name
type Measurements map[string]Point

// Metadata to store for measurements
type metadata struct {
	device      string
	groups      []string
	measurement string
}
