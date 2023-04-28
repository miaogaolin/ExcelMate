package main

import (
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type excelStat struct {
	Data [][]string `json:"data"`
	Path string     `json:"path"`
}

// OpenExcelFile excel file
func (a *App) OpenExcelFile() (*excelStat, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select excel file",
		Filters: []runtime.FileFilter{
			{DisplayName: "Excel Files", Pattern: "*.csv;*.xlam;*xlsm;*xlsx;*xltx"},
		},
	})
	if err != nil {
		return nil, err
	}

	if path != "" {
		a.StoreSettings.ExcelFile = path
		a.saveSettings()
	}
	data, err := a.ReadExcel(path)
	if err != nil {
		return nil, err
	}
	return &excelStat{
		Data: data,
		Path: path,
	}, nil
}

// OpenConfigFile
func (a *App) OpenConfigFile() ([]string, error) {
	dir := ""
	if a.StoreSettings.ConfigFile != "" {
		dir = filepath.Dir(a.StoreSettings.ConfigFile)
	}
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select config file",
		Filters: []runtime.FileFilter{
			{DisplayName: "Config Files", Pattern: "*.json"},
		},
		DefaultDirectory: dir,
	})
	if err != nil {
		return nil, err
	}

	var data string
	if path != "" {
		data, err = a.ReadConfig(path)
		if err != nil {
			return nil, err
		}

		a.StoreSettings.ConfigFile = path
		err = a.saveSettings()
		if err != nil {
			return nil, err
		}
	}
	return []string{path, data}, nil
}

func (a *App) DialogError(msg string) error {
	_, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    "DialogType",
		Message: msg,
	})
	return err
}
