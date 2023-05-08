package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/extrame/xls"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/transform"
)

type StoreSettings struct {
	ConfigFile string `json:"config_file"`
	ExcelFile  string `json:"excel_file"`
}

// storeConfig 与前端数据格式同步，用于验证 config 数据格式是否正确
type storeConfig struct {
	Current string `json:"current"`
	Files   []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		List []struct {
			Condition string `json:"condition"`
			Template  string `json:"template"`
			Color     string `json:"color"`
		} `json:"list"`
	}
}

func (a *App) ReadSettings() (*StoreSettings, error) {
	path, err := a.createHomeDir()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path + string(filepath.Separator) + "settings.json")
	if err != nil || len(b) == 0 {
		a.StoreSettings.ConfigFile = path + string(filepath.Separator) + "config.json"
		a.saveSettings()
		return &a.StoreSettings, nil
	}

	var setttings StoreSettings
	err = json.Unmarshal(b, &setttings)
	if err != nil {
		return nil, err
	}

	a.StoreSettings = setttings
	return &a.StoreSettings, nil
}

func (a *App) ReadExcel(path string) ([][]string, error) {
	if path == "" {
		settings, err := a.ReadSettings()
		if err != nil {
			return nil, err
		}
		path = settings.ExcelFile
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".csv":
		e, _, _, err := a.determineEncodingUtf8OrGBK(file)
		if err != nil {
			return nil, err
		}

		file.Seek(0, 0)
		r := transform.NewReader(file, e.NewDecoder())
		csvReader := csv.NewReader(r)

		csvReader.FieldsPerRecord = -1

		data, err := csvReader.ReadAll()
		return a.appendColumn(data), err
	case ".xls":
		_, name, _, err := a.determineEncodingUtf8OrGBK(file)
		if err != nil {
			return nil, err
		}

		file.Seek(0, 0)

		workbook, err := xls.OpenReader(file, name)
		if err != nil {
			return nil, err
		}

		var data [][]string
		if sheet1 := workbook.GetSheet(0); sheet1 != nil {
			maxRow := sheet1.MaxRow
			for i := 0; i < int(maxRow); i++ {
				var rows []string
				row := sheet1.Row(i)
				for index := row.FirstCol(); index < row.LastCol(); index++ {
					rows = append(rows, row.Col(index))
				}
				data = append(data, rows)
			}
		}
		return a.appendColumn(data), err
	default:
		f, err := excelize.OpenReader(file)
		if err != nil {
			return nil, err
		}
		defer func() {
			// Close the spreadsheet.
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		// Get all the rows in the Sheet1.
		data, err := f.GetRows(f.GetSheetName(0))
		return a.appendColumn(data), err
	}

}

func (a *App) SaveConfig(data string) error {
	settings, err := a.ReadSettings()
	if err != nil {
		return err
	}
	return os.WriteFile(settings.ConfigFile, []byte(data), 0755)
}

func (a *App) ReadConfig(path string) (string, error) {
	if path == "" {
		settings, err := a.ReadSettings()
		if err != nil {
			return "", err
		}
		path = settings.ConfigFile
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "path:%v", path)
	}

	var conf storeConfig
	if len(b) == 0 {
		return "", nil
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		return "", errors.Wrap(err, "Config file exception")
	}
	return string(b), nil
}

func (a *App) createHomeDir() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	path := u.HomeDir + string(filepath.Separator) + a.DataDirName
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func (a *App) saveSettings() error {
	path, err := a.createHomeDir()
	if err != nil {
		return err
	}

	b, _ := json.MarshalIndent(a.StoreSettings, "", "  ")
	return os.WriteFile(path+string(filepath.Separator)+"settings.json", b, 0755)
}

func (a *App) appendColumn(data [][]string) [][]string {
	maxColumn := 0
	for _, v := range data {
		if len(v) > maxColumn {
			maxColumn = len(v)
		}
	}

	for i, v := range data {
		diff := maxColumn - len(v)

		if diff > 0 {
			diffData := make([]string, diff)
			v = append(v, diffData...)
			data[i] = v
		}
	}

	a.exprEnv = make(map[string]interface{})
	if len(data) > 0 {
		for i := 0; i < maxColumn; i++ {
			a.exprEnv[string(rune(i+65))] = data[0][i]
		}
	}
	return data
}
