package tasklist

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/tasklist/list"
	"github.com/haykh/gobrain/ui/window/tasklist/task"
)

func (m *model) Sync() error {
	filenames_tasklists, err := m.app.GetMarkdownFilenames(m.app.Tasks)
	if err != nil {
		return err
	}
	m.tasklists = []list.List{}
	for _, filename := range filenames_tasklists {
		if title, tasks, checked, importances, dueDates, err := backend.ParseMarkdownTasklist(m.app.Tasks, filename); err != nil {
			return err
		} else {
			tl := list.New(title, filename, m.app.Tasks)
			for ti := range tasks {
				tl.AppendTask(task.New(tasks[ti], checked[ti], importances[ti], dueDates[ti]))
			}
			m.tasklists = append(m.tasklists, tl)
		}
	}
	return nil
}

func (m *model) Navigate(dir ui.Direction) tea.Cmd {
	if len(m.tasklists) == 0 {
		return nil
	}
	switch dir {
	case ui.DirUp:
		if m.cursor.TaskIndex == 0 {
			if m.cursor.TasklistIndex > 0 {
				m.cursor.TasklistIndex--
				m.cursor.TaskIndex = m.tasklists[m.cursor.TasklistIndex].NumTasks()
			}
		} else {
			m.cursor.TaskIndex--
		}
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))

	case ui.DirDown:
		if m.cursor.TaskIndex >= m.tasklists[m.cursor.TasklistIndex].NumTasks() {
			if m.cursor.TasklistIndex < len(m.tasklists)-1 {
				m.cursor.TasklistIndex++
				m.cursor.TaskIndex = 0
			}
		} else {
			m.cursor.TaskIndex++
		}
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))

	case ui.DirPgUp:
		if m.cursor.TasklistIndex > 0 && m.cursor.TaskIndex == 0 {
			m.cursor.TasklistIndex--
			m.cursor.TaskIndex = 0
		} else if m.cursor.TaskIndex > 0 {
			m.cursor.TaskIndex = 0
		}
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))

	case ui.DirPgDown:
		if m.cursor.TasklistIndex < len(m.tasklists)-1 {
			m.cursor.TasklistIndex++
			m.cursor.TaskIndex = 0
		}
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))
	}
	return func() tea.Msg {
		return ui.DebugMsg{
			Message: fmt.Sprintf(
				"Moved cursor to tasklist %d, task %d -- view_index %d: view_range.IMin: %d",
				m.cursor.TasklistIndex,
				m.cursor.TaskIndex,
				m.cursor.ViewIndex(m.tasklists),
				m.view_range.IMin,
			),
		}
	}
}

func (m *model) ToggleTask() tea.Cmd {
	if len(m.tasklists) == 0 {
		return nil
	}
	if m.cursor.TaskIndex == 0 {
		return nil
	}
	m.tasklists[m.cursor.TasklistIndex].Tasks[m.cursor.TaskIndex-1].Toggle()
	if err := m.tasklists[m.cursor.TasklistIndex].Sync(); err != nil {
		return func() tea.Msg {
			return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after toggling task:\n%v", err)}
		}
	}
	if err := m.Sync(); err != nil {
		return func() tea.Msg {
			return ui.ErrorMsg{Error: fmt.Errorf("Could not sync model after toggling task:\n%v", err)}
		}
	}
	return func() tea.Msg {
		return ui.DebugMsg{
			Message: fmt.Sprintf("Toggled task %d in tasklist %d", m.cursor.TaskIndex, m.cursor.TasklistIndex),
		}
	}
}

func (m *model) DeleteItem() tea.Cmd {
	if len(m.tasklists) == 0 {
		return nil
	}
	if m.cursor.TaskIndex == 0 {
		// delete tasklist
		deleted_ntasks := m.tasklists[m.cursor.TasklistIndex].NumTasks()
		if err := m.tasklists[m.cursor.TasklistIndex].Delete(m.app); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not delete tasklist:\n%v", err)}
			}
		}
		if m.cursor.TasklistIndex >= len(m.tasklists)-1 {
			m.cursor.TasklistIndex--
		}
		m.cursor.TaskIndex = 0
		if err := m.Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync model after deleting tasklist:\n%v", err)}
			}
		}
		m.view_range.RemoveItems(1+deleted_ntasks, m.tasklists)
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))
		return func() tea.Msg {
			return ui.DebugMsg{Message: "Deleted tasklist"}
		}
	} else {
		// delete task
		tl := &m.tasklists[m.cursor.TasklistIndex]
		tl.RemoveTask(m.cursor.TaskIndex - 1)
		if m.cursor.TaskIndex >= tl.NumTasks() && m.cursor.TaskIndex > 0 {
			m.cursor.TaskIndex--
		}
		if err := tl.Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after deleting task:\n%v", err)}
			}
		}
		if err := m.Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync model after deleting task:\n%v", err)}
			}
		}
		m.view_range.RemoveItems(1, m.tasklists)
		m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))
		return func() tea.Msg {
			return ui.DebugMsg{
				Message: fmt.Sprintf("Deleted task in tasklist %d", m.cursor.TasklistIndex),
			}
		}
	}
}

func (m *model) EditItem() tea.Cmd {
	if len(m.tasklists) == 0 {
		return nil
	}
	if m.cursor.TaskIndex == 0 {
		// edit tasklist
		// @TODO: edit tasklist title
		return nil
	} else {
		// edit task
		m.active_action = EditingTask
		m.active_item = ItemIndex{m.cursor.TasklistIndex, m.cursor.TaskIndex - 1}
		m.tasklists[m.active_item.ListIndex].Tasks[m.active_item.TaskIndex].StartEditing(&m.app.TypingInput)
		return tea.Batch(
			func() tea.Msg {
				return ui.TypingStartMsg{}
			},
			func() tea.Msg {
				return ui.DebugMsg{
					Message: fmt.Sprintf(
						"Editing task %d of tasklist %d",
						m.active_item.TaskIndex,
						m.active_item.ListIndex,
					),
				}
			},
		)
	}
}

func (m *model) AddTask() tea.Cmd {
	if len(m.tasklists) == 0 {
		return nil
	}
	m.active_action = AddingTask
	m.active_item = ItemIndex{m.cursor.TasklistIndex, m.cursor.TaskIndex}
	m.cursor.TaskIndex++

	newTask := task.New("", false, 0, time.Time{})
	newTask.StartEditing(&m.app.TypingInput)
	m.tasklists[m.active_item.ListIndex].AddTask(newTask, m.active_item.TaskIndex)

	m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))

	return tea.Batch(
		func() tea.Msg {
			return ui.TypingStartMsg{}
		},
		func() tea.Msg {
			return ui.DebugMsg{
				Message: fmt.Sprintf(
					"Adding task to tasklist %d at position %d cursor at %d",
					m.active_item.ListIndex,
					m.active_item.TaskIndex,
					m.cursor.TaskIndex,
				),
			}
		},
	)
}

func (m *model) AddTasklist() tea.Cmd {
	newlist := list.New("", "", m.app.Tasks)
	newlist.StartEditing(&m.app.TypingInput)

	new_tasklist_idx := len(m.tasklists)

	m.tasklists = append(m.tasklists, newlist)
	m.active_action = AddingList
	m.active_item = ItemIndex{new_tasklist_idx, -1}
	m.cursor.TasklistIndex = new_tasklist_idx
	m.cursor.TaskIndex = 0

	m.view_range.EnsureCursorVisible(m.cursor.ViewIndex(m.tasklists))

	return tea.Batch(
		func() tea.Msg {
			return ui.TypingStartMsg{}
		},
		func() tea.Msg {
			return ui.DebugMsg{Message: "Adding new tasklist"}
		},
	)
}

func (m *model) CancelAction() tea.Cmd {
	switch m.active_action {

	case AddingList, EditingList:
		// cancel tasklist editing/adding
		m.active_action = NoAction
		if m.active_item.ListIndex < 0 || m.active_item.ListIndex >= len(m.tasklists) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.ListIndex %d", m.active_item.ListIndex)}
			}
		}
		m.tasklists[m.active_item.ListIndex].StopEditing(false)
		m.tasklists = append(m.tasklists[:m.active_item.ListIndex], m.tasklists[m.active_item.ListIndex+1:]...)
		m.active_item = ItemIndex{-1, -1}

	case AddingTask, EditingTask:
		// cancel task editing/adding
		action_was := m.active_action
		m.active_action = NoAction
		if m.active_item.ListIndex < 0 || m.active_item.ListIndex >= len(m.tasklists) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.ListIndex %d", m.active_item.ListIndex)}
			}
		}
		if m.active_item.TaskIndex < 0 || m.active_item.TaskIndex >= len(m.tasklists[m.active_item.ListIndex].Tasks) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.TaskIndex %d", m.active_item.TaskIndex)}
			}
		}
		m.tasklists[m.active_item.ListIndex].Tasks[m.active_item.TaskIndex].StopEditing(false)
		if action_was == AddingTask {
			m.tasklists[m.active_item.ListIndex].RemoveTask(m.active_item.TaskIndex)
		}
		m.active_item = ItemIndex{-1, -1}

	}
	return tea.Batch(
		func() tea.Msg {
			return ui.TypingEndMsg{}
		},
		func() tea.Msg {
			return ui.DebugMsg{Message: "Cancelled action"}
		},
	)
}

func (m *model) AcceptAction() tea.Cmd {
	switch m.active_action {

	case AddingList, EditingList:
		m.active_action = NoAction
		title := strings.TrimSpace(m.app.TypingInput.Value())
		if _, err := m.app.CreateNew_Tasklist(title); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("could not create tasklist:\n%v", err)}
			}
		}

		if m.active_item.ListIndex < 0 || m.active_item.ListIndex >= len(m.tasklists) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.ListIndex %d", m.active_item.ListIndex)}
			}
		}
		m.tasklists[m.active_item.ListIndex].StopEditing(true)

		m.app.TypingInput.SetValue("")
		m.active_item = ItemIndex{-1, -1}
		if err := m.Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync model after adding tasklist:\n%v", err)}
			}
		}

	case AddingTask, EditingTask:
		m.active_action = NoAction
		if m.active_item.ListIndex < 0 || m.active_item.ListIndex >= len(m.tasklists) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.ListIndex %d", m.active_item.ListIndex)}
			}
		}
		if m.active_item.TaskIndex < 0 || m.active_item.TaskIndex >= len(m.tasklists[m.active_item.ListIndex].Tasks) {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Invalid active_item.TaskIndex %d", m.active_item.TaskIndex)}
			}
		}
		m.tasklists[m.active_item.ListIndex].Tasks[m.active_item.TaskIndex].StopEditing(true)
		if err := m.tasklists[m.active_item.ListIndex].Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync tasklist after adding task:\n%v", err)}
			}
		}
		if err := m.Sync(); err != nil {
			return func() tea.Msg {
				return ui.ErrorMsg{Error: fmt.Errorf("Could not sync model after adding task:\n%v", err)}
			}
		}
		m.active_item = ItemIndex{-1, -1}

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
