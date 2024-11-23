package matcher

import (
	"fmt"
	"strconv"
)

type intelligentType int

const (
	itString intelligentType = iota
	itNumber
	itVersion
	itFileSize
)

var intelligentMapping = map[intelligentType]map[string]string{
	itNumber: {
		">":  "greaterThan",
		">=": "notLessThan",
		"<":  "lessThan",
		"<=": "notGreaterThan",
	},
	itVersion: {
		">":  "versionGreaterThan",
		">=": "versionNotLessThan",
		"<":  "versionLessThan",
		"<=": "versionNotGreaterThan",
	},
	itFileSize: {
		">":  "sizeGreaterThan",
		">=": "sizeNotLessThan",
		"<":  "sizeLessThan",
		"<=": "sizeNotGreaterThan",
	},
	itString: {
		">":  "strGreaterThan",
		">=": "strNotLessThan",
		"<":  "strLessThan",
		"<=": "strNotGreaterThan",
	},
}

// newIntelligentMatcherFunc 根据待匹配的目标推断数据类型, 传入的目标只能是1个 !!!
// 支持智能推类型:
// ==, != 按照String的in/notIn处理
//
// >, >=, <, <= 支持下以类型：
//
//	Number (int64/float64)
//	Version (目标版本号至少写3段，不足3段一定要补0,否则视为float处理)
//	FileSize
//	其他情况视为String比较
func newIntelligentMatcherFunc(comparisonSymbol string) newMatcherFunc {
	return func(args []string, dataSource ...func(string) interface{}) (Matcher, error) {
		if len(args) != 1 {
			return nil, ErrArgsSize
		}

		switch comparisonSymbol {
		case "==":
			return newIn(args, dataSource...)
		case "!=":
			return not(newIn)(args, dataSource...)
		case ">", ">=", "<", "<=":
		default:
			return nil, fmt.Errorf("invalid comparison symbol:%s", comparisonSymbol)
		}

		s := args[0]
		var err error

		// Number check
		_, err = strconv.ParseFloat(s, 64)
		if err == nil {
			return internalMatcherFunctions[intelligentMapping[itNumber][comparisonSymbol]](args, dataSource...)
		}

		// Version check
		_, err = compileVersionWithErr(s)
		if err == nil {
			return internalMatcherFunctions[intelligentMapping[itVersion][comparisonSymbol]](args, dataSource...)
		}

		// FileSize check
		_, err = NewFileSize(s)
		if err == nil {
			return internalMatcherFunctions[intelligentMapping[itFileSize][comparisonSymbol]](args, dataSource...)
		}

		return internalMatcherFunctions[intelligentMapping[itString][comparisonSymbol]](args, dataSource...)
	}
}
