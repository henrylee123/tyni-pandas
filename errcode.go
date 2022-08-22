package tynipandas

import "strconv"

type TyniPandasErr struct {
	Code   int
	Desc   string
	Detail string
}

func (e TyniPandasErr) Error() string {
	errStr := strconv.Itoa(e.Code) + ": " + e.Desc
	if e.Detail != "" {
		errStr = errStr + "\n" + e.Detail
	}
	return errStr
}

func ErrColDulipcatedValue(detail string) error {
	return TyniPandasErr{
		Code:   10001,
		Desc:   "dataframe column has dulipcated values",
		Detail: detail,
	}
}
