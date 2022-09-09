package tynipandas

import (
	"math"

	"github.com/shopspring/decimal"
)

/*
Why to write this package:
Data-api use lots of operator between []map[string]interface{}, means need lots of loop
and duplicated code. I need a library like python pandas to operator data, but only find gota
and its doc is hard to learn and use. So I suppose to write a tiny-pandas just has tiny functions
of pandas in go.
*/

type Series struct {
	DType Type
	Name  string
	Nums  []N
	S     []string
	L     float64
}

type Type string

type NumStatus int

const (
	Normal NumStatus = iota
	Inf
	Nan
)

type N struct {
	T NumStatus
	V decimal.Decimal
}

func (n *N) SetZero() {
	n.T = NumStatus(0)
	n.V = decimal.Zero
}

// ******** num 2 num op ********
func (s *Series) Add(v interface{}) {
	if vf, ok := v.(float64); ok {
		for i := 0; i < int(s.L); i++ {
			s.Nums[i].V = s.Nums[i].V.Add(decimal.NewFromFloat(vf))
		}

	} else if s1, ok := v.(*Series); ok {
		l := int(math.Min(s.L, s1.L))
		for i := 0; i < l; i++ {
			s.Nums[i].V = s.Nums[i].V.Add(s1.Nums[i].V)
		}
	}
}

func (s *Series) Div(v interface{}) {
	if vf, ok := v.(float64); ok {
		for i := 0; i < int(s.L); i++ {
			if vf == 0 {
				if s.Nums[i].V.IsZero() {
					s.Nums[i].T = Nan
				} else {
					s.Nums[i].T = Inf
				}
			}
			s.Nums[i].V = s.Nums[i].V.Div(decimal.NewFromFloat(vf))
		}
	} else if s1, ok := v.(*Series); ok {
		l := int(math.Min(s.L, s1.L))
		for i := 0; i < l; i++ {
			if vf == 0 {
				if s.Nums[i].V.IsZero() {
					s.Nums[i].T = Nan
				} else {
					s.Nums[i].T = Inf
				}
			}
			s.Nums[i].V = s.Nums[i].V.Div(s1.Nums[i].V)
		}
	}
}

// TODO Add more function

// ******** num 2 str op ********
func (s *Series) Format(f func(n decimal.Decimal) string, inplace bool) []string {
	var sl = make([]string, 0, int(s.L))

	if f != nil {
		for i := 0; i < int(s.L); i++ {
			sl = append(sl, f(s.Nums[i].V))
		}
	} else {
		for i := 0; i < int(s.L); i++ {
			sl = append(sl, s.Nums[i].V.String())
		}
	}

	if inplace {
		s.S = sl
	}
	return sl
}

// ******** str 2 num op ********
func (s *Series) Parse(f func(s string) decimal.Decimal) {
	var sl = make([]N, 0, int(s.L))
	for i := 0; i < int(s.L); i++ {
		decNum := f(s.S[i])
		var n = N{
			T: Normal,
			V: decNum,
		}
		sl = append(sl, n)
	}
}
