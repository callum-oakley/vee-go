package state

type State struct {
	FileName string
	TabWidth int
	Text     []string
	Cursor   struct{ X, Y int }
}

func NewState(fileName string, tabWidth int, text []byte) State {
	s := State{FileName: fileName, TabWidth: tabWidth}
	i := 0
	for j, c := range text {
		if c == '\n' {
			s.Text = append(s.Text, string(text[i:j]))
			i = j + 1
		}
	}
	if i < len(text) {
		s.Text = append(s.Text, string(text[i:]))
	}
	return s
}

func (s *State) Left() {
	if s.Cursor.X > 0 && s.xMax(s.Cursor.Y) >= 0 {
		s.snapToEnd()
		s.Cursor.X--
		s.snapToTab(-1)
	}
}

func (s *State) Right() {
	if s.Cursor.X < s.xMax(s.Cursor.Y) {
		s.snapToEnd()
		s.Cursor.X++
		s.snapToTab(1)
	}
}

func (s *State) Up() {
	if s.Cursor.Y != 0 {
		s.Cursor.Y--
	}
}

func (s *State) Down() {
	if s.Cursor.Y != len(s.Text)-1 {
		s.Cursor.Y++
	}
}

func (s *State) snapToEnd() {
	if s.Cursor.X > s.xMax(s.Cursor.Y) {
		s.Cursor.X = s.xMax(s.Cursor.Y)
	}
}

func (s *State) snapToTab(dir int) {
	if s.charAt(s.Cursor.X, s.Cursor.Y) == '\t' {
		for s.Cursor.X%s.TabWidth > 0 {
			s.Cursor.X += dir
		}
	}
}

func (s *State) LogicalX(x, y int) int {
	visualX := 0
	logicalX := -1
	for logicalX = range s.Text[y] {
		visualX++
		if s.Text[y][logicalX] == '\t' {
			for visualX%s.TabWidth > 0 {
				visualX++
			}
		}
		if visualX > x {
			break
		}
	}
	return logicalX
}

func (s *State) charAt(x, y int) byte {
	return s.Text[y][s.LogicalX(x, y)]
}

func (s *State) xMax(y int) int {
	res := 0
	for _, c := range s.Text[y] {
		res++
		if c == '\t' {
			for res%s.TabWidth > 0 {
				res++
			}
		}
	}
	res--
	if len(s.Text[y]) > 0 && s.Text[y][len(s.Text[y])-1] == '\t' {
		for res%s.TabWidth > 0 {
			res--
		}
	}
	return res
}
