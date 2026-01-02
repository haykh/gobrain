package tasklist

import (
	"fmt"
	"strings"
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
				if m.is_adding_list && m.adding_list_idx >= 0 && m.adding_list_idx < len(m.tasklists) {
					m.tasklists[m.adding_list_idx].StopEditing(false)
					m.tasklists = append(m.tasklists[:m.adding_list_idx], m.tasklists[m.adding_list_idx+1:]...)
				}
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
					title := strings.TrimSpace(m.app.TypingInput.Value())
					if _, err := m.app.CreateNew_Tasklist(title); err != nil {
						return func() tea.Msg {
							return ui.ErrorMsg{Error: fmt.Errorf("could not create tasklist:\n%v", err)}
						}
					}

					if m.adding_list_idx >= 0 && m.adding_list_idx < len(m.tasklists) {
						m.tasklists[m.adding_list_idx].StopEditing(true)
					}

					m.app.TypingInput.SetValue("")
					m.is_adding_list = false
					m.adding_list_idx = -1
					m.Sync()
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
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch {
				case key.Matches(msg, keys.AddList):
					return m.addTasklist()
				}
			}
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

			case key.Matches(msg, keys.AddList):
				return m.addTasklist()

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
					listIdx := m.filtered_tasks[m.cursor.Index].Index()
					return m.addTaskToList(listIdx, len(m.tasklists[listIdx].Tasks))
				}

				if m.cursor.Index == m.task_view_idx_max-1 {
					m.task_view_idx_min++
					m.task_view_idx_max++
				}

				taskIdx := m.filtered_tasks[m.cursor.Index].Index() + 1
				listIdx := m.filtered_tasks[m.cursor.Index].(*task.Task).ListIndex()
				return m.addTaskToList(listIdx, taskIdx)
			}
		}
	}
	return nil
}

func (m *model) Filter() {
	m.filtered_tasks = []Selector{}
	if len(m.tasklists) == 0 {
		m.cursor.Index = 0
		m.task_view_idx_min = 0
		m.task_view_idx_max = 0
		return
	}
	for i := range m.tasklists {
		m.filtered_tasks = append(m.filtered_tasks, &m.tasklists[i])
		for j := range m.tasklists[i].Tasks {
			m.filtered_tasks = append(m.filtered_tasks, &m.tasklists[i].Tasks[j])
		}
	}
	m.task_view_idx_max = min(m.task_view_idx_min+int(ui.Height_Window), len(m.filtered_tasks))
	if m.cursor.Index >= len(m.filtered_tasks) {
		m.cursor.Index = len(m.filtered_tasks) - 1
	}
	if m.cursor.Index < 0 {
		m.cursor.Index = 0
	}
	m.ensureCursorVisible()
}

func (m *model) addTasklist() tea.Cmd {
	newList := list.New("", "", m.app.Tasks, len(m.tasklists))
	newList.StartEditing(&m.app.TypingInput, "New tasklist")

	m.tasklists = append(m.tasklists, newList)
	m.is_adding_list = true
	m.adding_list_idx = len(m.tasklists) - 1
	m.Filter()

	m.cursor.Index = m.indexForList(m.adding_list_idx)
	m.ensureCursorVisible()

	return tea.Batch(
		func() tea.Msg {
			return ui.TypingStartMsg{}
		},
		func() tea.Msg {
			return ui.DebugMsg{Message: "Adding new tasklist"}
		},
	)
}

func (m *model) addTaskToList(listIdx, insertAt int) tea.Cmd {
	if listIdx < 0 || listIdx >= len(m.tasklists) {
		return nil
	}

	insertAt = max(0, min(insertAt, len(m.tasklists[listIdx].Tasks)))

	m.is_adding_task = true
	m.adding_task_idx = insertAt
	m.adding_list_idx = listIdx

	newTask := task.New("", false, 0, time.Time{}, -1, m.adding_list_idx)
	newTask.StartEditing(&m.app.TypingInput)
	m.tasklists[m.adding_list_idx].AddTask(newTask, m.adding_task_idx)
	m.Filter()

	m.cursor.Index = m.indexForTask(m.adding_list_idx, m.adding_task_idx)
	m.ensureCursorVisible()

	return tea.Batch(
		func() tea.Msg {
			return ui.TypingStartMsg{}
		},
		func() tea.Msg {
			return ui.DebugMsg{Message: fmt.Sprintf("Adding task to tasklist %d at position %d", m.adding_list_idx, m.adding_task_idx)}
		},
	)
}

func (m *model) indexForList(listIdx int) int {
	for i, item := range m.filtered_tasks {
		if item.Type() == "tasklist" && item.Index() == listIdx {
			return i
		}
	}
	return 0
}

func (m *model) indexForTask(listIdx, taskIdx int) int {
	for i, item := range m.filtered_tasks {
		if item.Type() == "task" {
			if t, ok := item.(*task.Task); ok {
				if t.ListIndex() == listIdx && t.Index() == taskIdx {
					return i
				}
			}
		}
	}
	return 0
}

func (m *model) ensureCursorVisible() {
	if len(m.filtered_tasks) == 0 {
		m.task_view_idx_min = 0
		m.task_view_idx_max = 0
		return
	}

	if m.cursor.Index < m.task_view_idx_min {
		m.task_view_idx_min = m.cursor.Index
	}
	if m.cursor.Index >= m.task_view_idx_max {
		m.task_view_idx_max = m.cursor.Index + 1
	}

	m.task_view_idx_min = max(0, min(m.task_view_idx_min, len(m.filtered_tasks)-1))
	m.task_view_idx_max = min(max(m.task_view_idx_min+1, m.task_view_idx_max), len(m.filtered_tasks))
}

func (m *model) Sync() {
	filenames_tasklists, err := m.app.GetMarkdownFilenames(m.app.Tasks)
	if err != nil {
		panic("Could not get tasklist filenames: " + err.Error())
	}
	m.tasklists = []list.List{}
	for fi, filename := range filenames_tasklists {
		if title, tasks, checked, importances, dueDates, err := backend.ParseMarkdownTasklist(m.app.Tasks, filename); err != nil {
			panic("Could not parse tasklist: " + err.Error())
		} else {
			tl := list.New(title, filename, m.app.Tasks, fi)
			for ti := range tasks {
				tl.AppendTask(task.New(tasks[ti], checked[ti], importances[ti], dueDates[ti], ti, fi))
			}
			m.tasklists = append(m.tasklists, tl)
		}
	}
	m.Filter()
}
