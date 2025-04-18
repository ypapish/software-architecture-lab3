package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/draw"
)

type Operation interface {
	Do(t screen.Texture) bool
}

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) bool {
	ready := false
	for _, op := range ol {
		if op.Do(t) {
			ready = true
		}
	}
	return ready
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (updateOp) Do(t screen.Texture) bool { return true }

type WhiteFill struct{}

func (WhiteFill) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.White, draw.Src)
	return false
}

type GreenFill struct{}

func (GreenFill) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.RGBA{G: 255, A: 255}, draw.Src)
	return false
}

type Reset struct{}

func (Reset) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.Black, draw.Src)
	return false
}

type BgRect struct {
	X1, Y1, X2, Y2 float64
}

func (op BgRect) Do(t screen.Texture) bool {
	x1 := int(op.X1 * float64(t.Bounds().Dx()))
	y1 := int(op.Y1 * float64(t.Bounds().Dy()))
	x2 := int(op.X2 * float64(t.Bounds().Dx()))
	y2 := int(op.Y2 * float64(t.Bounds().Dy()))
	t.Fill(image.Rect(x1, y1, x2, y2), color.Black, draw.Src)
	return false
}

type Figure struct {
	X, Y float64
}

func (f Figure) Do(t screen.Texture) bool {
	x := int(f.X * float64(t.Bounds().Dx()))
	y := int(f.Y * float64(t.Bounds().Dy()))

	barWidth := 300 / 5
	halfSize := 300 / 2
	yellow := color.RGBA{255, 255, 0, 255}

	vertical := image.Rect(x-barWidth/2, y-halfSize, x+barWidth/2, y+halfSize)
	horizontal := image.Rect(x-barWidth/2, y-halfSize, x+halfSize, y-halfSize+barWidth)

	t.Fill(vertical, yellow, draw.Over)
	t.Fill(horizontal, yellow, draw.Over)

	return false
}

type Move struct {
    DX, DY  float64
    Figures *[]Figure
}

func (m Move) Do(t screen.Texture) bool {
    if m.Figures == nil {
        return false
    }
    for i := range *m.Figures {
        (*m.Figures)[i].X = clamp((*m.Figures)[i].X + m.DX, 0, 1)
        (*m.Figures)[i].Y = clamp((*m.Figures)[i].Y + m.DY, 0, 1)
    }
    return false
}

func clamp(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
