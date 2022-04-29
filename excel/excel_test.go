package excel

import (
	"github.com/miaocansky/go-tool/util/file"
	"testing"
)

func TestCreateExcel(t *testing.T) {
	head := []string{"请选择序号", "批次", "ID"}
	rows := make([][]interface{}, 0, 10)

	row1 := make([]interface{}, 0, 2)
	row1 = append(row1, "序号1", "批次1", 1)
	row2 := make([]interface{}, 0, 2)

	row2 = append(row2, "序号2", "批次2", 2)

	rows = append(rows, row1, row2)

	file, err := ExportExcelFile(head, rows)
	t.Log(file)
	t.Log(err)
}

func TestGetColumnName(t *testing.T) {
	name := GetColumnName(1, 1)
	t.Log(name)
	rowName := GetColumnRowName(name, 2)
	t.Log(rowName)
	s := string(name)
	t.Log(s)
}
func TestReadExcel(t *testing.T) {
	path, _ := file.GetPath()
	path = path + "/upload/excel"

	excelName := "Book3.xlsx"
	path = path + "/" + excelName
	rows, err := ReadExcel(path)
	t.Log(err)
	t.Log(rows)

}
