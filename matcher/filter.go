package matcher

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Filter 过滤器
// yaml examples：
//
// filter:
//
//	match: [key, in , 12, xy]
//
// filter:
//
//	matchAll:
//	- match: [key, in , 12, xy]
//	- match: [key, notIn , 34, xy]
//
// filter:
//
//	matchAny:
//	- match: [key, in , 12]
//	- match: [key, in , xy]
//
// filter:
//
//	matchAll:
//	- match: [key, in , 12, xy]
//	- matchAny:
//	  - match: [key, notIn , 34, xy]
//	  - match: [key, notIn , 56, xy]
type Filter struct {
	// Description 规则描述，不影响判断结果
	Description string `json:"desc,omitempty" yaml:"desc,omitempty"`

	// MergeResponseRecursively 合并Response对象时，对深层的map也执行数据合并操作
	MergeResponseRecursively bool `json:"mergeResponseRecursively,omitempty" yaml:"mergeResponseRecursively,omitempty"`

	// ResponseAlways 无论如何都返回
	ResponseAlways map[string]interface{} `json:"responseAlways,omitempty" yaml:"responseAlways,omitempty"`

	// ResponseOnMatch 当Filter匹配上后，返回这些信息, 支持使用单/双引号返回常量字符串
	ResponseOnMatch map[string]interface{} `json:"responseOnMatch,omitempty" yaml:"responseOnMatch,omitempty"`
	// ResponseOnNotMatch 当Filter匹配上后，返回这些信息, 支持使用单/双引号返回常量字符串
	ResponseOnNotMatch map[string]interface{} `json:"responseOnNotMatch,omitempty" yaml:"responseOnNotMatch,omitempty"`

	Match    *Expression `json:"match,omitempty" yaml:"match,omitempty"`       // 最基础的表达式
	NotMatch *Filter     `json:"notMatch,omitempty" yaml:"notMatch,omitempty"` // 只包含一个Filter,不匹配才算通过
	MatchAll []*Filter   `json:"matchAll,omitempty" yaml:"matchAll,omitempty"` // 一组Filters，每个都匹配才算通过
	MatchAny []*Filter   `json:"matchAny,omitempty" yaml:"matchAny,omitempty"` // 一组Filters，任意一个匹配都算通过
}

// FilterWithResponse 执行Filter并返回数据结果
func (f *Filter) FilterWithResponse(data map[string]string) (bool, map[string]interface{}) {
	if f == nil {
		return false, map[string]interface{}{}
	}
	if data == nil {
		data = map[string]string{}
	}

	resp := make(map[string]interface{})
	ok := f.doFilter(data, resp, f.MergeResponseRecursively)
	if !ok {
		//resp = nil
	}
	return ok, resp
}

// Walk 和Filter功能不同, 使用Walk时, matchAny和matchAll会走完所有子选项, 遇到false也不会立即返回
// 目的是将所有能匹配的分支中的responseOnMatch都返回出来
// 如:
//
//	data := map[string]string {
//	    "key1": "1",
//	}
//	matchAll:
//	  - match: ["key1", "in", "2"]
//	    responseOnMatch:
//	      xx: "xx"
//	  - match: ["key1", "in", "1"]
//	    responseOnMatch:
//	      yy: "yy"
//
// 返回的结果为: false, map[string]interface{}{ "yy": "yy" }
// 不会因为第一个match为false而退出
func (f *Filter) Walk(data map[string]string) (bool, map[string]interface{}) {
	if f == nil {
		return false, map[string]interface{}{}
	}
	if data == nil {
		data = map[string]string{}
	}

	resp := make(map[string]interface{})
	ok := f.walkFilter(data, resp, f.MergeResponseRecursively)
	return ok, resp
}

func (f *Filter) walkFilter(data map[string]string, resp map[string]interface{}, mergeValueRecursively bool) (result bool) {
	res := -1
	mergeValue(resp, f.ResponseAlways, data, mergeValueRecursively)
	defer func() {
		if result {
			mergeValue(resp, f.ResponseOnMatch, data, mergeValueRecursively)
		}
	}()

	if len(f.MatchAll) > 0 {
		pass := true
		for _, sub := range f.MatchAll {
			if sub == nil {
				continue
			}
			if !sub.walkFilter(data, resp, mergeValueRecursively) {
				pass = false
			}
		}
		if pass {
			res = 1
		} else {
			res = 0
		}
	}

	if len(f.MatchAny) > 0 {
		pass := false
		for _, sub := range f.MatchAny {
			if sub == nil {
				continue
			}
			if sub.walkFilter(data, resp, mergeValueRecursively) {
				pass = true
				if res != 0 {
					res = 1
				}
			}
		}
		if !pass {
			res = 0
		}
	}

	if f.Match != nil {
		if res == 0 {
			return false
		}
		return f.Match.Match(data)
	}

	return res == 1
}

// Filter 匹配Filter中所有条件，通过返回true，不通过返回false
// 未写任何条件的Filter返回true
func (f *Filter) Filter(data map[string]string) bool {
	if f == nil {
		return false
	}
	ok, _ := f.FilterWithResponse(data)
	return ok
}

func getValue(data map[string]string, key string) string {
	if strings.HasPrefix(key, "{{") && strings.HasSuffix(key, "}}") {
		return data[key[2:len(key)-2]]
	}
	return key
}

func mergeValue(target map[string]interface{}, from map[string]interface{}, data map[string]string, recursively bool) {
	if len(from) == 0 {
		return
	}

	// yaml.Unmarshal 解析出来的map是map[interface{}]interface{}类型， 这里强制改为map[string]interface{}
	m, _ := interfaceMapToStringMap(from).(map[string]interface{})
	for k, v := range m {
		if v == nil || v == "" {
			delete(target, k)
			continue
		}

		switch v := v.(type) {
		case string:
			function, arg, found := findFunction(v)
			if found {
				target[k] = function(target[k], getValue(data, arg))
			} else {
				target[k] = getValue(data, v)
			}

		case []interface{}:
			target[k] = copySliceForMerge(v, data, recursively)

		case map[string]interface{}:
			if recursively {
				targetChild, _ := target[k]
				targetChildMap, _ := targetChild.(map[string]interface{})
				if targetChildMap == nil {
					targetChildMap = make(map[string]interface{})
					target[k] = targetChildMap // 类型不同时，创建新的map继续merge，因为v的子map中可能还有函数需要解析
				}
				mergeValue(targetChildMap, v, data, recursively)
			} else {
				// 创建新的map继续merge，因为v的子map中可能还有函数需要解析
				// mergeResponseRecursively: false 这种模式下，append函数获取不到旧的值
				targetChildMap := make(map[string]interface{})
				target[k] = targetChildMap
				mergeValue(targetChildMap, v, data, recursively)
			}

		default:
			target[k] = v
		}
	}
}

func copySliceForMerge(slice []interface{}, data map[string]string, recursively bool) []interface{} {
	if len(slice) == 0 {
		return slice
	}

	news := make([]interface{}, 0, len(slice))
	for _, v := range slice {
		switch v := v.(type) {
		case string:
			function, arg, found := findFunction(v)
			if found {
				news = append(news, function(nil, getValue(data, arg)))
			} else {
				news = append(news, getValue(data, v))
			}
		case []interface{}:
			news = append(news, copySliceForMerge(v, data, recursively))
		case map[string]interface{}:
			m := make(map[string]interface{})
			mergeValue(m, v, data, recursively)
			news = append(news, m)
		default:
			news = append(news, v)
		}
	}
	return news
}

func interfaceMapToStringMap(src interface{}) interface{} {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case map[interface{}]interface{}:
		dst := make(map[string]interface{}, len(src))
		for k, v := range src {
			dst[fmt.Sprint(k)] = interfaceMapToStringMap(v)
		}
		return dst

	case map[string]interface{}:
		dst := make(map[string]interface{}, len(src))
		for k, v := range src {
			dst[k] = interfaceMapToStringMap(v)
		}
		return dst

	case []interface{}:
		a := make([]interface{}, 0, len(src))
		for _, v := range src {
			a = append(a, interfaceMapToStringMap(v))
		}
		return a

	default:
		return src
	}
}

func findFunction(s string) (function Function, arg string, found bool) {
	if len(s) < 3 {
		return
	}
	if s[len(s)-1] != ')' {
		return
	}
	begin := strings.IndexByte(s, '(')
	if begin <= 0 {
		return
	}

	name := strings.ToLower(s[:begin]) // 函数名不区分大小写
	function, found = registeredFunctions[name]
	if !found {
		return
	}
	arg = s[begin+1 : len(s)-1]
	return
}

func (f *Filter) doFilter(data map[string]string, resp map[string]interface{}, mergeValueRecursively bool) (result bool) {
	mergeValue(resp, f.ResponseAlways, data, mergeValueRecursively)
	defer func() {
		if !result {
			mergeValue(resp, f.ResponseOnNotMatch, data, mergeValueRecursively)
		}
	}()

	if f.Match != nil {
		if !f.Match.Match(data) {
			return false
		}
	}

	if f.NotMatch != nil {
		if f.NotMatch.doFilter(data, resp, mergeValueRecursively) {
			return false
		}
	}

	if len(f.MatchAll) > 0 {
		for _, sub := range f.MatchAll {
			if sub == nil {
				continue
			}
			if !sub.doFilter(data, resp, mergeValueRecursively) {
				return false
			}
		}
	}

	if len(f.MatchAny) > 0 {
		pass := false
		for _, sub := range f.MatchAny {
			if sub == nil {
				continue
			}
			if sub.doFilter(data, resp, mergeValueRecursively) {
				pass = true
				break
			}
		}
		if !pass {
			return false
		}
	}
	mergeValue(resp, f.ResponseOnMatch, data, mergeValueRecursively)

	return true
}

// String 返回json格式字符串
func (f *Filter) String() string {
	bs, err := json.Marshal(f)
	if err != nil {
		fmt.Printf("filter:%#v, json.Marshal err:%v\n", f, err)
	}
	return string(bs)
}

// ModifyValues 修改底层Expr中的参数, 一次性修改之后都可以使用
// 部分修改失败时不会回退
func (f *Filter) ModifyValues(modifier func(value string) []string) error {
	if f == nil {
		return nil
	}

	if f.Match != nil {
		if err := f.Match.ModifyValues(modifier); err != nil {
			return err
		}
	}
	if f.NotMatch != nil {
		if err := f.NotMatch.ModifyValues(modifier); err != nil {
			return err
		}
	}
	for _, v := range f.MatchAll {
		if err := v.ModifyValues(modifier); err != nil {
			return err
		}
	}
	for _, v := range f.MatchAny {
		if err := v.ModifyValues(modifier); err != nil {
			return err
		}
	}
	return nil
}

// ModifySource 修改底层Expr，由外部提供参数来源
func (f *Filter) ModifySource(modifier func(value string) func(string) interface{}) error {
	if f == nil {
		return nil
	}

	if f.Match != nil {
		if err := f.Match.ModifyDataSource(modifier); err != nil {
			return err
		}
	}
	if f.NotMatch != nil {
		if err := f.NotMatch.ModifySource(modifier); err != nil {
			return err
		}
	}
	for _, v := range f.MatchAll {
		if err := v.ModifySource(modifier); err != nil {
			return err
		}
	}
	for _, v := range f.MatchAny {
		if err := v.ModifySource(modifier); err != nil {
			return err
		}
	}
	return nil
}
