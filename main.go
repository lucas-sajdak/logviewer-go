package main

import (
	"fmt"
	"log"

	"github.com/marcusolsson/tui-go"
)

const displayedLines = 3

var testText [100]string

var lineIndex = 0

var labels map[int]*tui.Label
var input *tui.Entry
var tf *TextFile

func init() {
	for i := 0; i != 100; i++ {
		testText[i] = fmt.Sprintf("Line: %v", i)
	}

	labels = make(map[int]*tui.Label)
}

func createRootUIWidget() tui.Widget {
	tf.gotoLine(0)
	t := tui.NewTable(0, 0)
	t.AppendRow(tui.NewLabel("First line........."))
	for i := 0; i != displayedLines; i++ {
		l := tui.NewLabel("[empty line]")
		t.AppendRow(l)
		labels[i] = l
	}
	t.AppendRow(tui.NewLabel("Last line........."))
	t.OnSelectionChanged(func(t *tui.Table) {
		if t.Selected() == 0 {
			t.SetSelected(1)
			lineIndex--
			tf.gotoLine(lineIndex)
		}

		if t.Selected() == displayedLines+1 {
			t.SetSelected(displayedLines)
			lineIndex++
			tf.gotoLine(lineIndex)
		}
		updateDisplayedLines()
	})

	t.SetFocused(true)
	input = tui.NewEntry()
	input.SetText("enter text here")
	input.SetFocused(true)
	t.SetSelected(1)
	uiRoot := tui.NewVBox(
		t,
		tui.NewSpacer(),
		input)
	return uiRoot
}

func updateDisplayedLines() {
	for i := 0; i != displayedLines; i++ {
		l, _ := labels[i]
		text, ok := tf.Lines[i+lineIndex]
		if ok {
			l.SetText(text)
		} else {
			l.SetText("")
		}
	}
}

func main() {

	fmt.Print("Started...")

	tf = NewTextFile("./log.txt", 9)

	//	tf.gotoLine(0)
	//	fmt.Printf("%v", tf)

	//	tf.gotoLine(1)
	//	fmt.Printf("%v", tf)

	ui, err := tui.New(createRootUIWidget())
	if err != nil {
		log.Panicf("unable to create UI: %s", err.Error())
	}

	updateDisplayedLines()

	ui.SetKeybinding("Ctrl+C", func() { ui.Quit() })
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}
