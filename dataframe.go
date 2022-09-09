package tynipandas

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"github.com/noaway/dateparse"
	"github.com/shopspring/decimal"
)

type DataFrame struct {
	ColNames []string
	V        []*Series
	nameMap  map[string]int
}

func (d *DataFrame) AddCol(name string, series *Series) {
	d.nameMap[name] = len(d.ColNames)
	d.ColNames = append(d.ColNames, name)
	d.V = append(d.V, series)
}

func (d *DataFrame) AddColVal(name string, val interface{}) {
	series := d.GetCol(name)
	switch series.DType {
	case Int64:
		n := N{Normal, decimal.NewFromInt(val.(int64))}
		series.Nums = append(series.Nums, n)
	case Float64:
		n := N{Normal, decimal.NewFromFloat(val.(float64))}
		series.Nums = append(series.Nums, n)
	case String:
		s, _ := val.(string)
		series.S = append(series.S, s)
	case Time:
		t, _ := dateparse.ParseLocal(val.(string))
		n := N{Normal, decimal.NewFromInt(t.Unix())}
		series.Nums = append(series.Nums, n)
	}
}

func (d *DataFrame) GetCol(name string) *Series {
	idx := d.nameMap[name]
	return d.V[idx]
}

func (d *DataFrame) FromMaps(maps []map[string]interface{}, names []string) {
	if len(maps) == 0 {
		return
	}

	// 1. get type of each columns from the first map
	if len(names) == 0 {
		for name, value := range maps[0] {

			dType := getValueType(value)

			series := &Series{
				DType: dType,
				Name:  name,
			}

			d.AddCol(name, series)
		}
	}

	d.ColNames = names

	// 2. transform data
	for _, m := range maps {
		for _, name := range d.ColNames {
			v := m[name]
			d.AddColVal(name, v)
		}
	}
}

const (
	MergeTypeInner = "inner"
	MergeTypeOuter = "outer"
	MergeTypeLeft  = "left"
)

const (
	InDf1 = iota + 1
	InDf2
	InDf1_Df2
)

func (d *DataFrame) UniqueMerge(d2 *DataFrame, col string, mergeType string) {
	// check len
	if (len(d.ColNames) == 0) || (len(d2.ColNames) == 0) || d.V[0].L == 0 || d2.V[0].L == 0 {
		return
	}

	// check col unique
	if !d.CheckDistinctCol(col) || !d2.CheckDistinctCol(col) {
		panic(ErrColDulipcatedValue(""))
	}

	// 1. due with columns
	var colsMap = make(map[string]int, len(d.ColNames)+len(d2.ColNames))

	for _, col := range d.ColNames {
		colsMap[col] = InDf1
	}

	for _, col := range d2.ColNames {
		if _, ok := d.nameMap[col]; !ok {
			colsMap[col] = InDf2
		} else {
			colsMap[col] = InDf1_Df2
		}
	}

	var renames = make(map[string]string)
	for col, status := range colsMap {
		if status == 3 {
			renames[col] = col + "_right"
		}
	}

	// 2. due with lines
	type Status struct {
		Code   int
		Df2idx int
	}

	series := d.GetCol(col)
	series2 := d2.GetCol(col)

	var valMap = linkedhashmap.New() // orderMap

	var s1, s2 []string
	switch series.DType {
	case Float64, Int64, Time:
		s1 = series.Format(nil, false)
		s2 = series2.Format(nil, false)
	case String:
		s1 = series.S
		s2 = series2.S
	}

	for _, s := range s1 {
		// col in d1
		if _, exist := valMap.Get(s); !exist {
			valMap.Put(s, Status{Code: InDf1})
		}
	}
	for idx, s := range s2 {
		if statusInf, exist := valMap.Get(s); !exist {
			// col in d2
			valMap.Put(s, Status{Code: InDf2})
		} else {
			status := statusInf.(Status)
			if status.Code == InDf1 {
				// col both in d1 & d2
				status.Code = InDf1_Df2
				status.Df2idx = idx
				valMap.Put(s, status)
			}
		}
	}

	// 3. merge
	switch mergeType {
	case MergeTypeLeft:
		for _, colname := range d2.ColNames {
			series := d2.GetCol(colname)
			switch series.DType {
			case Float64, Int64, Time:
				var ns = make([]N, 0, d.V[0].L)

				it := valMap.Iterator()
				var n N
				for idx := 0; it.Next() && (idx < int(d.V[0].L)); idx++ {
					status := it.Value().(Status)
					switch status.Code {
					case InDf1:
						// set zero
						n.T = Normal
						n.V = decimal.Zero
					case InDf1_Df2:
						// add
						n = series.Nums[status.Df2idx]
						ns = append(ns, n)
					}

					n.SetZero()

					if newName, ok := renames[series.Name]; ok {
						d.AddCol(newName, &Series{DType: series.DType, Name: newName, Nums: ns, L: float64(len(ns))})
					} else {
						d.AddCol(series.Name, &Series{DType: series.DType, Name: series.Name, Nums: ns, L: float64(len(ns))})
					}
				}
			case String:

			}
		}

	case MergeTypeInner:
	case MergeTypeOuter:
	}
}

func (d *DataFrame) CheckDistinctCol(col string) bool {
	var ulen int
	series := d.GetCol(col)
	switch series.DType {
	case Float64, Int64, Time:
		var m = make(map[string]struct{}, series.L)
		for _, n := range series.Nums {
			m[n.V.String()] = struct{}{}
		}
		ulen = len(m)
	case String:
		ulen = len(UniqueArrayString(series.S))
	}
	if ulen != int(series.L) {
		return false
	}
	return true
}

func (d *DataFrame) Sort(col string) {}

func Format() {}

func Round() {}

func FillLost() {}

func (d *DataFrame) Rename(renames map[string]string) {
	for idx := 0; idx < len(d.ColNames); idx++ {
		currentName := d.ColNames[idx]
		if newName, ok := renames[currentName]; ok {
			d.ColNames[idx] = newName
			val := d.nameMap[currentName]
			delete(d.nameMap, currentName)
			d.nameMap[newName] = val
		}
	}
}

func RemainKeys() {}

func CalKeysValue() {}

type Operation string

// numeric operation
const (
	Sum        Operation = "sum"
	Percentage Operation = "percentage"
)

// format operation
const (
	ToPercentage Operation = "to_percentage"
	ToString     Operation = "to_string"
)
