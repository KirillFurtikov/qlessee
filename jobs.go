package main

import (
	"fmt"

	qless "github.com/KirillFurtikov/qlessee/pkg/qless"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type job struct {
	*qless.Job
}

type jobs struct {
	*tview.Table
	group      string
	filterWord string
}

func (f *jobs) name() string {
	return string("jobs")
}

func newJobs(g *Gui, group string) *jobs {
	jobs := &jobs{
		Table: tview.NewTable().SetSelectable(false, false).Select(0, 0).SetFixed(1, 1),
		group: group,
	}

	jobs.SetTitle("Jobs").
		SetBorder(false).
		SetBorderColor(tcell.ColorLightSkyBlue.TrueColor()).
		SetTitleColor(tcell.ColorSeashell.TrueColor())
	jobs.SetBorderPadding(0, 0, 1, 1)
	jobs.setEntries(g)
	jobs.setKeybinding(g)
	jobs.SetSelectedFunc(func(row, column int) {
		selectedJob := jobs.selectedJob(g)

		text := tview.TranslateANSI(string(selectedJob.Load().Pretty()))
		preview := tview.NewTextView().
			SetText(text).
			SetDynamicColors(true).
			SetWrap(true).
			SetWordWrap(true)
		preview.SetTitle("Job " + selectedJob.JID)
		preview.SetTitleAlign(tview.AlignCenter)
		preview.SetBorder(true)
		preview.SetBorderAttributes(tcell.AttrBold)
		preview.SetBorderColor(tcell.ColorLightSkyBlue.TrueColor()).SetBorderPadding(0, 0, 1, 1)

		viewName := "jobInfo"
		closePreview := func() {
			g.closeAndSwitchPanel(viewName, g.state.panels.panel[g.state.panels.currentPanel].name())
		}

		preview.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEsc {
				closePreview()
			}
			return event
		})

		g.pages.AddAndSwitchToPage(viewName, g.modal(preview, 150, 50), true).ShowPage("main")

	})
	logger.Print("jobs created: ", fmt.Sprint(jobs))

	return jobs
}

func (f *jobs) canSelect() bool {
	return true
}

func (g *Gui) jobsPanel() *jobs {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "jobs" {
			return panel.(*jobs)
		}
	}
	return nil
}

func (f *jobs) entries(g *Gui) {
	logger.Println("Entries")
	g.state.resources.jobs = make([]*job, 0)
	logger.Println("Forming entries for: " + fmt.Sprint(qlessClients))

	for _, qlessClient := range qlessClients {
		for _, failedJobs := range qlessClient.Jobs().GetFailedJobs() {
			for _, qJob := range failedJobs {
				g.state.resources.jobs = append(g.state.resources.jobs, &job{qJob.Load()})
			}
		}
	}
}

func (j *jobs) selectedJob(g *Gui) *job {
	row, _ := j.GetSelection()

	if len(g.state.resources.jobs) == 0 || len(g.state.resources.jobs) < row {
		return nil
	}
	if row-1 < 0 {
		return nil
	}

	return g.jobsPanel().GetCell(row, 0).GetReference().(*job)
}

func (f *jobs) setEntries(g *Gui) {
	f.entries(g)
	logger.Printf("setEntries: %s", fmt.Sprint(f))

	table := f.Clear()
	logger.Printf("setEntries [Clear]: %s", fmt.Sprint(table))
	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke).
		Foreground(tcell.ColorBlack)).
		SetBorder(true)

	headers := []string{
		"JID",
		"Klass",
		"Queue",
		"Worker",
		"Retries",
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

	for i, item := range g.state.resources.jobs {
		table.SetCell(i+1, 0, tview.NewTableCell(item.JID).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))

		table.SetCell(i+1, 1, tview.NewTableCell(item.Klass).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))

		table.SetCell(i+1, 2, tview.NewTableCell(fmt.Sprint(item.Queue)).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))

		table.SetCell(i+1, 3, tview.NewTableCell(fmt.Sprint(item.Worker)).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))

		table.SetCell(i+1, 4, tview.NewTableCell(fmt.Sprint(item.Retries)).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(0).
			SetExpansion(1).
			SetReference(item))
	}
}

func (f *jobs) updateEntries(g *Gui) {
	g.app.QueueUpdateDraw(func() {
		logger.Println("Set entries")
		f.setEntries(g)
	})
}

func (f *jobs) setKeybinding(g *Gui) {
	f.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		return event
	})
}

func (f *jobs) focus(g *Gui) {
	f.SetSelectable(true, false)
	g.app.SetFocus(f)
}

func (f *jobs) unfocus() {
	f.SetSelectable(false, false)
}

func (f *jobs) setFilterWord(word string) {
	f.filterWord = word
}

func (f *jobs) updateTitle(g *Gui) {
}
