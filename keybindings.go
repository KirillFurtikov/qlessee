package main

import (
	"context"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (g *Gui) setGlobalKeybinding(event *tcell.EventKey) {
	switch event.Rune() {
	case 'q':
		g.Stop()
	case '/':
		g.filter()
	}

	switch event.Key() {
	case tcell.KeyEsc:
		g.clearFilter()
	case tcell.KeyTab:
		g.nextPanel()
	case tcell.KeyBacktab:
		g.prevPanel()
	case tcell.KeyRight:
		g.nextPanel()
	case tcell.KeyLeft:
		g.prevPanel()
	}
}

func (g *Gui) clearFilter() {
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	currentPanel.setFilterWord("")
}

func (g *Gui) filter() {
	logger.Println("Filtering")
	currentPanel := g.state.panels.panel[g.state.panels.currentPanel]
	if currentPanel.name() == "tasks" {
		logger.Println("Current panel is 'tasks' - return")
		return
	}

	if currentPanel.name() == "leftMenu" {
		logger.Println("Current panel is 'leftMenu' - return")
		return
	}
	logger.Printf("Configure search input for panel: %s", currentPanel)
	currentPanel.setFilterWord("")
	logger.Println("Filter Word set to empty string")

	viewName := "filter"
	searchInput := tview.NewInputField().SetLabel("Word")
	searchInput.SetLabelColor(tcell.ColorBlue)
	searchInput.SetFieldBackgroundColor(tcell.ColorBlack)
	searchInput.SetLabelWidth(6)
	searchInput.SetPlaceholder("Enter queue name")
	searchInput.SetTitle("Filter")
	searchInput.SetTitleAlign(tview.AlignLeft)
	searchInput.SetBorder(true)
	searchInput.SetBorderAttributes(tcell.AttrMask(1))
	logger.Println("Search input configured")

	closeSearchInput := func() {
		g.closeAndSwitchPanel(viewName, g.state.panels.panel[g.state.panels.currentPanel].name())
	}

	searchInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			closeSearchInput()
		}
	})

	searchInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			closeSearchInput()
		}
		return event
	})

	searchInput.SetChangedFunc(func(text string) {
		currentPanel.setFilterWord(text)
	})
	logger.Printf("Switch to %v", viewName)

	g.pages.AddAndSwitchToPage(viewName, g.modal(searchInput, 80, 3), true).ShowPage("main")
}

func (g *Gui) nextPanel() {
	logger.Println("nextPanel")

	idx := (g.state.panels.currentPanel + 1) % len(g.state.panels.panel)

	logger.Printf("nextPanel [idx: %s]", fmt.Sprint(idx))
	logger.Printf("nextPanel [canSelect: %s]", fmt.Sprint(g.state.panels.panel))

	g.switchPanel(g.state.panels.panel[idx].name())
	if g.state.panels.panel[idx].canSelect() != true {
		logger.Printf("nextPanel: selected panel [%s] can not be selected, call it again", fmt.Sprint(idx))
		g.nextPanel()
	}
}

func (g *Gui) prevPanel() {
	g.state.panels.currentPanel--

	if g.state.panels.currentPanel < 0 {
		g.state.panels.currentPanel = len(g.state.panels.panel) - 1
	}

	idx := (g.state.panels.currentPanel) % len(g.state.panels.panel)

	logger.Printf("nextPanel [idx: %s]", fmt.Sprint(idx))
	logger.Printf("nextPanel [canSelect: %s]", fmt.Sprint(g.state.panels.panel))

	g.switchPanel(g.state.panels.panel[idx].name())
	if g.state.panels.panel[idx].canSelect() != true {
		logger.Printf("nextPanel: selected panel [%s] can not be selected, call it again", fmt.Sprint(idx))
		g.prevPanel()
	}
}

func (g *Gui) pauseQueue() {
	queue := g.selectedQueue()
	if queue == nil {
		return
	}
	logger.Println("Pause queue starting: " + fmt.Sprint(queue))

	g.startTask(fmt.Sprintf("Pause queue [blue::b]%s [yellow::b][[::b]%s]", queue.RedisName, queue.Name), func(ctx context.Context) error {
		if err := qlessClients[queue.RedisName].GetQueue(queue.Name).Pause().Err(); err != nil {
			return err
		}

		g.queuePanel().updateEntries(g)

		return nil
	})
	logger.Println("Pause queue completed: " + fmt.Sprint(queue))
}

func (g *Gui) continueQueue() {
	queue := g.selectedQueue()
	if queue == nil {
		return
	}

	g.startTask(fmt.Sprintf("Continue queue [blue::b]%s [yellow::b][[::b]%s]", queue.RedisName, queue.Name), func(ctx context.Context) error {
		if err := qlessClients[queue.RedisName].GetQueue(queue.Name).Continue().Err(); err != nil {
			return err
		}

		g.queuePanel().updateEntries(g)

		return nil
	})
}
