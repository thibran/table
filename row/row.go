package row

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Row is a slice of strings.
type Row []string

// ColumnCap holds the number of maximum runs for each column.
type ColumnCap []int

var errItemCountNotEqual = errors.New(`Number of items in the head and body []string must be
equal if there are more than 0 of either of each.`)

// New Row object. Line break will be striped from the strings
// and n whitespace added to the end of the cell.
func New(c ColumnCap, row []string, n uint8) Row {
	// trim too long rows
	if len(row) > len(c) {
		row = row[:len(c)]
	}
	for i, cell := range row {
		cell = purgeRunes(cell)
		cell = trimCell(cell, c[i]-int(n))
		row[i] = fmt.Sprintf("%s%s", cell, strings.Repeat(" ", int(n)))
	}
	return Row(row)
}

// trimCell if it is longer than max runes.
func trimCell(cell string, max int) string {
	if max < 0 {
		max = 0
	}
	count := utf8.RuneCountInString(cell)
	switch {
	case count <= max: // do nothing
		return cell
	case max <= 3: // trim without adding dots
		return cell[:max]
	default: // trim + dots
		return fmt.Sprintf("%s...", cell[:max-3])
	}
}

// NewColumnCap calculate ColumnCap for the rows with n whitespace added.
func NewColumnCap(rows [][]string, n uint8) ColumnCap {
	c := make(ColumnCap, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			cell = purgeRunes(cell)
			count := utf8.RuneCountInString(cell) + int(n)
			if count > c[i] {
				c[i] = count
			}
			if c[i] == 0 {
				c[i] = 1
			}
		}
	}
	return c
}

// ConvRunesPerColumn calculates ColumnCap (incorporate added whitespace n).
func ConvRunesPerColumn(runesPerColumn []int, n uint8) ColumnCap {
	c := make(ColumnCap, len(runesPerColumn))
	for i, count := range runesPerColumn {
		if count <= 0 {
			count = 0
		}
		c[i] = count + int(n)
		if c[i] == 0 {
			c[i] = 1
		}
	}
	return c
}

func purgeRunes(cell string) string {
	cell = strings.Replace(cell, "\n", "", -1)
	cell = strings.Replace(cell, "\r", "", -1)
	return cell
}

// TrimTextToMaxLength enlarges too short stings.
func TrimTextToMaxLength(s string, n int) string {
	size := utf8.RuneCountInString(s)
	if size == n {
		return s
	}
	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", n-size))
}
