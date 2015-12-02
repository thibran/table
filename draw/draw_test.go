package draw

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/thibran/table/row"
)

type bufDrawer struct {
	Draw
	*bytes.Buffer
}

func newBufDrawer() *bufDrawer {
	var buf bytes.Buffer
	return &bufDrawer{
		Draw: Draw{
			ColumnCap: []int{3, 3},
			Writer:    &buf,
		},
		Buffer: &buf,
	}
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
	arr := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.body(arr)
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
	arr := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.body(arr)
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
	arr := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.body(arr)
	s := d.String()
	exp := "------\na1 a2 \n------\nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestBody_topLineFrom_LineHeadBotTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.LineBodyTop = true
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	arr := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.body(arr)
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
	arr := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.body(arr)
	s := d.String()
	exp := "Ia1 Ia2 I\nIb1 Ib2 I"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_onlyText(t *testing.T) {
	d := newBufDrawer()
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "h1 h2 \na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_headerTopLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "======\nh1 h2 \na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_headerTopAndBottomLine(t *testing.T) {
	d := newBufDrawer()
	d.LineHeadTop = true
	d.LineHeadBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "======\nh1 h2 \n======\na1 a2 \nb1 b2 "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_bodyTopAndBottomLine(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyTop = true
	d.LineBodyBot = true
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "h1 h2 \n------\na1 a2 \n------\nb1 b2 \n------"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_bodyTopAndBottomLine_vlineBodyTrue(t *testing.T) {
	d := newBufDrawer()
	d.LineBodyTop = true
	d.LineBodyBot = true
	d.LineBodyV = true
	d.LineHeadV = false
	d.HeadEdge = '+'
	d.HeadLineH = '='
	d.HeadLineV = '|'
	d.BodyEdge = '#'
	d.BodyLineH = '-'
	d.BodyLineV = 'I'
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := " h1  h2  \n#---#---#\nIa1 Ia2 I\n#---#---#\nIb1 Ib2 I\n#---#---#"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_headTopAndBottomLine_vlineHeadTrue(t *testing.T) {
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
	head := row.Row{"h1 ", "h2 "}
	body := []row.Row{
		{"a1 ", "a2 "},
		{"b1 ", "b2 "},
	}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "+===+===+\n|h1 |h2 |\n+===+===+\n a1  a2  \n b1  b2  "
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_headVlineFalse_bodyVlineTrue(t *testing.T) {
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
	head := row.Row{"Name:", "Count:"}
	body := []row.Row{{"Kenobi ", "1"}}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "================\n Name:   Count: \n================\n|Kenobi |1     |"
	if s != exp {
		t.Error(err(exp, s))
	}
}

func TestDraw_headVlineTrue_bodyVlineFalse(t *testing.T) {
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
	head := row.Row{"Name:", "Count:"}
	body := []row.Row{{"Kenobi ", "1"}}
	d.Draw.Draw(head, body...)
	s := d.String()
	exp := "|Name:  |Count:|\n----------------\n Kenobi  1      \n----------------"
	if s != exp {
		t.Error(err(exp, s))
	}
}
