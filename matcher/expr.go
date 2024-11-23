package matcher

import (
	"errors"
	"strings"
)

// ErrUnknownOperator indicates that the operator is notMatcher registered.
var ErrUnknownOperator = errors.New("unknown operator")

// Expression basic expression with matcher.
type Expression struct {
	key      string
	operator string
	values   []string

	matcher Matcher
}

// NewExpression create a compiled filter.
func NewExpression(key, operator string, values []string) (*Expression, error) {
	if key == "" {
		return nil, errors.New("empty key for filter")
	}

	f, _ := registeredMatchers.Load(strings.ToLower(operator))
	if f == nil {
		return nil, ErrUnknownOperator
	}

	m, err := f.(newMatcherFunc)(values)
	if err != nil {
		return nil, err
	}

	return &Expression{
		key:      key,
		operator: operator,
		values:   values,
		matcher:  m,
	}, nil
}

// Match return true if f.key in data match f.matcher.
func (exp *Expression) Match(data map[string]string) bool {
	src := data[exp.key]
	return exp.matcher.Match(src)
}

// String implements the Stringer interface, return json format string.
func (exp *Expression) String() string {
	bs, _ := exp.MarshalJSON()
	return string(bs)
}

// ModifyValues update exp.matcher with new args, and do not update exp.values.
func (exp *Expression) ModifyValues(modifier func(value string) []string) error {
	if exp == nil || len(exp.values) == 0 {
		return nil
	}

	var newValues []string
	updated := false
	for _, v := range exp.values {
		vs := modifier(v)
		if !updated && (len(vs) != 1 || vs[0] != v) {
			updated = true
		}
		newValues = append(newValues, vs...)
	}

	if !updated {
		return nil
	}

	f, _ := registeredMatchers.Load(strings.ToLower(exp.operator))
	if f == nil {
		return ErrUnknownOperator
	}

	m, err := f.(newMatcherFunc)(newValues)
	if err != nil {
		return err
	}
	exp.matcher = m
	return nil
}

// ModifyDataSource update exp.matcher with new data source, and do not update exp.values.
func (exp *Expression) ModifyDataSource(modifier func(value string) func(string) interface{}) error {
	if exp == nil || len(exp.values) == 0 {
		return nil
	}

	var dataSource []func(string) interface{}
	var values []string
	for _, v := range exp.values {
		if ds := modifier(v); ds != nil {
			dataSource = append(dataSource, ds)
		} else {
			values = append(values, v)
		}
	}

	if len(dataSource) == 0 {
		return nil
	}

	f, _ := registeredMatchers.Load(strings.ToLower(exp.operator))
	if f == nil {
		return ErrUnknownOperator
	}

	m, err := f.(newMatcherFunc)(values, dataSource...)
	if err != nil {
		return err
	}
	exp.matcher = m

	return nil
}
