package row

import "testing"

func TestNew(t *testing.T) {
	r := New([]string{"a", "b"})
	if len(r[0]) != 2 {
		t.Fail()
	}
	if len(r[1]) != 1 {
		t.Fail()
	}
}

func TestMaxRuneCount(t *testing.T) {
	arr := make(ColumnCap, 2)
	r := Row{"a", "bä"}
	arr = r.MaxRuneCount(arr)
	if arr[0] != 1 {
		t.Fail()
	}
	if arr[1] != 2 {
		t.Fail()
	}
}

func TestMaxRuneCount_twoRows(t *testing.T) {
	arr := make(ColumnCap, 2)
	r := New([]string{"a", "bb"})
	arr = r.MaxRuneCount(arr)
	r = New([]string{"cc", "ddd"})
	arr = r.MaxRuneCount(arr)
	if arr[0] != 3 {
		t.Fail()
	}
	if arr[1] != 3 {
		t.Fail()
	}
}

func TestMaxRuneCount_linebreak(t *testing.T) {
	arr := make(ColumnCap, 2)
	r := New([]string{"a\n", "b"})
	arr = r.MaxRuneCount(arr)
	if arr[0] != 2 {
		t.Errorf("arr[%d] - '%s' should be %d but is %d", 0, r[0], 2, arr[0])
	}
	if arr[1] != 1 {
		t.Errorf("arr[%d] - '%s' should be %d but is %d", 1, r[1], 1, arr[1])
	}
}

func TestMaxRuneCount_carriageReturn(t *testing.T) {
	arr := make(ColumnCap, 2)
	r := New([]string{"a\r", "b"})
	arr = r.MaxRuneCount(arr)
	if arr[0] != 2 {
		t.Errorf("arr[%d] - '%s' should be %d but is %d", 0, r[0], 2, arr[0])
	}
	if arr[1] != 1 {
		t.Errorf("arr[%d] - '%s' should be %d but is %d", 1, r[1], 1, arr[1])
	}
}

func TestCheckBodyRow(t *testing.T) {
	headElements := 1
	bodyElements := 1
	r := New([]string{"a"})
	if err := r.CheckBodyRow(headElements, bodyElements); err != nil {
		t.Fail()
	}
}

func TestCheckBodyRowNoHeadOkay(t *testing.T) {
	headElements := 0
	bodyElements := 1
	r := New([]string{"a"})
	if err := r.CheckBodyRow(headElements, bodyElements); err != nil {
		t.Fail()
	}
}

func TestCheckHeadNotEqualBody(t *testing.T) {
	headElements := 2
	bodyElements := 1
	r := New([]string{"a"})
	if err := r.CheckBodyRow(headElements, bodyElements); err == nil {
		t.Fail()
	}
}

func TestCheckRowError(t *testing.T) {
	headElements := 1
	bodyElements := 1
	r := New([]string{})
	if err := r.CheckBodyRow(headElements, bodyElements); err == nil {
		t.Fail()
	}
}

func TestTrimTextToMaxLength(t *testing.T) {
	s := TrimTextToMaxLength("a", 2)
	if s != "a " {
		t.Fail()
	}
}
