package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestLoop_PostAndStop(t *testing.T) {
	var (
		loop Loop
		tr   testReceiver
	)
	loop.Receiver = &tr
	loop.Start(mockScreen{})

	var executed []string

	loop.Post(OperationFunc(func(t screen.Texture) {
		executed = append(executed, "op 1")
	}))
	loop.Post(WhiteFill{})
	loop.Post(UpdateOp)
	loop.Post(OperationFunc(func(t screen.Texture) {
		executed = append(executed, "op 2")
	}))
	loop.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type")
	}

	if len(mt.Colors) == 0 {
		t.Error("Expected texture to be filled at least once")
	}

	wantOrder := []string{"op 1", "op 2"}
	if !reflect.DeepEqual(executed, wantOrder) {
		t.Errorf("Unexpected execution order: got %v, want %v", executed, wantOrder)
	}
}

func TestLoop_UpdateSwitchTextures(t *testing.T) {
	var (
		loop Loop
		tr   testReceiver
	)
	loop.Receiver = &tr
	loop.Start(mockScreen{})

	loop.Post(WhiteFill{}) 
	loop.Post(UpdateOp)

	loop.StopAndWait()

	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Texture not valid")
	}
	if len(mt.Colors) == 0 || mt.Colors[0] != color.White {
		t.Errorf("Expected white fill, got: %v", mt.Colors)
	}
}

func TestLoop_StopWithoutOperations(t *testing.T) {
	var loop Loop
	loop.Receiver = &testReceiver{}
	loop.Start(mockScreen{})

	loop.StopAndWait()
}

// ==== mocks =====

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return &mockTexture{}, nil
}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("not implemented")
}
func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("not implemented")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point {
	return image.Pt(800, 800)
}

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rect(0, 0, 800, 800)
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
