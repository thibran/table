package main

import (
	"bytes"
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
	head               row.Row
	body               []row.Row
	columnCap          row.ColumnCap // max characters in column
}

// New Table object.
func New(head []string, body ...[]string) (*Table, error) {
	headRow := row.New(head)
	bodyLen := len(body)
	headRowLen := len(headRow)
	var bodyRowLen int
	if bodyLen > 0 {
		bodyRowLen = len(body[0])
	}
	if err := row.CheckElemets(headRowLen, bodyRowLen, bodyLen); err != nil {
		return nil, err
	}
	// collect maximum rune count per column
	columnCap := row.NewColumnCap(headRowLen, bodyRowLen)
	columnCap = headRow.MaxRuneCount(columnCap)
	arr, err := row.BodyToRow(body, bodyLen, headRowLen, bodyRowLen, columnCap)
	if err != nil {
		return nil, err
	}
	return &Table{
		head:               headRow,
		HeadStyle:          StyleSquare(),
		body:               arr,
		BodyStyle:          StyleEmpty(),
		columnCap:          columnCap,
		TopLine:            true,
		HeadOnlyBottomLine: true,
	}, nil
}

func (t *Table) newDraw(w io.Writer) *draw.Draw {
	return &draw.Draw{
		LineHeadTop: t.TopLine && !t.HeadStyle.isAllEmptyH() && !t.HeadOnlyBottomLine,
		LineHeadBot: t.TopLine && !t.HeadStyle.isAllEmptyH(),
		LineHeadV:   t.HeadStyle.vLines && !t.HeadStyle.isAllEmpty(),

		HeadEdge:  t.HeadStyle.edge,
		HeadLineH: t.HeadStyle.lineH,
		HeadLineV: t.HeadStyle.lineV,

		LineBodyTop: t.TopLine && !t.BodyStyle.isAllEmptyH(),
		LineBodyBot: t.BottomLine && !t.BodyStyle.isAllEmptyH(),
		LineBodyV:   t.BodyStyle.vLines && !t.BodyStyle.isAllEmpty(),

		BodyEdge:  t.BodyStyle.edge,
		BodyLineH: t.BodyStyle.lineH,
		BodyLineV: t.BodyStyle.lineV,

		ColumnCap: t.columnCap,
		Writer:    w,
	}
}

func (t *Table) String() string {
	var buf bytes.Buffer
	d := t.newDraw(io.Writer(&buf))
	d.Draw(t.head, t.body...)
	return buf.String()
}

func (t *Table) WriteTo(w io.Writer) {
	d := t.newDraw(w)
	d.Draw(t.head, t.body...)
}
