package client

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	window screen.Window

	list   *List
	next   *Next
	player *Player
}

func NewWindow(t string, w, h int) (window *Window, err error) {

	window = &Window{
		window: nil,
	}

	driver.Main(func(s screen.Screen) {
		opt := &screen.NewWindowOptions{
			Title:  t,
			Width:  w,
			Height: h,
		}

		window.window, err = s.NewWindow(opt)
		if err != nil {
			return
		}

		l, err := NewList(window.window, s)
		if err != nil {
			return
		}

		n, err := NewNext(window.window, s)
		if err != nil {
			return
		}

		p, err := NewPlayer(window.window, s)
		if err != nil {
			return
		}

		window.list = l
		window.next = n
		window.player = p

	})

	return
}

func (w *Window) Release() {
	w.window.Release()
	w.list.Release()
	w.next.Release()
	w.player.Release()
}
