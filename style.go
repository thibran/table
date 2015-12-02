package main

// Style of the line or edge between vertical & horizontal lines.
type Style struct {
	edge   rune // connection between a vertical & horizontal line
	lineH  rune // vertical line
	lineV  rune // horizontal line
	vLines bool // draw only vertical lines
}

// NewStyle object.
func NewStyle(edge, lineH, lineV rune, vLines bool) *Style {
	return &Style{
		edge:   edge,
		lineH:  lineH,
		lineV:  lineV,
		vLines: vLines,
	}
}

func (s *Style) isAllEmpty() bool {
	return s.edge == ' ' && s.lineH == ' ' && s.lineV == ' '
}

func (s *Style) isEmptyH() bool {
	return s.edge == ' ' && s.lineH == ' '
}

// StyleSquare
// ======
// h1 h2
// ======
//
// #===#===#
// |h1 |h2 |
// #===#===#
func StyleSquare() *Style {
	return NewStyle('#', '=', '|', false)
}

// StyleBoring
// ======
// h1 h2
// ======
//
// =========
// |h1 |h2 |
// =========
func StyleBoring() *Style {
	return NewStyle('=', '=', '|', false)
}

// StyleDot
// •••••••••
// |h1 |h2 |
// •••••••••
//
// ••••••
// h1 h2
// ••••••
func StyleDot() *Style {
	return NewStyle('•', '•', '|', true)
}

// StyleEmpty
// h1 h2
//
// |h1 |h2 |
func StyleEmpty() *Style {
	return NewStyle(' ', ' ', '|', false)
}
