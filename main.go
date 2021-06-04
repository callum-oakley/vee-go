package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/callum-oakley/vee/state"
	"github.com/callum-oakley/vee/ui"
	"github.com/gdamore/tcell/v2"
)

func main() {
	text, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(text), "\n")
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	s := state.State{
		FilePath: os.Args[2],
		TabWidth: 4,
		Text:     lines,
	}

	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	ui.Render(s, screen)

	for {
		switch e := screen.PollEvent().(type) {
		case *tcell.EventResize:
			ui.Render(s, screen)
		case *tcell.EventKey:
			switch e.Key() {
			case tcell.KeyRune:
				switch e.Rune() {
				case 'h':
					s.Left()
				case 'j':
					s.Down()
				case 'k':
					s.Up()
				case 'l':
					s.Right()
				}
				ui.Render(s, screen)
			case tcell.KeyEscape:
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}
