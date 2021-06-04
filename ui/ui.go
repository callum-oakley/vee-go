package ui

import (
	"fmt"

	"github.com/callum-oakley/vee/state"
	"github.com/gdamore/tcell/v2"
)

// TODO renderer struct with x, y, h, w, s, screen

func renderChar(screen tcell.Screen, x, y *int, c rune) {
	w, _ := screen.Size()
	screen.SetContent(*x%w, *y, c, nil, tcell.StyleDefault)
	*x++
	if *x == w {
		*y++
	}
}

func renderLine(s state.State, screen tcell.Screen, x, y *int, line int) {
	w, _ := screen.Size()
	for char := range s.Text[line] {
		if char == s.Cursor.X && line == s.Cursor.Y {
			screen.ShowCursor(*x%w, *y)
		}
		if s.Text[line][char] == '\t' {
			renderChar(screen, x, y, ' ')
			for *x%s.TabWidth > 0 {
				renderChar(screen, x, y, ' ')
			}
		} else {
			renderChar(screen, x, y, rune(s.Text[line][char]))
		}
	}
	if s.Cursor.Y == line && len(s.Text[line]) == 0 {
		screen.ShowCursor(0, *y)
	}
}

func renderText(s state.State, screen tcell.Screen) {
	_, h := screen.Size()
	x := 0
	y := 0
	for line := range s.Text {
		if y >= h-1 {
			break
		}
		renderLine(s, screen, &x, &y, line)
		y++
		x = 0
	}
}

var statusStyle = tcell.StyleDefault.Background(tcell.ColorSilver)

func renderStatus(s state.State, screen tcell.Screen) {
	lineNumber := []rune(fmt.Sprintf("%v", s.Cursor.Y+1))
	fileStatus := []rune(s.FilePath)
	w, h := screen.Size()
	for i := 0; i < w; i++ {
		if i < len(fileStatus) {
			screen.SetContent(i, h-1, fileStatus[i], nil, statusStyle)
		} else if i >= w-len(lineNumber) {
			screen.SetContent(
				i,
				h-1,
				lineNumber[len(lineNumber)-(w-i)],
				nil,
				statusStyle,
			)
		} else {
			screen.SetContent(i, h-1, ' ', nil, statusStyle)
		}
	}
}

func Render(s state.State, screen tcell.Screen) {
	screen.Clear()
	renderText(s, screen)
	renderStatus(s, screen)
	screen.Show()
}
