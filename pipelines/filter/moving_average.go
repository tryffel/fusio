package filter

import (
	"encoding/json"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/tryffel/fusio/pipelines/block"
)

type MovingAverageConfig struct {
	Length    uint    `json:"length" mapstructure:"length"`
	InitVal   float64 `json:"initial_value" mapstructure:"initial_value"`
	NextBlock int     `json:"next_block" mapstructure:"next_block"`
	Id        int     `json:"id" mapstructure:"id"`
}

func (m *MovingAverageConfig) SetId(id int) {
	m.Id = id
}

func (m *MovingAverageConfig) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, *m)
}

func (m *MovingAverageConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(m)
}

type movingAverage struct {
	Values    []float64
	index     int
	nextBlock int
	id        int
}

func (m *movingAverage) Put(msg *block.Message) error {
	m.index += 1
	m.index %= len(m.Values)
	m.Values[m.index] = msg.Value

	var avg float64
	for _, v := range m.Values {
		avg += v
	}
	avg /= float64(len(m.Values))

	msg.ReturnValue = avg
	msg.SetNext(m.nextBlock)
	return nil
}

func (m *movingAverage) Id() int {
	return m.id
}

func (m *movingAverage) Next() int {
	return m.nextBlock
}

func (m *movingAverage) SetId(id int) {
	m.id = id
}

func NewMovingAverage(conf map[string]interface{}) (block.Messager, error) {

	c := MovingAverageConfig{}
	err := mapstructure.Decode(conf, &c)

	if err != nil {
		return nil, errors.New("invalid configuration type")
	}

	if c.Length < 2 {
		return nil, errors.New("length must be > 1")
	}

	m := &movingAverage{
		Values:    make([]float64, c.Length),
		index:     0,
		nextBlock: c.NextBlock,
		id:        c.Id,
	}

	if c.InitVal != 0 {
		for i := range m.Values {
			m.Values[i] = c.InitVal
		}
	}
	return m, nil
}
