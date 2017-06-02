package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/sg3des/fizzgui"
)

var (
	lost  bool
	table *Table
)

func main() {
	log.SetFlags(log.Lshortfile)

	if err := NewWindow("2048", 600, 600); err != nil {
		log.Fatalln(err)
	}

	NewGame(nil)

	RenderLoop()
}

func NewGame(_ *fizzgui.Widget) {
	if table != nil {
		table.Container.Close()
	}

	table = NewTable()
	lost = false
}

//Table is main struct contains matrix 4x4
type Table struct {
	Container *fizzgui.Container
	Items     [16]*Item

	rand  *rand.Rand
	empty []int
}

//NewTable initialize table
func NewTable() *Table {
	t := &Table{
		Container: fizzgui.NewContainer("table", "0px", "0px", "100%", "100%"),
		rand:      rand.New(rand.NewSource(time.Now().Unix())),
	}
	t.Container.Style.BackgroundColor = fizzgui.Color(187, 173, 160, 255)

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			item := t.NewItem()
			item.transSrc = TransSrc{X: float32(col * 25), Y: float32(row * 25)}

			t.Items[4*row+col] = item
		}
	}

	t.NextMove(1)

	return t
}

type Item struct {
	N          int
	btn        *fizzgui.Widget
	transition bool
	transSrc   TransSrc
}

type TransSrc struct {
	X, Y, W float32
}

//NewItem created new number square button contains 2 or 4 in any random empty position
func (t *Table) NewItem() *Item {
	item := &Item{}

	item.btn = t.Container.NewButton(strconv.Itoa(item.N), nil)
	item.btn.Hidden = true
	item.btn.Layout.PositionFixed = true

	item.btn.StyleHover = fizzgui.Style{}
	item.btn.Style.BorderWidth = 0

	item.btn.Layout.SetHeight("0%")
	item.btn.Layout.SetWidth("0%")

	return item
}

func (t *Table) NextMove(count int) {
	for i := 0; i < count; i++ {
		t.FillRandomItem()
	}

	for i, item := range t.Items {
		if item.N == 0 {
			item.btn.Hidden = true
			continue
		}

		item.btn.Hidden = false

		if !item.transition {
			row := fmt.Sprintf("%d%%", i/4*25)
			col := fmt.Sprintf("%d%%", i%4*25)

			item.btn.Layout.SetX(col)
			item.btn.Layout.SetY(row)
			item.btn.Layout.SetWidth("25%")
			item.btn.Layout.SetHeight("25%")
		}

		switch item.N {
		case 2:
			item.btn.Style.BackgroundColor = fizzgui.Color(238, 228, 218, 255)
			item.btn.Style.TextColor = fizzgui.Color(80, 80, 80, 255)
		case 4:
			item.btn.Style.BackgroundColor = fizzgui.Color(236, 224, 200, 255)
			item.btn.Style.TextColor = fizzgui.Color(80, 80, 80, 255)
		case 8:
			item.btn.Style.BackgroundColor = fizzgui.Color(242, 177, 121, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 16:
			item.btn.Style.BackgroundColor = fizzgui.Color(245, 149, 99, 255)
		case 32:
			item.btn.Style.BackgroundColor = fizzgui.Color(245, 124, 95, 255)
		case 64:
			item.btn.Style.BackgroundColor = fizzgui.Color(246, 93, 59, 255)
		case 128:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 206, 113, 255)
		case 256:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 204, 97, 255)
		case 512:
			item.btn.Style.BackgroundColor = fizzgui.Color(236, 200, 80, 255)
		case 1024:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 197, 63, 255)
		case 2048:
			item.btn.Style.BackgroundColor = fizzgui.Color(236, 196, 0, 255)
		}
	}
}

func Transitions(dt float32) {
	if table == nil {
		return
	}

	dt = dt * 512

	for i, item := range table.Items {
		if !item.transition {
			continue
		}

		row := float32(i / 4 * 25)
		col := float32(i % 4 * 25)

		var rowEqual bool
		if row > item.transSrc.Y+dt {
			item.transSrc.Y += dt
			row = item.transSrc.Y
		} else if row < item.transSrc.Y-dt {
			item.transSrc.Y -= dt
			row = item.transSrc.Y
		} else {
			rowEqual = true
		}

		var colEqual bool
		if col > item.transSrc.X+dt {
			item.transSrc.X += dt
			col = item.transSrc.X
		} else if col < item.transSrc.X-dt {
			item.transSrc.X -= dt
			col = item.transSrc.X
		} else {
			colEqual = true
		}

		var widthEqual bool
		if item.transSrc.W < 25-dt {
			item.transSrc.W += dt
		} else {
			item.transSrc.W = 25
			widthEqual = true
		}

		if colEqual && rowEqual && widthEqual {
			item.transition = false
		}

		if !widthEqual {
			col += 12.5 - item.transSrc.W/2
			row += 12.5 - item.transSrc.W/2
		}

		item.btn.Layout.SetX(fmt.Sprintf("%.0f%%", col))
		item.btn.Layout.SetY(fmt.Sprintf("%.0f%%", row))
		item.btn.Layout.SetWidth(fmt.Sprintf("%0.0f%%", item.transSrc.W))
		item.btn.Layout.SetHeight(fmt.Sprintf("%0.0f%%", item.transSrc.W))
	}
}

func (t *Table) FillRandomItem() {
	t.empty = t.empty[:0]
	for i, item := range t.Items {
		if item.N == 0 {
			t.empty = append(t.empty, i)
		}
	}

	if len(t.empty) == 0 {
		return
	}

	i := t.empty[t.rand.Intn(len(t.empty))]
	item := t.Items[i]

	item.SetValue(t.newNum())
	item.btn.Hidden = false
	item.transition = true
	item.transSrc.W = 0
	item.transSrc.X = float32(i % 4 * 25)
	item.transSrc.Y = float32(i / 4 * 25)
}

// //newNum return new number 2 or 4
func (t *Table) newNum() int {
	if t.rand.Uint32()%2 == 0 {
		return 2
	}
	return 4
}

func (item *Item) SetValue(n int) {
	item.N = n
	item.btn.Text = strconv.Itoa(n)
}

func (item *Item) Hide() {
	item.N = 0
}

func (t *Table) MoveItem(r0, c0, r1, c1 int) {
	t.Items[4*r0+c0].transition = true
	t.Items[4*r0+c0].transSrc = TransSrc{X: float32(c0 * 25), Y: float32(r0 * 25), W: 25}
	// t.Items[4*r0+c0].transDst = [2]int{r1 % 4 * 25, c1 / 4 * 25}
	t.Items[4*r0+c0], t.Items[4*r1+c1] = t.Items[4*r1+c1], t.Items[4*r0+c0]
}

func (t *Table) LookupPrevItem(row, col int) (int, int, *Item) {
	for c := col - 1; c >= 0; c-- {
		if item := t.Items[4*row+c]; item.N > 0 {
			return row, c, item
		}
	}

	return row, 0, nil
}

func (t *Table) MoveLeft() (moves int) {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			item := t.Items[4*row+col]
			if item.N == 0 {
				continue
			}

			r1, c1, prev := t.LookupPrevItem(row, col)
			if prev == nil {
				if col != c1 {
					moves++
					table.MoveItem(row, col, r1, c1)
				}
				continue
			}

			if prev.N != item.N {
				if col != c1+1 {
					moves++
					table.MoveItem(row, col, r1, c1+1)
				}
				continue
			}

			moves++
			prev.Hide()
			item.SetValue(item.N * 2)
			table.MoveItem(row, col, r1, c1)
		}
	}

	return
}

func (t *Table) LookupNextItem(row, col int) (int, int, *Item) {
	for c := col + 1; c <= 3; c++ {
		if item := t.Items[4*row+c]; item.N > 0 {
			return row, c, item
		}
	}

	return row, 3, nil
}

func (t *Table) MoveRight() (moves int) {
	for row := 0; row < 4; row++ {
		for col := 3; col >= 0; col-- {
			item := t.Items[4*row+col]
			if item.N == 0 {
				continue
			}

			r1, c1, prev := t.LookupNextItem(row, col)
			if prev == nil {
				if col != c1 {
					moves++
					table.MoveItem(row, col, r1, c1)
				}
				continue
			}

			if prev.N != item.N {
				if col != c1-1 {
					moves++
					table.MoveItem(row, col, r1, c1-1)
				}
				continue
			}

			moves++
			prev.Hide()
			item.SetValue(item.N * 2)
			table.MoveItem(row, col, r1, c1)
		}
	}

	return
}

func (t *Table) LookupPrevItemCol(row, col int) (int, int, *Item) {
	for r := row - 1; r >= 0; r-- {
		if item := t.Items[4*r+col]; item.N > 0 {
			return r, col, item
		}
	}

	return 0, col, nil
}

func (t *Table) MoveUp() (moves int) {
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			item := t.Items[4*row+col]
			if item.N == 0 {
				continue
			}

			r1, c1, prev := t.LookupPrevItemCol(row, col)
			if prev == nil {
				if row != r1 {
					moves++
					table.MoveItem(row, col, r1, c1)
				}
				continue
			}

			if prev.N != item.N {
				if row != r1+1 {
					moves++
					table.MoveItem(row, col, r1+1, c1)
				}

				continue
			}

			moves++
			prev.Hide()
			item.SetValue(item.N * 2)
			table.MoveItem(row, col, r1, c1)
		}
	}

	return
}

func (t *Table) LookupNextItemCol(row, col int) (int, int, *Item) {
	for r := row + 1; r <= 3; r++ {
		if item := t.Items[4*r+col]; item.N > 0 {
			return r, col, item
		}
	}

	return 3, col, nil
}

func (t *Table) MoveDown() (moves int) {
	for col := 0; col < 4; col++ {
		for row := 3; row >= 0; row-- {

			item := t.Items[4*row+col]
			if item.N == 0 {
				continue
			}

			r1, c1, prev := t.LookupNextItemCol(row, col)
			if prev == nil {
				if row != r1 {
					moves++
					table.MoveItem(row, col, r1, c1)
				}
				continue
			}

			if prev.N != item.N {
				if row != r1-1 {
					moves++
					table.MoveItem(row, col, r1-1, c1)
				}

				continue
			}

			moves++
			prev.Hide()
			item.SetValue(item.N * 2)
			table.MoveItem(row, col, r1, c1)
		}
	}

	return
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	if key == glfw.KeyEscape {
		w.SetShouldClose(true)
		return
	}

	if !lost {
		var moves int

		switch key {
		case glfw.KeyLeft:
			moves = table.MoveLeft()
		case glfw.KeyRight:
			moves = table.MoveRight()
		case glfw.KeyUp:
			moves = table.MoveUp()
		case glfw.KeyDown:
			moves = table.MoveDown()
		case glfw.KeyEscape:

		}

		if moves > 0 {
			table.NextMove(1)
		} else {

			for _, item := range table.Items {
				if item.N == 0 {
					lost = false
					break
				} else {
					lost = true

				}
			}
		}
	}

	if lost {
		Lost()
	}

}

func Lost() {
	if lost {
		lostBtn := table.Container.NewButton("YOU LOSE! RESTART?", NewGame)
		lostBtn.Layout.PositionFixed = true
		lostBtn.Font = TextFont
		lostBtn.Layout.SetWidth("80%")
		lostBtn.Layout.SetX("10%")
		lostBtn.Layout.SetY("40%")
	}
}

//Close it`s callback from renderLoop, should close application
func Close() {
	os.Exit(0)
}
