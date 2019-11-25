package main

import (
	"fmt"

	"github.com/marcusolsson/tui-go"
)

const filepath = "d:/files/log.txt"
const visibleLinesCount = 1
const marginLinesCount = 1

func newUI(filename string) tui.Widget {

	filenameLabel := tui.NewLabel(filename)
	headersBox := tui.NewHBox(filenameLabel)
	headersBox.SetBorder(true)
	//	headersBox.si

	fileLines := tui.NewTable(0, 0)
	for i := 0; i != 10; i++ {
		fileLines.AppendRow(tui.NewLabel(fmt.Sprintf("line %v", i)))
	}

	debugLines := tui.NewTable(0, 0)
	for i := 0; i != 5; i++ {
		fileLines.AppendRow(tui.NewLabel(fmt.Sprintf("debug %v", i)))
	}

	return tui.NewVBox(headersBox, fileLines, debugLines)
}

func main() {
	fmt.Println("Started...")

	rootWidget := newUI(filepath)
	ui, _ := tui.New(rootWidget)

	ui.SetKeybinding("Ctrl+C", func() { ui.Quit() })
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	ui.Run()
}
