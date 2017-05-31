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
	loose bool
	table *Table
)

func init() {

	log.SetFlags(log.Lshortfile)
}

func main() {
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
	loose = false
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
			t.Items[4*row+col] = t.NewItem()
		}
	}

	t.NextMove(3)

	return t
}

type Item struct {
	N   int
	btn *fizzgui.Widget
}

//NewItem created new number square button contains 2 or 4 in any random empty position
func (t *Table) NewItem() *Item {
	item := &Item{}

	item.btn = t.Container.NewButton(strconv.Itoa(item.N), nil)
	item.btn.Hidden = true
	item.btn.Layout.PositionFixed = true

	item.btn.StyleHover = fizzgui.Style{}
	item.btn.Style.BorderWidth = 0

	item.btn.Layout.SetHeight("25%")
	item.btn.Layout.SetWidth("25%")

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

		col := fmt.Sprintf("%d%%", i%4*25)
		row := fmt.Sprintf("%d%%", i/4*25)

		item.btn.Layout.SetX(col)
		item.btn.Layout.SetY(row)

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
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 32:
			item.btn.Style.BackgroundColor = fizzgui.Color(245, 124, 95, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 64:
			item.btn.Style.BackgroundColor = fizzgui.Color(246, 93, 59, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 128:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 206, 113, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 256:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 204, 97, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 512:
			item.btn.Style.BackgroundColor = fizzgui.Color(236, 200, 80, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 1024:
			item.btn.Style.BackgroundColor = fizzgui.Color(237, 197, 63, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		case 2048:
			item.btn.Style.BackgroundColor = fizzgui.Color(236, 196, 0, 255)
			item.btn.Style.TextColor = fizzgui.Color(249, 246, 241, 255)
		}
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

	rc := t.empty[t.rand.Intn(len(t.empty))]
	item := t.Items[rc]

	item.SetValue(t.newNum())
	item.btn.Hidden = false
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
	// log.Println(r0, c0, r1, c1)
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
			prev.SetValue(prev.N * 2)
			item.Hide()
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
			prev.SetValue(prev.N * 2)
			item.Hide()
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
			prev.SetValue(prev.N * 2)
			item.Hide()
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
			prev.SetValue(prev.N * 2)
			item.Hide()
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

	if !loose {
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
			table.NextMove(2)
		} else {

			for _, item := range table.Items {
				if item.N == 0 {
					loose = false
					break
				} else {
					loose = true

				}
			}
		}
	}

	if loose {
		Loose()
	}

}

func Loose() {
	if loose {
		loose := table.Container.NewButton("YOU LOOSE! RESTART?", NewGame)
		loose.Layout.PositionFixed = true
		loose.Layout.SetWidth("80%")
		loose.Layout.SetX("10%")
		loose.Layout.SetY("40%")
	}
}

//Close it`s callback from renderLoop, should close application
func Close() {
	os.Exit(0)
}
