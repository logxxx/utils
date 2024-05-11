package matcher

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Matcher implements an interface to match src and dst.
type Matcher interface {
	Description() string
	Match(src string) bool
}

type newMatcherFunc func(args []string, datadataSource ...func(string) interface{}) (Matcher, error)

var registeredMatchers = new(sync.Map)

// ErrArgsSize return when args for function is wrong.
var ErrArgsSize = errors.New("error args size")

// ErrArgsType return when args type is wrong.
//var ErrArgsType = errors.New("error args type")

var internalMatcherFunctions = map[string]newMatcherFunc{
	"in":    newIn,
	"notIn": not(newIn),

	"lessThan":       newLessThan,
	"notLessThan":    not(newLessThan),
	"greaterThan":    newGreaterThan,
	"notGreaterThan": not(newGreaterThan),

	"versionLessThan":       newVersionLessThan,
	"versionNotLessThan":    newVersionNotLessThan,
	"versionGreaterThan":    newVersionGreaterThan,
	"versionNotGreaterThan": newVersionNotGreaterThan,

	"regexMatch":    newRegexMatch,
	"regexNotMatch": not(newRegexMatch),

	"hasPrefix":    newStringMatcherFunc("has prefix", strings.HasPrefix),
	"notHasPrefix": not(newStringMatcherFunc("has prefix", strings.HasPrefix)),

	"hasSuffix":    newStringMatcherFunc("has suffix", strings.HasSuffix),
	"notHasSuffix": not(newStringMatcherFunc("has suffix", strings.HasSuffix)),

	"contains":    newStringMatcherFunc("contains", strings.Contains),
	"notContains": not(newStringMatcherFunc("contains", strings.Contains)),

	"strLessThan":       newStringMatcherFunc("string less than", func(src, dst string) bool { return src < dst }),
	"strNotLessThan":    newStringMatcherFunc("string not less than", func(src, dst string) bool { return src >= dst }),
	"strGreaterThan":    newStringMatcherFunc("string greater than", func(src, dst string) bool { return src > dst }),
	"strNotGreaterThan": newStringMatcherFunc("string not greater than", func(src, dst string) bool { return src <= dst }),

	"sizeLessThan":       newFileSizeMatcherFunc("string less than", func(src, dst FileSize) bool { return src < dst }),
	"sizeNotLessThan":    newFileSizeMatcherFunc("string not less than", func(src, dst FileSize) bool { return src >= dst }),
	"sizeGreaterThan":    newFileSizeMatcherFunc("string greater than", func(src, dst FileSize) bool { return src > dst }),
	"sizeNotGreaterThan": newFileSizeMatcherFunc("string not greater than", func(src, dst FileSize) bool { return src <= dst }),
}

func init() {
	for k, v := range internalMatcherFunctions {
		RegisterMatcher(k, v)
	}
	for _, v := range []string{
		"==",
		"!=",
		">",
		">=",
		"<",
		"<=",
	} {
		RegisterMatcher(v, newIntelligentMatcherFunc(v))
	}
}

// RegisterMatcher register matcher, if already exists, over write it.
// name case insensitive.
func RegisterMatcher(name string, f func(args []string, datadataSource ...func(string) interface{}) (Matcher, error)) {
	registeredMatchers.Store(strings.ToLower(name), newMatcherFunc(f))
}

// ------not some------

func not(ori newMatcherFunc) newMatcherFunc {
	return func(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
		m, err := ori(args, dataSource...)
		if err != nil {
			return nil, err
		}
		return notMatcher{m}, nil
	}
}

type notMatcher struct {
	Matcher
}

func (n notMatcher) Description() string {
	return "[NOT] " + n.Matcher.Description()
}

func (n notMatcher) Match(src string) bool {
	return !n.Matcher.Match(src)
}

// ------in------

func newIn(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	// allow empty args
	mp := make(map[string]struct{}, len(args))
	for _, v := range args {
		mp[v] = struct{}{}
	}
	return &in{data: mp, dataSource: dataSource}, nil
}

type in struct {
	data       map[string]struct{}
	dataSource []func(string) interface{}
}

func (m *in) Description() string {
	return "in"
}

func (m *in) Match(src string) bool {
	for _, ds := range m.dataSource {
		data, ok := ds(src).(map[string]struct{})
		if !ok {
			continue
		}
		if _, ok = data[src]; ok {
			return ok
		}
	}
	_, ok := m.data[src]
	return ok
}

// ------lessThan------

func newLessThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	f, err := compileNumber(args)
	if err != nil {
		return nil, err
	}
	return &lessThan{data: f}, nil
}

type lessThan struct {
	data float64
}

func (m *lessThan) Description() string {
	return "less than"
}

func (m *lessThan) Match(src string) bool {
	return toNumber(src) < m.data
}

// ------greaterThan------

func newGreaterThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	f, err := compileNumber(args)
	if err != nil {
		return nil, err
	}
	return &greaterThan{data: f}, nil
}

type greaterThan struct {
	data float64
}

func (m *greaterThan) Description() string {
	return "greater than"
}

func (m *greaterThan) Match(src string) bool {
	return toNumber(src) > m.data
}

// ------versionLessThan------

func newVersionLessThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	if len(args) != 1 {
		return nil, ErrArgsSize
	}
	vs, err := compileVersionWithErr(args[0])
	if err != nil {
		return nil, err
	}
	return &versionLessThan{data: vs}, nil
}

type versionLessThan struct {
	data []int
}

func (m *versionLessThan) Description() string {
	return "version less than"
}

func (m *versionLessThan) Match(src string) bool {
	v1 := compileVersion(src)
	return versionCompare(v1, m.data) < 0
}

// ------versionNotLessThan------

func newVersionNotLessThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	if len(args) != 1 {
		return nil, ErrArgsSize
	}
	vs, err := compileVersionWithErr(args[0])
	if err != nil {
		return nil, err
	}
	return &versionNotLessThan{data: vs}, nil
}

type versionNotLessThan struct {
	data []int
}

func (m *versionNotLessThan) Description() string {
	return "version not less than"
}

func (m *versionNotLessThan) Match(src string) bool {
	v1 := compileVersion(src)
	return versionCompare(v1, m.data) >= 0
}

// ------versionGreaterThan------

func newVersionGreaterThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	if len(args) != 1 {
		return nil, ErrArgsSize
	}
	vs, err := compileVersionWithErr(args[0])
	if err != nil {
		return nil, err
	}
	return &versionGreaterThan{data: vs}, nil
}

type versionGreaterThan struct {
	data []int
}

func (m *versionGreaterThan) Description() string {
	return "version greater than"
}

func (m *versionGreaterThan) Match(src string) bool {
	v1 := compileVersion(src)
	return versionCompare(v1, m.data) > 0
}

// ------versionNotGreaterThan------

func newVersionNotGreaterThan(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	if len(args) != 1 {
		return nil, ErrArgsSize
	}
	vs, err := compileVersionWithErr(args[0])
	if err != nil {
		return nil, err
	}
	return &versionNotGreaterThan{data: vs}, nil
}

type versionNotGreaterThan struct {
	data []int
}

func (m *versionNotGreaterThan) Description() string {
	return "version not greater than"
}

func (m *versionNotGreaterThan) Match(src string) bool {
	v1 := compileVersion(src)
	return versionCompare(v1, m.data) <= 0
}

// ------regexMatch------

func newRegexMatch(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
	if len(args) != 1 {
		return nil, ErrArgsSize
	}
	r, err := regexp.Compile(args[0])
	if err != nil {
		return nil, err
	}
	return &regexMatch{data: r}, nil
}

type regexMatch struct {
	data *regexp.Regexp
}

func (m *regexMatch) Description() string {
	return "regex match"
}

func (m *regexMatch) Match(src string) bool {
	return m.data.MatchString(src)
}

// ------fileSizeMatcher------
type fileSizeMatcher struct {
	dst         FileSize
	description string
	matchFunc   func(src, dst FileSize) bool
}

func (m *fileSizeMatcher) Description() string {
	return m.description
}

func (m *fileSizeMatcher) Match(src string) bool {
	size, _ := NewFileSize(src)
	return m.matchFunc(size, m.dst)
}

func newFileSizeMatcherFunc(desc string, matchFunc func(src, dst FileSize) bool) newMatcherFunc {
	return func(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
		if len(args) != 1 {
			return nil, ErrArgsSize
		}
		dst, err := NewFileSize(args[0])
		if err != nil {
			return nil, err
		}
		return &fileSizeMatcher{
			dst:         dst,
			description: desc,
			matchFunc:   matchFunc,
		}, nil
	}
}

// ------stringMatcher------

type stringMatcher struct {
	desc      string
	dst       string
	matchFunc func(src, dst string) bool
}

func newStringMatcherFunc(desc string, matchFunc func(src, dst string) bool) newMatcherFunc {
	return func(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
		if len(args) != 1 {
			return nil, ErrArgsSize
		}

		return &stringMatcher{
			desc:      desc,
			dst:       args[0],
			matchFunc: matchFunc,
		}, nil
	}
}

func (m *stringMatcher) Description() string {
	return m.desc
}

func (m *stringMatcher) Match(src string) bool {
	return m.matchFunc(src, m.dst)
}

type FileSize int64

func NewFileSize(path string) (FileSize, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return FileSize(stat.Size()), nil
}
