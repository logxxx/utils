package matcher

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Do(t *testing.T) {
	data := map[string]string{
		"string":   "Abbbbc",
		"int":      "-123",
		"float":    "123.456789",
		"version":  "1.2.3.4",
		"size":     "2GB",
		"size_int": strconv.Itoa(2 * 1024 * 1024 * 1024),
	}

	type config struct {
		key      string
		operator string
		values   []string
	}

	testCases := []struct {
		expression config
		expectErr  bool
		expectRet  bool
	}{
		{
			expression: config{"", "in", []string{"", "Abbbbc", "xx"}},
			expectErr:  true,
		},
		{
			expression: config{"string", "", []string{"", "Abbbbc", "xx"}},
			expectErr:  true,
		},
		{
			expression: config{"string", "in", []string{}},
			expectErr:  false,
			expectRet:  false,
		},

		{
			expression: config{"string", "in", []string{"", "Abbbbc", "xx"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "NotIn", []string{"", "abbbbc", "xx"}},
			expectRet:  true,
		},

		{
			expression: config{"int", "LessThan", []string{"-120"}},
			expectRet:  true,
		},
		{
			expression: config{"int", "LessThan", []string{"-125"}},
			expectRet:  false,
		},
		{
			expression: config{"float", "GreaterThan", []string{"100.001"}},
			expectRet:  true,
		},
		{
			expression: config{"version", "versionLessThan", []string{"1.100"}},
			expectRet:  true,
		},
		{
			expression: config{"version", "versionLessThan", []string{"1.1.100"}},
			expectRet:  false,
		},
		{
			expression: config{"version", "versionGreaterThan", []string{"1.1.1.1.1"}},
			expectRet:  true,
		},
		{
			expression: config{"version", "versionGreaterThan", []string{"1.100.100"}},
			expectRet:  false,
		},
		{
			expression: config{"string", "hasPrefix", []string{"Ab"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "hasPrefix", []string{"ab"}},
			expectRet:  false,
		},
		{
			expression: config{"string", "hasSuffix", []string{"bc"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "hasSuffix", []string{"xx"}},
			expectRet:  false,
		},
		{
			expression: config{"string", "contains", []string{"bbbb"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "contains", []string{"bbxbb"}},
			expectRet:  false,
		},
		{
			expression: config{"string", "regexMatch", []string{"^Ab{4}c$"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "regexMatch", []string{"^Ab{3}c$"}},
			expectRet:  false,
		},

		// intelligent matcher
		{
			expression: config{"string", "==", []string{"Abbbbc"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "!=", []string{"Ac"}},
			expectRet:  true,
		},
		{
			expression: config{"string", ">", []string{"Aa"}},
			expectRet:  true,
		},
		{
			expression: config{"string", ">=", []string{"Abbbbc"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "<", []string{"Ac"}},
			expectRet:  true,
		},
		{
			expression: config{"string", "<=", []string{"Ac"}},
			expectRet:  true,
		},

		{
			expression: config{"int", ">", []string{"-999999"}},
			expectRet:  true,
		},
		{
			expression: config{"int", ">=", []string{"-123"}},
			expectRet:  true,
		},
		{
			expression: config{"int", "<", []string{"0"}},
			expectRet:  true,
		},
		{
			expression: config{"int", "<=", []string{"-122.99999"}},
			expectRet:  true,
		},

		{
			expression: config{"float", ">", []string{"123"}},
			expectRet:  true,
		},
		{
			expression: config{"float", ">=", []string{"123.456789"}},
			expectRet:  true,
		},
		{
			expression: config{"float", "<", []string{"124"}},
			expectRet:  true,
		},
		{
			expression: config{"float", "<=", []string{"124.123"}},
			expectRet:  true,
		},

		{
			expression: config{"version", ">", []string{"1.2.0"}},
			expectRet:  true,
		},
		{
			expression: config{"version", ">=", []string{"1.2.3.4"}},
			expectRet:  true,
		},
		{
			expression: config{"version", "<", []string{"1.2.4"}},
			expectRet:  true,
		},
		{
			expression: config{"version", "<=", []string{"1.3.0.0"}},
			expectRet:  true,
		},

		{
			expression: config{"size", ">", []string{"2047MB"}},
			expectRet:  true,
		},
		{
			expression: config{"size", ">=", []string{"2GB"}},
			expectRet:  true,
		},
		{
			expression: config{"size", "<", []string{"1TB"}},
			expectRet:  true,
		},
		{
			expression: config{"size", "<=", []string{"2048GB"}},
			expectRet:  true,
		},
		{
			expression: config{"size_int", ">", []string{"2047MB"}},
			expectRet:  true,
		},
		{
			expression: config{"size_int", ">=", []string{"2GB"}},
			expectRet:  true,
		},
		{
			expression: config{"size_int", "<", []string{"1TB"}},
			expectRet:  true,
		},
		{
			expression: config{"size_int", "<=", []string{"2048GB"}},
			expectRet:  true,
		},
	}

	for _, v := range testCases {
		filter, err := NewExpression(v.expression.key, v.expression.operator, v.expression.values)
		if v.expectErr {
			assert.NotNil(t, err, "new filter error:%v, case: %v", v.expression, v)
		} else {
			assert.Nil(t, err, "new filter error:%v, case: %v", v.expression, v)
			ret := filter.Match(data)
			assert.Equal(t, v.expectRet, ret, v)
		}
	}
}

func Benchmark_RegexMatch(b *testing.B) {
	pattern := `(^1[0-9]{10}$)|(^.*@.*\.com$)|(^.*@.*\.cn$)`
	data := map[string]string{"key": "someone@logxxx.com"}
	expect := true

	filter, err := NewExpression("key", "regexMatch", []string{pattern})
	if err != nil {
		b.Fatal(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ret := filter.Match(data); ret != expect {
				b.Fatalf("ret:%v, err:%v", ret, err)
			}
		}
	})
}

func Benchmark_Match(b *testing.B) {
	data := map[string]string{"key": "someone@logxxx.com"}

	filter, err := NewExpression("key", "hasSuffix", []string{".com"})
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		filter.Match(data)
	}
}

func Benchmark_MatchNotCompile(b *testing.B) {
	data := map[string]string{"key": "someone@logxxx.com"}

	for i := 0; i < b.N; i++ {
		filter, err := NewExpression("key", "hasSuffix", []string{".com"})
		if err != nil {
			b.Fatal(err)
		}
		filter.Match(data)
	}
}

func Benchmark_OnlyCompile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewExpression("key", "hasSuffix", []string{".com"})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark_MatchRaw(b *testing.B) {
	data := map[string]string{"key": "someone@logxxx.com"}
	for i := 0; i < b.N; i++ {
		strings.HasSuffix(data["key"], ".com")
	}
}
