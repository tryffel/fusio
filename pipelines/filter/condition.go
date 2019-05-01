package filter

import (
	"errors"
	"github.com/knetic/govaluate"
	"github.com/mitchellh/mapstructure"
	"github.com/tryffel/fusio/err"
	"github.com/tryffel/fusio/pipelines/block"
)

type ConditionConfig struct {
	If        string `json:"if" mapstructure:"if"`
	Id        int    `json:"id" mapstructure:"id"`
	IfBlock   int    `json:"then" mapstructure:"then"`
	ElseBlock int    `json:"else" mapstructure:"else"`
}

type Condition struct {
	exp       *govaluate.EvaluableExpression
	refVal    float64
	id        int
	ifBlock   int
	elseBlock int
}

//func NewCondition(id int, ifId int, elseId int, refval float64) block.Messager {
func NewCondition(conf map[string]interface{}) (block.Messager, error) {

	c := ConditionConfig{}
	err := mapstructure.Decode(conf, &c)

	if err != nil {
		return nil, errors.New("invalid configuration type")
	}

	valuator, err := govaluate.NewEvaluableExpression(c.If)
	if err != nil {
		return nil, Err.UserError("invalid conditional clause", err)
	}

	variables := valuator.Vars()
	if len(variables) == 0 {
		return nil, Err.UserError("no parameters in condition", nil)
	}

	// Check it has {{value}} at least once
	valueExists := false
	for _, v := range variables {
		if v == "value" {
			valueExists = true
		}
	}
	if !valueExists {
		return nil, Err.UserError("invalid condition: 'value' has to exists at least once", nil)
	}

	return &Condition{
		exp:       valuator,
		refVal:    0,
		id:        c.Id,
		ifBlock:   c.IfBlock,
		elseBlock: c.ElseBlock,
	}, nil
}

func (c *Condition) Put(msg *block.Message) error {
	vars := map[string]interface{}{}
	vars["value"] = msg.Value

	status, err := c.exp.Evaluate(vars)
	if err != nil {
		return Err.InternalError("failed to evaluate condition", err)
	}

	if status == false {
		msg.SetNext(c.elseBlock)
	} else if status == true {
		msg.SetNext(c.ifBlock)
	}
	msg.ReturnValue = msg.Value
	return nil
}

func (c *Condition) Id() int {
	return c.id
}

func (c *Condition) Next() int {
	return c.ifBlock
}

func (c *Condition) SetId(id int) {
	c.id = id
}
