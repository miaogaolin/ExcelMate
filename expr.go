package main

import (
	"errors"
	"fmt"
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
				param := fmt.Sprintf("%v", params[0])
				res, err := strconv.ParseInt(param, 10, 64)
				if err == nil {
					return res, nil
				}
				return strconv.ParseFloat(param, 64)
			},
		),
	}
}
