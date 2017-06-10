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

//NewGame - create new table and fill it with 2 items
func NewGame(_ *fizzgui.Widget) {
	if table != nil {
		table.Container.Close()
	}

	table = NewTable()
	table.FillRandomItem()
	table.FillRandomItem()
	table.ReDraw()
	lost = false
}

//Table is main struct contains matrix 4x4
type Table struct {
	Container *fizzgui.Container
	Items     [16]*Item

	rand *rand.Rand
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

	return t
}

//Item is square with number on table
type Item struct {
	N          int
	btn        *fizzgui.Widget
	transition bool
	transSrc   TransSrc
}

//TransSrc contains X,Y and Size values for transitions
type TransSrc struct {
	Y, X, S float32
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

//ReDraw func update values, positions and styles of items
func (t *Table) ReDraw() {
	for i, item := range t.Items {

		if item.N == 0 {
			item.btn.Hidden = true
			continue
		}

		item.btn.Hidden = false
		item.btn.Text = strconv.Itoa(item.N)

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

//Dump is print value of items 4x4 to stdout
func (t *Table) Dump() {
	for i, item := range t.Items {
		for _i, _item := range t.Items {
			if item == _item && i != _i {
				log.Fatalln("item equal %d == %d", i, _i)
			}
		}
	}

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			i := r*4 + c
			fmt.Printf("%2d:%d ", i, t.Items[i].N)
		}
		fmt.Println()
	}
}

//Transitions is handle animations
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
		if item.transSrc.S < 25-dt {
			item.transSrc.S += dt / 4
		} else {
			item.transSrc.S = 25
			widthEqual = true
		}

		if colEqual && rowEqual && widthEqual {
			item.transition = false
		}

		if !widthEqual {
			col += 12.5 - item.transSrc.S/2
			row += 12.5 - item.transSrc.S/2
		}

		item.btn.Layout.SetX(fmt.Sprintf("%.0f%%", col))
		item.btn.Layout.SetY(fmt.Sprintf("%.0f%%", row))
		item.btn.Layout.SetWidth(fmt.Sprintf("%0.0f%%", item.transSrc.S))
		item.btn.Layout.SetHeight(fmt.Sprintf("%0.0f%%", item.transSrc.S))
	}
}

//FillRandomItem - fill random empty position on table with number 2 or 4
func (t *Table) FillRandomItem() {
	var empty []int
	for i, item := range t.Items {
		if item.N == 0 {
			empty = append(empty, i)
		}
	}

	if len(empty) == 0 {
		return
	}

	i := empty[t.rand.Intn(len(empty))]
	table.FillItem(i, t.newNum())
}

//FillItem set value to item on table
func (t *Table) FillItem(i, num int) {
	item := t.Items[i]

	item.N = num
	item.transition = true
	item.transSrc.S = 0
	item.transSrc.Y = float32(i / 4 * 25)
	item.transSrc.X = float32(i % 4 * 25)
}

// //newNum return new number 2 or 4
func (t *Table) newNum() int {
	if t.rand.Uint32()%2 == 0 {
		return 2
	}
	return 4
}

func (t *Table) MoveLeft() (moves int) {
	for r := 0; r < 4; r++ {
		l := t.GetRow(r).Calculate()
		moves += t.PutRow(r, l)
	}
	return
}

func (t *Table) MoveRight() (moves int) {
	for r := 0; r < 4; r++ {
		l := t.GetRow(r).Reverse().Calculate().Reverse()
		moves += t.PutRow(r, l)
	}
	return
}

func (t *Table) MoveUp() (moves int) {
	for c := 0; c < 4; c++ {
		l := t.GetCol(c).Calculate()
		moves += t.PutCol(c, l)
	}
	return
}

func (t *Table) MoveDown() (moves int) {
	for c := 0; c < 4; c++ {
		l := t.GetCol(c).Reverse().Calculate().Reverse()
		moves += t.PutCol(c, l)
	}
	return
}

//Line contains in from one row or column, Src it original position
type Line struct {
	Items [4]*Item
	Src   [4]int
}

//GetRow get 4 items from specify row and return Line
func (t *Table) GetRow(r int) (l *Line) {
	l = new(Line)
	r *= 4
	for i := 0; i < 4; i++ {
		l.Items[i] = t.Items[r+i]
		l.Src[i] = r + i

	}
	return
}

//PutRow put line items to table by specify row
func (t *Table) PutRow(r int, l *Line) (moves int) {
	r *= 4
	for i := 0; i < 4; i++ {
		if t.Items[r+i] != l.Items[i] {
			t.Items[r+i] = l.Items[i]
			moves++
		}
	}
	return
}

//GetCol get 4 items for specify column and return Line
func (t *Table) GetCol(c int) (l *Line) {
	l = new(Line)
	for i := 0; i < 4; i++ {
		l.Items[i] = t.Items[c+i*4]
		l.Src[i] = c + i*4
	}
	return
}

//PutCol put line items to table by specify column
func (t *Table) PutCol(c int, l *Line) (moves int) {
	for i := 0; i < 4; i++ {
		if t.Items[c+i*4] != l.Items[i] {
			t.Items[c+i*4] = l.Items[i]
			moves++
		}
	}
	return
}

//Reverse line
func (l *Line) Reverse() *Line {
	l.Items[0], l.Items[1], l.Items[2], l.Items[3] = l.Items[3], l.Items[2], l.Items[1], l.Items[0]
	l.Src[0], l.Src[1], l.Src[2], l.Src[3] = l.Src[3], l.Src[2], l.Src[1], l.Src[0]
	return l
}

//Calculate is important function it move line items, always to left, calcluate item positions and values.
func (l *Line) Calculate() *Line {
	var offset int
	var prev *Item

	for i, item := range l.Items {
		if item.N == 0 {
			continue
		}

		offset, prev = l.LookupPrev(offset, i)
		if prev == nil {
			l.Move(offset, i)
			continue
		}

		if prev.N != item.N {
			offset++
			l.Move(offset, i)
			continue
		}

		prev.N = 0
		item.N *= 2
		l.Move(offset, i)
		offset++
	}

	return l
}

//LookupPrev lookup previous items in this line
func (l *Line) LookupPrev(offset, count int) (int, *Item) {
	for i := offset; i < count; i++ {
		if l.Items[i].N != 0 {
			return i, l.Items[i]
		}
	}
	return offset, nil
}

//Move - swap 2 items and prepare transition values
func (l *Line) Move(dst, src int) {
	if dst == src {
		return
	}
	i := l.Src[src]
	row := float32(i / 4 * 25)
	col := float32(i % 4 * 25)

	l.Items[src].transition = true
	l.Items[src].transSrc = TransSrc{Y: row, X: col, S: 25}
	l.Items[dst], l.Items[src] = l.Items[src], l.Items[dst]
	// l.I[dst], l.I[src] = l.I[src], l.I[dst]
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
		}

		if moves > 0 {
			table.FillRandomItem()
			table.ReDraw()
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

//Lost create lost/restart button
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
