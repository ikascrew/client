package window

import (
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	Owner screen.Window

	List   *List
	Next   *Next
	Player *Player
}

func New(t string, w, h int) (window *Window, err error) {

	window = &Window{
		Owner: nil,
	}

	driver.Main(func(s screen.Screen) {
		opt := &screen.NewWindowOptions{
			Title:  t,
			Width:  w,
			Height: h,
		}

		window.Owner, err = s.NewWindow(opt)
		if err != nil {
			return
		}

		l, err := NewList(window.Owner, s)
		if err != nil {
			return
		}

		n, err := NewNext(window.Owner, s)
		if err != nil {
			return
		}

		p, err := NewPlayer(window.Owner, s)
		if err != nil {
			return
		}

		window.List = l
		window.Next = n
		window.Player = p

	})

	return
}

func (w *Window) Release() {
	w.Owner.Release()
	w.List.Release()
	w.Next.Release()
	w.Player.Release()
}
