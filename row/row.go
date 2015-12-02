package row

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Row is a slice of strings.
type Row []string

// ColumnCap holds the number of maximum runs for each column.
type ColumnCap []int

const itemCountNotEqual = `Number of items in the head and body []string must be
equal if there are more than 0 of either of each.`

func NewColumnCap(headRowLen, bodyRowLen int) ColumnCap {
	if headRowLen > 0 {
		return make(ColumnCap, headRowLen)
	}
	return make(ColumnCap, bodyRowLen)
}

// New Row object. Line break will be striped from the strings
// and a whitespace added to the end.
func New(arr []string) Row {
	size := len(arr)
	r := make(Row, size)
	for i, s := range arr {
		s = strings.Replace(s, "\n", "", -1)
		s = strings.Replace(s, "\r", "", -1)
		// whitespace, ignore last column
		if i < size-1 {
			r[i] = fmt.Sprintf("%s ", s)
		} else {
			r[i] = s
		}
	}
	return r
}

func BodyToRow(
	body [][]string,
	bodyLen, headRowLen, bodyRowLen int,
	columnCap ColumnCap,
) ([]Row, error) {
	arr := make([]Row, bodyLen)

	for i := 0; i < bodyLen; i++ {
		r := New(body[i])
		if err := r.CheckBodyRow(headRowLen, bodyRowLen); err != nil {
			return nil, fmt.Errorf(itemCountNotEqual)
		}
		columnCap = r.MaxRuneCount(columnCap)
		arr[i] = r
	}
	return arr, nil
}

func (r Row) CheckBodyRow(headRowLen, bodyRowLen int) error {
	// check if row has same amount of items than head
	if headRowLen > 0 && headRowLen != bodyRowLen {
		return fmt.Errorf(itemCountNotEqual)
	}
	// check if row has equal itmes
	if bodyRowLen != len(r) {
		return fmt.Errorf(itemCountNotEqual)
	}
	return nil
}

func (r Row) MaxRuneCount(arr ColumnCap) ColumnCap {
	var count int
	for i, cell := range r {
		count = utf8.RuneCountInString(cell)
		if count > arr[i] {
			arr[i] = count
		}
	}
	return arr
}

func CheckElemets(headRowLen, bodyRowLen, bodyLen int) error {
	if headRowLen > 0 && bodyLen > 0 && headRowLen != bodyRowLen {
		return fmt.Errorf(itemCountNotEqual)
	}
	return nil
}

// TrimTextToMaxLength enlarges too short stings.
func TrimTextToMaxLength(s string, n int) string {
	size := utf8.RuneCountInString(s)
	if size == n {
		return s
	}
	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", n-size))
}
