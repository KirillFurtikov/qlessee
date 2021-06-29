package main

import (
	"github.com/rivo/tview"
)

type panels struct {
	currentPanel int
	panel        []panel
}

type resources struct {
	tasks           []*task
	queues          []*queue
	menuItems       []*menuItem
	failedJobGroups []*failedJobGroup
	jobs            []*job
	hotkeyPrompts   []*hotkeyPrompt
}

type state struct {
	resources resources
	panels    panels
	stopChans map[string]chan int
}

// Gui have all panels
type Gui struct {
	app   *tview.Application
	pages *tview.Pages
	state *state
}

// NewGui create new gui
func NewGui() *Gui {
	return &Gui{
		app:   tview.NewApplication().EnableMouse(true),
		state: newState(),
	}
}

// Start start application
func (g *Gui) Start() error {
	logger.Println("Gui start")
	g.showQueuesPage()

	if err := g.app.Run(); err != nil {
		g.app.Stop()
		return err
	}

	return nil
}

func (g *Gui) stopMonitoring() {
	for c := range g.state.stopChans {
		logger.Printf("%s -> %d", c, g.state.stopChans[c])
		g.state.stopChans[c] <- 1
	}
}

func (g *Gui) switchPanel(panelName string) {
	for i, panel := range g.state.panels.panel {
		if panel.name() == panelName {
			panel.focus(g)
			g.state.panels.currentPanel = i
		} else {
			panel.unfocus()
		}
	}
}

func newState() *state {
	return &state{
		stopChans: make(map[string]chan int),
	}
}

// Stop stop application
func (g *Gui) Stop() error {
	g.stopMonitoring()
	g.app.Stop()
	return nil
}

func (g *Gui) currentPanel() panel {
	return g.state.panels.panel[g.state.panels.currentPanel]
}

func (g *Gui) closeAndSwitchPanel(removePanel, switchPanel string) {
	g.pages.RemovePage(removePanel).ShowPage("main")
	g.switchPanel(switchPanel)
}

func (g *Gui) modal(p tview.Primitive, width, height int) tview.Primitive {
	logger.Println("Init new modal primitive")
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (g *Gui) showQueuesPage() *tview.Flex {
	g.state.stopChans = make(map[string]chan int)
	header := newHeader(g)
	leftMenu := newLeftMenu(g)
	queues := newQueues(g)
	tasks := newTasks(g)

	g.state.panels.panel = []panel{
		header, leftMenu, queues, tasks,
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 5, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(leftMenu, 30, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(queues, 0, 3, true).
					AddItem(tasks, 0, 1, false), 0, 2, false), 0, 7, false), 0, 1, false)

	g.pages = tview.NewPages().
		AddAndSwitchToPage("main", flex, true)

	g.app.SetRoot(g.pages, true)
	g.switchPanel("leftMenu")
	leftMenu.Select(1, 0)

	stop := make(chan int, 1)
	g.state.stopChans["queue"] = stop
	g.state.stopChans["task"] = stop
	go g.monitoringTask()
	go g.queuePanel().monitoringQueues(g)

	return flex
}

func (g *Gui) showFailedJobGroupsPage() *tview.Flex {
	g.state.stopChans = make(map[string]chan int)
	header := newHeader(g)
	leftMenu := newLeftMenu(g)
	failedJobs := newFailedJobGroups(g)
	tasks := newTasks(g)

	g.state.panels.panel = []panel{
		header, leftMenu, failedJobs, tasks,
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 5, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(leftMenu, 30, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(failedJobs, 0, 3, true).
					AddItem(tasks, 0, 1, false), 0, 2, false), 0, 7, false), 0, 1, false)

	g.pages = tview.NewPages().
		AddAndSwitchToPage("main", flex, true)

	g.app.SetRoot(g.pages, true)
	g.switchPanel("leftMenu")
	leftMenu.Select(2, 0)

	stop := make(chan int, 1)
	g.state.stopChans["task"] = stop
	go g.monitoringTask()

	return flex
}

func (g *Gui) showJobsPage(group string) *tview.Flex {
	g.state.stopChans = make(map[string]chan int)
	header := newHeader(g)
	leftMenu := newLeftMenu(g)
	jobs := newJobs(g, group)
	tasks := newTasks(g)

	g.state.panels.panel = []panel{
		header, leftMenu, jobs, tasks,
	}

	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(header, 5, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(leftMenu, 30, 0, false).
				AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
					AddItem(jobs, 0, 3, true).
					AddItem(tasks, 0, 1, false), 0, 2, false), 0, 7, false), 0, 1, false)

	g.pages = tview.NewPages().
		AddAndSwitchToPage("main", flex, true)

	g.app.SetRoot(g.pages, true)
	g.switchPanel("leftMenu")
	leftMenu.Select(2, 0)

	stop := make(chan int, 1)
	g.state.stopChans["task"] = stop
	go g.monitoringTask()

	return flex
}
