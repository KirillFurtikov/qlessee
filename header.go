package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type hotkeyPrompt struct {
	action, key string
}

type header struct {
	*tview.Table
	filterWord string
}

func (h *header) name() string {
	return string("header")
}

func newHeader(g *Gui) *header {
	header := &header{
		Table: tview.NewTable().SetSelectable(false, false).Select(0, 0).SetFixed(1, 1),
	}

	header.SetBorder(true).
		SetBorderColor(tcell.ColorLightSkyBlue.TrueColor())
	header.SetBorderPadding(0, 0, 1, 1)
	header.setEntries(g)
	header.setKeybinding(g)
	logger.Print("header created: ", fmt.Sprint(header))

	return header
}

func (h *header) canSelect() bool {
	return false
}

func (g *Gui) headerPanel() *leftMenu {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "header" {
			return panel.(*leftMenu)
		}
	}
	return nil
}

func (h *header) entries(g *Gui) {
	g.state.resources.hotkeyPrompts = make([]*hotkeyPrompt, 0)
	hotkeys := []*hotkeyPrompt{
		{action: "Exit", key: "<q>"},
		{action: "Filter:", key: "</>"},
		{action: "Reset filter:", key: "<Esc>"},
		{action: "Pause/unpause queue", key: "<p>/<u>"},
		{action: "Select next panel:", key: "<Tab>/<Right>"},
		{action: "Select previous panel:", key: "<Shift+Tab>/<Left>"},
	}
	g.state.resources.hotkeyPrompts = append(g.state.resources.hotkeyPrompts, hotkeys...)
}

func (h *header) setEntries(g *Gui) {
	h.entries(g)
	logger.Printf("setEntries: %s", fmt.Sprint(h))

	table := h.Clear()

	logger.Printf("setEntries [Clear]: %s", fmt.Sprint(table))

	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke).
		Foreground(tcell.ColorBlack))

	groups := make([][]*hotkeyPrompt, 0)

	for 3 < len(g.state.resources.hotkeyPrompts) {
		g.state.resources.hotkeyPrompts, groups = g.state.resources.hotkeyPrompts[3:], append(groups, g.state.resources.hotkeyPrompts[0:3:3])
	}

	groups = append(groups, g.state.resources.hotkeyPrompts)

	column := 0
	for i, prompts := range groups {
		if i != 0 {
			column = column + 2
		}
		for j, item := range prompts {
			table.SetCell(j, column, tview.NewTableCell(item.action).
				SetTextColor(tcell.ColorSlateGray.TrueColor()).
				SetAlign(tview.AlignRight).
				SetMaxWidth(0).
				SetExpansion(0))

			table.SetCell(j, column+1, tview.NewTableCell(item.key).
				SetTextColor(tcell.ColorDeepSkyBlue.TrueColor()).
				SetMaxWidth(0).
				SetExpansion(0))
		}
	}

}

func (h *header) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		logger.Println("Set entries")
		h.setEntries(g)
	})
}

func (h *header) setKeybinding(g *Gui) {
}

func (h *header) focus(g *Gui) {

}

func (h *header) unfocus() {

}

func (h *header) setFilterWord(word string) {

}

func (h *header) updateTitle(g *Gui) {
}
