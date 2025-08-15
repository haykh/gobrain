package backend

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type Backend struct {
	MainPath        string
	DailyNotesPath  string
	TasksPath       string
	RandomNotesPath string

	TrashPath string
}

func New(path ...string) *Backend {
	main_path := ""
	if len(path) > 0 {
		main_path = filepath.Join(path...)
	} else {
		usr, err := user.Current()
		if err != nil {
			panic("Could not get current user: " + err.Error())
		}
		homedir := usr.HomeDir
		main_path = filepath.Join(homedir, ".gobrain")
	}
	return &Backend{
		MainPath:        main_path,
		DailyNotesPath:  filepath.Join(main_path, "daily_notes"),
		TasksPath:       filepath.Join(main_path, "tasks"),
		RandomNotesPath: filepath.Join(main_path, "random_notes"),
		TrashPath:       filepath.Join(main_path, ".trash"),
	}
}

func (b Backend) Init() {
	for _, path := range []string{b.MainPath, b.DailyNotesPath, b.TasksPath, b.RandomNotesPath, b.TrashPath} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, os.ModePerm); err != nil {
				panic("Could not create child directory: " + err.Error())
			}
		}
	}
}

func (b Backend) GetFilepaths_RandomNotes() ([]string, error) {
	files, err := os.ReadDir(b.RandomNotesPath)
	if err != nil {
		return nil, err
	}
	var filenames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".md" {
			filenames = append(filenames, file.Name())
		}
	}
	return filenames, nil
}

func (b Backend) CreateNew_RandomNote() (string, error) {
	now := time.Now()
	datePrefix := now.Format("2006-Jan-02")

	var filename string
	var filename_full string
	for i := 1; i < 100; i++ {
		filename = fmt.Sprintf("%s-%02d.md", datePrefix, i)
		filename_full = filepath.Join(b.RandomNotesPath, filename)
		if _, err := os.Stat(filename_full); os.IsNotExist(err) {
			break
		}
		if i == 99 {
			return "", fmt.Errorf("could not create new random note, too many files for today")
		}
	}

	file, err := os.Create(filename_full)
	if err != nil {
		return "", err
	}
	file.Close()

	if err := AddMetadataMarkdownNote(b.RandomNotesPath, filename, "icon", "\"\""); err != nil {
		return "", fmt.Errorf("could not add icon to new random note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.RandomNotesPath, filename, "tags", []string{}); err != nil {
		return "", fmt.Errorf("could not add tags to new random note: %v", err)
	}

	if err := AddMetadataMarkdownNote(b.RandomNotesPath, filename, "created", now.Format(time.RFC3339)); err != nil {
		return "", fmt.Errorf("could not add created metadata to new random note: %v", err)
	}

	return filename, nil
}

func (b Backend) Trash_RandomNote(filename string) error {
	filename_full := filepath.Join(b.RandomNotesPath, filename)
	if _, err := os.Stat(filename_full); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename_full)
	}
	if err := AddMetadataMarkdownNote(b.RandomNotesPath, filename, "trashed", time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("could not add trashed metadata to random note: %v", err)
	}
	for i := 1; i < 100; i++ {
		filename_trash_full := filepath.Join(b.TrashPath, filename)
		if i > 1 {
			ext := filepath.Ext(filename_trash_full)
			filename_trash_full = filename_trash_full[:len(filename_trash_full)-len(ext)]
			filename_trash_full = fmt.Sprintf("%s.v%02d%s", filename_trash_full, i, ext)
		}
		if _, err := os.Stat(filename_trash_full); os.IsNotExist(err) {
			if err := os.Rename(filename_full, filename_trash_full); err != nil {
				return fmt.Errorf("could not move file to trash: %v", err)
			} else {
				return nil
			}
		}
	}
	return fmt.Errorf("could not move file to trash, too many versions already exist: %s", filename)
}
