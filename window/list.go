package window

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"

	"strings"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/client/tool"
)

var max = 0

type List struct {
	cursor   int
	idx      int
	images   []image.Image
	resource []string
	*Part
}

func NewList(w screen.Window, s screen.Screen) (*List, error) {

	l := &List{}

	r := image.Rect(0, 0, 320, 720)
	l.Part = &Part{}
	l.Init(w, s, r)

	work := "./.client/images"

	paths, err := tool.Search(work, "_thumb.jpg", nil)
	if err != nil {
		return nil, err
	}

	l.images = make([]image.Image, len(paths)+1)
	l.resource = make([]string, len(paths)+1)

	for idx, path := range paths {
		l.images[idx], _ = tool.LoadImage(path)

		id := strings.Replace(path, "_thumb.jpg", "", -1)
		l.resource[idx] = id + ".jpg"
	}

	fmt.Println(len(paths))

	max = len(paths) * 100 * 100

	return l, nil
}

func (l *List) Draw() {

	fmt.Println("List Draw")

	m := l.Part.buffer.RGBA()

	lox := 0
	loy := 0
	hix := 320
	hiy := 720

	ver := l.cursor / 200

	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	h := 64
	hf := h / 2
	cur := h*2 + hf
	start := (ver / h)

	for y := loy; y < hiy; y++ {

		var img image.Image

		d := y / h
		idx := start + d

		if idx >= 0 && idx < len(l.images) {
			img = l.images[start+d]
		}

		dy := y - (d * h)

		flag := false
		yflag := false

		if (y+hf) > cur && (y-hf) < cur {
			l.idx = idx
			if dy <= 5 || dy >= (h-5) {
				flag = true
			} else {
				yflag = true
			}
		}

		for x := lox; x < hix; x++ {

			if yflag {
				if x <= 5 || x >= (hix-5) {
					flag = true
				} else {
					flag = false
				}
			}

			go func(img image.Image, x, y, dy int, flag bool) {
				if img == nil {
					m.Set(x, y, black)
				} else if flag {
					m.Set(x, y, white)
				} else {
					m.Set(x, y, img.At(x, dy))
				}
			}(img, x, y, dy, flag)
		}
	}
}

func (l *List) Get() string {
	if l.idx < 0 || l.idx >= len(l.resource) {
		return ""
	}
	return l.resource[l.idx]
}

func (l *List) SetCursor(d int) {
	l.cursor = l.cursor + d
	l.Draw()
}

func (l *List) ZeroCursor() {
	l.cursor = 0
	l.Draw()
}

func (l *List) MaxCursor() {
	l.cursor = l.cursor + max
	l.Draw()
}
