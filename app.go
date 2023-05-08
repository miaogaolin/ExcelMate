package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/sprig/v3"
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// App struct
type App struct {
	ctx              context.Context
	DataDirName      string
	StoreSettings    StoreSettings
	conditionProgram map[string]interface{}
	templateError    map[string]error
	exprEnv          map[string]interface{}
}

type File struct {
	Data []byte `json:"data"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

// NewApp creates a new App application struct
func NewApp(dataDirName string) *App {
	a := &App{
		DataDirName:      dataDirName,
		conditionProgram: make(map[string]interface{}),
		templateError:    make(map[string]error),
	}
	return a
}

// Validate 验证条件和数据是否匹配
func (a *App) Validate(data interface{}, conditionExpr string) (bool, error) {
	rowData := a.getExcelRow(data)

	// 没有条件代表，代表所有匹配
	if conditionExpr == "" {
		return true, nil
	}

	if _, ok := a.conditionProgram[conditionExpr]; !ok {
		options := append(ExprOptions(), expr.AsBool())
		// if v, ok := rowData["A"]; ok {
		// 	options = append(options, expr.Env(v))
		// }
		program, err := expr.Compile(conditionExpr, options...)
		if err != nil {
			a.conditionProgram[conditionExpr] = err
			return false, err
		}
		a.conditionProgram[conditionExpr] = program
	}

	if len(rowData) == 0 {
		return false, nil
	}

	val := a.conditionProgram[conditionExpr]
	if err, ok := val.(error); ok {
		return false, err
	}

	output, err := expr.Run(val.(*vm.Program), rowData)
	if err != nil {
		return false, err
	}

	if _, ok := output.(bool); !ok {
		return false, errors.New("result type is bool")
	}

	return output.(bool), nil
}

func (a *App) ValidateExpr(conditionExpr string) (bool, error) {
	// 没有条件代表，代表所有匹配
	if conditionExpr == "" {
		return true, nil
	}

	options := append(ExprOptions(), expr.AsBool())
	if len(a.exprEnv) > 0 {
		options = append(options, expr.Env(a.exprEnv))
	}
	_, err := expr.Compile(conditionExpr, options...)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Template 使用模板渲染数据
func (a *App) Template(data interface{}, tpl string) (string, error) {
	rowData := a.getExcelRow(data)
	if tpl != "" {
		if err := a.templateError[tpl]; err != nil {
			return "", errors.Wrap(err, "template")
		}

		res := bytes.NewBuffer(nil)
		t, err :=
			template.New("base").Funcs(sprig.FuncMap()).Parse(tpl)
		if err != nil {
			a.templateError[tpl] = err
			return "", errors.Wrap(err, "template")
		}
		err = t.Execute(res, rowData)
		if err != nil {
			return "", errors.Wrap(err, "template")
		}
		return res.String(), nil
	}
	return "", nil
}

func (a *App) getExcelRow(data interface{}) map[string]interface{} {
	var excelData []interface{}
	if v, ok := data.([]interface{}); ok {
		excelData = v
	}

	rowData := make(map[string]interface{})
	for i, v := range excelData {
		val := fmt.Sprintf("%v", v)
		key := fmt.Sprintf("%c", 'A'+i)
		if num, err := a.getMoneyNum(val); err == nil {
			rowData[key] = num
		} else {
			v = strings.TrimSpace(val)
			rowData[key] = v
		}
	}
	return rowData
}

func (a *App) determineEncodingUtf8OrGBK(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	rd := bufio.NewReader(r)
	b, err := rd.Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	if name != "utf-8" {
		e = simplifiedchinese.GBK
		name = "gbk"
	}
	return
}

func (a *App) getMoneyNum(money string) (float64, error) {
	money = strings.Trim(money, `"`)

	var sign float64 = 1
	if len(money) >= 1 && money[0] == '-' {
		sign = -1
		money = money[1:]
	}

	// 匹配人民币、美元等等格式：￥123.45
	reCurrency := regexp.MustCompile(`^(\$|¥|€|£)(((\d{1,3}(,\d{3})*)|\d{4,})(\.\d+)?)$`)
	// 匹配普通数字格式：123.45 或 1,234.56
	reNum := regexp.MustCompile(`^((\d{1,3}(,\d{3})*)|\d{4,})(\.\d{1,2})?$`)

	// 按优先级尝试匹配不同的金额格式
	if reCurrency.MatchString(money) {
		m := []rune(money)
		money = string(m[1:])
	}

	if reNum.MatchString(money) {
		numStr := strings.ReplaceAll(money, ",", "") // 去掉逗号分隔符
		num, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, err
		}
		return math.Round(num*sign*100) / 100, nil
	}

	// 如果没有匹配成功，则返回错误
	return 0, errors.New("invalid money format")
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
