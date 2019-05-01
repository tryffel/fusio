package output

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/tryffel/fusio/pipelines/block"
)

type LogConfig struct {
	NextBlock int `json:"next"`
	Id        int `json:"id"`
}

func (l *LogConfig) SetId(id int) {
	l.Id = id
}

type log struct {
	NextBlock int
	id        int
}

func (l *log) Id() int {
	return l.id
}

func (l *log) Next() int {
	return l.NextBlock
}

func (l *log) SetId(id int) {
	l.id = id
}

func (l *log) Put(msg *block.Message) error {
	fmt.Print(msg.Value)
	msg.SetNext(l.NextBlock)
	return nil
}

func NewLog(conf map[string]interface{}) (block.Messager, error) {
	c := LogConfig{}
	err := mapstructure.Decode(conf, &c)
	if err != nil {
		return nil, err
	}
	return &log{NextBlock: c.NextBlock}, nil
}
