package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marcusolsson/tui-go"
)

const filepath = "d:/files/log.txt"
const displayedLines = 3

var lineIndex uint
var labels map[int]*tui.Label
var input *tui.Entry
var tf *TextFile
var debugText *tui.Label

func init() {
	labels = make(map[int]*tui.Label)
}

func createRootUIWidget() tui.Widget {
	tf.goTo(lineIndex)
	t := tui.NewTable(0, 0)
	t.AppendRow(tui.NewLabel("First line........."))
	for i := 0; i != displayedLines; i++ {
		l := tui.NewLabel("[empty line]")
		t.AppendRow(l)
		labels[i] = l
	}
	t.AppendRow(tui.NewLabel("Last line........."))
	t.OnSelectionChanged(func(t *tui.Table) {
		now := time.Now()
		if t.Selected() == 0 {
			t.SetSelected(1)
			lineIndex--
			tf.goTo(lineIndex)
		}
		if t.Selected() == displayedLines+1 {
			t.SetSelected(displayedLines)
			lineIndex++
			tf.goTo(lineIndex)
		}
		updateDisplayedLines()
		debugText.SetText(fmt.Sprintf("OnSelectionChanged took %v", time.Since(now)))
	})

	t.SetFocused(true)

	debugText = tui.NewLabel("debug text")
	debugBox := tui.NewVBox(debugText)
	debugBox.SetBorder(true)

	input = tui.NewEntry()
	input.SetText("enter text here")
	input.SetFocused(true)
	t.SetSelected(1)
	uiRoot := tui.NewVBox(
		debugBox,
		t,
		tui.NewSpacer(),
		input)
	return uiRoot
}

func updateDisplayedLines() {
	for i := 0; i != displayedLines; i++ {
		l, _ := labels[i]
		_, ok := tf.CachedLines[uint(i)+lineIndex]
		if ok {
			line, _ := tf.CachedLines[uint(i)+lineIndex]
			l.SetText(line.Contents)
		} else {
			l.SetText("")
		}
	}

	input.SetText(fmt.Sprintf("Cache %v - %v", tf.startingLineIndex, tf.cacheSize))
}

func printLinesCount() {
	f, err := os.Open(filepath)
	if err != nil {
		log.Panicf("Open() failed: %v", err.Error())
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var lines int
	for s.Scan() {
		lines++
	}
	debugText.SetText(fmt.Sprintf("Lines in file: %v", lines))
}

func main() {
	fmt.Println("Started...")
	tf = NewTextFile(filepath, 3)
	//	fmt.Println(tf)

	tf.goTo(3000000)
	tf.goTo(0)
	tf.goTo(3000000)
	//	fmt.Println(tf)

	tf.goTo(1)
	//	fmt.Println(tf)

	tf.goTo(3)
	//	fmt.Println(tf)

	tf.goTo(1)
	//	fmt.Println(tf)

	/*
		ui, err := tui.New(createRootUIWidget())
		if err != nil {
			log.Panicf("unable to create UI: %s", err.Error())
		}

		printLinesCount()

		ui.SetKeybinding("Ctrl+C", func() { ui.Quit() })
		ui.SetKeybinding("Esc", func() { ui.Quit() })

		if err := ui.Run(); err != nil {
			log.Fatal(err)
		}
	*/

}
