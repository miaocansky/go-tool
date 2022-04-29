package excel

import (
	"github.com/miaocansky/go-tool/util/file"
	"github.com/xuri/excelize/v2"
	"strconv"
)

// maxCharCount 最多26个字符A-Z
const maxCharCount = 26

func init() {
	createExcelPath()
}

//
//  ExportExcelFile
//  @Description: 导出exclefile文件
//  @param headers 头信息
//  @param rows 数据
//  @return string 文件保存地址
//  @return error 错误信息
//
func ExportExcelFile(headers []string, rows [][]interface{}) (string, error) {
	path, _ := file.GetPath()
	path = path + "/upload/excel"

	excelName := "Book3.xlsx"
	path = path + "/" + excelName
	f, err := ExportExcel("Sheet1", headers, rows)
	if err != nil {
		return "", nil
	}
	err = f.SaveAs(path)
	if err != nil {
		return "", err
	}
	return path, nil

}

//
//  ExportExcel
//  @Description: 导出文件
//  @param sheetName 表名
//  @param headers 头信息
//  @param rows 数据
//  @return *excelize.File
//  @return error
//
func ExportExcel(sheetName string, headers []string, rows [][]interface{}) (*excelize.File, error) {
	f := excelize.NewFile()
	var sheetIndex int

	sheetIndex = f.NewSheet(sheetName)

	maxColumnRowNameLen := 1 + len(strconv.Itoa(len(rows)))
	columnCount := len(headers)
	if columnCount > maxCharCount {
		maxColumnRowNameLen++
	} else if columnCount > maxCharCount*maxCharCount {
		maxColumnRowNameLen += 2
	}
	columnNames := make([][]byte, 0, columnCount)
	for i, header := range headers {
		columnName := GetColumnName(i, maxColumnRowNameLen)
		columnNames = append(columnNames, columnName)
		// 初始化excel表头，这里的index从1开始要注意
		curColumnName := GetColumnRowName(columnName, 1)
		err := f.SetCellValue(sheetName, curColumnName, header)
		if err != nil {
			return nil, err
		}
	}
	for rowIndex, row := range rows {
		for columnIndex, columnName := range columnNames {
			// 从第二行开始
			err := f.SetCellValue(sheetName, GetColumnRowName(columnName, rowIndex+2), row[columnIndex])
			if err != nil {
				return nil, err
			}
		}
	}
	if sheetName != "Sheet1" {
		f.SetActiveSheet(sheetIndex)
	}

	return f, nil
}

func ReadExcel(filePath string) ([][]string, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	//var slices [][]string = make([][]string, 0, 2000)
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}
	rows = rows[1:]
	return rows, nil
}

// getColumnName 生成列名
// Excel的列名规则是从A-Z往后排;超过Z以后用两个字母表示，比如AA,AB,AC;两个字母不够以后用三个字母表示，比如AAA,AAB,AAC
// 这里做数字到列名的映射：0 -> A, 1 -> B, 2 -> C
// maxColumnRowNameLen 表示名称框的最大长度，假设数据是10行，1000列，则最后一个名称框是J1000(如果有表头，则是J1001),是4位
// 这里根据 maxColumnRowNameLen 生成切片，后面生成名称框的时候可以复用这个切片，而无需扩容
//
//  GetColumnName
//  @Description: 生成列名
//  @param column
//  @param maxColumnRowNameLen
//  @return []byte
//
func GetColumnName(column, maxColumnRowNameLen int) []byte {
	const A = 'A'
	if column < maxCharCount {
		// 第一次就分配好切片的容量
		slice := make([]byte, 0, maxColumnRowNameLen)
		return append(slice, byte(A+column))
	} else {
		// 递归生成类似AA,AB,AAA,AAB这种形式的列名
		return append(GetColumnName(column/maxCharCount-1, maxColumnRowNameLen), byte(A+column%maxCharCount))
	}
}

//
//  GetColumnRowName
//  @Description: 生成单元格名称 Excel的名称框是用A1,A2,B1,B2来表示的，这里需要传入前一步生成的列名切片，然后直接加上行索引来生成名称框，就无需每次分配内存
//  @param columnName 单元格所在列名称
//  @param rowIndex 行数
//  @return columnRowName
//
func GetColumnRowName(columnName []byte, rowIndex int) (columnRowName string) {
	l := len(columnName)
	columnName = strconv.AppendInt(columnName, int64(rowIndex), 10)
	columnRowName = string(columnName)
	// 将列名恢复回去
	columnName = columnName[:l]
	return
}

//
//  createExcelPath
//  @Description: 创建地址
//  @return error
//
func createExcelPath() error {
	path, _ := file.GetPath()
	path = path + "/upload/excel"
	exists, _ := file.PathExists(path)
	if !exists {
		err := file.MakeDir(path)
		return err
	}
	return nil

}
