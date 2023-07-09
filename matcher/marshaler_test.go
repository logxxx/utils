package matcher

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestExpression_UnmarshalYAML(t *testing.T) {
	data := map[string]string{"key": "12"}
	exp := new(Expression)
	in := []byte(`
- key
- in
- 12
- xy
`)
	err := yaml.Unmarshal(in, exp)
	if err != nil {
		t.Fatalf("unmarshal from %s, err:%v", in, err)
	}

	t.Log(exp.String())

	if exp.Match(data) != true {
		t.Fatal("result not match")
	}
}

func TestExpression_UnmarshalJSON(t *testing.T) {
	data := map[string]string{"key": "12"}
	exp := new(Expression)
	in := []byte(`["key","in","12","xy"]`)
	err := json.Unmarshal(in, exp)
	if err != nil {
		t.Fatalf("unmarshal from %s, err:%v", in, err)
	}

	t.Log(exp.String())

	if exp.Match(data) != true {
		t.Fatal("result not match")
	}
}
