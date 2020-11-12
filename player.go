package client

import (
	"C"

	"image"

	"golang.org/x/exp/shiny/screen"

	"github.com/ikascrew/client/tool"
)

type Player struct {
	target image.Image
	*Part
}

func NewPlayer(w screen.Window, s screen.Screen) (*Player, error) {
	p := &Player{}

	r := image.Rect(320, 180, 1536, 720)
	p.Part = &Part{}
	p.Init(w, s, r)

	return p, nil
}

func (p *Player) setFile(f string) error {

	img, err := tool.LoadImage(f)
	if err != nil {
		return err
	}

	p.target = img

	return nil

}

func (p *Player) Draw() {

	if p.target == nil {
		return
	}

	//m := p.Part.buffer.RGBA()
	//img := p.target

	/*
		for y := 0; y < height; y++ {
			for x := 0; x < step; x = x + channels {
				c.B = uint8(data[y*step+x])
				c.G = uint8(data[y*step+x+1])
				c.R = uint8(data[y*step+x+2])
				if channels == 4 {
					c.A = uint8(data[y*step+x+3])
				}
				m.Set(int(x/channels), y, c)
			}
		}
	*/
	return
}
