package state

import (
	"strings"

	"github.com/atotto/clipboard"
)

func (s *State) copy() {
	from, to := s.normalisedSelection()
	selectedText := ""
	for y := from.Y; y <= to.Y; y++ {
		if y == from.Y && from.X == -1 || y == to.Y && to.X == -1 {
			selectedText += "\n"
		} else if y == from.Y && y == to.Y {
			selectedText += s.Text[y][from.X:s.xRightOf(&to)]
		} else if y == from.Y {
			selectedText += s.Text[y][from.X:] + "\n"
		} else if y == to.Y {
			selectedText += s.Text[y][:s.xRightOf(&to)]
		} else {
			selectedText += s.Text[y] + "\n"
		}
	}
	if err := clipboard.WriteAll(selectedText); err != nil {
		panic(err)
	}
}

func (s *State) paste() {
	_, to := s.normalisedSelection()
	text, err := clipboard.ReadAll()
	if err != nil {
		panic(err)
	}
	after := strings.Split(text, "\n")
	after[0] = s.Text[to.Y][:s.xRightOf(&to)] + after[0]
	after[len(after)-1] += s.Text[to.Y][s.xRightOf(&to):]
	s.applyDiff(diff{
		start:  s.Cursor.Y,
		before: s.Text[s.Cursor.Y : s.Cursor.Y+1],
		after:  after,
	})
}
