package pipelines

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestPipelineRunner(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	c := `
{
	"name": "test",
	"type": "plain",
	"on_error":"stop",
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
				"if": "value % 10 > 8",
				"then": 0,
				"else": 3
			}
		}
	}
}`

	p, err := NewPipeline(c, nil)
	if err != nil {
		t.Error(err)
		return
	}

	d := "d1"

	go p.Run()

	time.Sleep(time.Millisecond * 10)

	for i := 0; i < 100; i++ {
		p.Ingest(&d, float64(i))
		fmt.Printf("\n")

	}

	p.Stop()
}
