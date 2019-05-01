package pipelines

import (
	"github.com/tidwall/gjson"
	"github.com/tryffel/fusio/pipelines/block"
	"testing"
)

func TestLoadBlock(t *testing.T) {

	config :=
		`
{
	"name": "test",
	"type": "plain",
	"blocks":
	{
		"3": {
			"type":"filter.moving_average",
			"config": {
				"length": 5,
				"initial_value": 5,
				"next_block": 2
			}
		},
		"2": {
			"type": "output.log",
			"config": {
				"next_block": 0
			}
			
		},
		"1": {
			"type": "filter.condition",
			"config": {
				"if": "valu{{asd == 10",
				"then": 2,
				"else": 3
			}
		}
	
	}
}`

	s, err := NewStreamer(gjson.Parse(config).Get("blocks"))
	if err != nil {
		t.Error(err)
		return
	}

	d := "device"
	m := "temperature"

	msg := &block.Message{
		Device:      &d,
		Measurement: &m,
		Value:       2,
	}

	err = s.Push(msg)
	if err != nil {
		t.Error(err)
		return
	}

}

func BenchmarkPipeline(b *testing.B) {

	config :=
		`
{
	"name": "test",
	"type": "plain",
	"blocks":
	{
		"1": {
			"type":"filter.moving_average",
			"config": {
				"length": 5,
				"initial_value": 5,
				"next_block": 0
			}
		}
	}
}`

	s, err := NewStreamer(gjson.Parse(config).Get("blocks"))
	if err != nil {
		b.Error(err.Error())
		return
	}

	d := "device"
	m := "temperature"

	msg := &block.Message{
		Device:      &d,
		Measurement: &m,
		Value:       2,
	}

	for i := 0; i < b.N; i++ {
		msg.Value = float64(i)

		s.Push(msg)
	}

}
