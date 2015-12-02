package row

import "testing"

func TestNew_default(t *testing.T) {
	c := ColumnCap{3, 2}
	arr := []string{"abcde", "1234", "zzz"}
	postSpace := uint8(1)
	r := New(c, arr, postSpace)
	if len(r) != 2 {
		t.Fail()
	}
	if r[0] != "ab " {
		t.Errorf("r[0] should be %q but is %q\n", "abc ", r[0])
	}
	if r[1] != "1 " {
		t.Errorf("r[0] should be %q but is %q\n", "12 ", r[1])
	}
}

func TestNew_twoSpaces(t *testing.T) {
	c := ColumnCap{3, 2}
	arr := []string{"abcde", "1234", "zzz"}
	postSpace := uint8(2)
	r := New(c, arr, postSpace)
	if len(r) != 2 {
		t.Fail()
	}
	if r[0] != "a  " {
		t.Errorf("r[0] should be %q but is %q\n", "a  ", r[0])
	}
	if r[1] != "  " {
		t.Errorf("r[0] should be %q but is %q\n", "  ", r[1])
	}
}

func TestNewColumnCap(t *testing.T) {
	postSpace := uint8(1)
	c := NewColumnCap([][]string{{"a", "bä"}}, postSpace)
	if c[0] != 2 {
		t.Errorf("c[0] should be %d but is %d", 2, c[0])
	}
	if c[1] != 3 {
		t.Errorf("c[1] should be %d but is %d", 3, c[0])
	}
}

func TestNewColumnCap_twoRows(t *testing.T) {
	postSpace := uint8(1)
	c := NewColumnCap([][]string{
		{"a", "bb"},
		{"cc", "ddd"},
		{"ee", "fff"},
	}, postSpace)
	if c[0] != 3 {
		t.Errorf("c[0] should be %d but is %d", 3, c[0])
	}
	if c[1] != 4 {
		t.Errorf("c[1] should be %d but is %d", 4, c[1])
	}
}

func TestNewColumnCap_linebreak(t *testing.T) {
	postSpace := uint8(1)
	c := NewColumnCap([][]string{{"a\n", "b"}}, postSpace)
	if c[0] != 2 {
		t.Errorf("c[0] should be %d but is %d", 2, c[0])
	}
	if c[1] != 2 {
		t.Errorf("c[1] should be %d but is %d", 1, c[1])
	}
}

func TestNewColumnCap_carriageReturn(t *testing.T) {
	postSpace := uint8(1)
	c := NewColumnCap([][]string{{"a\r", "b"}}, postSpace)
	if c[0] != 2 {
		t.Errorf("c[0] should be %d but is %d", 2, c[0])
	}
	if c[1] != 2 {
		t.Errorf("c[1] should be %d but is %d", 1, c[1])
	}
}

func TestNewColumnCap_empty(t *testing.T) {
	postSpace := uint8(0)
	c := NewColumnCap([][]string{{""}}, postSpace)
	if c[0] != 1 {
		t.Fail()
	}
}

func TestConvRunesPerColumn(t *testing.T) {
	postSpace := uint8(1)
	c := ConvRunesPerColumn([]int{2, 3}, postSpace)
	if c[0] != 3 {
		t.Fail()
	}
	if c[1] != 4 {
		t.Fail()
	}
}

func TestConvRunesPerColumn_noSpace(t *testing.T) {
	postSpace := uint8(0)
	c := ConvRunesPerColumn([]int{2, 3}, postSpace)
	if c[0] != 2 {
		t.Fail()
	}
	if c[1] != 3 {
		t.Fail()
	}
}

func TestConvRunesPerColumn_empty(t *testing.T) {
	postSpace := uint8(0)
	c := ConvRunesPerColumn([]int{0}, postSpace)
	if c[0] != 1 {
		t.Fail()
	}
}

func TestTrimTextToMaxLength(t *testing.T) {
	s := TrimTextToMaxLength("a", 2)
	if s != "a " {
		t.Fail()
	}
}

func TestTrimCell_noDots(t *testing.T) {
	cell := "abcd"
	max := 3
	s := trimCell(cell, max)
	if s != "abc" {
		t.Fail()
	}
}

func TestTrimCell_returnInput(t *testing.T) {
	cell := "abcd"
	max := 4
	s := trimCell(cell, max)
	if s != "abcd" {
		t.Errorf("should be %q but is %q", "ab...", s)
	}
}

func TestTrimCell_withDots(t *testing.T) {
	cell := "abcde"
	max := 4
	s := trimCell(cell, max)
	if s != "a..." {
		t.Errorf("should be %q but is %q", "a...", s)
	}
}

func TestTrimCell_zero(t *testing.T) {
	cell := "abcde"
	max := 0
	s := trimCell(cell, max)
	exp := ""
	if s != exp {
		t.Errorf("should be %q but is %q", exp, s)
	}
}
