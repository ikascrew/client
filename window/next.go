package window

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/client/tool"
)

type Next struct {
	cursor   int
	idx      int
	targets  []image.Image
	resource []string

	*Part
}

func NewNext(w screen.Window, s screen.Screen) (*Next, error) {

	n := &Next{}
	r := image.Rect(320, 0, 1536, 180)

	n.Part = &Part{}
	n.Init(w, s, r)

	n.targets = make([]image.Image, 0)
	n.resource = make([]string, 0)

	return n, nil
}

func (n *Next) Draw() {

	m := n.Part.buffer.RGBA()

	lox := 0
	loy := 0
	hix := 1280
	hiy := 180

	hor := n.cursor / 175

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	w := 320
	h := hiy
	start := (hor / w)

	for y := loy; y < hiy; y++ {
		var img image.Image
		for x := lox; x < hix; x++ {

			d := x / w
			idx := start + d

			flag := false
			if x >= 0 && x < w {
				n.idx = idx
				flag = true
				if x > 5 && x < (w-5) {
					if y > 5 && y < (h-5) {
						flag = false
					}
				}
			}

			if idx >= 0 && idx < len(n.targets) {
				img = n.targets[idx]
			} else {
				img = nil
			}

			dx := x - (d * w)
			go func(img image.Image, x, y, dx int, flag bool) {
				if img == nil {
					m.Set(x, y, black)
				} else if flag {
					m.Set(x, y, white)
				} else {
					m.Set(x, y, img.At(dx, y))
				}
			}(img, x, y, dx, flag)
		}
	}
}

func (n *Next) Get() string {
	sz := len(n.resource)
	if n.idx > sz-1 || sz == 0 || n.idx < 0 {
		return ""
	}

	rtn := n.resource[n.idx]
	return rtn
}

func (n *Next) Delete() error {

	sz := len(n.resource)
	if n.idx > sz-1 || sz == 0 || n.idx < 0 {
		return fmt.Errorf("Pusher Index Error")
	}

	newres := make([]string, 0)
	newtar := make([]image.Image, 0)
	for idx, elm := range n.resource {
		if idx != n.idx {
			newres = append(newres, elm)
			newtar = append(newtar, n.targets[idx])
		}
	}
	n.resource = newres
	n.targets = newtar

	n.cursor = 0
	n.Draw()

	return nil
}

func (n *Next) Add(f string) error {

	for _, elm := range n.resource {
		if f == elm {
			return fmt.Errorf("Resource[" + f + "] exist")
		}
	}
	n.resource = append(n.resource, f)

	img, err := tool.LoadImage(f)
	if err != nil {
		return err
	}
	n.targets = append(n.targets, img)

	n.cursor = 0
	n.Draw()

	return nil
}

func (n *Next) SetCursor(d int) {
	n.cursor = n.cursor + d
	n.Draw()
}
