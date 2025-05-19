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
	centerX := int(f.X * float64(t.Bounds().Dx()))
	centerY := int(f.Y * float64(t.Bounds().Dy()))

	const fixedFigureSize float64 = 300.0
	figureSize := fixedFigureSize

	halfSize := figureSize / 2
	barWidth := figureSize / 5
	yellow := color.RGBA{R: 255, G: 255, A: 255}

	verticalPart := image.Rect(
		centerX-int(halfSize),               
		centerY-int(halfSize),               
		centerX-int(halfSize)+int(barWidth),
		centerY+int(halfSize),             
	)

	horizontalPart := image.Rect(
		centerX-int(halfSize),          
		centerY-int(barWidth/2),      
		centerX+int(halfSize),         
		centerY+int(barWidth/2),       
	)

	t.Fill(verticalPart, yellow, draw.Src)
	t.Fill(horizontalPart, yellow, draw.Src)

	return false
}
