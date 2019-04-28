package filter

import (
	"github.com/tryffel/fusio/pipeline/block"
)

type MovingAverage struct {
	Values []float64
	index  int
}

func (m *MovingAverage) Put(msg *block.Message) error {
	m.index += 1
	m.index %= len(m.Values)
	m.Values[m.index] = msg.Value

	var avg float64
	for _, v := range m.Values {
		avg += v
	}
	avg /= float64(len(m.Values))

	msg.ReturnValue = avg
	return nil
}

func NewMovingAverage(n int, initialVal float64) block.Messager {
	m := &MovingAverage{
		Values: make([]float64, n),
		index:  0,
	}

	if initialVal != 0 {
		for i := range m.Values {
			m.Values[i] = initialVal
		}
	}
	return m
}
