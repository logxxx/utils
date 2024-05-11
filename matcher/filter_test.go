package matcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestFilter_Filter(t *testing.T) {
	data := map[string]string{
		"key1":  "1",
		"key2":  "2",
		"keyAA": "AA",
		"keyBB": "BB",
	}

	testCase := []struct {
		filter string
		expect bool
	}{
		{
			filter: `
match: [key1, in, 1]
`,
			expect: true,
		},

		{
			filter: `
match: [key1, in, 2]
`,
			expect: false,
		},

		{
			filter: `
notMatch:
  match: [key1, in, 2]
`,
			expect: true,
		},

		{
			filter: `
notMatch:
  matchAll:
  - match: [key2, in, 2]
  - match: [key2, greaterThan, 0]
`,
			expect: false,
		},

		{
			filter: `
matchAll:
- match: [key2, in, 2]
- match: [key2, greaterThan, 0]
`,
			expect: true,
		},

		{
			filter: `
matchAny:
- match: [key2, in, 2]
- match: [key2, lessThan, 0]
`,
			expect: true,
		},

		{
			filter: `
match: [key1, in, 1]
matchAll:
- match: [key2, in, 2]
- match: [key2, greaterThan, 0]
matchAny:
- match: [key2, lessThan, 0]
- match: [key2, in, 2]
`,
			expect: true,
		},

		{
			filter: `
matchAll:
- match: [key2, in, 2]
- matchAny:
  - match: [key2, lessThan, 0]
  - match: [key2, in, 2]
`,
			expect: true,
		},

		{
			filter: `
matchAll:
- match: [key2, in, 2]
- matchAll:
  - match: [key2, greaterThan, 0]
  - match: [key2, in, 2]
  - matchAny:
    - match: [key2, lessThan, 0]
    - match: [key2, in, 2]
`,
			expect: true,
		},
		{
			filter: `
matchAll:
- match: [key2, in, 2]
- matchAll:
  - match: [key2, greaterThan, 0]
  - match: [key2, in, 2]
  - matchAll:
    - match: [key2, lessThan, 0]
    - match: [key2, in, 2]
`,
			expect: false,
		},
		{
			filter: `
matchAll:
- match: [key2, in, 2]
- matchAll:
  - match: [key2, greaterThan, 0]
  - match: [key2, in, 2]
  - matchAll:
    - matchAny:
      - match: [key2, lessThan, 0]
      - match: [key2, greaterThan, 0]
    - matchAll:
      - match: [key2, lessThan, 0]
      - match: [key2, greaterThan, 0]
`,
			expect: false,
		},
	}

	for _, v := range testCase {
		f := new(Filter)
		if err := yaml.Unmarshal([]byte(v.filter), f); err != nil {
			t.Fatalf("unmarshal: %s, err: %v", v.filter, err)
		}

		got := f.Filter(data)
		assert.Equal(t, v.expect, got, "filter: %s", f)
		t.Logf("filter:%s, got:%v", f, got)

		// NotMatch
		{
			f := &Filter{NotMatch: f}
			got := f.Filter(data)
			assert.Equal(t, !v.expect, got, "filter: %s", f)
			t.Logf("filter:%s, got:%v", f, got)
		}
	}
}

func TestFilter_ModifyValues(t *testing.T) {
	data := map[string]string{
		"key1": "1",
		"key2": "2",
	}
	f := new(Filter)
	bs := []byte(`
match: [key1, in, "1|2|3"]
matchAll:
- match: [key1, in, "1|2|3"]
- match: [key2, in, "1|2|3"]
matchAny:
- match: [key2, in]
- match: [key1, in, "1|2|3"]
`)

	if err := yaml.Unmarshal(bs, f); err != nil {
		t.Fatal(err)
	}

	if f.Filter(data) != false {
		t.Fatal("error test case")
	}

	f.ModifyValues(func(value string) []string {
		return strings.Split(value, "|")
	})

	if f.Filter(data) != true {
		t.Fatal("error test case after ModifyValues")
	}

	t.Log(f.String())
}

func TestFilter_ModifySource(t *testing.T) {
	data := map[string]string{
		"key1": "1",
		"key2": "2",
	}
	f := new(Filter)
	bs := []byte(`
match: [key1, in, "xxx"]
matchAll:
- match: [key1, in, "xxx"]
- match: [key2, in, "xxx"]
matchAny:
- match: [key2, in]
- match: [key1, in, "xxx"]
`)

	if err := yaml.Unmarshal(bs, f); err != nil {
		t.Fatal(err)
	}

	if f.Filter(data) != false {
		t.Fatal("error test case")
	}

	f.ModifySource(func(value string) func(string) interface{} {
		return func(string) interface{} {
			return map[string]struct{}{
				"1": {},
				"2": {},
			}
		}
	})

	if f.Filter(data) != true {
		t.Fatal("error test case after ModifySource")
	}

	t.Log(f.String())
}

func TestFilter_FilterWithResponse(t *testing.T) {
	data := map[string]string{
		"key1":  "1",
		"key2":  "2",
		"keyAA": "AA",
		"keyBB": "BB",
	}

	testCase := []struct {
		filter           string
		expectedRet      bool
		expectedResponse map[string]interface{}
	}{
		{
			filter: `
match: [key1, in, 1]
responseOnMatch:
  key: "{{key1}}"
  key1_constant: "key1"
  key2_constant: "key2"
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"key":           "1",
				"key1_constant": "key1",
				"key2_constant": "key2",
			},
		},
		{
			filter: `
match: [key1, in, 1]
responseAlways:
  keyAlways: "{{key1}}"
responseOnMatch:
  key: "{{key1}}"
  key1_constant: key1
  key2_constant: key2
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"keyAlways":     "1",
				"key":           "1",
				"key1_constant": "key1",
				"key2_constant": "key2",
			},
		},
		{
			filter: `
match: [key1, in, 2]
responseAlways:
  keyAlways: "{{key1}}"
responseOnNotMatch:
  key: "{{key1}}"
  key1_constant: "key1"
  key2_constant: "key2"
`,
			expectedRet: false,
			expectedResponse: map[string]interface{}{
				"keyAlways":     "1",
				"key":           "1",
				"key1_constant": "key1",
				"key2_constant": "key2",
			},
		},

		{
			filter: `
responseOnMatch:
  key: "{{key1}}"
matchAll:
- match: [key2, in, 2]
  responseOnMatch:
    key: "{{key2}}"
    key1_constant: key1
- matchAny:
  - match: [key2, lessThan, 0]
  - match: [key2, in, 2]
    responseOnMatch:
      key: "{{keyAA}}"
      key2_constant: key2
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"key":           "1",
				"key1_constant": "key1",
				"key2_constant": "key2",
			},
		},
	}

	for _, v := range testCase {
		f := new(Filter)
		if err := yaml.Unmarshal([]byte(v.filter), f); err != nil {
			t.Fatalf("unmarshal: %s, err: %v", v.filter, err)
		}

		ret, resp := f.FilterWithResponse(data)
		assert.Equal(t, v.expectedRet, ret, "filter: %s", f)
		assert.EqualValues(t, v.expectedResponse, resp, "filter: %s", f)
		t.Logf("filter:%s, got:%v, gotResp:%v", f, ret, resp)
	}
}

func TestFilter_Walk(t *testing.T) {
	data := map[string]string{
		"key1":  "1",
		"key2":  "2",
		"keyAA": "AA",
		"keyBB": "BB",
	}

	testCase := []struct {
		filter         string
		expectRet      bool
		expectResponse map[string]interface{}
	}{
		{
			filter: `
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         yyx: first
     - match: ["key1", "in", "1"]
       responseOnMatch:
         yy: yyyy
 - match: ["key1", "in", "1"]
   responseOnMatch:
     xx: xxx
responseOnMatch:
 onMatchx: xxxxm
`,
			expectRet: true,
			expectResponse: map[string]interface{}{
				"yyx":      "first",
				"yy":       "yyyy",
				"xx":       "xxx",
				"onMatchx": "xxxxm",
			},
		},
		{
			filter: `
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         yyx: first
     - match: ["key1", "in", "2"]
       responseOnMatch:
         yy: yyyy
 - match: ["key1", "in", "1"]
   responseOnMatch:
     xx: xxx
responseOnMatch:
 onMatchx: xxxxm
`,
			expectRet: false,
			expectResponse: map[string]interface{}{
				"yyx": "first",
				"xx":  "xxx",
			},
		},
		{
			filter: `
matchAny:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         yyx: first
     - match: ["key1", "in", "2"]
       responseOnMatch:
         yy: yyyy
 - match: ["key1", "in", "1"]
   responseOnMatch:
     xx: xxx
responseOnMatch:
  onMatchx: xxxxm
`,
			expectRet: true,
			expectResponse: map[string]interface{}{
				"yyx":      "first",
				"xx":       "xxx",
				"onMatchx": "xxxxm",
			},
		},
		{
			filter: `
matchAny:
  - matchAll:
      - match: ["key1", "in", "2"]
        responseOnMatch:
          yyx: first
      - match: ["key1", "in", "2"]
        responseOnMatch:
          yy: yyyy
  - match: ["key1", "in", "1"]
    responseOnMatch:
      xx: xxx
responseOnMatch:
  onMatchx: xxxxm
`,
			expectRet: true,
			expectResponse: map[string]interface{}{
				"xx":       "xxx",
				"onMatchx": "xxxxm",
			},
		},
		{
			filter: `
matchAny:
 - matchAll:
     - match: ["key1", "in", "2"]
       responseOnMatch:
         yyx: first
     - match: ["key1", "in", "2"]
       responseOnMatch:
         yy: yyyy
 - match: ["key1", "in", "2"]
   responseOnMatch:
     xx: xxx
responseOnMatch:
  onMatchx: xxxxm
`,
			expectRet:      false,
			expectResponse: map[string]interface{}{},
		},
		{
			filter: `
matchAll:
 - match: ["key1", "in", "1"]
   responseOnMatch:
     a: 11
 - match: ["key1", "in", "2"]
   responseOnMatch:
     b: 22
matchAny:
 - match: ["key1", "in", "1"]
   responseOnMatch:
     c: 33
match: ["key1", "in", "1"]
responseOnMatch:
  xx: xxx
`,
			expectRet: false,
			expectResponse: map[string]interface{}{
				"a": 11,
				"c": 33,
			},
		},
		{
			filter: `
matchAny:
 - matchAny:
   - match: ["key1", "in", "1"]
   responseOnMatch:
     xx: "xxxx"
 - match: ["key1", "in", "1"]
   responseOnMatch:
     yy: "yyyy"
`,
			expectRet: true,
			expectResponse: map[string]interface{}{
				"xx": "xxxx",
				"yy": "yyyy",
			},
		},
		{
			filter: `
matchAny:
 - match: ["key1", "in", "1"]
   responseOnMatch:
     yy: "yyyy"
     innerJSON:
       xx: "xx"
       zz: "zz"
`,
			expectRet: true,
			expectResponse: map[string]interface{}{
				"yy": "yyyy",
				"innerJSON": map[string]interface{}{
					"xx": "xx",
					"zz": "zz",
				},
			},
		},
	}

	for _, v := range testCase {
		filter := new(Filter)
		err := yaml.Unmarshal([]byte(v.filter), &filter)
		if err != nil {
			t.Fatal(err)
		}

		ret, resp := filter.Walk(data)
		//fmt.Println(ret, resp)
		assert.Equal(t, v.expectRet, ret, "filter: %s", filter)
		assert.EqualValues(t, toJSON(v.expectResponse), toJSON(resp), "filter: %s", filter)
		t.Logf("filter:%s, got:%v, gotResp:%v", filter, ret, toJSON(resp))
	}
}

func TestFilter_Append(t *testing.T) {
	data := map[string]string{
		"key1":  "1",
		"key2":  "2",
		"keyAA": "AA",
		"keyBB": "BB",
	}

	testCase := []struct {
		filter           string
		expectedRet      bool
		expectedResponse map[string]interface{}
	}{
		{
			filter: `
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         xx: "{{keyAA}}"
     - match: ["key1", "in", "1"]
       responseOnMatch:
         xx: append({{keyBB}})
 - match: ["key1", "in", "1"]
   responseOnMatch:
     xx: append(,third,)
responseOnMatch:
  xx: append(true)
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"xx": "AABB,third,true",
			},
		},
		{
			filter: `
matchAny:
  - match: ["key1", "in", "1"]
    responseOnMatch:
      xx: append(world)
  - match: ["key1", "in", "1"]
    responseOnMatch:
      xx: append({{key2}})
responseAlways:
  xx: ["hello"]
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"xx": []interface{}{"hello", "world", "2"},
			},
		},
	}

	for _, v := range testCase {
		filter := new(Filter)
		err := yaml.Unmarshal([]byte(v.filter), &filter)
		if err != nil {
			t.Fatal(err)
		}

		ret, resp := filter.Walk(data)
		fmt.Println(ret, resp)
		assert.Equal(t, v.expectedRet, ret, "filter: %s", filter)
		assert.EqualValues(t, v.expectedResponse, resp, "filter: %s", filter)
		t.Logf("filter:%s, got:%v, gotResp:%v", filter, ret, resp)
	}
}

func TestFilter_Function(t *testing.T) {
	data := map[string]string{
		"key1":  "1",
		"key2":  "2",
		"keyAA": "AA",
		"keyBB": "BB",
	}

	testCase := []struct {
		filter           string
		expectedRet      bool
		expectedResponse map[string]interface{}
	}{
		{
			filter: `
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         xx: "{{keyAA}}"
         string_1: "{{key1}}"
         number_1: "toNumber({{key1}})"
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"xx":       "AA",
				"string_1": "1",
				"number_1": float64(1),
			},
		},

		{
			filter: `
mergeResponseRecursively: true
responseAlways:
  xx: "BB"
  struct:
    append: "{{key1}}"
    number: "toNumber({{key2}})"
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         xx: "{{keyAA}}"
         struct:
           append: 'append({{key2}})'
           number: "toNumber({{key1}})"
           string: "{{key1}}"
           struct2:
             string: "{{key1}}"
           
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"xx": "AA",
				"struct": map[string]interface{}{
					"append": "12",
					"number": 1,
					"string": "1",
					"struct2": map[string]interface{}{
						"string": "1",
					},
				},
			},
		},

		{
			filter: `
mergeResponseRecursively: false
responseAlways:
  xx: "BB"
  struct:
    not_use_1: xxx
    not_use_2: xxx
    append: "{{key1}}"
    number: "toNumber({{key2}})"
matchAll:
 - matchAll:
     - match: ["key1", "in", "1"]
       responseOnMatch:
         xx: "{{keyAA}}"
         struct:
           append: 'append({{key2}})'
           number: "toNumber({{key1}})"
           string: "{{key1}}"
           struct2:
             string: "{{key1}}"
           
`,
			expectedRet: true,
			expectedResponse: map[string]interface{}{
				"xx": "AA",
				"struct": map[string]interface{}{
					"append": "2", // mergeResponseRecursively: false 这种模式下，append函数获取不到旧的值
					"number": 1,
					"string": "1",
					"struct2": map[string]interface{}{
						"string": "1",
					},
				},
			},
		},
	}

	for _, v := range testCase {
		filter := new(Filter)
		err := yaml.Unmarshal([]byte(v.filter), &filter)
		if err != nil {
			t.Fatal(err)
		}

		ret, resp := filter.Walk(data)
		assert.Equal(t, v.expectedRet, ret, "filter: %s", filter)
		assert.EqualValues(t, toJSON(v.expectedResponse), toJSON(resp), "filter: %s", filter)

		t.Logf("filter:%s, got:%v, gotResp:%v", filter, ret, toJSON(resp))
	}
}

// recursively
func TestMergeValueRecursively(t *testing.T) {
	strToMap := func(s string) (map[string]interface{}, error) {
		bs := []byte(strings.TrimSpace(s))
		m := make(map[string]interface{})
		var err error
		if bytes.HasPrefix(bs, []byte("{")) && json.Valid(bs) {
			err = json.Unmarshal(bs, &m)
		} else {
			err = yaml.Unmarshal(bs, &m)
			//t.Logf("unmarshal yaml:%s, result:%v", bs, m)
		}
		if err != nil {
			return nil, err
		}

		m2 := interfaceMapToStringMap(m).(map[string]interface{})
		return m2, nil
	}

	mapToStr := func(v map[string]interface{}) string {
		m := interfaceMapToStringMap(v).(map[string]interface{})
		bs, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return fmt.Sprintf("ERROR: %s, raw input:%+v", err.Error(), v)
		}
		return string(bs)
	}

	var testCases = []struct {
		targetStr string
		fromStr   string
		data      map[string]string

		expectStr string
	}{
		{
			targetStr: `
{
  "str": "v_str",
  "int": 100,
  "float": 99.999,
  "bool": true,
  "str2": "v_str_2"
}
`,
			fromStr: `
{
  "str": "{{test}}",
  "int": 0,
  "float": -100000.1111,
  "bool": false
}
`,
			data: map[string]string{
				"test": "testVal",
			},
			expectStr: `
{
  "str": "testVal",
  "int": 0,
  "float": -100000.1111,
  "bool": false,
  "str2": "v_str_2"
}
`,
		},

		{
			targetStr: `
{
  "str": "testVal",
  "int": 100,
  "float": 99.999,
  "bool": true,
  "str2": "v_str_2",
  "map": {
    "str": "v_str",
    "map": {
      "str": "v_str",
      "map": {
        "str": "v_str"
      }
    }
  }
}
`,
			fromStr: `
{
  "str": "v_str_2",
  "int": 0,
  "float": -100000.1111,
  "bool": false,
  "map": {
    "str": "v_str",
    "str2": "v_str2",
    "map": {
      "str": "v_str",
      "str2": "v_str2",
      "map": "XXX"
    }
  }
}
`,
			expectStr: `
{
  "str": "v_str_2",
  "int": 0,
  "float": -100000.1111,
  "bool": false,
  "str2": "v_str_2",
  "map": {
    "str": "v_str",
    "str2": "v_str2",
    "map": {
      "str": "v_str",
      "str2": "v_str2",
      "map": "XXX"
    }
  }
}
`,
		},

		{
			targetStr: `
map:
  str: v_str
  100: 999
  bool: true
  map:
    str: v_str
    map:
      str: v_str
  array: [1, 2, 3]
`,
			fromStr: `
map:
  str: "{{test}}"
  100: 777
  bool: false
  101: 888
  bool2: false
  map:
    str: v_str2
    map: xxx
  array: [ "{{test}}", "test", "toNumber({{num}})" ]
  arrayStruct:
    - key1: [ "{{test}}", "test", "toNumber({{num}})" ]
    - key2: "{{test}}"
    - key3: "toNumber(10.0)"
    - key4: "xxx"
`,
			data: map[string]string{
				"test": "testVal",
				"num":  "10.000",
			},
			expectStr: `
map:
  str: testVal
  100: 777
  bool: false
  101: 888
  bool2: false
  map:
    str: v_str2
    map: xxx
  array: [ "testVal", "test", 10 ]
  arrayStruct:
    - key1: [ "testVal", "test", 10 ]
    - key2: "testVal"
    - key3: 10
    - key4: "xxx"
`,
		},
	}

	for i, v := range testCases {
		target, err := strToMap(v.targetStr)
		if err != nil {
			t.Fatal(err)
		}
		from, err := strToMap(v.fromStr)
		if err != nil {
			t.Fatal(err)
		}
		expectM, err := strToMap(v.expectStr)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("\ncase:%d --------------------------", i)
		t.Logf("\n---target:%s,\n---raw target:%#v", mapToStr(target), target)
		t.Logf("\n---from:%s", mapToStr(from))
		t.Logf("\n---expect:%s", mapToStr(expectM))

		mergeValue(target, from, v.data, true)

		expect := mapToStr(expectM)
		got := mapToStr(target)

		t.Logf("\n---got:%s", got)
		if expect != got {
			t.Fatal("result not match")
		}
	}
}
