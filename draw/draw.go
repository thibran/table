package draw

import (
	"io"
	"unicode/utf8"

	"github.com/thibran/table/row"
)

// Drawer objct.
type Drawer struct {
	LineHeadTop bool
	LineHeadBot bool
	LineHeadV   bool

	LineBodyTop bool
	LineBodyBot bool
	LineBodyV   bool

	HeadEdge  rune
	HeadLineH rune
	HeadLineV rune

	BodyEdge     rune
	BodyLineH    rune
	BodyLineV    rune
	ColumnCap    row.ColumnCap
	BytesWritten int64
	Err          error
	io.Writer
}

// Row writes a single row into the io.Writer.
func (d *Drawer) Row(r row.Row, isHeader, firstBodyRow bool) {
	if isHeader {
		d.head(r)
		return
	}
	if len(r) > 0 {
		d.body(r, firstBodyRow)
	}
	return
}

func (d *Drawer) head(r row.Row) {
	edge := d.edgeRune(false)
	hline := d.hlineRune(false)
	vLine := d.isLineV(false)
	opositVlineTrue := d.isOpositVlineTrue(false)
	writeEdge := func() {
		d.writeEdge(vLine, opositVlineTrue, edge, hline)
	}
	// top line
	if d.LineHeadTop {
		d.lineH(hline, writeEdge)
		d.writeRune('\n')
	}
	// header
	d.writeRow(r, false)
	// bottom line
	if d.LineHeadBot {
		d.writeRune('\n')
		d.lineH(hline, writeEdge)
	}
}

func (d *Drawer) body(r row.Row, firstBodyRow bool) {
	edge := d.edgeRune(true)
	hline := d.hlineRune(true)
	vLine := d.isLineV(true)
	opositVlineTrue := d.isOpositVlineTrue(true)
	writeEdge := func() {
		d.writeEdge(vLine, opositVlineTrue, edge, hline)
	}
	writeHline := func() {
		d.lineH(hline, writeEdge)
	}
	d.bodyRow(r, firstBodyRow, writeHline)
}

func (d *Drawer) WriteBottomBodyLine() {
	edge := d.edgeRune(true)
	hline := d.hlineRune(true)
	vLine := d.isLineV(true)
	opositVlineTrue := d.isOpositVlineTrue(true)
	writeEdge := func() {
		d.writeEdge(vLine, opositVlineTrue, edge, hline)
	}
	if d.LineBodyBot {
		d.writeRune('\n')
		d.lineH(hline, writeEdge)
	}
}

func (d *Drawer) WriteNewline() {
	d.writeRune('\n')
}

func (d *Drawer) bodyRow(r row.Row, firstRow bool, writeHline func()) {
	// dont print a topline if there is already one
	printFirstRow := firstRow && !d.LineHeadBot && d.LineBodyTop
	printNonFirstRow := !firstRow && d.LineBodyTop
	if printFirstRow || printNonFirstRow {
		writeHline()
		d.writeRune('\n')
	}
	d.writeRow(r, true)
}

func (d *Drawer) writeRow(r row.Row, isBody bool) {
	d.lineV(isBody)
	for i, cell := range r {
		cell = row.TrimTextToMaxLength(cell, d.ColumnCap[i])
		d.writeString(cell)
		d.lineV(isBody)
	}
}

// writeEdge rune if either head or body vline is true (rune differs).
// In the case that both are false, nothing is written.
func (d *Drawer) writeEdge(vLine, opositVlineTrue bool, edge, hline rune) {
	if vLine {
		d.writeRune(edge)
	} else if opositVlineTrue {
		d.writeRune(hline)
	}
}

func (d *Drawer) lineH(hline rune, writeEdge func()) {
	// edge or hline rune, or nothing if vline head & body are false
	writeEdge()
	for _, count := range d.ColumnCap {
		for i := 0; i < count; i++ {
			d.writeRune(hline)
		}
		writeEdge()
	}
}

// lineV prints the vline rune or a space if.
func (d *Drawer) lineV(isBody bool) {
	if d.isLineV(isBody) {
		r := d.vlineRune(isBody)
		d.writeRune(r)
	} else if d.isOpositVlineTrue(isBody) {
		d.writeRune(' ')
	}
}

// isLineV returns true, if a vertical line should be drawn.
func (d *Drawer) isLineV(isBody bool) bool {
	if isBody {
		return d.LineBodyV
	}
	return d.LineHeadV
}

// edgeRune returns the head or body edge rune.
func (d *Drawer) edgeRune(isBody bool) rune {
	if isBody {
		return d.BodyEdge
	}
	return d.HeadEdge
}

// hlineRune returns the head or body hline rune.
func (d *Drawer) hlineRune(isBody bool) rune {
	if isBody {
		return d.BodyLineH
	}
	return d.HeadLineH
}

// vlineRune returns the head or body vline rune.
func (d *Drawer) vlineRune(isBody bool) rune {
	if isBody {
		return d.BodyLineV
	}
	return d.HeadLineV
}

// isOpositVlineTrue true, e.g. if HeadLineV True & BodyLineV False
func (d *Drawer) isOpositVlineTrue(isBody bool) bool {
	if isBody {
		return d.LineHeadV == true
	}
	return d.LineBodyV == true
}

func (d *Drawer) writeString(s string) {
	d.addBytesWritten(io.WriteString(d, s))
}

func (d *Drawer) writeRune(r rune) {
	n := utf8.RuneLen(r)
	buf := make([]byte, n)
	utf8.EncodeRune(buf, r)
	d.addBytesWritten(d.Write(buf))
}

func (d *Drawer) addBytesWritten(n int, err error) {
	d.BytesWritten += int64(n)
	if err != nil {
		d.Err = err
	}
}

// func (d *Drawer) Debug() string {
// 	var buf bytes.Buffer
// 	buf.WriteString(fmt.Sprintf("LineHeadTop: %v\n", d.LineHeadTop))
// 	buf.WriteString(fmt.Sprintf("LineHeadBot: %v\n", d.LineHeadBot))
// 	buf.WriteString(fmt.Sprintf("LineHeadV:   %v\n\n", d.LineHeadV))
//
// 	buf.WriteString(fmt.Sprintf("LineBodyTop: %v\n", d.LineBodyTop))
// 	buf.WriteString(fmt.Sprintf("LineBodyBot: %v\n", d.LineBodyBot))
// 	buf.WriteString(fmt.Sprintf("LineBodyV:   %v\n\n", d.LineBodyV))
//
// 	buf.WriteString(fmt.Sprintf("HeadEdge:    '%s'\n", string(d.HeadEdge)))
// 	buf.WriteString(fmt.Sprintf("HeadLineH:   '%s'\n", string(d.HeadLineH)))
// 	buf.WriteString(fmt.Sprintf("HeadLineV:   '%s'\n\n", string(d.HeadLineV)))
//
// 	buf.WriteString(fmt.Sprintf("BodyEdge:    '%s'\n", string(d.BodyEdge)))
// 	buf.WriteString(fmt.Sprintf("BodyLineH:   '%s'\n", string(d.BodyLineH)))
// 	buf.WriteString(fmt.Sprintf("BodyLineV:   '%s'\n\n", string(d.BodyLineV)))
// 	buf.WriteString(fmt.Sprintf("ColumnCap:   %v", d.ColumnCap))
// 	return buf.String()
// }
