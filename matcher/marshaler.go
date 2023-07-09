package matcher

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
)

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (exp *Expression) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var a []string
	if err := unmarshal(&a); err != nil {
		return err
	}
	return exp.unmarshal(a)
}

// MarshalYAML implements the yaml.Marshaler interface.
func (exp Expression) MarshalYAML() (interface{}, error) {
	a := make([]string, 0, len(exp.values)+2)
	a = append(a, exp.key)
	a = append(a, exp.operator)
	a = append(a, exp.values...)
	return yaml.Marshal(a)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (exp *Expression) UnmarshalJSON(b []byte) error {
	var a []string
	err := json.Unmarshal(b, &a)
	if err != nil {
		return err
	}
	return exp.unmarshal(a)
}

// MarshalJSON implements the json.Marshaler interface.
func (exp Expression) MarshalJSON() ([]byte, error) {
	a := make([]string, 0, len(exp.values)+2)
	a = append(a, exp.key)
	a = append(a, exp.operator)
	a = append(a, exp.values...)
	return json.Marshal(a)
}

func (exp *Expression) unmarshal(in []string) error {
	if len(in) < 2 {
		return errors.New("array size must greater than 1, format: [key, operator, values...]")
	}
	v, err := NewExpression(in[0], in[1], in[2:])
	if err != nil {
		return err
	}

	*exp = *v
	return nil
}
