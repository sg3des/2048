package main

import (
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func TestNewWindow(t *testing.T) {
	err := NewWindow("2048", 600, 600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewTable(t *testing.T) {
	table = NewTable()
}

func TestFillItem(t *testing.T) {
	table.FillItem(1, 2)
	if table.Items[1].N != 2 {
		t.Fatal("failed fill item")
	}
}

func TestMoveLeft(t *testing.T) {
	table = NewTable()
	table.FillItem(1, 2)

	table.MoveLeft()
	if table.Items[0].N != 2 {
		t.Error("failed move left")
	}
	if table.Items[1].N != 0 {
		t.Error("failed move left")
	}

	table.FillItem(3, 2)

	table.MoveLeft()
	if table.Items[0].N != 4 {
		t.Fatal("failed summ on 0 position")
	}
	for i := 1; i < 16; i++ {
		if table.Items[i].N != 0 {
			t.Fatalf("item %d should be 0", i)
		}
	}

	table.FillItem(2, 4)
	table.MoveLeft()
	if table.Items[0].N != 8 {
		t.Fatal("failed summ on 0 position to 8")
	}

	for i := 1; i < 16; i++ {
		if table.Items[i].N != 0 {
			t.Fatalf("item %d should be 0", i)
		}
	}
}

func TestMoveRight(t *testing.T) {
	table = NewTable()
	table.FillItem(0, 2)
	table.FillItem(1, 2)
	table.MoveRight()
	if table.Items[3].N != 4 {
		t.Error("failed move right, sum on 3 position should be 4")
	}

	for i := 0; i < 3; i++ {
		if table.Items[i].N != 0 {
			t.Errorf("failed move right, position %d should be 0, but whis is %d", i, table.Items[i].N)
		}
	}
}
