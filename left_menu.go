package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/kyokomi/emoji/v2"
	"github.com/rivo/tview"
)

type menuItem struct {
	icon string
	text string
	page string
}

type leftMenu struct {
	*tview.Table
	filterWord string
}

func (m *leftMenu) name() string {
	return string("leftMenu")
}

func newLeftMenu(g *Gui) *leftMenu {
	leftMenu := &leftMenu{
		Table: tview.NewTable().SetSelectable(false, false).Select(0, 0).SetFixed(1, 1),
	}

	leftMenu.SetTitle("Left Menu").
		SetBorder(false)
	leftMenu.SetBorderPadding(1, 1, 1, 1)
	leftMenu.setEntries(g)
	leftMenu.setKeybinding(g)

	leftMenu.SetSelectedFunc(func(row, column int) {
		g.stopMonitoring()
		switch g.state.resources.menuItems[row-1].page {
		case "showQueuesPage":
			g.showQueuesPage()
		case "showFailedJobGroupsPage":
			g.showFailedJobGroupsPage()
		}
	})
	logger.Print("leftMenu created: ", fmt.Sprint(leftMenu))

	return leftMenu
}

func (m *leftMenu) canSelect() bool {
	return true
}

func (g *Gui) leftMenuPanel() *leftMenu {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "leftMenu" {
			return panel.(*leftMenu)
		}
	}
	return nil
}

func (m *leftMenu) entries(g *Gui) {
	g.state.resources.menuItems = []*menuItem{
		{icon: emoji.Sprint(":repeat:"), text: "Queues", page: "showQueuesPage"},
		{icon: emoji.Sprint(":ladybug:"), text: "Failed jobs", page: "showFailedJobGroupsPage"},
	}
}

func (m *leftMenu) setEntries(g *Gui) {
	m.entries(g)
	logger.Printf("setEntries: %s", fmt.Sprint(m))

	table := m.Clear()
	logger.Printf("setEntries [Clear]: %s", fmt.Sprint(table))
	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke).
		Foreground(tcell.ColorBlack))

	headers := []string{
		emoji.Sprint(":beer: MENU"),
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignCenter,
			Color:           tcell.ColorSeashell.TrueColor(),
			BackgroundColor: tcell.ColorBlack,
			Attributes:      tcell.AttrBold,
		})
	}

	for i, item := range g.state.resources.menuItems {
		table.SetCell(i+1, 0, tview.NewTableCell(item.icon+item.text).
			SetTextColor(tcell.ColorWhiteSmoke.TrueColor()).
			SetMaxWidth(0).
			SetExpansion(1))
	}
}

func (m *leftMenu) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		logger.Println("Set entries")
		m.setEntries(g)
	})
}

func (m *leftMenu) setKeybinding(g *Gui) {
	m.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		return event
	})
}

func (m *leftMenu) focus(g *Gui) {
	m.SetSelectable(true, false)
	g.app.SetFocus(m)
}

func (m *leftMenu) unfocus() {
	m.SetSelectable(false, false)
}

func (m *leftMenu) setFilterWord(word string) {
	m.filterWord = word
}

func (m *leftMenu) updateTitle(g *Gui) {
}
