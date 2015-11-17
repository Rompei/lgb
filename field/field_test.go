package field

import (
	"testing"
)

const (
	Width  = 10
	Height = 10
	Text   = "Test"
)

func TestAddPoint(t *testing.T) {

	f := InitField(Width, Height)

	err := f.AddPoint(2, 3, Text)
	if err != nil {
		t.Error(err)
	}
	err = f.AddPoint(0, 0, Text)
	if err != nil {
		t.Error(err)
	}
	err = f.AddPoint(9, 9, Text)
	if err != nil {
		t.Error(err)
	}

	if !f.Points[3][2].IsAlive {
		t.Error("Cell was dead")
	}

	if f.Points[3][2].Str != Text {
		t.Error("Text went wrong.")
	}
	if !f.Points[0][0].IsAlive {
		t.Error("Cell was dead")
	}

	if f.Points[0][0].Str != Text {
		t.Error("Text went wrong.")
	}
	if !f.Points[9][9].IsAlive {
		t.Error("Cell was dead")
	}

	if f.Points[9][9].Str != Text {
		t.Error("Text went wrong.")
	}

}

func TestAddRandomPoint(t *testing.T) {

	f := InitField(Width, Height)

	f.AddRandomPoint(Text)

	numAlive := 0
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if f.Points[y][x].IsAlive {
				if f.Points[y][x].Str != Text {
					t.Error("Text went wrong.")
				}
				numAlive++
			}
		}
	}

	if numAlive != 1 {
		t.Error("The number of alive cells went wrong.")
	}
}

func TestDelPoint(t *testing.T) {
	f := InitField(Width, Height)

	f.AddPoint(2, 3, Text)
	f.AddPoint(0, 0, Text)
	f.AddPoint(9, 9, Text)

	f.DelPoint(2, 3)
	f.DelPoint(0, 0)
	f.DelPoint(9, 9)

	if f.Points[3][2].IsAlive {
		t.Error("Cell was not dead.")
	}
}

func TestIsAlive(t *testing.T) {
	f := InitField(Width, Height)

	f.AddPoint(2, 3, Text)
	f.AddPoint(0, 0, Text)
	f.AddPoint(9, 9, Text)

	if !f.IsAlive(2, 3) {
		t.Error("Cell was dead actuary alive.")
	}
	if !f.IsAlive(0, 0) {
		t.Error("Cell was dead actuary alive.")
	}
	if !f.IsAlive(9, 9) {
		t.Error("Cell was dead actuary alive.")
	}

	if f.IsAlive(1, 1) {
		t.Error("Cell was alive actuary dead.")
	}
}

func TestGetNumberOfAlive(t *testing.T) {
	f := InitField(Width, Height)

	if f.GetNumberOfAlive() != 0 {
		t.Error("The number of alive cells was went wrong actuary 0")
	}

	f.AddPoint(2, 3, Text)
	f.AddPoint(2, 2, Text)
	f.AddPoint(5, 7, Text)

	if f.GetNumberOfAlive() != 3 {
		t.Error("The number of alive cells was went wrong actuary 3")
	}

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			f.AddPoint(x, y, Text)
		}
	}

	if f.GetNumberOfAlive() != 100 {
		t.Error("The number of alive cells was went wrong actuary 100")
	}
}

func TestGetAliveCells(t *testing.T) {
	f := InitField(Width, Height)

	f.AddPoint(0, 0, Text)
	f.AddPoint(0, 2, Text)
	f.AddPoint(0, 1, Text)
	f.AddPoint(0, 3, Text)

	p1 := f.GetAliveCells(1, 1)
	if p1 != 3 {
		t.Errorf("The number of alive cells %v went wrong actuary 3.", p1)
	}

	p2 := f.GetAliveCells(2, 2)
	if p2 != 0 {
		t.Errorf("The number of alive cells %v went wrong actuary 0.", p2)
	}

	p3 := f.GetAliveCells(1, 3)
	if p3 != 2 {
		t.Errorf("The number of alive cells %v went wrong actuary 2.", p3)
	}

	p4 := f.GetAliveCells(1, 4)
	if p4 != 1 {
		t.Errorf("The number of alive cells %v went wrong actuary 1.", p3)
	}
}
