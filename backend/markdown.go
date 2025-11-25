package backend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func ParseMarkdownNote(path, name string) (string, string, []string, time.Time, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return "", "", nil, time.Time{}, err
	}
	defer file.Close()

	var (
		icon    string
		tags    []string
		title   string
		created time.Time
	)

	scanner := bufio.NewScanner(file)
	var metadata strings.Builder
	inMetadata := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "+++" {
			if !inMetadata {
				inMetadata = true
				continue
			} else {
				inMetadata = false
				continue
			}
		}

		if inMetadata {
			metadata.WriteString(line + "\n")
			continue
		}

		if strings.HasPrefix(line, "# ") && title == "" {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return "", "", nil, time.Time{}, err
	}

	metaDataMap := make(map[string]any)
	if err := toml.Unmarshal([]byte(metadata.String()), &metaDataMap); err != nil {
		return "", "", nil, time.Time{}, err
	}

	if val, ok := metaDataMap["icon"].(string); ok {
		icon = val
	}
	if val, ok := metaDataMap["tags"].([]any); ok {
		for _, tag := range val {
			if tagStr, ok := tag.(string); ok {
				tags = append(tags, tagStr)
			}
		}
	}
	if val, ok := metaDataMap["created"].(time.Time); ok {
		created = val
	}

	if icon == "" {
		icon = "󰎞"
	}

	if title == "" {
		title = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return title, icon, tags, created, nil
}

func ParseMarkdownTasklist(path, name string) (string, []string, []bool, []int, []time.Time, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return "", nil, nil, nil, nil, err
	}
	defer file.Close()

	var (
		title       string
		tasks       []string
		checked     []bool
		importances []int
		dueDates    []time.Time
	)

	scanner := bufio.NewScanner(file)
	inMetadata := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "+++" {
			inMetadata = !inMetadata
			continue
		}

		if strings.HasPrefix(line, "# ") && title == "" {
			title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}

		pattern := regexp.MustCompile(`^- \[( |x)\] (.*?)( \{(\!*)\})?( \{([0-9]{4}-[0-9]{2}-[0-9]{2})?\})?$`)
		matches := pattern.FindStringSubmatch(line)

		var (
			taskText   string
			isChecked  bool
			importance int
			dueDate    time.Time
		)
		if len(matches) > 3 {
			if matches[4] != "" {
				importance = len(matches[4])
			}
		}
		if len(matches) > 5 {
			if matches[6] != "" {
				dueDate, _ = time.Parse("2006-01-02", matches[6])
			}
		}
		if len(matches) > 2 {
			taskText = matches[2]
			isChecked = matches[1] == "x" || matches[1] == "X"

			tasks = append(tasks, taskText)
			checked = append(checked, isChecked)
			importances = append(importances, importance)
			dueDates = append(dueDates, dueDate)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, nil, nil, nil, err
	}

	if title == "" {
		title = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return title, tasks, checked, importances, dueDates, nil
}

func WriteMarkdownTasklist(path, name, title string, tasks []string, checked []bool, importances []int, dueDates []time.Time) error {
	file, err := os.Create(filepath.Join(path, name))
	if err != nil {
		return fmt.Errorf("could not create markdown tasklist: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if _, err := fmt.Fprintf(writer, "# %s\n\n", title); err != nil {
		return fmt.Errorf("could not write title to markdown tasklist: %v", err)
	}

	for i, task := range tasks {
		checkMark := " "
		if checked[i] {
			checkMark = "x"
		}
		importanceStr := ""
		if importances[i] > 0 {
			importanceStr = " {" + strings.Repeat("!", importances[i]) + "}"
		}
		dueDateStr := ""
		if !dueDates[i].IsZero() {
			dueDateStr = " {" + dueDates[i].Format("2006-01-02") + "}"
		}
		if _, err := fmt.Fprintf(writer, "- [%s] %s%s%s\n", checkMark, task, importanceStr, dueDateStr); err != nil {
			return fmt.Errorf("could not write task to markdown tasklist: %v", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("could not flush markdown tasklist writer: %v", err)
	}

	return nil
}

func ReadMarkdownNote(path, name string) (string, error) {
	file, err := os.Open(filepath.Join(path, name))
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	content := string(data)

	start := strings.Index(content, "+++")
	if start != -1 {
		end := strings.Index(content[start+3:], "+++")
		if end != -1 {
			content = content[:start] + content[start+3+end+3:]
		}
	}

	return strings.TrimSpace(content), nil
}

func AddTitleMarkdownNote(path, name, title string) error {
	file, err := os.OpenFile(filepath.Join(path, name), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "\n# %s\n", title); err != nil {
		return fmt.Errorf("could not write title to markdown note: %v", err)
	}
	return nil
}

func AddMetadataMarkdownNote(path, name, metadata string, value any) error {
	file, err := os.OpenFile(filepath.Join(path, name), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}
	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil && err != io.EOF {
		return err
	}

	parts := strings.SplitN(string(content), "+++", 3)
	var metadataBlock, body string
	if len(parts) == 3 {
		metadataBlock = strings.TrimSpace(parts[1])
		body = strings.TrimSpace(parts[2])
	} else {
		body = string(content)
	}

	metadataMap := make(map[string]string)
	if metadataBlock != "" {
		for line := range strings.SplitSeq(metadataBlock, "\n") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				metadataMap[key] = val
			}
		}
	}

	metadataMap[metadata] = fmt.Sprintf("%v", value)

	var newMetadataBlock strings.Builder
	newMetadataBlock.WriteString("+++\n")

	keys := make([]string, 0, len(metadataMap))
	for key := range metadataMap {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		val := metadataMap[key]
		newMetadataBlock.WriteString(fmt.Sprintf("%s = %s\n", key, val))
	}
	newMetadataBlock.WriteString("+++\n\n")

	newContent := newMetadataBlock.String() + body

	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(newContent))
	if err != nil {
		return err
	}

	return nil
}
