package tasklist

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/tasklist/list"
	"github.com/haykh/gobrain/ui/window/tasklist/task"
)

func (m *model) Update(msg tea.Msg) tea.Cmd {
	if m.is_adding_list || m.is_adding_task {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, ui.Key_Cancel):
				if m.is_adding_task {
					m.tasklists[m.adding_list_idx].Tasks[m.adding_task_idx].StopEditing(false)
					if m.tasklists[m.adding_list_idx].Tasks[m.adding_task_idx].Index() == -1 {
						m.tasklists[m.adding_list_idx].RemoveTask(m.adding_task_idx)
						if m.cursor.Index == m.task_view_idx_max-2 {
							m.task_view_idx_min--
							m.task_view_idx_max--
						}
					}
				}
				m.Filter()
				m.is_adding_list = false
				m.is_adding_task = false
				m.adding_list_idx = -1
				m.adding_task_idx = -1
				return tea.Batch(
					func() tea.Msg {
						return ui.TypingEndMsg{}
					},
					func() tea.Msg {
						return ui.DebugMsg{Message: "Cancelled adding tasklist"}
					},
				)

			case key.Matches(msg, ui.Key_Accept):
				if m.is_adding_list {
					m.is_adding_list = false
					// @TODO: add empty tasklist
				} else if m.is_adding_task {
					m.is_adding_task = false
					m.tasklists[m.adding_list_idx].Tasks[m.adding_task_idx].StopEditing(true)
					if err := m.tasklists[m.adding_list_idx].Sync(); err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after adding task:\n%v", err)}
						}
					}
					if m.tasklists[m.adding_list_idx].Tasks[m.adding_task_idx].Index() == -1 {
						m.cursor.Index++
					}
					m.Sync()
					m.adding_list_idx = -1
					m.adding_task_idx = -1
				}
				return tea.Batch(
					func() tea.Msg {
						return ui.TypingEndMsg{}
					},
					func() tea.Msg {
						return ui.DebugMsg{Message: "Added new item"}
					},
				)
			}
		}
	} else {
		if len(m.filtered_tasks) == 0 {
			return nil
		}
		adjust_minmax_idx := func() {
			m.task_view_idx_min = min(max(0, m.task_view_idx_min), len(m.filtered_tasks)-1)
			m.task_view_idx_max = min(max(1, m.task_view_idx_max), len(m.filtered_tasks))
		}
		idx_minus := func() {
			m.cursor.Index = max(0, m.cursor.Index-1)
			if m.cursor.Index < m.task_view_idx_min {
				m.task_view_idx_min--
				m.task_view_idx_max--
			}
			adjust_minmax_idx()
		}
		idx_plus := func() {
			m.cursor.Index = min(len(m.filtered_tasks)-1, m.cursor.Index+1)
			if m.cursor.Index >= m.task_view_idx_max {
				m.task_view_idx_min++
				m.task_view_idx_max++
			}
			adjust_minmax_idx()
		}
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch {

			case key.Matches(msg, keys.Up):
				idx_minus()
				return nil

			case key.Matches(msg, keys.Down):
				idx_plus()
				return nil

			case key.Matches(msg, keys.PgUp):
				idx_minus()
				for {
					if m.cursor.Index > 0 {
						if m.filtered_tasks[m.cursor.Index].Type() == "tasklist" {
							break
						}
						idx_minus()
					} else {
						break
					}
				}
				return nil

			case key.Matches(msg, keys.PgDown):
				idx_plus()
				for {
					if m.cursor.Index < len(m.filtered_tasks)-1 {
						if m.filtered_tasks[m.cursor.Index].Type() == "tasklist" {
							break
						}
						idx_plus()
					} else {
						break
					}
				}
				return nil

			case key.Matches(msg, keys.Toggle):
				if m.filtered_tasks[m.cursor.Index].Type() == "task" {
					task_idx := m.filtered_tasks[m.cursor.Index].Index()
					tasklist_idx := m.filtered_tasks[m.cursor.Index].(*task.Task).ListIndex()
					m.tasklists[tasklist_idx].Tasks[task_idx].Toggle()
					if err := m.tasklists[tasklist_idx].Sync(); err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after toggling task:\n%v", err)}
						}
					}
					m.Sync()
					return func() tea.Msg {
						return ui.DebugMsg{
							Message: fmt.Sprintf("Toggled task %d in tasklist %d", m.filtered_tasks[m.cursor.Index].Index(), tasklist_idx),
						}
					}
				}
				return nil

			case key.Matches(msg, keys.Edit):
				if m.filtered_tasks[m.cursor.Index].Type() == "task" {
					m.is_adding_task = true
					m.adding_task_idx = m.filtered_tasks[m.cursor.Index].Index()
					m.adding_list_idx = m.filtered_tasks[m.cursor.Index].(*task.Task).ListIndex()
					m.tasklists[m.adding_list_idx].Tasks[m.adding_task_idx].StartEditing(&m.app.TypingInput)
					m.Filter()
					return tea.Batch(
						func() tea.Msg {
							return ui.TypingStartMsg{}
						},
						func() tea.Msg {
							return ui.DebugMsg{Message: fmt.Sprintf("Editing task %d of tasklist %d", m.adding_task_idx, m.adding_list_idx)}
						},
					)
				} else if m.filtered_tasks[m.cursor.Index].Type() == "tasklist" {
					// edit tasklist title
				}
				// @TODO: synchronize with backend
				return nil

			case key.Matches(msg, keys.Delete):
				if m.filtered_tasks[m.cursor.Index].Type() == "task" {
					task_idx := m.filtered_tasks[m.cursor.Index].Index()
					tasklist_idx := m.filtered_tasks[m.cursor.Index].(*task.Task).ListIndex()
					m.tasklists[tasklist_idx].RemoveTask(task_idx)
					if err := m.tasklists[tasklist_idx].Sync(); err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after deleting task:\n%v", err)}
						}
					}
					m.Sync()
					return func() tea.Msg {
						return ui.DebugMsg{
							Message: fmt.Sprintf("Deleted task %d from tasklist %d", task_idx, tasklist_idx),
						}
					}
				} else if m.filtered_tasks[m.cursor.Index].Type() == "tasklist" {
					tasklist_idx := m.filtered_tasks[m.cursor.Index].Index()
					if err := m.tasklists[tasklist_idx].Delete(m.app); err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: fmt.Errorf("Could not delete tasklist:\n%v", err)}
						}
					}
					m.Sync()
					return func() tea.Msg {
						return ui.DebugMsg{
							Message: fmt.Sprintf("Deleted tasklist %d", tasklist_idx),
						}
					}
				}
				return nil

			case key.Matches(msg, keys.Add):
				if m.filtered_tasks[m.cursor.Index].Type() == "tasklist" {
					m.is_adding_list = true
					// @TODO: record position to insert new tasklist
					return func() tea.Msg {
						return ui.TypingStartMsg{}
					}
				} else {
					if m.cursor.Index == m.task_view_idx_max-1 {
						m.task_view_idx_min++
						m.task_view_idx_max++
					}
					m.is_adding_task = true
					m.adding_task_idx = m.filtered_tasks[m.cursor.Index].Index() + 1
					m.adding_list_idx = m.filtered_tasks[m.cursor.Index].(*task.Task).ListIndex()
					new_task := task.New("", false, 0, time.Time{}, -1, m.adding_list_idx)
					new_task.StartEditing(&m.app.TypingInput)
					m.tasklists[m.adding_list_idx].AddTask(new_task, m.adding_task_idx)
					m.Filter()
					return tea.Batch(
						func() tea.Msg {
							return ui.TypingStartMsg{}
						},
						func() tea.Msg {
							return ui.DebugMsg{Message: fmt.Sprintf("Adding task to tasklist %d at position %d", m.adding_list_idx, m.adding_task_idx)}
						},
					)
				}
			}
		}
	}
	return nil
}

func (m *model) Filter() {
	m.filtered_tasks = []Selector{}
	if len(m.tasklists) == 0 {
		return
	}
	for _, l := range m.tasklists {
		m.filtered_tasks = append(m.filtered_tasks, &l)
		for _, t := range l.Tasks {
			m.filtered_tasks = append(m.filtered_tasks, &t)
		}
	}
	m.task_view_idx_max = min(m.task_view_idx_min+int(ui.Height_Window), len(m.filtered_tasks))
	if m.cursor.Index >= len(m.filtered_tasks) {
		m.cursor.Index = len(m.filtered_tasks) - 1
	}
}

func (m *model) Sync() {
	filenames_tasklists, err := m.app.GetMarkdownFilenames(m.app.TasksPath)
	if err != nil {
		panic("Could not get tasklist filenames: " + err.Error())
	}
	m.tasklists = []list.List{}
	for fi, filename := range filenames_tasklists {
		if title, tasks, checked, importances, dueDates, err := backend.ParseMarkdownTasklist(m.app.TasksPath, filename); err != nil {
			panic("Could not parse tasklist: " + err.Error())
		} else {
			tl := list.New(title, filename, m.app.TasksPath, fi)
			for ti := range tasks {
				tl.AppendTask(task.New(tasks[ti], checked[ti], importances[ti], dueDates[ti], ti, fi))
			}
			m.tasklists = append(m.tasklists, tl)
		}
	}
	m.Filter()
}
