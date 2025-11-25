package backend

import "path/filepath"

// Paths holds all filesystem locations used by the application. Keeping them in
// one place makes it easy to plug in alternative storage layers (e.g. git
// syncing or cloud-backed libraries) without touching callers.
type Paths struct {
	Root            string
	DailyNotes      string
	Tasks           string
	RandomNotes     string
	Library         string
	TrashRoot       string
	TrashDailyNotes string
	TrashTasks      string
	TrashRandom     string
}

// NewPaths builds the default directory layout under the supplied root. Extra
// locations (Library, TrashRandom) are included up front to keep future
// features straightforward.
func NewPaths(root string) Paths {
	trashRoot := filepath.Join(root, ".trash")

	return Paths{
		Root:            root,
		DailyNotes:      filepath.Join(root, "daily_notes"),
		Tasks:           filepath.Join(root, "tasks"),
		RandomNotes:     filepath.Join(root, "random_notes"),
		Library:         filepath.Join(root, "library"),
		TrashRoot:       trashRoot,
		TrashDailyNotes: filepath.Join(trashRoot, "daily_notes"),
		TrashTasks:      filepath.Join(trashRoot, "tasks"),
		TrashRandom:     filepath.Join(trashRoot, "random_notes"),
	}
}

// All returns every directory that should exist on disk. This keeps Init()
// declarative and easy to extend.
func (p Paths) All() []string {
	return []string{
		p.Root,
		p.DailyNotes,
		p.Tasks,
		p.RandomNotes,
		p.Library,
		p.TrashRoot,
		p.TrashDailyNotes,
		p.TrashTasks,
		p.TrashRandom,
	}
}
