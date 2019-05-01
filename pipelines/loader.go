package pipelines

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/pipelines/block"
	"github.com/vmihailenco/msgpack"
	"strconv"
)

// LoadPipelineBlocks loads pipeline from byte array
func LoadPipelineBlocks(data *[]byte) (*map[int]block.Messager, error) {
	blocks := &map[int]block.Messager{}
	return blocks, msgpack.Unmarshal(*data, &blocks)
}

// UnloadPipelineBlocks unloads blocks into byte array
func UnloadPipelineBlocks(blocks *map[int]block.Messager) (*[]byte, error) {
	val, err := msgpack.Marshal(*blocks)
	return &val, err
}

func NewStreamer(blocks gjson.Result) (streamer, error) {
	s := &stream{}
	s.blocks = make(map[int]block.Messager)

	arr := blocks.Map()

	// Create blocks
	for i, v := range arr {
		id, err := strconv.ParseInt(i, 10, 32)
		if err != nil {
			return s, &Err.Error{Code: Err.Einvalid, Message: fmt.Sprintf("pipeline id '%s' not a number", i)}
		}

		b, err := LoadBlock(int(id), &v)
		if err != nil {
			return s, err
		}
		if b == nil {
			return s, &Err.Error{Code: Err.Einternal, Err: errors.New("got nil block")}
		} else {
			s.blocks[b.Id()] = b
		}
	}

	first := s.blocks[1]
	if first == nil {
		return s, errors.New("no 1st block defined. \"1\" must be defined")
	}

	s.firstBlock = 1

	return s, nil
}

// LoadBlock tries to parse configuration for a block and return a functioning block messager
// root: {"type: "test", "config": {...}}
func LoadBlock(id int, root *gjson.Result) (block.Messager, error) {

	Type := root.Get("type").String()
	if Type == "" {
		return nil, errors.New("No block type defined")
	}

	blockType := Blocks[Type]
	if blockType.Config == nil {
		return nil, &Err.Error{Code: Err.Enotfound, Message: fmt.Sprintf("block '%s' not found", Type)}
	}
	conf := root.Get("config").Value()

	c, ok := conf.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid configuration block")
	}

	messager, err := blockType.Init(c)
	if err != nil || messager == nil {
		return nil, Err.Wrap(&err, fmt.Sprintf("failed to intiialize block %s (%d)", Type, id))
	}

	messager.SetId(id)
	return messager, nil
}
