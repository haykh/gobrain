package randnotes

import (
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/haykh/gobrain/backend"
	"github.com/haykh/gobrain/ui"
	"github.com/haykh/gobrain/ui/window/randnotes/note"
)

func navigate(c Cursor, dir ui.Direction, num_notes int, num_tags int) Cursor {
	switch dir {
	case ui.DirUp:
		if c.NoteIndex > 0 {
			return Cursor{c.NoteIndex - 1, c.TagIndex, c.TagFilter}
		}
	case ui.DirDown:
		if c.NoteIndex < num_notes-1 {
			return Cursor{c.NoteIndex + 1, c.TagIndex, c.TagFilter}
		}
	case ui.DirLeft:
		if c.TagIndex > -1 {
			return Cursor{c.NoteIndex, c.TagIndex - 1, c.TagFilter}
		}
	case ui.DirRight:
		if c.TagIndex < num_tags-1 {
			return Cursor{c.NoteIndex, c.TagIndex + 1, c.TagFilter}
		}
	}
	return c
}

func (m *model) Update(msg tea.Msg) tea.Cmd {
	if len(m.filtered_notes) == 0 {
		return nil
	}
	switch msg := msg.(type) {

	case tea.KeyMsg:
		num_notes := len(m.filtered_notes)
		note_idx := m.cursor.NoteIndex
		tag_idx := m.cursor.TagIndex
		num_tags_on_current_note := len(m.filtered_notes[note_idx].Tags)

		adjustTagIndex := func(new_tag_idx, new_num_tags int) int {
			if new_tag_idx > -1 {
				if new_num_tags == 0 {
					return -1
				} else {
					return min(max(0, new_num_tags-(num_tags_on_current_note-tag_idx)), new_num_tags)
				}
			}
			return new_tag_idx
		}

		switch {

		case key.Matches(msg, keys.Up):
			m.cursor = navigate(m.cursor, ui.DirUp, num_notes, num_tags_on_current_note)
			if m.cursor.NoteIndex < m.note_view_idx_min {
				m.note_view_idx_min--
				m.note_view_idx_max--
			}
			m.cursor.TagIndex = adjustTagIndex(m.cursor.TagIndex, len(m.filtered_notes[m.cursor.NoteIndex].Tags))
			return nil

		case key.Matches(msg, keys.Down):
			m.cursor = navigate(m.cursor, ui.DirDown, num_notes, num_tags_on_current_note)
			if m.cursor.NoteIndex >= m.note_view_idx_max {
				m.note_view_idx_min++
				m.note_view_idx_max++
			}
			m.cursor.TagIndex = adjustTagIndex(m.cursor.TagIndex, len(m.filtered_notes[m.cursor.NoteIndex].Tags))
			return nil

		case key.Matches(msg, keys.Left):
			m.cursor = navigate(m.cursor, ui.DirLeft, num_notes, num_tags_on_current_note)
			return nil

		case key.Matches(msg, keys.Right):
			m.cursor = navigate(m.cursor, ui.DirRight, num_notes, num_tags_on_current_note)
			return nil

		case key.Matches(msg, keys.Select):
			return func() tea.Msg {
				return ui.NewViewer{
					Filepath: m.app.RandomNotesPath,
					Filename: m.filtered_notes[note_idx].Filename,
				}
			}

		case key.Matches(msg, keys.Edit):
			return func() tea.Msg {
				return ui.OpenEditorMsg{
					Filename: filepath.Join(m.app.RandomNotesPath, m.filtered_notes[note_idx].Filename),
				}
			}

		case key.Matches(msg, keys.Filter):
			old_filter := m.cursor.TagFilter
			if m.cursor.TagIndex == -1 {
				m.cursor.TagFilter = ""
			} else {
				new_tag := m.filtered_notes[m.cursor.NoteIndex].Tags[m.cursor.TagIndex]
				if new_tag == m.cursor.TagFilter {
					m.cursor.TagFilter = ""
				} else {
					m.cursor.TagFilter = new_tag
				}
			}
			if old_filter != m.cursor.TagFilter {
				m.Filter()
			}
			return nil

		case key.Matches(msg, keys.Delete):
			return func() tea.Msg {
				return ui.TrashRandomNoteMsg{
					Filename: m.filtered_notes[m.cursor.NoteIndex].Filename,
				}
			}
		}

	}
	return nil
}

func (m *model) Filter() {
	m.filtered_notes = []*note.Note{}
	for ni, n := range m.notes {
		if ni == m.cursor.NoteIndex {
			m.cursor.NoteIndex = len(m.filtered_notes)
		}
		if m.cursor.TagFilter != "" {
			for _, t := range n.Tags {
				if t == m.cursor.TagFilter {
					m.filtered_notes = append(m.filtered_notes, &n)
				}
			}
		} else {
			m.filtered_notes = append(m.filtered_notes, &n)
		}
	}
	m.note_view_idx_min = 0
	m.note_view_idx_max = min(int(ui.Height_Window/2), len(m.filtered_notes))
	if m.cursor.NoteIndex >= m.note_view_idx_max {
		m.note_view_idx_min = m.note_view_idx_max - m.note_view_idx_min
		m.note_view_idx_max = m.cursor.NoteIndex + 1
	}
}

func (m *model) Sync() {
	filenames, err := m.app.GetMarkdownFilenames(m.app.RandomNotesPath)
	if err != nil {
		panic("Could not get random notes filenames: " + err.Error())
	}
	m.notes = []note.Note{}
	for _, filename := range filenames {
		if title, icon, tags, created, err := backend.ParseMarkdownNote(m.app.RandomNotesPath, filename); err != nil {
			panic("Could not parse random note: " + err.Error())
		} else {
			n := note.New(filename, title, icon, tags, created)
			m.notes = append(m.notes, n)
		}
	}
	m.Filter()
}
