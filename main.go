package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	s := state.NewState(os.Args[2], 4, text)

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
