package pipelines

import (
	"errors"
	"github.com/tryffel/fusio/pipelines/block"
	"github.com/tryffel/fusio/pipelines/filter"
	"github.com/tryffel/fusio/pipelines/output"
)

type blockType struct {
	Config interface{}
	BlockConstructor
	Init func(map[string]interface{}) (block.Messager, error)
}

type BlockConstructor interface {
	New() (block.Messager, error)
}

func (b *blockType) New() (block.Messager, error) {
	constructor, ok := b.Config.(BlockConstructor)
	if !ok {
		return nil, errors.New("not a valid block constructor")
	}
	return constructor.New()
}

// All pipeline blocks that application is aware of
var Blocks = map[string]blockType{
	"filter.moving_average": {Config: filter.MovingAverageConfig{}, Init: filter.NewMovingAverage},
	"filter.condition":      {Config: filter.ConditionConfig{}, Init: filter.NewCondition},
	"output.log":            {Config: output.LogConfig{}, Init: output.NewLog},
}
