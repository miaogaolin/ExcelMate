package main

import (
	"errors"
	"strconv"

	"github.com/antonmedv/expr"
)

func ExprOptions() []expr.Option {
	return []expr.Option{
		expr.Function(
			"number",
			func(params ...any) (any, error) {
				if len(params) == 0 {
					return nil, errors.New("missing param")
				}
				param := params[0].(string)
				res, err := strconv.ParseInt(param, 10, 64)
				if err == nil {
					return res, nil
				}
				return strconv.ParseFloat(param, 64)
			},
		),
	}
}
