package matcher

import (
	"encoding/json"
	"strconv"
	"strings"
)

func toNumber(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func compileNumber(args []string) (float64, error) {
	if len(args) != 1 {
		return 0, ErrArgsSize
	}
	return strconv.ParseFloat(args[0], 64)
}

func compileVersion(ver string) []int {
	if len(ver) > 0 && (ver[0] == 'v' || ver[0] == 'V') {
		ver = ver[1:]
	}

	ss := strings.Split(ver, ".")
	a := make([]int, len(ss))

	for i, v := range ss {
		a[i], _ = strconv.Atoi(v)
	}
	return a
}

func compileVersionWithErr(ver string) ([]int, error) {
	ss := strings.Split(ver, ".")
	a := make([]int, len(ss))

	var err error
	for i, v := range ss {
		a[i], err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func versionCompare(v1, v2 []int) int {
	max := len(v1)
	if len(v2) > max {
		max = len(v2)
	}
	for i := 0; i < max; i++ {
		var i1, i2 int
		if i < len(v1) {
			i1 = v1[i]
		}
		if i < len(v2) {
			i2 = v2[i]
		}

		if i1 > i2 {
			return 1
		}
		if i2 > i1 {
			return -1
		}
	}
	return 0
}

// VersionCompare 版本号解析为数字比较
// v1 > v2 return 1
// v1 = v2 return 0
// v1 < v2 return -1
func VersionCompare(v1, v2 string) int {
	a1 := compileVersion(v1)
	a2 := compileVersion(v2)
	return versionCompare(a1, a2)
}

func toJSON(v interface{}) string {
	bs, _ := json.Marshal(v)
	return string(bs)
}
