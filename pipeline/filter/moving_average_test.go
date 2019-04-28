package filter

import (
	"github.com/tryffel/fusio/pipeline/block"
	"testing"
)

func TestMovingAverage_Put(t *testing.T) {

	ma := NewMovingAverage(5, 1)

	msg := &block.Message{}
	msg.Value = 2
	err := ma.Put(msg)
	if err != nil {
		t.Error("putting to moving average failed: ", err)
	}
	msg.Value = 3
	_ = ma.Put(msg)

	msg.Value = 4
	_ = ma.Put(msg)

	msg.Value = -4
	_ = ma.Put(msg)

	if msg.ReturnValue != 1.2 {
		t.Errorf("invalid return value from initial round, exptected %f, got %f",
			1.2, msg.ReturnValue)
	}

	msg.Value = 5
	_ = ma.Put(msg)

	msg.Value = 10
	_ = ma.Put(msg)

	if msg.ReturnValue != 3.6 {
		t.Errorf("invalid return value from second round, exptected %f, got %f",
			3.6, msg.ReturnValue)
	}
}

func BenchmarkMovingAverage10_Put(b *testing.B) {
	ma := NewMovingAverage(10, 1)
	msg := &block.Message{Value: 2}

	for i := 0; i < b.N; i++ {
		msg.Value = float64(i % 250)
		_ = ma.Put(msg)
	}
}

func BenchmarkMovingAverage1000_Put(b *testing.B) {
	ma := NewMovingAverage(1000, 1)
	msg := &block.Message{Value: 2}

	for i := 0; i < b.N; i++ {
		msg.Value = float64(i % 250)
		_ = ma.Put(msg)
	}
}
