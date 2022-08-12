package window

import (
	"fmt"
	"log"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

type Window struct {
	Owner screen.Window

	List   *List
	Next   *Next
	Player *Player

	Title  string
	Width  int
	Height int
}

func New(t string, width, height int) *Window {
	var w Window
	w.Owner = nil
	w.Title = t
	w.Width = width
	w.Height = height
	return &w
}

func (w *Window) Start() error {

	var err error

	driver.Main(func(s screen.Screen) {

		opt := &screen.NewWindowOptions{
			Title:  w.Title,
			Width:  w.Width,
			Height: w.Height,
		}

		w.Owner, err = s.NewWindow(opt)
		if err != nil {
			log.Printf("NewList() error: %+v\n", err)
			return
		}
		defer w.Owner.Release()

		l, err := NewList(w.Owner, s)
		if err != nil {
			log.Printf("NewList() error: %+v\n", err)
			return
		}
		defer l.Release()
		w.List = l

		n, err := NewNext(w.Owner, s)
		if err != nil {
			log.Printf("NewNext() error: %+v\n", err)
			return
		}
		defer n.Release()
		w.Next = n

		p, err := NewPlayer(w.Owner, s)
		if err != nil {
			log.Printf("NewPlayer() error: %+v\n", err)
			return
		}
		defer p.Release()
		w.Player = p

		//クライアント描画
		for {
			e := w.Owner.NextEvent()
			switch e := e.(type) {
			case *Part:
				e.Redraw()
			default:
			}
		}
	})

	fmt.Println("Window Return")

	return nil
}
