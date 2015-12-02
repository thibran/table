package main

import (
	"bytes"
	"io"
	"strings"
	"sync"
	"testing"
)

func TestTable_defaults(t *testing.T) {
	ta, _ := New(true, [][]string{
		{"Name:", "Count:"},
		{"Obi Wan Kenobi", "1"},
		{"Banana", "80"},
		{"Hinz & Kunz", "2"}}...)
	s := ta.String()
	exp := "Name:          Count: \n======================\nObi Wan Kenobi 1      \nBanana         80     \nHinz & Kunz    2      "
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

func TestTable_allEmptyStyle_vlinesTrue(t *testing.T) {
	ta, _ := New(true, [][]string{
		{"Name:", "Count:"},
		{"Obi Wan", "1"},
		{"Banana", "80"}}...)
	ta.HeadStyle = StyleEmpty()
	ta.BodyStyle = StyleEmpty()
	ta.HeadStyle.vLines = true
	ta.BodyStyle.vLines = true
	s := ta.String()
	exp := "|Name:   |Count: |\n|Obi Wan |1      |\n|Banana  |80     |"
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

func TestTable_multipleStyles(t *testing.T) {
	ta, _ := New(true, [][]string{
		{"Fruits:", "Count:"},
		{"Apple", "4"},
		{"Banana", "25"}}...)
	ta.HeadStyle = NewStyle('o', '=', '|', true)
	ta.BodyStyle = NewStyle('•', '–', '|', true)
	ta.HeadOnlyBottomLine = false
	ta.BottomLine = true
	s := ta.String()
	exp := "o========o=======o\n|Fruits: |Count: |\no========o=======o\n|Apple   |4      |\n•––––––––•–––––––•\n|Banana  |25     |\n•––––––––•–––––––•"
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

func TestTable_headNoLines_bodyWithLines(t *testing.T) {
	rows := [][]string{
		{"h1", "h2"},
		{"a1", "a2"},
		{"b1", "b2"},
	}
	ta, _ := New(true, rows...)
	ta.HeadStyle = NewStyle(' ', ' ', '|', false)
	ta.BodyStyle = NewStyle('#', '-', 'I', true)
	ta.TopLine = true
	ta.HeadOnlyBottomLine = true
	ta.BottomLine = true
	s := ta.String()
	exp := " h1  h2  \n#---#---#\nIa1 Ia2 I\n#---#---#\nIb1 Ib2 I\n#---#---#"
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

func TestTable_empty(t *testing.T) {
	ta, err := New(true, [][]string{{"", ""}, {"", ""}}...)
	if err != nil {
		t.Fail()
	}
	ta.String()
}

func TestWriteTo_default(t *testing.T) {
	ta, _ := New(true, [][]string{
		{"Name:", "Count:"},
		{"Obi Wan Kenobi", "1"},
		{"Banana", "80"},
		{"Hinz & Kunz", "2"}}...)
	var b bytes.Buffer
	ta.WriteTo(io.Writer(&b))
	s := b.String()
	exp := "Name:          Count: \n======================\nObi Wan Kenobi 1      \nBanana         80     \nHinz & Kunz    2      "
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

func TestWriteTo_byteCount(t *testing.T) {
	rows := [][]string{
		{"h1", "h2"},
		{"a1", "a2"},
		{"b1", "b2"},
	}
	ta, _ := New(true, rows...)
	ta.HeadStyle = NewStyle(' ', ' ', '|', false)
	ta.BodyStyle = NewStyle('#', '-', 'I', true)
	ta.TopLine = true
	ta.HeadOnlyBottomLine = true
	ta.BottomLine = true
	var b bytes.Buffer
	res, err := ta.WriteTo(&b)
	if err != nil {
		t.Error(err)
	}
	exp := int64(len(" h1  h2  \n#---#---#\nIa1 Ia2 I\n#---#---#\nIb1 Ib2 I\n#---#---#"))
	if res != exp {
		t.Errorf("written bytes should be %d but are %d", exp, res)
	}
	if int64(b.Len()) != res {
		t.Fail()
	}
}

func TestRaw(t *testing.T) {
	ta, _ := New(true, [][]string{
		{"Name:", "Count:"},
		{"Obi Wan Kenobi", "1"},
		{"Hinz & Kunz", "2"}}...)
	s := ta.String()
	exp := "Name:          Count: \n======================\nObi Wan Kenobi 1      \nHinz & Kunz    2      "
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}

type stub struct {
	i int
	c chan result
	b bytes.Buffer
}

type result struct {
	i int
	v string
}

func (st *stub) Read(p []byte) (int, error) {
	// result is send before the receiver-function obtains the return values
	defer func() {
		st.c <- result{i: st.i, v: st.b.String()}
		st.i++
	}()
	switch st.i {
	case 0:
		return copy(p, "11,22\n"), nil
	case 1:
		return copy(p, "33,44\n"), nil
	default:
		return 0, io.EOF
	}
}

func (st *stub) check(t *testing.T, i int, exp, res string) {
	if res != exp {
		t.Errorf("\n[%d] should be:\n%s\n\nbut is:\n%s", i, exp, res)
	}
}

func TestReadFrom_stream(t *testing.T) {
	st := stub{c: make(chan result)}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		for r := range st.c {
			switch r.i {
			case 1:
				st.check(t, r.i, "11 22 \n======", r.v)
			case 2:
				st.check(t, r.i, "11 22 \n======\n33 44 ", r.v)
			}
			wg.Done()
		}
		close(st.c)
	}()
	ta, _ := ReadFrom(&st, true, []int{2, 2})
	ta.WriteTo(&st.b)
	wg.Wait()
}

func TestReadFrom_compareWithNew(t *testing.T) {
	text := "\"Name:\",\"Count:\"\n\"Obi Wan Kenobi\",\"1\"\n"
	r := strings.NewReader(text)
	t1, _ := ReadFrom(r, true, []int{14, 6})
	t2, _ := New(true, [][]string{{"Name:", "Count:"},
		{"Obi Wan Kenobi", "1"}}...)
	s1 := t1.String()
	s2 := t2.String()
	if s1 == "" {
		t.Error(`s1 == ""`)
	}
	if s1 != s2 {
		t.Fail()
	}
}

func TestReadFrom_cutCell(t *testing.T) {
	text := "\"Name:\",\"Count:\"\n\"Obi Wan Kenobi\",\"1\"\n"
	r := strings.NewReader(text)
	ta, _ := ReadFrom(r, true, []int{2, 2})
	s := ta.String()
	exp := "Na Co \n======\nOb 1  "
	if s != exp {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", exp, s)
	}
}
