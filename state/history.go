package state

type diff struct {
	start         int
	before, after []string
}

type change struct {
	diff
	anchorBefore, anchorAfter, cursorBefore, cursorAfter cursor
}

func (d *diff) isEmpty() bool {
	return len(d.before) == 0 && len(d.after) == 0
}

func apply(d diff, text []string) []string {
	if d.isEmpty() {
		return text
	}
	return append(
		text[:d.start],
		append(d.after, text[d.start+len(d.before):]...)...,
	)
}

func revert(d diff, text []string) []string {
	if d.isEmpty() {
		return text
	}
	return append(
		text[:d.start],
		append(d.before, text[d.start+len(d.after):]...)...,
	)
}

func compose(b diff, a diff) diff {
	if a.isEmpty() {
		return b
	}
	start := min(a.start, b.start)
	end := max(a.start+len(a.after), b.start+len(b.before))
	intermediateA := make([]string, end-start)
	intermediateB := make([]string, end-start)
	for i := start; i < end; i++ {
		if i >= a.start && i < a.start+len(a.after) {
			intermediateA[i-start] = a.after[i-a.start]
			intermediateB[i-start] = a.after[i-a.start]
		} else {
			intermediateA[i-start] = b.before[i-b.start]
			intermediateB[i-start] = b.before[i-b.start]
		}
	}
	return diff{
		start: start,
		before: revert(
			diff{start: a.start - start, before: a.before, after: a.after},
			intermediateA,
		),
		after: apply(
			diff{start: b.start - start, before: b.before, after: b.after},
			intermediateB,
		),
	}
}

func (s *State) startChange() {
	s.change = change{
		anchorBefore: s.Anchor,
		cursorBefore: s.Cursor,
	}
}

func (s *State) endChange() {
	if s.change.isEmpty() {
		return
	}
	s.change.anchorAfter = s.Anchor
	s.change.cursorAfter = s.Cursor
	s.history = append(s.history[:s.historyHead], s.change)
	s.historyHead++
}

func (s *State) applyDiff(d diff) {
	before := d.before
	d.before = make([]string, len(before))
	copy(d.before, before)

	after := d.after
	d.after = make([]string, len(after))
	copy(d.after, after)

	s.Text = apply(d, s.Text)
	s.change.diff = compose(d, s.change.diff)
}

func (s *State) undo() {
	if s.historyHead == 0 {
		return
	}
	s.historyHead--
	s.Text = revert(s.history[s.historyHead].diff, s.Text)
	s.Anchor = s.history[s.historyHead].anchorBefore
	s.Cursor = s.history[s.historyHead].cursorBefore
}

func (s *State) redo() {
	if s.historyHead >= len(s.history) {
		return
	}
	s.Text = apply(s.history[s.historyHead].diff, s.Text)
	s.Anchor = s.history[s.historyHead].anchorAfter
	s.Cursor = s.history[s.historyHead].cursorAfter
	s.historyHead++
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
