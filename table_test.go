package main

import (
	"fmt"
	"testing"
)

func TestTable_defaults(t *testing.T) {
	ta, _ := New([]string{"Name:", "Count:"}, [][]string{
		{"Obi Wan Kenobi", "1"},
		{"Banana", "80"},
		{"Hinz & Kunz", "2"}}...)
	s := ta.String()
	res := "Name:          Count:\n=====================\nObi Wan Kenobi 1     \nBanana         80    \nHinz & Kunz    2     "
	if s != res {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", res, s)
	}
}

func TestTable_allEmptyStyle_vlinesTrue(t *testing.T) {
	ta, _ := New([]string{"Name:", "Count:"}, [][]string{
		{"Obi Wan", "1"},
		{"Banana", "80"}}...)
	ta.HeadStyle = StyleEmpty()
	ta.BodyStyle = StyleEmpty()
	ta.HeadStyle.vLines = true
	ta.BodyStyle.vLines = true
	s := ta.String()
	res := "|Name:   |Count:|\n|Obi Wan |1     |\n|Banana  |80    |"
	if s != res {
		t.Errorf("\n\nshould be:\n%s\n\nbut is:\n%s", res, s)
	}
}

func TestTable_foo(t *testing.T) {
	ta, _ := New(
		[]string{"Fruits:", "Count:"},
		[][]string{
			{"Apple", "4"},
			{"Banana", "25"}}...)
	ta.HeadStyle = NewStyle('o', '=', '|', true)
	ta.BodyStyle = NewStyle('•', '–', '|', true)
	ta.HeadOnlyBottomLine = false
	ta.BottomLine = true
	fmt.Println(ta.String())
}

// func TestTable_foo(t *testing.T) {
// 	ta, _ := New([]string{"Name:", "Count:"}, [][]string{
// 		{"Obi Wan Kenobi", "1"},
// 		{"Banana", "80"}}...)
// 	ta.HeadStyle = StyleSquare()
// 	ta.BodyStyle = StyleEmpty()
// 	//ta.HeadOnlyBottomLine = false
// 	s := ta.String()
// 	fmt.Println(s)
// }
