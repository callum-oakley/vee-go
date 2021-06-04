package state

type State struct {
	FilePath string
	TabWidth int
	Text     []string
	Cursor   struct{ X, Y, Col int }
}

func (s *State) updateCol() {
	s.Cursor.Col = 0
	for _, c := range s.Text[s.Cursor.Y][:s.Cursor.X] {
		s.Cursor.Col++
		for c == '\t' && s.Cursor.Col%s.TabWidth > 0 {
			s.Cursor.Col++
		}
	}
}

func (s *State) updateX() {
	s.Cursor.X = -1
	col := 0
	for s.Cursor.X = range s.Text[s.Cursor.Y] {
		col++
		for s.Text[s.Cursor.Y][s.Cursor.X] == '\t' && col%s.TabWidth > 0 {
			col++
		}
		if col > s.Cursor.Col {
			break
		}
	}
}

func (s *State) Left() {
	if s.Cursor.X > 0 {
		s.Cursor.X--
		s.updateCol()
	}
}

func (s *State) Right() {
	if s.Cursor.X < len(s.Text[s.Cursor.Y])-1 {
		s.Cursor.X++
		s.updateCol()
	}
}

func (s *State) Up() {
	if s.Cursor.Y != 0 {
		s.Cursor.Y--
		s.updateX()
	}
}

func (s *State) Down() {
	if s.Cursor.Y < len(s.Text)-1 {
		s.Cursor.Y++
		s.updateX()
	}
}
