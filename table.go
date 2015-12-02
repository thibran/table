package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/thibran/table/draw"
	"github.com/thibran/table/row"
)

// Table object.
type Table struct {
	HeadStyle          *Style
	BodyStyle          *Style
	TopLine            bool // if true, a line is drawn above the header or first entry
	BottomLine         bool // if true, a line is drawn below the last entry
	HeadOnlyBottomLine bool
	hasHeader          bool
	postfixSpace       uint8         // whitespace after every cell
	columnCap          row.ColumnCap // max characters in column
	r                  *csv.Reader
}

// ReadFrom reader r rows unitl io.EOF. If header is true, the first row is
// treated as header. The number of columns and the column-width-per-rune
// is specified by runesPerColumn. The returned csv.Reader can be used to
// set e.g. the delimiter rune for the passed io.Reader.
func ReadFrom(r io.Reader, header bool, runesPerColumn []int) (*Table, *csv.Reader) {
	postfixSpace := uint8(1)
	c := row.ConvRunesPerColumn(runesPerColumn, postfixSpace)
	rd := csv.NewReader(r)
	return newTable(rd, header, c, postfixSpace), rd
}

// New table from slice of rows.
// If hasHeader is true the first row is treated as header-row.
func New(hasHeader bool, rows ...[]string) (*Table, error) {
	if len(rows) == 0 {
		return new(Table), nil
	}
	var b bytes.Buffer
	for _, row := range rows {
		for i, cell := range row {
			fmt.Fprintf(&b, "%q", cell)
			if i < len(row)-1 {
				b.WriteRune(',')
			}
		}
		b.WriteRune('\n')
	}
	postfixSpace := uint8(1)
	r := csv.NewReader(&b)
	c := row.NewColumnCap(rows, postfixSpace)
	return newTable(r, hasHeader, c, postfixSpace), nil
}

// WriteTo returns the bytes written.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	d := t.newDraw(w)
	return t.drawRow(d)
}

func (t *Table) String() string {
	var buf bytes.Buffer
	d := t.newDraw(io.Writer(&buf))
	if _, err := t.drawRow(d); err != nil {
		panic(err)
	}
	return buf.String()
}

func newTable(
	r *csv.Reader,
	hasHeader bool,
	c row.ColumnCap,
	postfixSpace uint8) *Table {
	return &Table{
		HeadStyle:          StyleSquare(),
		BodyStyle:          StyleEmpty(),
		columnCap:          c,
		TopLine:            true,
		HeadOnlyBottomLine: true,
		postfixSpace:       postfixSpace,
		r:                  r,
		hasHeader:          hasHeader,
	}
}

func (t *Table) newDraw(w io.Writer) *draw.Drawer {
	return &draw.Drawer{
		LineHeadTop: t.TopLine && !t.HeadStyle.isEmptyH() && !t.HeadOnlyBottomLine,
		LineHeadBot: t.TopLine && !t.HeadStyle.isEmptyH(),
		LineHeadV:   t.HeadStyle.vLines && !t.HeadStyle.isAllEmpty(),

		HeadEdge:  t.HeadStyle.edge,
		HeadLineH: t.HeadStyle.lineH,
		HeadLineV: t.HeadStyle.lineV,

		LineBodyTop: t.TopLine && !t.BodyStyle.isEmptyH(),
		LineBodyBot: t.BottomLine && !t.BodyStyle.isEmptyH(),
		LineBodyV:   t.BodyStyle.vLines && !t.BodyStyle.isAllEmpty(),

		BodyEdge:  t.BodyStyle.edge,
		BodyLineH: t.BodyStyle.lineH,
		BodyLineV: t.BodyStyle.lineV,

		ColumnCap: t.columnCap,
		Writer:    w,
	}
}

func (t *Table) drawRow(d *draw.Drawer) (int64, error) {
	var i int
	for {
		b, err := t.r.Read()
		if err == io.EOF {
			d.WriteBottomBodyLine()
			break
		}
		if err != nil {
			return d.BytesWritten, err
		}
		isHeader := i == 0 && t.hasHeader
		firstBodyRow := t.isFirstBodyRow(i)
		row := row.New(t.columnCap, []string(b), t.postfixSpace)
		if i != 0 {
			d.WriteNewline()
		}
		d.Row(row, isHeader, firstBodyRow)
		i++
	}
	return d.BytesWritten, d.Err
}

func (t *Table) isFirstBodyRow(i int) bool {
	if i == 0 && !t.hasHeader {
		return true
	} else if i == 1 && t.hasHeader {
		return true
	}
	return false
}
