package matcher

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Function 简单的函数支持, 用于各类response取数据
// 错误不返回, 请在函数中处理错误情况!!!
type Function func(old interface{}, arg interface{}) (new interface{})

var registeredFunctions = map[string]Function{
	"append":   funcAppend,
	"tostring": funcToString,
	"tonumber": funcToNumber,
}

// RegisterFunction 注册函数
func RegisterFunction(name string, f Function) error {
	if f == nil {
		return errors.New("function should not be nil")
	}

	name = strings.ToLower(name) // 函数名不区分大小写
	registeredFunctions[name] = f
	return nil
}

func funcAppend(old interface{}, arg interface{}) (new interface{}) {
	switch oldVal := old.(type) {
	case string:
		return fmt.Sprintf("%s%v", oldVal, arg)

	case []interface{}:
		return append(oldVal, arg)

	default:
		// 默认返回arg
		return arg
	}
}

// return string
func funcToString(_ interface{}, arg interface{}) (new interface{}) {
	return fmt.Sprint(arg)
}

// return float64 number
func funcToNumber(_ interface{}, arg interface{}) (new interface{}) {
	switch v := arg.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case bool:
		if v {
			return float64(1)
		} else {
			return float64(0)
		}
	default:
		return float64(0)
	}
}
