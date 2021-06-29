package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	success   = "Success"
	executing = "Executing"
	cancel    = "canceled"
)

type task struct {
	Name    string
	Status  string
	Created string
	Func    func(ctx context.Context) error
	Ctx     context.Context
	Cancel  context.CancelFunc
}

type tasks struct {
	*tview.Table
	tasks chan *task
}

func (g *Gui) updateTask() {
	go g.app.QueueUpdateDraw(func() {
		g.taskPanel().setEntries(g)
		g.taskPanel().ScrollToEnd()
		logger.Println("Tasks updated")
	})
}

func (t *tasks) name() string {
	return "tasks"
}

func (t *tasks) canSelect() bool {
	return true
}

func (g *Gui) taskPanel() *tasks {
	for _, panel := range g.state.panels.panel {
		if panel.name() == "tasks" {
			return panel.(*tasks)
		}
	}
	return nil
}

func newTasks(g *Gui) *tasks {
	logger.Println("New tasks")
	tasks := &tasks{
		Table: tview.NewTable().SetSelectable(true, false).Select(0, 0).SetFixed(1, 1),
		tasks: make(chan *task),
	}

	tasks.SetTitle(" [::b]Tasks ").SetTitleAlign(tview.AlignCenter).SetTitleColor(tcell.ColorSeashell.TrueColor())
	tasks.SetBorder(true)
	tasks.setEntries(g)
	tasks.setKeybinding(g)
	logger.Println("Task created")
	return tasks
}

func (t *tasks) setKeybinding(g *Gui) {
	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		g.setGlobalKeybinding(event)

		return event
	})
}

func (g *Gui) monitoringTask() {
LOOP:
	for {
		select {
		case task := <-g.taskPanel().tasks:
			go func() {
				if err := task.Func(task.Ctx); err != nil {
					task.Status = err.Error()
				} else {
					task.Status = success
				}
				g.updateTask()
			}()
		case <-g.state.stopChans["task"]:
			break LOOP
		}
	}
}

func (g *Gui) startTask(taskName string, f func(ctx context.Context) error) {
	logger.Println("Start task " + taskName)
	ctx, cancel := context.WithCancel(context.Background())

	task := &task{
		Name:    taskName,
		Status:  executing,
		Created: time.Now().Format("2006/01/02 15:04:05"),
		Func:    f,
		Ctx:     ctx,
		Cancel:  cancel,
	}
	logger.Println("Task: " + fmt.Sprint(task))

	g.state.resources.tasks = append(g.state.resources.tasks, task)
	g.updateTask()
	logger.Println("Send task goroutine into tasks chanel: start")
	g.taskPanel().tasks <- task
	logger.Println("Send task goroutine into tasks chanel: complete")
}

func (t *tasks) entries(g *Gui) {
	// do nothing
}

func (t *tasks) setEntries(g *Gui) {
	t.entries(g)
	table := t.Clear()
	table.SetSelectedStyle(tcell.Style{}.
		Background(tcell.ColorWhiteSmoke.TrueColor()).
		Foreground(tcell.ColorBlack))

	headers := []string{
		"Name",
		"Status",
		"Created",
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

	for i, task := range g.state.resources.tasks {
		table.SetCell(i+1, 0, tview.NewTableCell(task.Name).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 1, tview.NewTableCell(task.Status).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(1).
			SetExpansion(1))

		table.SetCell(i+1, 2, tview.NewTableCell(task.Created).
			SetTextColor(tcell.ColorWhiteSmoke).
			SetMaxWidth(1).
			SetExpansion(1))

	}
}

func (t *tasks) focus(g *Gui) {
	t.SetSelectable(true, false)
	g.app.SetFocus(t)
	t.Select(t.GetRowCount()-1, 0)
}

func (t *tasks) unfocus() {
	t.SetSelectable(false, false)
}

func (t *tasks) setFilterWord(word string) {
	// do nothings
}

func (t *tasks) updateEntries(g *Gui) {
	// do nothings
}

func (t *tasks) updateTitle(g *Gui) {
	// do nothings
}
