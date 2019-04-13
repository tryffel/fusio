package alarm

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

var regexSplitAggregation, _ = regexp.Compile("(\\w+)\\(([a-z]+),?([0-9]+)?\\)")

//  Currently only support two fields, more advanced queries will require more
type Aggregation struct {
	Function string
	Key      string
}

// GetAggregation Split Aggregation functions into logical parts:
// e.g. derivative(temperature, 10) -> derivative, temperature, 10
func GetAggregations(s string) (*[]Aggregation, error) {
	match := regexSplitAggregation.FindAllStringSubmatch(s, 2)
	aggr := make([]Aggregation, len(match))
	for i, el := range match {
		if (len(el) <= 2) || (len(el) > 4) {
			return &aggr, errors.New("Invalid aggregation")
		}
		aggr[i] = Aggregation{
			Function: el[1],
			Key:      el[2],
		}
	}
	return &aggr, nil
}

func (a *Aggregation) ToString() string {
	return fmt.Sprintf("%s(\"%s\")", a.Function, a.Key)
}

//func (a *Aggregation) ToInput() models.Input {
//	//i := models.Input{
//		Aggregation: a.Function,
//		Key:         a.Key,
//	}
//	return i
//}
