package ui

import (
	"fmt"
	"strings"

	"github.com/callum-oakley/vee/state"
	"github.com/gdamore/tcell/v2"
	rw "github.com/mattn/go-runewidth"
)

var (
	statusStyle    = tcell.StyleDefault.Background(tcell.ColorSilver)
	selectionStyle = tcell.StyleDefault.Background(tcell.ColorSilver)
)

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

func (r *Renderer) renderText(height int) {
	// Discard lines we definitely won't be rendering.
	start := max(0, min(r.S.Cursor.Y-(height-1)/2, len(r.S.Text)-height))
	rawLines := r.S.Text[start:min(start+height, len(r.S.Text))]

	// Expand tabs and wrap.
	var lines []string
	var cursor struct{ x, y int }
	var anchor struct{ x, y int }
	if r.S.Anchor.Y >= start+height {
		anchor.y = height
	}
	for y, line := range rawLines {
		lines = append(lines, "")
		if len(line) == 0 {
			line = " "
		}
		for x, char := range line {
			if x == max(0, r.S.Cursor.X) && y == r.S.Cursor.Y-start {
				cursor.x = rw.StringWidth(lines[len(lines)-1])
				cursor.y = len(lines) - 1
			}
			if x == max(0, r.S.Anchor.X) && y == r.S.Anchor.Y-start {
				anchor.x = rw.StringWidth(lines[len(lines)-1])
				anchor.y = len(lines) - 1
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
	}

	// Discard the lines that didn't make the cut after wrapping.
	start = max(0, min(cursor.y-(height-1)/2, len(lines)-height))
	lines = lines[start:min(start+height, len(lines))]
	cursor.y -= start
	anchor.y -= start
	r.Screen.ShowCursor(cursor.x, cursor.y)

	if cursor.x == anchor.x && cursor.y == anchor.y {
		for y, line := range lines {
			puts(r.Screen, tcell.StyleDefault, 0, y, line)
		}
		return
	}

	if cursor.y < anchor.y || cursor.y == anchor.y && cursor.x < anchor.x {
		cursor, anchor = anchor, cursor
	}

	for y, line := range lines {
		if y < anchor.y || y > cursor.y {
			puts(r.Screen, tcell.StyleDefault, 0, y, line)
		} else if y == anchor.y && y == cursor.y {
			x := puts(r.Screen, tcell.StyleDefault, 0, y, line[:anchor.x])
			x += puts(r.Screen, selectionStyle, x, y, line[anchor.x:cursor.x+1])
			puts(r.Screen, tcell.StyleDefault, x, y, line[cursor.x+1:])
		} else if y == anchor.y {
			x := puts(r.Screen, tcell.StyleDefault, 0, y, line[:anchor.x])
			puts(r.Screen, selectionStyle, x, y, line[anchor.x:])
		} else if y == cursor.y {
			x := puts(r.Screen, selectionStyle, 0, y, line[:cursor.x+1])
			puts(r.Screen, tcell.StyleDefault, x, y, line[cursor.x+1:])
		} else { // y > anchor.y && y < cursor.y
			puts(r.Screen, selectionStyle, 0, y, line)
		}
	}
}

func padBetween(left, right string, width int) string {
	return left + strings.Repeat(
		" ",
		width-rw.StringWidth(left)-rw.StringWidth(right),
	) + right
}

func (r *Renderer) renderStatus(y int) {
	puts(r.Screen, statusStyle, 0, y, padBetween(
		r.S.FilePath,
		fmt.Sprintf("%v", r.S.Cursor.Y+1),
		r.w,
	))
	puts(r.Screen, tcell.StyleDefault, 0, y+1, r.S.Msg)
}

func (r *Renderer) Render() {
	r.w, r.h = r.Screen.Size()
	r.Screen.Clear()
	r.renderText(r.h - 2)
	r.renderStatus(r.h - 2)
	r.Screen.Show()
}
