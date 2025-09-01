package service

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func NewExcelService(templatesDir string) *ExcelService {
	return &ExcelService{templatesDir: templatesDir}
}

func (s *ExcelService) RenderExcel(code string, data map[string]interface{}) ([]byte, error) {
	path := fmt.Sprintf("%s/excel/%s.xlsx", s.templatesDir, code)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("open template: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheets := f.GetSheetList()

	for _, sheet := range sheets {
		if err := replaceScalarsInSheet(f, sheet, data); err != nil {
			return nil, err
		}
	}

	for _, sheet := range sheets {
		if err := expandArrayOnSheet(f, sheet, "transactions", data); err != nil {
			return nil, err
		}
		if err := expandArrayOnSheet(f, sheet, "waitTransactions", data); err != nil {
			return nil, err
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func replaceScalarsInSheet(f *excelize.File, sheet string, data map[string]interface{}) error {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}
	for rIdx, row := range rows {
		for cIdx, cell := range row {
			if !strings.Contains(cell, "{{") {
				continue
			}
			newText := cell
			for key, val := range data {

				if strings.Contains(key, ".") {
					continue
				}
				switch val.(type) {
				case []interface{}:
					continue
				}
				ph := "{{" + key + "}}"
				if strings.Contains(newText, ph) {
					newText = strings.ReplaceAll(newText, ph, fmt.Sprint(val))
				}
			}
			if newText != cell {
				addr, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
				if err := f.SetCellStr(sheet, addr, newText); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func expandArrayOnSheet(f *excelize.File, sheet, arrayKey string, data map[string]interface{}) error {
	rows, err := f.GetRows(sheet)
	if err != nil {
		return err
	}

	tplRow0 := -1
	var cols []int
	var fields []string
	var styles []int

	markerPrefix := "{{" + arrayKey + "."

	for rIdx, row := range rows {
		for cIdx, cell := range row {
			if !strings.Contains(cell, markerPrefix) {
				continue
			}
			if tplRow0 == -1 {
				tplRow0 = rIdx
			}
			if tplRow0 == rIdx {
				field := extractArrayField(cell, arrayKey)
				if field == "" {
					continue
				}
				cols = append(cols, cIdx)
				fields = append(fields, field)
				addr, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
				stID, _ := f.GetCellStyle(sheet, addr)
				styles = append(styles, stID)
			}
		}
	}
	if tplRow0 == -1 {
		return nil
	}

	raw := data[arrayKey]
	arr, ok := raw.([]interface{})
	if !ok || len(arr) == 0 {

		for _, c := range cols {
			addr, _ := excelize.CoordinatesToCellName(c+1, tplRow0+1)
			_ = f.SetCellStr(sheet, addr, "")
		}
		return nil
	}

	n := len(arr)
	if n > 1 {
		if err := f.InsertRows(sheet, tplRow0+2, n-1); err != nil {
			return err
		}
	}

	for i := 0; i < n; i++ {
		row1 := tplRow0 + 1 + i // 1-based
		item, ok := arr[i].(map[string]interface{})
		if !ok {
			return fmt.Errorf("%s[%d] must be an object", arrayKey, i)
		}

		for j, c := range cols {
			addr, _ := excelize.CoordinatesToCellName(c+1, row1)
			val := fmt.Sprint(item[fields[j]])
			if err := f.SetCellStr(sheet, addr, val); err != nil {
				return err
			}
			if styles[j] != 0 {
				_ = f.SetCellStyle(sheet, addr, addr, styles[j])
			}
		}
	}

	return nil
}

func extractArrayField(cell, arrayKey string) string {
	prefix := "{{" + arrayKey + "."
	i := strings.Index(cell, prefix)
	if i == -1 {
		return ""
	}
	j := strings.Index(cell[i+len(prefix):], "}}")
	if j == -1 {
		return ""
	}
	return cell[i+len(prefix) : i+len(prefix)+j]
}
