package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type queue struct {
	RedisName string
	Name      string
	Status    string
	Stalled   string
	Work      string
	Depends   string
	Scheduled string
	Recurring string
	Failed    string
}

type queues struct {
	*tview.Table
	filterWord, title string
}

func newQueues(g *Gui) *queues {
	logger.Println("new queues")
	queues := &queues{
		Table:      tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		title:      "Queues List",
		filterWord: "",
	}

	queues.SetTitle(queues.title).
		SetTitleAlign(tview.AlignCenter).
		SetTitleColor(tcell.ColorSeashell.TrueColor()).
		SetBorder(true).
		SetBorderColor(tcell.ColorLightSkyBlue.TrueColor()).
		SetBorderPadding(0, 0, 1, 1)
	queues.setEntries(g)
	queues.setKeybinding(g)
	queues.SetSelectedFunc(func(row, column int) {
		queues.showJobs(g)
	})
	logger.Println("queues created")
	return queues
}

func (q *queues) showJobs(g *Gui) {

}

func (q *queues) name() string {
	return "queues"
}

func (q *queues) canSelect() bool {
	return true
}

func (g *Gui) queuePanel() *queues {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "queues" {
			return panel.(*queues)
		}
	}
	return nil
}

func (q *queues) entries(g *Gui) {
	logger.Println("Entries")
	g.state.resources.queues = make([]*queue, 0)
	logger.Println("Forming entries for: " + fmt.Sprint(qlessClients))

	for _, qlessClient := range qlessClients {
		for _, qu := range qlessClient.Queues {
			if strings.Index(qu.Name, q.filterWord) == -1 {
				continue
			}

			data := qu.Counts()
			queue := &queue{
				RedisName: qlessClient.Name,
				Name:      qu.Name,
				Status:    fmt.Sprintf("%t", qu.Paused()),
				Stalled:   fmt.Sprint(data["stalled"]),
				Work:      fmt.Sprint(data["work"]),
				Depends:   fmt.Sprint(data["depends"]),
				Scheduled: fmt.Sprint(data["scheduled"]),
				Recurring: fmt.Sprint(data["recurring"]),
				Failed:    fmt.Sprint(qu.FailedCount()),
			}
			g.state.resources.queues = append(g.state.resources.queues, queue)

			logger.Println("Queue entry formed: " + fmt.Sprint(queue))
		}
	}
}

func (q *queues) setEntries(g *Gui) {
	q.entries(g)
	q.updateTitle(g)
	logger.Printf("setEntries: %s", fmt.Sprint(q))
	table := q.Clear()
	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke).
		Foreground(tcell.ColorBlack))

	headers := []string{
		"Qless",
		"Name",
		"Status",
		"Stalled",
		"Work",
		"Depends",
		"Scheduled",
		"Recurring",
		"Failed",
	}

	for i, header := range headers {
		table.SetCell(0, i, &tview.TableCell{
			Text:            header,
			NotSelectable:   true,
			Align:           tview.AlignLeft,
			Color:           tcell.ColorWhite.TrueColor(),
			BackgroundColor: tcell.ColorBlack,
			Attributes:      tcell.AttrBold,
		})
	}

	// Sort queues by RedisName -> QueueName
	sort.Slice(g.state.resources.queues, func(i, j int) bool {
		if g.state.resources.queues[i].RedisName != g.state.resources.queues[j].RedisName {
			return g.state.resources.queues[i].RedisName < g.state.resources.queues[j].RedisName
		}

		return g.state.resources.queues[i].Name < g.state.resources.queues[j].Name
	})

	for i, queue := range g.state.resources.queues {
		textColor := tcell.ColorLightGreen.TrueColor()
		if queue.Status == "true" {
			textColor = tcell.ColorOrange.TrueColor()
		}
		table.SetCell(i+1, 0, tview.NewTableCell(queue.RedisName).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue))

		table.SetCell(i+1, 1, tview.NewTableCell(queue.Name).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(3).
			SetReference(queue.Name))

		pausedCell := tview.NewTableCell("active").
			SetMaxWidth(1).
			SetTextColor(textColor).
			SetBackgroundColor(tcell.ColorBlack).
			SetExpansion(1).
			SetReference(queue.Status)

		if queue.Status == "true" {
			pausedCell.SetText("paused")
		}

		table.SetCell(i+1, 2, pausedCell)

		table.SetCell(i+1, 3, tview.NewTableCell(queue.Work).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Work))

		table.SetCell(i+1, 4, tview.NewTableCell(queue.Scheduled).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Scheduled))

		table.SetCell(i+1, 5, tview.NewTableCell(queue.Stalled).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Stalled))

		table.SetCell(i+1, 6, tview.NewTableCell(queue.Depends).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Depends))

		table.SetCell(i+1, 7, tview.NewTableCell(queue.Recurring).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Recurring))

		table.SetCell(i+1, 8, tview.NewTableCell(queue.Failed).
			SetTextColor(textColor).
			SetMaxWidth(1).
			SetExpansion(1).
			SetReference(queue.Failed))
	}
}

func (q *queues) focus(g *Gui) {
	q.SetSelectable(true, false)
	g.app.SetFocus(q)
}

func (q *queues) unfocus() {
	q.SetSelectable(false, false)
}

func (q *queues) updateEntries(g *Gui) {
	logger.Println("Updating entries")
	g.app.QueueUpdateDraw(func() {
		logger.Println("Set entries")
		q.setEntries(g)
		g.leftMenuPanel().setEntries(g)
	})
}

func (q *queues) setFilterWord(word string) {
	q.filterWord = word
}

func (q *queues) updateTitle(g *Gui) {
	title := q.title
	title += fmt.Sprintf(" [green::b](%d)", q.GetRowCount()-1)
	if len(q.filterWord) > 0 {
		title += fmt.Sprintf(" [orange::b][[yellow::b]/%s[orange::b]]", q.filterWord)
	}
	q.SetTitle(title)
}

func (q *queues) monitoringQueues(g *Gui) {
	logger.Println("monitoring queues starting")
	ticker := time.NewTicker(1 * time.Second)

LOOP:
	for {
		select {
		case <-ticker.C:
			q.updateEntries(g)
			logger.Println("entries updated")
		case <-g.state.stopChans["container"]:
			ticker.Stop()
			logger.Println("ticker stopped")
			break LOOP
		}
	}
}

func (g *Gui) selectedQueue() *queue {
	row, _ := g.queuePanel().GetSelection()

	if len(g.state.resources.queues) == 0 || len(g.state.resources.queues) < row {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	logger.Println("Select queue -> row: " + fmt.Sprint(row))
	return g.queuePanel().GetCell(row, 0).GetReference().(*queue)
}

func (q *queues) setKeybinding(g *Gui) {
	q.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		switch event.Rune() {
		case 'u':
			g.continueQueue()
		case 'p':
			g.pauseQueue()
		}
		return event
	})
}
