package main

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ReadExcel(filename string) (err error) {
	var file *excelize.File
	file, err = excelize.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("unable to open excel file: %w", err)
	}
	defer file.Close()

	var currentSheet string
	for _, sheetName := range file.GetSheetList() {
		if file.GetSheetVisible(sheetName) {
			currentSheet = sheetName
		}
	}
	if len(currentSheet) == 0 {
		return fmt.Errorf("no visible worksheets found in %s", filename)
	}

	var rows *excelize.Rows
	rows, err = file.Rows(currentSheet)
	if err != nil {
		return fmt.Errorf("could not iterate over rows in %s [%s]: %w", currentSheet, filename, err)
	}

	offset := NewOffset()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("could not read row values: %w", err)
		}
		if !offset.Valid() {
			offset, err = headerOffset(row)
			if err != nil {
				return err
			}
			continue
		}
		fmt.Println(row)
	}
	fmt.Printf("header: %v", offset)

	return nil
}

type offset struct {
	ID           int `header:"Код"`
	Description  int `header:"Действие"`
	Dependencies int `header:"Предшествующие действия"`
	Duration     int `header:"Прогноз длительности"`
}

func (o *offset) Valid() bool {
	return o.ID >= 0 && o.Description >= 0 && o.Dependencies >= 0 && o.Duration >= 0
}

func NewOffset() offset {
	return offset{-1, -1, -1, -1}
}

func headerOffset(row []string) (o offset, err error) {
	o = NewOffset()
	var header string

	element := reflect.ValueOf(&o).Elem()
	structure := reflect.ValueOf(o).Type()

	for index, value := range row {
		value = strings.TrimSpace(value)
		for i := 0; i < structure.NumField(); i++ {
			field := element.Field(i)
			header = structure.Field(i).Tag.Get("header")
			if len(header) == 0 {
				continue
			}
			if field.IsValid() && field.CanSet() && strings.HasPrefix(value, header) {
				field.SetInt(int64(index))
				continue
			}
		}
	}
	if !o.Valid() {
		return NewOffset(), fmt.Errorf("could not parse header %v->%v", row, o)
	}
	return o, nil
}
