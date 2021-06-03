package ui

import (
	"fmt"

	"github.com/callum-oakley/vee/state"
	"github.com/gdamore/tcell/v2"
)

func renderChar(
	s state.State,
	screen tcell.Screen,
	cursorX *int,
	x, y *int,
	line int,
	c rune,
) {
	w, _ := screen.Size()
	if *x == s.Cursor.X && line == s.Cursor.Y {
		screen.ShowCursor(*x%w, *y)
		*cursorX = *x
	}
	screen.SetContent(*x%w, *y, c, nil, tcell.StyleDefault)
	*x++
	if *x == w {
		*y++
	}
}

func renderLine(
	s state.State,
	screen tcell.Screen,
	cursorX *int,
	x, y *int,
	line int,
) {
	for _, c := range s.Text[line] {
		if c == '\t' {
			renderChar(s, screen, cursorX, x, y, line, ' ')
			for *x%s.TabWidth > 0 {
				renderChar(s, screen, cursorX, x, y, line, ' ')
			}
		} else {
			renderChar(s, screen, cursorX, x, y, line, c)
		}
	}
	if line == s.Cursor.Y && *cursorX == -1 {
		if *x == 0 {
			*cursorX = 0
		} else {
			*cursorX = *x - 1
		}
		w, _ := screen.Size()
		screen.ShowCursor(*cursorX%w, *y)
	}
}

func renderText(s state.State, screen tcell.Screen, cursorX *int) {
	_, h := screen.Size()
	x := 0
	y := 0
	for line := range s.Text {
		if y >= h-1 {
			break
		}
		renderLine(s, screen, cursorX, &x, &y, line)
		y++
		x = 0
	}
}

var statusStyle = tcell.StyleDefault.Background(tcell.ColorSilver)

func renderStatus(s state.State, screen tcell.Screen, cursorX *int) {
	logicalX := s.LogicalX(*cursorX, s.Cursor.Y)
	var cursorStatus []rune
	if logicalX == *cursorX {
		cursorStatus = []rune(fmt.Sprintf("%v,%v", *cursorX+1, s.Cursor.Y+1))
	} else {
		cursorStatus = []rune(
			fmt.Sprintf("%v-%v,%v", logicalX+1, *cursorX+1, s.Cursor.Y+1),
		)
	}
	fileStatus := []rune(s.FileName)

	w, h := screen.Size()
	for i := 0; i < w; i++ {
		if i < len(fileStatus) {
			screen.SetContent(i, h-1, fileStatus[i], nil, statusStyle)
		} else if i >= w-len(cursorStatus) {
			screen.SetContent(
				i,
				h-1,
				cursorStatus[len(cursorStatus)-(w-i)],
				nil,
				statusStyle,
			)
		} else {
			screen.SetContent(i, h-1, ' ', nil, statusStyle)
		}
	}
}

func Render(s state.State, screen tcell.Screen) {
	cursorX := -1
	screen.Clear()
	renderText(s, screen, &cursorX)
	renderStatus(s, screen, &cursorX)
	screen.Show()
}
