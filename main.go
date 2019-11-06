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
const visibleLinesCount = 1
const marginLinesCount = 1

var firstVisibleLineIndex uint
var topLabel *tui.Label
var labels map[int]*tui.Label
var input *tui.Entry
var tf *TextFile
var debugText *tui.Label

type tabWidget struct {
	*tui.Box

	filename *tui.Label
	filter   *tui.Label
	selected *tui.Label
}

func init() {
	labels = make(map[int]*tui.Label)
}

func createRootUIWidget() tui.Widget {

	tw := &tabWidget{}
	tw.filename = tui.NewLabel("filename")

	tw.filter = tui.NewLabel("empty")

	header := tui.NewHBox(tw.filename, tw.filter, tui.NewSpacer())
	header.SetBorder(true)

	t := tui.NewTable(0, 0)
	topLabel = tui.NewLabel("First line.........")
	t.AppendRow(topLabel)
	for i := 0; i != visibleLinesCount; i++ {
		l := tui.NewLabel("[empty line]")
		t.AppendRow(l)
		labels[i] = l
	}
	t.AppendRow(tui.NewLabel("Last line........."))
	t.OnSelectionChanged(func(t *tui.Table) {
		if t.Selected() == 0 {
			t.SetSelected(1)
			if firstVisibleLineIndex > 0 {
				firstVisibleLineIndex--
			}
			updateCache()
		}
		if t.Selected() == visibleLinesCount+1 {
			t.SetSelected(visibleLinesCount)
			firstVisibleLineIndex++
			updateCache()
		}
		updateVisibleLines()
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
		header,
		debugBox,
		t,
		tui.NewSpacer(),
		input)

	return uiRoot
}

func updateVisibleLines() {
	topLabel.SetText(fmt.Sprintf("[Line %v/...]", firstVisibleLineIndex))
	for i := 0; i != visibleLinesCount; i++ {
		l, _ := labels[i]
		_, ok := tf.CachedLines[uint(i)+firstVisibleLineIndex]
		if ok {
			line, _ := tf.CachedLines[uint(i)+firstVisibleLineIndex]
			l.SetText(line.Contents)
		} else {
			l.SetText("")
		}
	}
	input.SetText(fmt.Sprintf("Cached [%v-%v]", tf.startingLineIndex, tf.startingLineIndex+tf.cacheSize-1))
}

func updateCache() {
	now := time.Now()
	if firstVisibleLineIndex < tf.startingLineIndex ||
		firstVisibleLineIndex+visibleLinesCount > tf.startingLineIndex+tf.cacheSize {
		newStartingLine := uint(Max(0, int64(firstVisibleLineIndex-marginLinesCount)))
		tf.goTo(newStartingLine)
	}

	elapsed := time.Since(now)
	debugText.SetText(fmt.Sprintf("took %s", elapsed))
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
}

func main() {
	fmt.Println("Started...")

	f, err := os.OpenFile(filepath, os.O_RDONLY, os.ModeCharDevice)
	if err != nil {
		log.Panicf("OpenFile() failed: %v", err.Error())
	}

	tf = NewTextFile(f, visibleLinesCount+2*marginLinesCount)

	ui, err := tui.New(createRootUIWidget())
	if err != nil {
		log.Panicf("unable to create UI: %s", err.Error())
	}

	updateVisibleLines()

	ui.SetKeybinding("Ctrl+C", func() { ui.Quit() })
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	ui.SetKeybinding("Ctrl+Down", func() {
		firstVisibleLineIndex += 10000
		updateCache()
		updateVisibleLines()
	})

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}

}
