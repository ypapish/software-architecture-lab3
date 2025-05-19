package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ypapish/software-architecture-lab3/painter"
)

type Parser struct {
	lastBgColor painter.Operation
	lastBgRect  *painter.BgRect
	figures     []painter.Figure // Цей зріз тепер зберігає поточний стан фігур
	hasUpdate   bool
}

func (p *Parser) initialize() {
	p.lastBgColor = nil
	p.lastBgRect = nil
	p.figures = nil
	p.hasUpdate = false
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if err := p.parseLine(line); err != nil {
			return nil, err
		}
	}

	return p.finalize(), nil
}

func (p *Parser) parseLine(line string) error {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "white":
		p.lastBgColor = painter.WhiteFill{}
	case "green":
		p.lastBgColor = painter.GreenFill{}
	case "bgrect":
		if len(parts) != 5 {
			return fmt.Errorf("bgrect requires 4 arguments")
		}
		x1, err1 := strconv.ParseFloat(parts[1], 64)
		y1, err2 := strconv.ParseFloat(parts[2], 64)
		x2, err3 := strconv.ParseFloat(parts[3], 64)
		y2, err4 := strconv.ParseFloat(parts[4], 64)
		if err := firstNonNil(err1, err2, err3, err4); err != nil {
			return fmt.Errorf("invalid number in bgrect: %v", err)
		}
		p.lastBgRect = &painter.BgRect{X1: x1, Y1: y1, X2: x2, Y2: y2}
	case "figure":
		if len(parts) != 3 {
			return fmt.Errorf("figure requires 2 arguments")
		}
		x, err1 := strconv.ParseFloat(parts[1], 64)
		y, err2 := strconv.ParseFloat(parts[2], 64)
		if err := firstNonNil(err1, err2); err != nil {
			return fmt.Errorf("invalid number in figure: %v", err)
		}
		p.figures = append(p.figures, painter.Figure{X: x, Y: y})
	case "move":
		if len(parts) != 3 {
			return fmt.Errorf("move requires 2 arguments")
		}
		dx, err1 := strconv.ParseFloat(parts[1], 64)
		dy, err2 := strconv.ParseFloat(parts[2], 64)
		if err := firstNonNil(err1, err2); err != nil {
			return fmt.Errorf("invalid number in move: %v", err)
		}
		for i := range p.figures {
			p.figures[i].X = clamp(p.figures[i].X+dx, 0, 1)
			p.figures[i].Y = clamp(p.figures[i].Y+dy, 0, 1)
		}
	case "reset":
		p.initialize()
		p.lastBgColor = painter.Reset{}
	case "update":
		p.hasUpdate = true
	default:
		return fmt.Errorf("unknown command: %s", parts[0])
	}
	return nil
}

func (p *Parser) finalize() []painter.Operation {
	var result []painter.Operation

	if p.lastBgRect != nil {
		result = append(result, *p.lastBgRect)
	}

	if p.lastBgColor != nil {
		result = append(result, p.lastBgColor)
	}

	for _, f := range p.figures {
		result = append(result, f)
	}

	if p.hasUpdate {
		result = append(result, painter.UpdateOp)
		p.hasUpdate = false
	}

	return result
}

func firstNonNil(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
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
