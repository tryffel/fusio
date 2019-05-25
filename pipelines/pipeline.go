package pipelines

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/pipelines/block"
	"github.com/tryffel/fusio/storage"
	"github.com/tryffel/fusio/storage/repository"
	"sync"
	"time"
)

const (
	OnErrorContinue = "continue"
	OnErrorStop     = "stop"
	OnErrorBreak    = "break"
)

// Configuration for low level pipeline
type PipelineConfig struct {
	Id      string      `json:"id"`
	OnError string      `json:"on_error"`
	Blocks  interface{} `json:"blocks"`
}

type streamer interface {
	Push(m *block.Message) error
}

type status int

const (
	pipelineActive  status = 1
	pipelineInCache status = 2
	pipelineOffline status = 3
)

type Data struct {
	Device string
	Value  float64
}

// Pipeline interface, which runs in background and can ingest new measurements
type Pipeline interface {
	Run()
	Stop()
	Ingest(device *string, value float64)
}

// Pipeline struct
type pipeline struct {
	id      string
	onError string
	sigStop chan bool
	input   chan Data
	status  status
	time    time.Ticker
	stream  streamer
	store   *storage.Store
	lock    sync.Mutex
	running bool
}

func (p *pipeline) Ingest(device *string, value float64) {
	data := Data{Device: *device, Value: value}
	p.input <- data
}

func (p *pipeline) Stop() {
	p.sigStop <- true
}

func (p *pipeline) Run() {
	p.lock.Lock()
	if p.running {
		logrus.Error("Pipeline is already running")
		p.lock.Unlock()
		return
	}
	p.running = true
	p.lock.Unlock()

	msg := &block.Message{}

	for {
		select {
		case in := <-p.input:
			logrus.Debugf("pipeline %s got input %f", p.id, in.Value)
			msg.Value = in.Value
			msg.Device = &in.Device
			err := p.stream.Push(msg)
			if err != nil {
				switch p.onError {
				case OnErrorStop:
					logrus.Infof("Error on pipeline %s, stopping, err: %s", p.id, err.Error())
					return
				default:
					logrus.Debugf("Error on pipeline %s, continuing", p.id)
				}
			}
		case <-p.sigStop:
			return
		}
	}
}

// Stream consists of blocks with their internal data and can be pushed to cache
type stream struct {
	// Where to begin from
	firstBlock int
	// all blocks
	blocks map[int]block.Messager
}

func (s *stream) Push(m *block.Message) error {

	if m == nil {
		return errors.New("message cannot be empty")
	}

	if s.blocks[s.firstBlock] == nil {
		return &Err.Error{Code: Err.Einternal, Err: errors.New("no pipeline blocks defined")}
	}

	err := s.blocks[s.firstBlock].Put(m)
	if err != nil {
		return err
	}

	// Process as long as there's new block to run
	for {
		// Stop processing
		if m.NextBlock < 1 {
			break
		}
		m.Value = m.ReturnValue
		err = s.blocks[m.NextBlock].Put(m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *stream) Dump(cache repository.Cache) {

	data, err := UnloadPipelineBlocks(&s.blocks)
	if err != nil {
		logrus.Error(err)
	}

	err = cache.Put(fmt.Sprintf("pipeline_%s", "test1"), data, time.Second*7200)
	if err != nil {
		logrus.Error(err)
	}
}

func (s *stream) Load(cache repository.Cache) {

	data := &[]byte{}
	err := cache.Get("pipeline_test1", *data)
	if err != nil {
		logrus.Error(err)
	}

	blocks, err := LoadPipelineBlocks(data)

	s.blocks = *blocks

}

func NewPipeline(config string, store *storage.Store) (Pipeline, error) {
	p := &pipeline{
		store:   store,
		sigStop: make(chan bool),
		input:   make(chan Data, 10),
	}

	if !gjson.Valid(config) {
		return p, Err.UserError("invalid json", nil)
	}

	gson := gjson.Parse(config)

	if gson.Get("name").String() == "" {
		return p, Err.UserError("name is required parameter", nil)
	}
	if gson.Get("blocks").String() == "" {
		return p, Err.UserError("blocks are required parameter", nil)
	}

	blocks := gson.Get("blocks")

	stream, err := NewStreamer(blocks)
	if err != nil {
		return p, Err.Wrap(&err, "failed to create pipeline")
	}

	p.stream = stream
	return p, nil
}

func ValidateBlocks(pipelineDto string) error {
	gson := gjson.Parse(pipelineDto)

	blocks := gson.Get("blocks")
	_, err := NewStreamer(blocks)
	if err != nil {
		return Err.Wrap(&err, "failed to validate pipeline blocks")
	}
	return nil
}

// NewStream returns working stream aka pipeline
func NewPipelineBytes(config []byte, store *storage.Store) (Pipeline, error) {

	c := string(config)
	return NewPipeline(c, store)
}
