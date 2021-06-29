package main

import (
	"fmt"

	qless "github.com/KirillFurtikov/qlessee/pkg/qless"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type failedJobGroup struct {
	redisName string
	klass     string
	counts    uint64
	jobs      []*qless.Job
}

type failedJobGroups struct {
	*tview.Table
	filterWord string
}

func (f *failedJobGroups) name() string {
	return string("failedJobGroups")
}

func newFailedJobGroups(g *Gui) *failedJobGroups {
	failedJobGroups := &failedJobGroups{
		Table: tview.NewTable().SetSelectable(false, false).Select(0, 0).SetFixed(1, 1),
	}

	failedJobGroups.SetTitle("Failed Jobs").
		SetBorder(false).
		SetBorderColor(tcell.ColorLightSkyBlue.TrueColor()).
		SetTitleColor(tcell.ColorSeashell.TrueColor())
	failedJobGroups.SetBorderPadding(0, 0, 1, 1)
	failedJobGroups.setEntries(g)
	failedJobGroups.setKeybinding(g)
	failedJobGroups.SetSelectedFunc(func(row, column int) {
		g.showJobsPage(failedJobGroups.GetCell(row, column).Reference.(*failedJobGroup).klass)
	})
	logger.Print("failedJobs created: ", fmt.Sprint(failedJobGroups))

	return failedJobGroups
}

func (f *failedJobGroups) canSelect() bool {
	return true
}

func (g *Gui) failedJobGroupsPanel() *failedJobGroups {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "failedJobGroups" {
			return panel.(*failedJobGroups)
		}
	}
	return nil
}

func (f *failedJobGroups) entries(g *Gui) {
	logger.Println("Entries")
	g.state.resources.failedJobGroups = make([]*failedJobGroup, 0)
	logger.Println("Forming entries for: " + fmt.Sprint(qlessClients))

	for _, qlessClient := range qlessClients {
		for klass, counts := range qlessClient.Jobs().GetFailedCounts() {
			job := &failedJobGroup{redisName: qlessClient.Name, klass: klass, counts: counts}
			g.state.resources.failedJobGroups = append(g.state.resources.failedJobGroups, job)

			logger.Println("failedJobGroups entry formed: " + fmt.Sprint(job))
		}
	}
}

func (f *failedJobGroups) setEntries(g *Gui) {
	f.entries(g)
	logger.Printf("setEntries: %s", fmt.Sprint(f))

	table := f.Clear()
	logger.Printf("setEntries [Clear]: %s", fmt.Sprint(table))
	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke).
		Foreground(tcell.ColorBlack)).
		SetBorder(true)

	headers := []string{
		"Qless",
		"Job Klass",
		"Count",
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

	for i, item := range g.state.resources.failedJobGroups {
		table.SetCell(i+1, 0, tview.NewTableCell(item.redisName).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))

		table.SetCell(i+1, 1, tview.NewTableCell(item.klass).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprint(item.counts)).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1))
	}
}

func (f *failedJobGroups) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		logger.Println("Set entries")
		f.setEntries(g)
	})
}

func (f *failedJobGroups) setKeybinding(g *Gui) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		return event
	})
}

func (f *failedJobGroups) focus(g *Gui) {
	f.SetSelectable(true, false)
	g.app.SetFocus(f)
}

func (f *failedJobGroups) unfocus() {
	f.SetSelectable(false, false)
}

func (f *failedJobGroups) setFilterWord(word string) {
	f.filterWord = word
}

func (f *failedJobGroups) updateTitle(g *Gui) {
}
