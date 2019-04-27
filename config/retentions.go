package config

import (
	"github.com/tryffel/fusio/util"
	"time"
)

var DefaultRetentions = []RetentionPolicy{
	{
		Duration:      util.Interval(24 * time.Hour),
		Interval:      util.Interval(time.Second),
		ChunkInterval: util.Interval(24 * time.Hour),
	},
	{
		Duration:      util.Interval(30 * 24 * time.Hour),
		Interval:      util.Interval(5 * time.Minute),
		ChunkInterval: util.Interval(7 * 24 * time.Hour),
	},
	{
		Duration:      util.Interval(365 * 24 * time.Hour),
		Interval:      util.Interval(30 * time.Minute),
		ChunkInterval: util.Interval(60 * 24 * time.Hour),
	},
	{
		Duration:      util.Interval(10 * 365 * 24 * time.Hour),
		Interval:      util.Interval(60 * time.Minute),
		ChunkInterval: util.Interval(365 * 24 * time.Hour),
	},
}
