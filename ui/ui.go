package ui

import (
	"fmt"
	"strings"

	"github.com/callum-oakley/vee/state"
	"github.com/gdamore/tcell/v2"
	rw "github.com/mattn/go-runewidth"
)

var statusStyle = tcell.StyleDefault.Background(tcell.ColorSilver)

type Renderer struct {
	S      *state.State
	Screen tcell.Screen
	w, h   int
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (r *Renderer) renderText() {
	// Discard lines we definitely won't be rendering.
	height := r.h - 1
	start := max(0, min(r.S.Cursor.Y-(height-1)/2, len(r.S.Text)-height))
	rawLines := r.S.Text[start:min(start+height, len(r.S.Text))]

	// Expand tabs and wrap.
	var lines []string
	var cursor struct{ x, y int }
	for y, line := range rawLines {
		lines = append(lines, "")
		for x, char := range line {
			if x == max(0, r.S.Cursor.X) && y == r.S.Cursor.Y-start {
				cursor.x = rw.StringWidth(lines[len(lines)-1])
				cursor.y = len(lines) - 1
			}
			if rw.StringWidth(lines[len(lines)-1]+string(char)) > r.w {
				lines = append(lines, "")
			}
			if char == '\t' {
				lines[len(lines)-1] +=
					strings.Repeat(
						" ",
						r.S.TabWidth-
							rw.StringWidth(lines[len(lines)-1])%r.S.TabWidth,
					)
			} else {
				lines[len(lines)-1] += string(char)
			}
		}
		if y == r.S.Cursor.Y-start && len(line) == 0 {
			cursor.y = len(lines) - 1
		}
	}

	// Finally discard the lines that didn't make the cut after wrapping.
	start = max(0, min(cursor.y-(height-1)/2, len(lines)-height))
	lines = lines[start:min(start+height, len(lines))]
	cursor.y -= start

	for y, line := range lines {
		puts(r.Screen, tcell.StyleDefault, 0, y, line)
	}
	r.Screen.ShowCursor(cursor.x, cursor.y)
}

func (r *Renderer) renderStatus() {
	left := r.S.FilePath
	right := fmt.Sprintf("%v", r.S.Cursor.Y+1)
	if r.S.Msg != "" {
		right = r.S.Msg
	}
	padding := strings.Repeat(
		" ",
		r.w-rw.StringWidth(left)-rw.StringWidth(right),
	)
	puts(r.Screen, statusStyle, 0, r.h-1, left+padding+right)
}

func (r *Renderer) Render() {
	r.w, r.h = r.Screen.Size()
	r.Screen.Clear()
	r.renderText()
	r.renderStatus()
	r.Screen.Show()
}
