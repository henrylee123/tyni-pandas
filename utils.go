package tynipandas

import "github.com/noaway/dateparse"

func copyMapStrInt(m map[string]int) map[string]int {
	res := make(map[string]int, len(m))
	for k, v := range m {
		res[k] = v
	}
	return res
}

// 去重数组
func UniqueArrayString(slices ...[]string) []string {
	m := make(map[string]struct{})
	for _, slice := range slices {
		for _, e := range slice {
			if _, ok := m[e]; !ok {
				m[e] = struct{}{}
			}
		}
	}
	var res = make([]string, 0, len(m))
	for e, _ := range m {
		res = append(res, e)
	}
	return res
}

const (
	Int64   Type = "int64"
	Float64 Type = "float64"
	String  Type = "string"
	Time    Type = "time"
)

func getValueType(value interface{}) string {
	var dType string
	if _, ok := value.(int); ok {
		dType = Int64
	} else if _, ok := value.(float64); ok {
		dType = Float64
	} else if valStr, ok := value.(string); ok {
		if _, err := dateparse.ParseLocal(valStr); err != nil {
			dType = Time
		} else {
			dType = String
		}
	}
	return dType
}
