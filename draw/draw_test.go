package draw

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/thibran/table/row"
)

type bufDrawer struct {
	Drawer
	*bytes.Buffer
}

func newBufDrawer() *bufDrawer {
	var buf bytes.Buffer
	return &bufDrawer{
		Drawer: Drawer{
			ColumnCap: []int{3, 3},
			Writer:    &buf,
		},
		Buffer: &buf,
	}
}

func (d *bufDrawer) writeAll(rows []row.Row, header bool) {
	var firstBodyRow bool
	for i, r := range rows {
		// check if its the first body row
		if i == 0 && !header {
			firstBodyRow = true
		} else if i == 1 && header {
			firstBodyRow = true
		} else {
			firstBodyRow = false
		}
		d.Row(r, i == 0 && header, firstBodyRow)
		// print linebreak between rows, except the last row
		if i < len(rows)-1 {
			fmt.Fprintln(d)
		}
	}
	// print last line after the last row, if any
	d.WriteBottomBodyLine()
}

func (d *bufDrawer) bodyRowTest(r row.Row, firstRow bool) {
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
	d.bodyRow(r, firstRow, writeHline)
}

func err(exp, res string) error {
	return fmt.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, res)
}

func TestHead_checkByteCount(t *testing.T) {
	d := newBufDrawer()
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	res := d.BytesWritten
	exp := int64(len("h1 h2 "))
	if res != exp {
		t.Errorf("written bytes should be %d but are %d", exp, res)
	}
}

func TestHead_checkByteCount_multiline(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.LineHeadV = true
	d.LineBodyV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	res := d.BytesWritten
	exp := int64(len("|h1 |h2 |\n+===+===+"))
	if res != exp {
		t.Errorf("written bytes should be %d but are %d", exp, res)
	}
}

func TestHead_onlyText(t *testing.T) {
	d := newBufDrawer()
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "h1 h2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead_topLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.HeadLineH = '='
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "======\nh1 h2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead_topLine_vlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.LineHeadV = true
	d.LineBodyV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "+===+===+\n|h1 |h2 |"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead_bottomLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.HeadLineH = '='
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "h1 h2 \n======"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead_bottomLine_vlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.LineHeadV = true
	d.LineBodyV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "|h1 |h2 |\n+===+===+"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead_vlineFalse_bodyVlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadV = false
	d.LineBodyV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := " h1  h2  "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestHead(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.LineHeadBot = true
	d.LineHeadV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"h1 ", "h2 "}
	d.head(r)
	s := d.String()
	exp := "+===+===+\n|h1 |h2 |\n+===+===+"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_onlyText_firstTrue(t *testing.T) {
	d := newBufDrawer()
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, true)
	s := d.String()
	exp := "a1 a2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_onlyText_firstFalse(t *testing.T) {
	d := newBufDrawer()
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, false)
	s := d.String()
	exp := "a1 a2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_topLine_firstFalse(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyTop = true
	d.BodyEdge = '+'
	d.BodyLineH = '='
	d.BodyLineV = '|'
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, false)
	s := d.String()
	exp := "======\na1 a2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_topLine_firstTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.LineBodyTop = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, true)
	s := d.String()
	exp := "a1 a2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_vlineFalse_headVlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadV = true
	d.LineBodyV = false
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, false)
	s := d.String()
	exp := " a1  a2  "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBodyRow_topLine_vlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyTop = true
	d.LineBodyV = true
	d.LineHeadV = true
	d.BodyEdge = '+'
	d.BodyLineH = '='
	d.BodyLineV = '|'
	r := row.Row{"a1 ", "a2 "}
	d.bodyRowTest(r, false)
	s := d.String()
	exp := "+===+===+\n|a1 |a2 |"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_onlyText(t *testing.T) {
	d := newBufDrawer()
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, false)
	s := d.String()
	exp := "a1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_bottomLine(t *testing.T) {
	d := newBufDrawer()
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.LineBodyBot = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, false)
	s := d.String()
	exp := "a1 a2 \nb1 b2 \n------"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_topLine(t *testing.T) {
	d := newBufDrawer()
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.LineBodyTop = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, false)
	s := d.String()
	exp := "------\na1 a2 \n------\nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_firstRowToplineFromBodyStyle(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.LineBodyTop = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, false)
	s := d.String()
	exp := "a1 a2 \n------\nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_topLineFrom_vlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyV = true
	d.LineHeadV = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, false)
	s := d.String()
	exp := "Ia1 Ia2 I\nIb1 Ib2 I"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_onlyText(t *testing.T) {
	d := newBufDrawer()
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "h1 h2 \na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_headerTopLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "======\nh1 h2 \na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_headerTopAndBottomLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.LineHeadBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "======\nh1 h2 \n======\na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_bodyTopAndBottomLine(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyTop = true
	d.LineBodyBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)

	s := d.String()
	exp := "h1 h2 \n------\na1 a2 \n------\nb1 b2 \n------"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_bodyTopAndBottomLine_vlineBodyTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadV = false
	d.LineBodyTop = true
	d.LineBodyBot = true
	d.LineBodyV = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := " h1  h2  \n#---#---#\nIa1 Ia2 I\n#---#---#\nIb1 Ib2 I\n#---#---#"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_headTopAndBottomLine_vlineHeadTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.LineHeadBot = true
	d.LineHeadV = true
	d.LineBodyV = false
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	rows := []row.Row{
		{"h1 ", "h2 "},
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "+===+===+\n|h1 |h2 |\n+===+===+\n a1  a2  \n b1  b2  "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_headVlineFalse_bodyVlineTrue(t *testing.T) {
	d := newBufDrawer()
	d.ColumnCap = row.ColumnCap{7, 6}
	d.LineHeadTop = true
	d.LineHeadBot = true
	d.LineBodyV = true
	d.HeadEdge = '#'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = 'o'
	d.BodyLineH = '-'
	d.BodyLineV = '|'
	rows := []row.Row{
		{"Name:", "Count:"},
		{"Kenobi ", "1"},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "================\n Name:   Count: \n================\n|Kenobi |1     |"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestWriteAll_headVlineTrue_bodyVlineFalse(t *testing.T) {
	d := newBufDrawer()
	d.ColumnCap = row.ColumnCap{7, 6}
	d.LineBodyTop = true
	d.LineBodyBot = true
	d.LineHeadV = true
	d.HeadEdge = '#'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = 'o'
	d.BodyLineH = '-'
	d.BodyLineV = '|'
	rows := []row.Row{
		{"Name:", "Count:"},
		{"Kenobi ", "1"},
	}
	d.writeAll(rows, true)
	s := d.String()
	exp := "|Name:  |Count:|\n----------------\n Kenobi  1      \n----------------"
	if s != exp {
		t.Error(err(exp, s))
	}
}
