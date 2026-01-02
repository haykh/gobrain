package backend

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type Backend struct {
	Paths
	storage Storage

	TypingInput textinput.Model
}

func New(path ...string) *Backend {
	root := ""
	if len(path) > 0 {
		root = filepath.Join(path...)
	} else {
		usr, err := user.Current()
		if err != nil {
			panic("Could not get current user: " + err.Error())
		}
		root = filepath.Join(usr.HomeDir, ".gobrain")
	}

	paths := NewPaths(root)
	return &Backend{
		Paths:       paths,
		storage:     localStorage{},
		TypingInput: textinput.New(),
	}
}

func (b Backend) Init() {
	if err := b.storage.Ensure(b.Paths); err != nil {
		panic("Could not create child directory: " + err.Error())
	}
}

func (b Backend) GetMarkdownFilenames(path string) ([]string, error) {
	filenames, err := b.storage.ListMarkdown(path)
	if err != nil {
		return nil, fmt.Errorf("could not read markdown directory: %w", err)
	}
	return filenames, nil
}

func (b Backend) CreateNew_DailyNote(date time.Time) (string, error) {
	filename := fmt.Sprintf("%s.md", date.Format("2006-Jan-02"))
	filaneme_full := filepath.Join(b.DailyNotes, filename)
	file, err := b.storage.Create(filaneme_full)
	if err != nil {
		return "", fmt.Errorf("could not create daily note: %v", err)
	}
	file.Close()

	if err := AddTitleMarkdownNote(b.DailyNotes, filename, date.Format("Jan 2 2006")); err != nil {
		return "", fmt.Errorf("could not add title to daily note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.DailyNotes, filename, "created", date.Format(time.RFC3339)); err != nil {
		return "", fmt.Errorf("could not add created metadata to daily note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.DailyNotes, filename, "icon", "\"󰃭\""); err != nil {
		return "", fmt.Errorf("could not add icon to daily note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.DailyNotes, filename, "tags", []string{"\"daily\""}); err != nil {
		return "", fmt.Errorf("could not add tags to daily note: %v", err)
	}

	return filename, nil
}

func (b Backend) CreateNew_RandomNote() (string, error) {
	now := time.Now()
	datePrefix := now.Format("2006-Jan-02")

	var filename string
	var filename_full string
	for i := 1; i < 100; i++ {
		filename = fmt.Sprintf("%s-%02d.md", datePrefix, i)
		filename_full = filepath.Join(b.RandomNotes, filename)
		if !b.storage.Exists(filename_full) {
			break
		}
		if i == 99 {
			return "", fmt.Errorf("could not create new random note, too many files for today")
		}
	}

	file, err := b.storage.Create(filename_full)
	if err != nil {
		return "", err
	}
	file.Close()

	if err := AddMetadataMarkdownNote(b.RandomNotes, filename, "icon", "\"\""); err != nil {
		return "", fmt.Errorf("could not add icon to new random note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.RandomNotes, filename, "tags", []string{}); err != nil {
		return "", fmt.Errorf("could not add tags to new random note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.RandomNotes, filename, "created", now.Format(time.RFC3339)); err != nil {
		return "", fmt.Errorf("could not add created metadata to new random note: %v", err)
	}

	return filename, nil
}

func (b Backend) CreateNew_Tasklist(title string) (string, error) {
	if title == "" {
		title = "New Tasklist"
	}

	prefix := time.Now().Format("2006-Jan-02-tasklist")

	var filename string
	var filenameFull string
	for i := 1; i < 100; i++ {
		filename = fmt.Sprintf("%s-%02d.md", prefix, i)
		filenameFull = filepath.Join(b.Tasks, filename)
		if !b.storage.Exists(filenameFull) {
			break
		}
		if i == 99 {
			return "", fmt.Errorf("could not create new tasklist, too many files for today")
		}
	}

	if err := WriteMarkdownTasklist(b.Tasks, filename, title, []string{}, []bool{}, []int{}, []time.Time{}); err != nil {
		return "", fmt.Errorf("could not write new tasklist: %v", err)
	}

	return filename, nil
}

func (b Backend) TrashNote(filename, path string) error {
	filename_full := filepath.Join(path, filename)
	if !b.storage.Exists(filename_full) {
		return fmt.Errorf("file does not exist: %s", filename_full)
	}
	if err := AddMetadataMarkdownNote(path, filename, "trashed", time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("could not add trashed metadata to random note: %v", err)
	}
	trash_path := b.TrashRoot
	if strings.HasSuffix(path, "daily_notes") {
		trash_path = b.TrashDailyNotes
	}
	if strings.HasSuffix(path, "random_notes") {
		trash_path = b.TrashRandom
	}
	for i := 1; i < 100; i++ {
		filename_trash_full := filepath.Join(trash_path, filename)
		if i > 1 {
			ext := filepath.Ext(filename_trash_full)
			filename_trash_full = filename_trash_full[:len(filename_trash_full)-len(ext)]
			filename_trash_full = fmt.Sprintf("%s.v%02d%s", filename_trash_full, i, ext)
		}
		if !b.storage.Exists(filename_trash_full) {
			if err := b.storage.Move(filename_full, filename_trash_full); err != nil {
				return fmt.Errorf("could not move file to trash: %v", err)
			} else {
				return nil
			}
		}
	}
	return fmt.Errorf("could not move file to trash, too many versions already exist: %s", filename)
}

func (b Backend) TrashTasklist(filename, path string) error {
	filename_full := filepath.Join(path, filename)
	if !b.storage.Exists(filename_full) {
		return fmt.Errorf("file does not exist: %s", filename_full)
	}
	if err := AddMetadataMarkdownNote(path, filename, "trashed", time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("could not add trashed metadata to random note: %v", err)
	}
	trash_path := b.TrashTasks
	for i := 1; i < 100; i++ {
		filename_trash_full := filepath.Join(trash_path, filename)
		if i > 1 {
			ext := filepath.Ext(filename_trash_full)
			filename_trash_full = filename_trash_full[:len(filename_trash_full)-len(ext)]
			filename_trash_full = fmt.Sprintf("%s.v%02d%s", filename_trash_full, i, ext)
		}
		if !b.storage.Exists(filename_trash_full) {
			if err := b.storage.Move(filename_full, filename_trash_full); err != nil {
				return fmt.Errorf("could not move file to trash: %v", err)
			} else {
				return nil
			}
		}
	}
	return fmt.Errorf("could not move file to trash, too many versions already exist: %s", filename)
}

// Storage exposes the storage implementation for tests.
func (b *Backend) Storage() Storage {
	return b.storage
}
